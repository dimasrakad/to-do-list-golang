package controllers

import (
	"net/http"
	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"github.com/gin-gonic/gin"
)

func GetColors(c *gin.Context) {
	var colors []models.CategoryColor

	config.DB.Model(&models.CategoryColor{}).Find(&colors)

	c.JSON(http.StatusOK, gin.H{"data": colors})
}
