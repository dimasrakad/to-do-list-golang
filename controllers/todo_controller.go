package controllers

import (
	"fmt"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"

	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"github.com/gin-gonic/gin"
)

func GetTodos(c *gin.Context) {
	var todos []models.Todo

	query := config.DB.Model(&models.Todo{})

	// search title
	if search := c.Query("search"); search != "" {
		query = query.Where("LOWER(title) LIKE ?", "%"+strings.ToLower(search)+"%")
	}

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

	// sorting
	sortBy := strings.Split(c.DefaultQuery("sortBy", "due"), ",")
	orders := strings.Split(c.DefaultQuery("order", "asc"), ",")

	for i, field := range sortBy {
		field = strings.TrimSpace(field)

		order := "asc"
		if i < len(orders) {
			if o := strings.ToLower(strings.TrimSpace(orders[i])); o == "desc" {
				order = "desc"
			}

		}

		switch field {
		case "priority":
			query = query.Order(fmt.Sprintf(`
			CASE
				WHEN priority = 'low' THEN 1
				WHEN priority = 'medium' THEN 2
				WHEN priority = 'high' THEN 3
			END %s`, order))
		case "status":
			query = query.Order(fmt.Sprintf(`
			CASE
				WHEN status = 'pending' THEN 1
				WHEN status = 'in progress' THEN 2
				WHEN status = 'done' THEN 3
			END %s`, order))
		default:
			query = query.Order(fmt.Sprintf("%s %s", field, order))
		}
	}

	// pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	if page < 1 {
		page = 1
	}
	if limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var total int64
	query.Count(&total)

	if err := query.Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":       todos,
		"page":       page,
		"limit":      limit,
		"total":      total,
		"totalPages": int(math.Ceil(float64(total) / float64(limit))),
	})
}

func CreateTodo(c *gin.Context) {
	cfg := config.LoadConfig()

	var input struct {
		Title       string `json:"title" binding:"required"`
		Priority    string `json:"priority" binding:"oneof=low medium high"`
		Description string `json:"description"`
		Due         string `json:"due" binding:"required"`
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

	c.JSON(http.StatusCreated, gin.H{"data": todo})
}

func UpdateTodo(c *gin.Context) {
	cfg := config.LoadConfig()
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Title       *string `json:"title"`
		Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
		Status      *string `json:"status" binding:"omitempty,oneof='pending' 'in progress' 'done'"`
		Description *string `json:"description"`
		Due         *string `json:"due"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var todo models.Todo
	if err := config.DB.First(&todo, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
		return
	}

	if input.Title != nil {
		todo.Title = *input.Title
	}
	if input.Priority != nil {
		todo.Priority = *input.Priority
	}
	if input.Status != nil {
		todo.Status = *input.Status
	}
	if input.Description != nil {
		todo.Description = *input.Description
	}

	if input.Due != nil {
		loc, _ := time.LoadLocation(cfg.AppLocation)
		due, err := time.ParseInLocation("2006-01-02 15:04:05", *input.Due, loc)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD hh:mm:ss"})
			return
		}

		if due.Before(time.Now()) {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Due cannot be in the past"})
			return
		}

		todo.Due = &due
	}

	if err := config.DB.Save(&todo).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": todo})
}

func DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	if err := config.DB.Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Todo deleted"})
}
