package controllers

import (
	"net/http"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
)

// "Get Category Colors" godoc
// @Summary "Get Category Colors"
// @Description Get all category colors
// @Tags Category
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /category-colors [get]
func GetColors(c *gin.Context) {
	var colors []models.CategoryColor

	if err := config.DB.Model(&models.CategoryColor{}).Find(&colors).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	res := dtos.SuccessResponse{
		Data:    colors,
		Message: "",
	}
	c.JSON(http.StatusOK, res)
}
