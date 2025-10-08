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
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// "Get Todos" godoc
// @Summary "Get Todos"
// @Description Get all todos
// @Tags Todo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param search query string true "Todo search filter"
// @Param status query string true "Todo status filter"
// @Param priority query string true "Todo priority filter"
// @Param category query string true "Todo category filter"
// @Param dueDate query string true "Todo due date filter"
// @Param dueFrom query string true "Todo due from filter"
// @Param dueTo query string true "Todo due to filter"
// @Param sortBy query string true "Sort todo by field(s)"
// @Param order query string true "Order sort by asc/desc"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /todos [get]
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

	// filter by category
	if category := c.Query("category"); category != "" {
		query = query.Joins("Category").Where("categories.name = ?", category)
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

	if err := todoWithRelations(query).Offset(offset).Limit(limit).Find(&todos).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res := dtos.PaginationResponse{
		Data:       todos,
		Page:       page,
		Limit:      limit,
		Total:      total,
		TotalPages: int(math.Ceil(float64(total) / float64(limit))),
	}
	c.JSON(http.StatusOK, res)
}

// "Create Todo" godoc
// @Summary "Create Todo"
// @Description Create new todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param Payload body dtos.CreateTodoRequest true "Create category input"
// @Success 201 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /todos [post]
func CreateTodo(c *gin.Context) {
	cfg := config.LoadConfig()
	var res any

	var input dtos.CreateTodoRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	loc, _ := time.LoadLocation(cfg.AppLocation)
	due, err := time.ParseInLocation("2006-01-02 15:04:05", input.Due, loc)

	if err != nil {
		res = dtos.ErrorResponse{
			Error: "Invalid date format. Use YYYY-MM-DD hh:mm:ss",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	if due.Before(time.Now()) {
		res = dtos.ErrorResponse{
			Error: "Due cannot be in the past",
		}
		c.JSON(http.StatusBadRequest, res)
		return
	}

	todo := models.Todo{
		Title:      input.Title,
		Priority:   input.Priority,
		Due:        &due,
		CategoryID: input.CategoryID,
	}

	if err := config.DB.Create(&todo).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := todoWithRelations(config.DB).First(&todo, todo.ID).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res = dtos.SuccessResponse{
		Data:    todo,
		Message: "",
	}
	c.JSON(http.StatusCreated, res)
}

// "Update Todo" godoc
// @Summary "Update Todo"
// @Description Update existing todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param ID path uint true "Todo id"
// @Param Authorization header string true "Bearer token"
// @Param Payload body dtos.UpdateTodoRequest true "Update todo input"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /todos [patch]
func UpdateTodo(c *gin.Context) {
	cfg := config.LoadConfig()
	id, _ := strconv.Atoi(c.Param("id"))
	var res any

	var input dtos.UpdateTodoRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	var todo models.Todo
	if err := config.DB.First(&todo, id).Error; err != nil {
		utils.HandleDBError(c, err)
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
			res := dtos.ErrorResponse{
				Error: "Invalid date format. Use YYYY-MM-DD hh:mm:ss",
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		if due.Before(time.Now()) {
			res := dtos.ErrorResponse{
				Error: "Due cannot be in the past",
			}
			c.JSON(http.StatusBadRequest, res)
			return
		}

		todo.Due = &due
	}

	if err := config.DB.Save(&todo).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := todoWithRelations(config.DB).First(&todo, todo.ID).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res = dtos.SuccessResponse{
		Data:    todo,
		Message: "",
	}
	c.JSON(http.StatusOK, res)
}

// "Delete Todo" godoc
// @Summary "Delete Todo"
// @Description Delete existing todo
// @Tags Todo
// @Accept json
// @Produce json
// @Param ID path uint true "Todo id"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /todos [delete]
func DeleteTodo(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var todo models.Todo
	if err := config.DB.First(&todo, id).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := config.DB.Delete(&models.Todo{}, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	res := dtos.SuccessResponse{
		Data:    nil,
		Message: "Todo deleted",
	}
	c.JSON(http.StatusOK, res)
}

func todoWithRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("Category").Preload("Category.CategoryColor")
}
