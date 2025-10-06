package middlewares

import (
	"net/http"
	"strings"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/utils"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authorization header required"})
			c.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid authorization format"})
			c.Abort()
			return
		}

		tokenStr := parts[1]

		var revoked models.RevokedToken
		if err := config.DB.Where("token = ?", tokenStr).First(&revoked).Error; err == nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token revoked"})
			c.Abort()
			return
		}

		claims := &utils.Claims{}

		token, err := jwt.ParseWithClaims(tokenStr, claims, func(token *jwt.Token) (any, error) {
			return utils.AccessSecret, nil
		})

		if err != nil || !token.Valid {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid or expired token"})
			c.Abort()
			return
		}

		if claims.ExpiresAt.Time.Before(time.Now()) {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Token expired"})
			c.Abort()
			return
		}

		c.Set("userId", claims.UserID)
		c.Set("tokenString", tokenStr)
		c.Set("claims", claims)
		c.Next()
	}
}
