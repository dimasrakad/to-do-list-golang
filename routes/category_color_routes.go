package routes

import (
	"to-do-list-golang/controllers"

	"github.com/gin-gonic/gin"
)

func CategoryColorRoute(r *gin.Engine) {
	{
		colors := r.Group("/category-colors")
		colors.GET("", controllers.GetColors)
	}
}
