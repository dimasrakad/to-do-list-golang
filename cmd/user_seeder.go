package cmd

import (
	"strconv"
	"to-do-list-golang/config"
	"to-do-list-golang/models"
	"to-do-list-golang/utils"
)

func UserSeed() {
	for i := range 3 {
		index := strconv.Itoa(i + 1)
		password, _ := utils.HashPassword("password" + index)
		config.DB.FirstOrCreate(&models.User{}, models.User{
			Name:     "Seeder User " + index,
			Email:    "seederuser" + index + "@yopmail.com",
			Password: password,
		})
	}
}
