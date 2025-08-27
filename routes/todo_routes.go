package routes

import (
	"to-do-list-golang/controllers"
	"to-do-list-golang/middleware"

	"github.com/gin-gonic/gin"
)

func TodoRoute(r *gin.Engine) {
	todos := r.Group("/todos")
	todos.Use(middleware.AuthMiddleware())
	{
		todos.GET("", controllers.GetTodos)
		todos.POST("", controllers.CreateTodo)
		todos.PATCH("/:id", controllers.UpdateTodo)
		todos.DELETE("/:id", controllers.DeleteTodo)
	}
}
