package routes

import "github.com/gin-gonic/gin"

func RouteIndex(r *gin.Engine) {
	v1 := r.Group("/api/v1")
	{
		TodoRoute(v1)
		CategoryRoute(v1)
		CategoryColorRoute(v1)
		AuthRoute(v1)
	}
}
