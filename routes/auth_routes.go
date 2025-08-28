package routes

import (
	"to-do-list-golang/controllers"

	"github.com/gin-gonic/gin"
)

func AuthRoute(r *gin.Engine) {
	{
		auths := r.Group("/auth")
		auths.POST("/register", controllers.Register)
		auths.POST("/login", controllers.Login)
		auths.POST("/refresh", controllers.RefreshAccessToken)
	}
}
