package todo

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
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

func (t *TodoHandler) List(ctx *gin.Context) {
	var todos []Todo

	result := t.db.Find(&todos)

	if err := result.Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, todos)
}

func (t *TodoHandler) Remove(ctx *gin.Context) {
	todoIdParam := ctx.Param("id")

	todoId, err := strconv.Atoi(todoIdParam)

	if err != nil {
		ctx.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})

		return
	}

	result := t.db.Delete(&Todo{}, todoId)

	if err := result.Error; err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})

		return
	}

	ctx.JSON(http.StatusOK, gin.H{
		"status": "success remove",
	})
}
