package main

import (
	"to-do-list-golang/config"
	"to-do-list-golang/goroutines"
	"to-do-list-golang/models"
	"to-do-list-golang/routes"
	"to-do-list-golang/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func main() {
	// Load config
	cfg := config.LoadConfig()

	// Connect DB
	config.ConnectDatabase()
	config.DB.AutoMigrate(
		&models.Category{},
		&models.Todo{},
		&models.CategoryColor{},
		&models.User{},
		&models.RefreshToken{},
		&models.RevokedToken{},
	)

	goroutines.StartTokenCleanup()

	// Flag
	// seed := flag.Bool("seed", false, "run database seeder")
	// count := flag.Int("count", 10, "number of records to seed")
	// truncate := flag.Bool("truncate", false, "truncate table before seeding")
	// flag.Parse()

	// if *seed {
	// 	cmd.CategoryColorSeed()
	// 	cmd.CategorySeed()
	// 	cmd.TodoSeed(*count, *truncate)
	// 	return
	// }

	// Init Gin
	r := gin.Default()

	// Router
	routes.RouteIndex(r)

	// Custom validator
	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validator.RegisterValidation("password", validators.PasswordValidator)
	}

	// Run server
	r.Run(":" + cfg.AppPort)
}
