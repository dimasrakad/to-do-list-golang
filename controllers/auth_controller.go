package controllers

import (
	"net/http"
	"strconv"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/models/dtos"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

// Register godoc
// @Summary Register a new user
// @Description Register a new user with name, email, and password
// @Tags Auth
// @Accept json
// @Produce json
// @Param Payload body dtos.RegisterRequest true "Register input"
// @Success 201 {object} dtos.SuccessResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /auth/register [post]
func Register(c *gin.Context) {
	var input dtos.RegisterRequest
	var err error
	var res any

	if !utils.ValidateInput(c, &input) {
		return
	}

	hashedPassword, err := utils.HashPassword(input.Password)
	if err != nil {
		res = dtos.ErrorResponse{
			Error: "Error while hashing password: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	verifyToken, err := utils.GenerateRandomToken(32)
	if err != nil {
		res = dtos.ErrorResponse{
			Error: "Error while generating verification token: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	user := models.User{
		Name:        input.Name,
		Email:       input.Email,
		Password:    hashedPassword,
		VerifyToken: &verifyToken,
	}

	if err = config.DB.Create(&user).Error; err != nil {
		utils.HandleDBError(c, err)
		return
	}

	if err = utils.SendVerificationEmail(user); err != nil {
		res = dtos.ErrorResponse{
			Error: "Error while sending verification email: " + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res = dtos.SuccessResponse{
		Data:    user,
		Message: "",
	}

	c.JSON(http.StatusCreated, res)
}

// Verify Email godoc
// @Summary Verify user's email
// @Description Verifies a user's email using a verification token sent via email.
// @Tags Auth
// @Accept text/html
// @Produce text/html
// @Param token query string true "Verification token"
// @Success 200 {string} string "Email verified successfully (HTML page)"
// @Failure 500 {string} string "Internal server error (HTML page)"
// @Router /auth/verify [get]
func VerifyEmail(c *gin.Context) {
	var token = c.Query("token")

	if token == "" {
		c.HTML(http.StatusInternalServerError, "failed-verify.html", nil)
		return
	}

	var user models.User
	if err := config.DB.Where("verify_token = ?", token).First(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "failed-verify.html", nil)
		return
	}

	user.IsVerified = true
	user.VerifyToken = nil
	now := time.Now()
	user.VerifiedAt = &now
	if err := config.DB.Save(&user).Error; err != nil {
		c.HTML(http.StatusInternalServerError, "failed-verify.html", nil)
		return
	}

	c.HTML(http.StatusOK, "success-verify.html", nil)
}

// Login godoc
// @Summary Login
// @Description Login with email and password to receive access and refresh tokens
// @Tags Auth
// @Accept json
// @Produce json
// @Param Payload body dtos.LoginRequest true "Login input"
// @Success 200 {object} dtos.TokenResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /auth/login [post]
func Login(c *gin.Context) {
	cfg := config.LoadConfig()
	var res any

	var input dtos.LoginRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	var user models.User
	if err := config.DB.Where("email = ?", input.Email).First(&user).Error; err != nil {
		res = dtos.ErrorResponse{
			Error: "Invalid email or password",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if !utils.CheckPasswordHash(input.Password, user.Password) {
		res = dtos.ErrorResponse{
			Error: "Invalid email or password",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if !user.IsVerified {
		res = dtos.ErrorResponse{
			Error: "Email not verified",
		}
		c.JSON(http.StatusUnauthorized, res)
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
		res = dtos.ErrorResponse{
			Error: "Could not save refresh token\n" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	res = dtos.TokenResponse{
		AccessToken:  accessToken,
		RefreshToken: refreshToken,
	}
	c.JSON(http.StatusOK, res)
}

// Refresh Token godoc
// @Summary Refresh token
// @Description Refresh access token & refresh token using a valid refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Payload body dtos.RefreshTokenRequest true "Refresh token input"
// @Success 200 {object} dtos.TokenResponse
// @Failure 400 {object} dtos.ErrorResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /auth/refresh [post]
func RefreshToken(c *gin.Context) {
	cfg := config.LoadConfig()

	var res any

	var input dtos.RefreshTokenRequest

	if !utils.ValidateInput(c, &input) {
		return
	}

	claims := &utils.Claims{}
	token, err := jwt.ParseWithClaims(input.RefreshToken, claims, func(token *jwt.Token) (any, error) {
		return utils.RefreshSecret, nil
	})

	if err != nil || !token.Valid {
		res = dtos.ErrorResponse{
			Error: "Invalid refresh token",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	userID := claims.UserID

	var refreshToken models.RefreshToken
	if err := config.DB.Where("user_id = ? AND token = ?", userID, input.RefreshToken).First(&refreshToken).Error; err != nil {
		res = dtos.ErrorResponse{
			Error: "Refresh token not found",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if refreshToken.ExpiresAt.Before(time.Now()) {
		config.DB.Delete(&refreshToken)
		res = dtos.ErrorResponse{
			Error: "Refresh token expired",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	config.DB.Delete(&refreshToken)

	newAccessToken, err := utils.GenerateAccessToken(userID)
	if err != nil {
		res = dtos.ErrorResponse{
			Error: "Could not generate new access token\n" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	newRefreshToken, err := utils.GenerateRefreshToken(userID)
	if err != nil {
		res = dtos.ErrorResponse{
			Error: "Could not generate new refresh token\n" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	expMinutes, _ := time.ParseDuration(cfg.JWTRefreshExpire + "m")

	newToken := models.RefreshToken{
		UserID:    userID,
		Token:     newRefreshToken,
		ExpiresAt: time.Now().Add(expMinutes),
	}

	config.DB.Create(&newToken)

	res = dtos.TokenResponse{
		AccessToken:  newAccessToken,
		RefreshToken: newRefreshToken,
	}
	c.JSON(http.StatusOK, res)
}

// Logout godoc
// @Summary Logout
// @Description Logout user by revoking the current access token and deleting the refresh token
// @Tags Auth
// @Accept json
// @Produce json
// @Param Authorization header string true "Bearer token"
// @Success 200 {object} dtos.SuccessResponse
// @Failure 401 {object} dtos.ErrorResponse
// @Failure 500 {object} dtos.ErrorResponse
// @Router /auth/logout [post]
func Logout(c *gin.Context) {
	userId, userIdExists := c.Get("userId") // from JWT claims
	tokenString, tokenStringExists := c.Get("tokenString")
	claims, claimsExists := c.Get("claims")
	var res any

	if !userIdExists || !tokenStringExists || !claimsExists {
		res = dtos.ErrorResponse{
			Error: "Unauthorized",
		}
		c.JSON(http.StatusUnauthorized, res)
		return
	}

	if err := config.DB.Where("user_id = ?", userId).Delete(&models.RefreshToken{}).Error; err != nil {
		res = dtos.ErrorResponse{
			Error: "Failed to logout\n" + err.Error(),
		}
		c.JSON(http.StatusInternalServerError, res)
		return
	}

	expTime := claims.(*utils.Claims).ExpiresAt.Time

	revoked := models.RevokedToken{
		Token:     tokenString.(string),
		UserID:    userId.(uint),
		ExpiresAt: expTime,
	}

	config.DB.Create(&revoked)

	res = dtos.SuccessResponse{
		Data:    nil,
		Message: "Logout success",
	}
	c.JSON(http.StatusOK, res)
}
