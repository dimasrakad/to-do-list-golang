package models

import "time"

type Todo struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	Title       string    `json:"title" gorm:"not null;unique"`
	Description *string   `json:"description" gorm:"type:text"`
	AssignedTo  *uint     `json:"assignedTo"`
	Status      string    `json:"status" gorm:"type:enum('pending','in progress','done');default:'pending';not null"`
	Priority    string    `json:"priority" gorm:"type:enum('low','medium','high');default:'low';not null"`
	CategoryID  uint      `json:"categoryId" gorm:"not null"` // foreign key
	Category    Category  `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Due         time.Time `json:"due" gorm:"not null"`
	CreatedAt   time.Time `json:"createdAt"`
	CreatedBy   uint      `json:"createdBy" gorm:"not null"`
	UpdatedAt   time.Time `json:"updatedAt"`
	UpdatedBy   *uint     `json:"updatedBy"`
}
