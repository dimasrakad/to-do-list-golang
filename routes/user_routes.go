package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func UserRoute(r *gin.RouterGroup) {
	users := r.Group("/users")
	users.Use(middlewares.AuthMiddleware())
	{
		users.GET("", controllers.GetUserNames)
	}
}
