package controllers

import (
	"net/http"
	"strconv"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// "Get Categories" godoc
// @Summary "Get Categories"
// @Description Get all todo categories
// @Tags Category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /categories [get]
func GetCategories(c *gin.Context) {
	var categories []models.Category

	query := config.DB.Model(&models.Category{})

	if err := categoryWithRelations(query).Find(&categories).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res := dtos.SuccessResponse{
		Data:    categories,
		Message: "",
	}
	c.JSON(http.StatusOK, res)
}

// "Create Category" godoc
// @Summary "Create Category"
// @Description Create new todo category
// @Tags Category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Param Payload body dtos.CreateCategoryRequest true "Create category input"
// @Success 201 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /categories [post]
func CreateCategory(c *gin.Context) {
	var input dtos.CreateCategoryRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	category := models.Category{
		Name:            input.Name,
		CategoryColorID: input.CategoryColorID,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := config.DB.Preload("CategoryColor").First(&category, category.ID).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res := dtos.SuccessResponse{
		Data:    category,
		Message: "",
	}
	c.JSON(http.StatusCreated, res)
}

// "Update Category" godoc
// @Summary "Update Category"
// @Description Update existing todo category
// @Tags Category
// @Accept json
// @Produce json
// @Param ID path uint true "Category id"
// @Param Authorization header string true "Bearer token"
// @Param Payload body dtos.UpdateCategoryRequest true "Update category input"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /categories [patch]
func UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input dtos.UpdateCategoryRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.CategoryColorID != nil {
		category.CategoryColorID = *input.CategoryColorID
	}

	if err := config.DB.Save(&category).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := config.DB.Preload("CategoryColor").First(&category, category.ID).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res := dtos.SuccessResponse{
		Data:    category,
		Message: "",
	}
	c.JSON(http.StatusOK, res)
}

// "Delete Category" godoc
// @Summary "Delete Category"
// @Description Delete existing todo category
// @Tags Category
// @Accept json
// @Produce json
// @Param ID path uint true "Category id"
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 404 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /categories [delete]
func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err := config.DB.Delete(&models.Category{}, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	res := dtos.SuccessResponse{
		Data:    nil,
		Message: "Category deleted",
	}
	c.JSON(http.StatusOK, res)
}

func categoryWithRelations(db *gorm.DB) *gorm.DB {
	return db.Preload("CategoryColor")
}
