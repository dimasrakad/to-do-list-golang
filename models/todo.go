package models

import "time"

type Todo struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Title       string     `json:"title" gorm:"unique"`
	Description string     `json:"description" gorm:"type:text"`
	Status      string     `json:"status" gorm:"type:enum('pending','in progress','completed');default:'pending'"`
	Priority    string     `json:"priority" gorm:"type:enum('low','medium','high');default:'low'"`
	Due         *time.Time `json:"due"` // not null
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
