package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middleware"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {
	{
		auths := r.Group("/auth")
		auths.POST("/register", controllers.Register)
		auths.POST("/login", controllers.Login)
		auths.POST("/refresh", controllers.RefreshToken)
		auths.POST("/logout", middleware.AuthMiddleware(), controllers.Logout)
	}
}
