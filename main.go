package main

import (
	"flag"
	"to-do-list-golang/cmd"
	"to-do-list-golang/config"
	_ "to-do-list-golang/docs"
	"to-do-list-golang/goroutines"
	"to-do-list-golang/models"
	"to-do-list-golang/routes"
	"to-do-list-golang/validators"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @title To Do List API
// @version 1.0
// @description This is a simple To Do List API server using Golang.
// @termsOfService http://swagger.io/terms/

// @contact.name Dimas Raka Dewanggana
// @contact.email dimasdewanggana@gmail.com

// @host localhost:8080
// @BasePath /api/v1

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
	seed := flag.Bool("seed", false, "run database seeder")
	count := flag.Int("count", 10, "number of records to seed")
	truncate := flag.Bool("truncate", false, "truncate table before seeding")
	flag.Parse()

	if *seed {
		cmd.UserSeed()
		cmd.CategoryColorSeed()
		cmd.CategorySeed()
		cmd.TodoSeed(*count, *truncate)
		return
	}

	// Init Gin
	r := gin.Default()

	// Router
	routes.RouteIndex(r)

	// Custom validator
	if validator, ok := binding.Validator.Engine().(*validator.Validate); ok {
		validator.RegisterValidation("password", validators.PasswordValidator)
	}

	// Swagger
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Run server
	r.Run(":" + cfg.AppPort)
}
