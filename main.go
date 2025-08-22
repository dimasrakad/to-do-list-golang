package main

import (
	"flag"
	"to-do-list-golang/cmd"
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

	// Flag
	seed := flag.Bool("seed", false, "run database seeder")
	count := flag.Int("count", 10, "number of records to seed")
	truncate := flag.Bool("truncate", false, "truncate table before seeding")
	flag.Parse()

	if *seed {
		cmd.Seed(*count, *truncate)
		return
	}

	// Init Gin
	r := gin.Default()

	// Router
	routes.TodoRoute(r)

	// Run server
	r.Run(":" + cfg.AppPort)
}
