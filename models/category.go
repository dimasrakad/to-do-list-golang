package models

import "time"

type Category struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	Name            string        `json:"name" gorm:"unique;not null"`
	Todos           []Todo        `json:"-" gorm:"foreignKey:CategoryID"`
	CategoryColorID uint          `json:"categoryColorId" gorm:"not null"`
	CategoryColor   CategoryColor `json:"categoryColor" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	CreatedAt       time.Time     `json:"createdAt"`
	UpdatedAt       time.Time     `json:"updatedAt"`
}
