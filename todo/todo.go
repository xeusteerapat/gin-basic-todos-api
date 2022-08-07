package todo

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/xeusteerapat/gin-basic-todos-api/auth"
	"gorm.io/gorm"
)

type Todo struct {
	Title string `json:"text"`
	gorm.Model
}

func (Todo) TableName() string {
	return "todos"
}

type TodoHandler struct {
	db *gorm.DB
}

func NewTodoHandler(db *gorm.DB) *TodoHandler {
	return &TodoHandler{db: db}
}

func (t *TodoHandler) NewTask(ctx *gin.Context) {
	// Validate token
	bearerToken := ctx.Request.Header.Get("Authorization")
	tokenString := strings.TrimPrefix(bearerToken, "Bearer ")

	if err := auth.Protect(tokenString); err != nil {
		ctx.AbortWithStatus(http.StatusUnauthorized) // If unauthorized, then no need to go to next handler
		return
	}

	var todo Todo
	if err := ctx.ShouldBindJSON(&todo); err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	r := t.db.Create(&todo)
	if err := r.Error; err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusCreated, gin.H{
		"ID": todo.ID,
	})
}
