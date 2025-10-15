package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middlewares"

	"github.com/gin-gonic/gin"
)

func TodoRoute(r *gin.RouterGroup) {
	todos := r.Group("/todos")
	todos.Use(middlewares.AuthMiddleware())
	{
		todos.GET("", controllers.GetTodos)
		todos.POST("", controllers.CreateTodo)
		todos.PATCH("/:id", controllers.UpdateTodo)
		todos.DELETE("/:id", controllers.DeleteTodo)
		todos.POST("/:id/attachments", controllers.UploadAttachments)
		todos.GET("/attachments/:id", controllers.GetAttachment)
		todos.DELETE("/attachments/:id", controllers.DeleteAttachment)
	}
}
