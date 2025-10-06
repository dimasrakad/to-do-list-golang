package controllers

import (
	"net/http"
	"strconv"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
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
	cfg := config.LoadConfig()

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

	accessToken, _ := utils.GenerateAccessToken(user.ID)
	refreshToken, _ := utils.GenerateRefreshToken(user.ID)

	refreshExpire, _ := strconv.Atoi(cfg.JWTRefreshExpire)

	refreshTokenDB := models.RefreshToken{
		UserID:    user.ID,
		Token:     refreshToken,
		ExpiresAt: time.Now().Add(time.Duration(refreshExpire) * time.Minute),
	}

	if err := config.DB.Create(&refreshTokenDB).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not save refresh token\n" + err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}

func RefreshToken(c *gin.Context) {
	cfg := config.LoadConfig()

	var input struct {
		RefreshToken string `json:"refreshToken" binding:"required"`
	}

	if !utils.ValidateInput(c, &input) {
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(input.RefreshToken, claims, func(token *jwt.Token) (any, error) {
		return utils.RefreshSecret, nil
	})

	if err != nil || !token.Valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid refresh token"})
		return
	}

	userID := claims.UserID

	var refreshToken models.RefreshToken
	if err := config.DB.Where("user_id = ? AND token = ?", userID, input.RefreshToken).First(&refreshToken).Error; err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Refresh token not found"})
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		config.DB.Delete(&refreshToken)
		c.JSON(http.StatusUnauthorized, gin.H{"error": "refresh token expired"})
		return
	}

	config.DB.Delete(&refreshToken)

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new access token\n" + err.Error()})
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not generate new refresh token\n" + err.Error()})
		return
	}

	expMinutes, _ := time.ParseDuration(cfg.JWTRefreshExpire + "m")

	newToken := models.RefreshToken{
		UserID:    userID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(expMinutes),
	}

	config.DB.Create(&newToken)

	c.JSON(http.StatusOK, gin.H{
		"accessToken":  newAccessToken,
		"refreshToken": newRefreshToken,
	})
}

func Logout(c *gin.Context) {
	userId, userIdExists := c.Get("userId") // from JWT claims
	tokenString, tokenStringExists := c.Get("tokenString")
	claims, claimsExists := c.Get("claims")

	if !userIdExists || !tokenStringExists || !claimsExists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	if err := config.DB.Where("user_id = ?", userId).Delete(&models.RefreshToken{}).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to logout\n" + err.Error()})
		return
	}

	expTime := claims.(*utils.Claims).ExpiresAt.Time

	revoked := models.RevokedToken{
		Token:     tokenString.(string),
		UserID:    userId.(uint),
		ExpiresAt: expTime,
	}

	config.DB.Create(&revoked)

	c.JSON(http.StatusOK, gin.H{"message": "Logout success"})
}
