package cmd

import (
	"math/rand"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
)

func CategorySeed() {
	categories := []string{"work", "personal", "study", "other"}

	for _, name := range categories {
		config.DB.FirstOrCreate(&models.Category{}, models.Category{
			Name:            name,
			CategoryColorID: uint(rand.Intn(10) + 1),
		})
	}

}
