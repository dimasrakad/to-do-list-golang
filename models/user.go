package models

import "time"

type User struct {
	ID          uint       `json:"id" gorm:"primaryKey"`
	Name        string     `json:"name" gorm:"not null"`
	Email       string     `json:"email" gorm:"unique;not null"`
	Password    string     `json:"-" gorm:"not null"`
	IsVerified  bool       `json:"isVerified" gorm:"not null;default:false"`
	VerifyToken *string    `json:"-"`
	VerifiedAt  *time.Time `json:"verifiedAt"`
	CreatedAt   time.Time  `json:"createdAt"`
	UpdatedAt   time.Time  `json:"updatedAt"`
}
