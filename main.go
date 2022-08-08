package main

import (
	"log"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/xeusteerapat/gin-basic-todos-api/auth"
	"github.com/xeusteerapat/gin-basic-todos-api/todo"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // gorm will add ID, CreatedAt, UpdatedAt, DeletedAt for us
	Name       string
}

var db *gorm.DB

func main() {
	err := godotenv.Load("local.env")
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

	r.Run()
}

type UsersHandler struct {
	db *gorm.DB
}

func (usersHandler *UsersHandler) User(ctx *gin.Context) {
	var u User
	usersHandler.db.First(&u)

	ctx.JSON(200, u)
}
