package routes

import (
	"to-do-list-golang/config"

	"github.com/gin-gonic/gin"
)

func RouteIndex(r *gin.Engine) {
	cfg := config.LoadConfig()
	path := r.Group(cfg.AppPath)
	{
		TodoRoute(path)
		CategoryRoute(path)
		CategoryColorRoute(path)
		AuthRoute(path)
		UserRoute(path)
	}
}
