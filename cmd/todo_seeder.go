package cmd

import (
	"fmt"
	"math/rand"
	"time"
	"to-do-list-golang/config"
	"to-do-list-golang/models"

	"github.com/go-faker/faker/v4"
)

func TodoSeed(count int, truncate bool) {
	cfg := config.LoadConfig()

	statuses := []string{"pending", "in progress", "done"}
	priorities := []string{"low", "medium", "high"}

	if truncate {
		config.DB.Exec("TRUNCATE TABLE todos")
		fmt.Println("Table todos truncated")
	}

	rand.Seed(time.Now().UnixNano())

	for i := 0; i < count; i++ {
		randomDays := rand.Intn(30) + 1
		randomHours := rand.Intn(24)
		randomMinutes := rand.Intn(60)

		due := time.Now().AddDate(0, 0, randomDays).
			Add(time.Duration(randomHours) * time.Hour).
			Add(time.Duration(randomMinutes) * time.Minute)

		loc, _ := time.LoadLocation(cfg.AppLocation)
		due = due.In(loc)

		todo := models.Todo{
			Title:       faker.Sentence(),
			Description: faker.Paragraph(),
			Status:      statuses[rand.Intn(len(statuses))],
			Priority:    priorities[rand.Intn(len(priorities))],
			Due:         &due,
			CategoryID:  uint(rand.Intn(4) + 1),
		}

		config.DB.Create(&todo)
		fmt.Println("Inserted: ", todo.Title)
	}
}
