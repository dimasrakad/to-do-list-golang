package cmd

import (
	"to-do-list-golang/config"
	"to-do-list-golang/models"
)

func CategoryColorSeed() {
	colors := map[string]string{
		"red":         "E74C3C",
		"orange":      "E67E22",
		"yellow":      "F1C40F",
		"green":       "2ECC71",
		"blue":        "3498DB",
		"purple":      "9B59B6",
		"dark grey":   "34495E",
		"light brown": "D35400",
		"light grey":  "BDC3C7",
		"pink":        "FF6B81",
	}

	for name, code := range colors {
		config.DB.FirstOrCreate(&models.CategoryColor{}, models.CategoryColor{
			Name: name,
			Code: code,
		})
	}
}
