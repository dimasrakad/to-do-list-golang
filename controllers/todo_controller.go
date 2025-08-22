package controllers

import (
	"net/http"
	"strconv"
	"time"

	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo

	query := config.DB.Model(&models.Todo{})

	// filter by status
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}

	// filter by priority
	if priority := c.Query("priority"); priority != "" {
		query = query.Where("priority = ?", priority)
	}

	// filter by due date
	if dueDate := c.Query("dueDate"); dueDate != "" {
		query = query.Where("DATE(due) = ?", dueDate)
	}

	// filter by due range
	if dueFrom := c.Query("dueFrom"); dueFrom != "" {
		query = query.Where("DATE(due) >= ?", dueFrom)
	}
	if dueTo := c.Query("dueTo"); dueTo != "" {
		query = query.Where("DATE(due) <= ?", dueTo)
	}

	if err := query.Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, todos)
}

func CreateTodo(c *gin.Context) {
	cfg := config.LoadConfig()

	var input struct {
		Title       string `json:"title" binding:"required"`
		Priority    string `json:"priority" binding:"oneof=low medium high"`
		Description string `json:"description"`
		Due         string `json:"due"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	loc, _ := time.LoadLocation(cfg.AppLocation)
	due, err := time.ParseInLocation("2006-01-02 15:04:05", input.Due, loc)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD hh:mm:ss"})
		return
	}

	if due.Before(time.Now()) {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Due cannot be in the past"})
		return
	}

	todo := models.Todo{
		Title:    input.Title,
		Priority: input.Priority,
		Due:      &due,
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, todo)
}

func UpdateTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	var todo models.Todo

	if err := config.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	config.DB.Save(&todo)
	c.JSON(http.StatusOK, todo)
}

func DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
