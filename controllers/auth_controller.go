package controllers

import (
	"net/http"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
)

func Register(c *gin.Context) {
	var input struct {
		Name     string `json:"name" binding:"required"`
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required,password"`
	}

	if !utils.ValidateInput(c, &input) {
		return
	}

	hashedPassword, _ := utils.HashPassword(input.Password)

	user := models.User{
		Name:     input.Name,
		Email:    input.Email,
		Password: hashedPassword,
	}

	if err := config.DB.Create(&user).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	c.JSON(http.StatusCreated, gin.H{"data": user})
}

func Login(c *gin.Context) {
	var input struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if !utils.ValidateInput(c, &input) {
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	token, _ := utils.GenerateJWT(user.ID)

	c.JSON(http.StatusOK, gin.H{
		"token": token,
		"user":  user,
	})
}
