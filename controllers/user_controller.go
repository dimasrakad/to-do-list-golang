package controllers

import (
	"net/http"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
)

// "Get User Names" godoc
// @Summary "Get User Names"
// @Description Get all user names
// @Tags User
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /users [get]
func GetUserNames(c *gin.Context) {
	var users []models.User

	if err := config.DB.Model(&models.User{}).Find(&users).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	var userNames []string

	for _, v := range users {
		userNames = append(userNames, v.Name)
	}

	res := dtos.SuccessResponse{
		Data:    userNames,
		Message: "",
	}
	c.JSON(http.StatusOK, res)
}
