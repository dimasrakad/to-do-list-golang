package models

import "time"

type Todo struct {
	ID          uint         `json:"id" gorm:"primaryKey"`
	Title       string       `json:"title" gorm:"not null;unique"`
	Description *string      `json:"description" gorm:"type:text"`
	Status      string       `json:"status" gorm:"type:enum('pending','in progress','done');default:'pending';not null"`
	Priority    string       `json:"priority" gorm:"type:enum('low','medium','high');default:'low';not null"`
	CategoryID  uint         `json:"categoryId" gorm:"not null"` // foreign key
	Category    Category     `json:"category" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	Due         time.Time    `json:"due" gorm:"not null"`
	CreatedAt   time.Time    `json:"createdAt"`
	CreatedByID uint         `json:"createdBy" gorm:"not null"`
	CreatedBy   User         `json:"createdByUser" gorm:"foreignKey:CreatedByID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	UpdatedAt   time.Time    `json:"updatedAt"`
	UpdatedByID *uint        `json:"updatedBy"`
	UpdatedBy   *User        `json:"updatedByUser"`
	Assignees   []User       `json:"assignees" gorm:"many2many:todo_assignees;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;"`
	Attachments []Attachment `json:"attachment" gorm:"constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}
