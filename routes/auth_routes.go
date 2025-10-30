package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.RouterGroup) {
	auths := r.Group("/auth")
	{
		auths.POST("/register", controllers.Register)
		auths.GET("/verify", controllers.VerifyEmail)
		// auths.GET("/resend-verification", controllers.ResendVerificationEmail)
		auths.POST("/login", controllers.Login)
		auths.POST("/refresh", controllers.RefreshToken)
		auths.POST("/logout", middlewares.AuthMiddleware(), controllers.Logout)
	}
}
