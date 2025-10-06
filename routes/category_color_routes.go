package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryColorRoute(r *gin.Engine) {
	colors := r.Group("/category-colors")
	colors.Use(middlewares.AuthMiddleware())
	{
		colors.GET("", controllers.GetColors)
	}
}
