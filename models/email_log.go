package models

import "time"

type EmailLog struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Recipient string    `json:"recipient" gorm:"not null"`
	Subject   string    `json:"type" gorm:"not null"`
	SentAt    time.Time `json:"sentAt" gorm:"not null"`
}
