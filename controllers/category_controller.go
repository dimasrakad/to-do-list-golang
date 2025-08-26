package controllers

import (
	"net/http"
	"strconv"
	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"github.com/gin-gonic/gin"
)

func GetCategories(c *gin.Context) {
	var categories []models.Category

	config.DB.Model(&models.Category{}).Find(&categories)

	c.JSON(http.StatusOK, gin.H{"data": categories})
}

func CreateCategory(c *gin.Context) {
	var input struct {
		Name            string `json:"name" binding:"required"`
		CategoryColorID uint   `json:"categoryColorId" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var categoryColor models.CategoryColor

	if err := config.DB.Model(&models.CategoryColor{}).Where("id = ?", input.CategoryColorID).First(&categoryColor).Error; err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Category color not found"})
		return
	}

	category := models.Category{
		Name:            input.Name,
		CategoryColorID: input.CategoryColorID,
	}

	if err := config.DB.Create(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Preload("CategoryColor").First(&category, category.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": category})
}

func UpdateCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var input struct {
		Name            *string `json:"name"`
		CategoryColorID *uint   `json:"categoryColorId"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if input.Name != nil {
		category.Name = *input.Name
	}
	if input.CategoryColorID != nil {
		var categoryColor models.CategoryColor

		if err := config.DB.Model(&models.CategoryColor{}).Where("id = ?", *input.CategoryColorID).First(&categoryColor).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Category color not found"})
			return
		}

		category.CategoryColorID = *input.CategoryColorID
	}

	if err := config.DB.Save(&category).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	if err := config.DB.Preload("CategoryColor").First(&category, category.ID).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"data": category})
}

func DeleteCategory(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))

	var category models.Category
	if err := config.DB.First(&category, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Category not found"})
		return
	}

	if err := config.DB.Delete(&models.Category{}, id); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Category deleted"})
}
