package dtos

import (
	"time"
	"to-do-list-golang/models"
)

type GetTodosResponse struct {
	ID         uint            `json:"id"`
	Title      string          `json:"title"`
	Status     string          `json:"status"`
	Priority   string          `json:"priority"`
	CategoryID uint            `json:"categoryId"`
	Category   models.Category `json:"category"`
	Due        time.Time       `json:"due"`
}
