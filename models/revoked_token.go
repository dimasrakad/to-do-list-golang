package models

import "time"

type RevokedToken struct {
	ID        uint      `gorm:"primaryKey"`
	Token     string    `gorm:"type:text;not null"`
	UserID    uint      `gorm:"not null"`
	ExpiresAt time.Time `gorm:"not null"`
	CreatedAt time.Time
}
