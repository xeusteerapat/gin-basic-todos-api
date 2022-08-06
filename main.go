package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

type User struct {
	gorm.Model // gorm will add ID, CreatedAt, UpdatedAt, DeletedAt for us
	Name       string
}

var db *gorm.DB

func main() {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{})
	if err != nil {
		panic("Database connection failed!")
	}

	db.AutoMigrate(&User{})

	db.Create(&User{
		Name: "Teerapat",
	})

	r := gin.Default()

	usersHandler := UsersHandler{db: db}
	r.GET("/users", usersHandler.User)

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
