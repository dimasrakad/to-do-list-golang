package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func CategoryRoute(r *gin.Engine) {
	categories := r.Group("/categories")
	categories.Use(middlewares.AuthMiddleware())
	{
		categories.GET("", controllers.GetCategories)
		categories.POST("", controllers.CreateCategory)
		categories.PATCH("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}
}
