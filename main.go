package main

import (
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/routes"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	config.ConnectDatabase()
	config.DB.AutoMigrate(&models.Todo{})

	// Init Gin
	r := gin.Default()

	// Router
	routes.TodoRoute(r)

	// Run server
	r.Run(":" + cfg.AppPort)
}
