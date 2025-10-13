package controllers

import (
	"net/http"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
)

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
