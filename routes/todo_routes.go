package routes

import (
	"to-do-list-golang/controllers"

	"github.com/gin-gonic/gin"
)

func TodoRoute(r *gin.Engine) {
	{
		todos := r.Group("/todos")
		todos.GET("", controllers.GetTodos)
		todos.POST("", controllers.CreateTodo)
		todos.PATCH("/:id", controllers.UpdateTodo)
		todos.DELETE("/:id", controllers.DeleteTodo)
	}
}
