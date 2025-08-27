package routes

import "github.com/gin-gonic/gin"

func RouteIndex(r *gin.Engine) {
	TodoRoute(r)
	CategoryRoute(r)
	CategoryColorRoute(r)
	AuthRoute(r)
}
