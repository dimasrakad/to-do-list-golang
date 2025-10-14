package models

import "time"

type Attachment struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	TodoID    uint      `json:"todoId" gorm:"not null;index"`
	FileName  string    `json:"filename" gorm:"not null"`
	FilePath  string    `json:"filePath" gorm:"not null"`
	FileType  string    `json:"fileType"`
	CreatedAt time.Time `json:"createdAt"`
}
