package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xeusteerapat/gin-basic-todos-api/auth"
	"github.com/xeusteerapat/gin-basic-todos-api/todo"
	"golang.org/x/time/rate"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	buildCommit = "dev"
	buildTime   = time.Now().String()
)

type User struct {
	gorm.Model // gorm will add ID, CreatedAt, UpdatedAt, DeletedAt for us
	Name       string
}

var db *gorm.DB

func main() {
	// Liveness probe
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Printf("Please consider environment variables: %s", err)
	}

	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Database connection failed!")
	}

	db.AutoMigrate(&User{})
	db.AutoMigrate(&todo.Todo{})

	db.Create(&User{
		Name: "Teerapat",
	})

	r := gin.Default()
	r.GET("/limit", limitedHandler)
	r.GET("/x", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"buildCommit": buildCommit,
			"buildTime":   buildTime,
		})
	})

	// health check
	r.GET("/health-check", func(ctx *gin.Context) {
		ctx.JSON(200, gin.H{
			"Status": "OK",
		})
	})

	r.GET("/token", auth.AccessToken(os.Getenv("SIGN")))

	protectedRoute := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	usersHandler := UsersHandler{db: db}
	r.GET("/users", usersHandler.User)

	handler := todo.NewTodoHandler(db)

	protectedRoute.POST("/todos", handler.NewTask)

	// Handle Graceful shutdown
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	server := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   20 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("Shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

type UsersHandler struct {
	db *gorm.DB
}

func (usersHandler *UsersHandler) User(ctx *gin.Context) {
	var u User
	usersHandler.db.First(&u)

	ctx.JSON(200, u)
}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(ctx *gin.Context) {
	if !limiter.Allow() {
		ctx.AbortWithStatus(http.StatusTooManyRequests)
		return
	}

	ctx.JSON(200, gin.H{
		"messaage": "pong",
	})
}
