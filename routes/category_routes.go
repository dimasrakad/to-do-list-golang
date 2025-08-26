package routes

import (
	"to-do-list-golang/controllers"

	"github.com/gin-gonic/gin"
)

func CategoryRoute(r *gin.Engine) {
	{
		categories := r.Group("/categories")
		categories.GET("", controllers.GetCategories)
		categories.POST("", controllers.CreateCategory)
		categories.PATCH("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}
}
