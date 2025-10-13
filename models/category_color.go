package models

type CategoryColor struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"unique;not null"`
	Code       string     `json:"code" gorm:"unique;not null"`
	Categories []Category `json:"-" gorm:"foreignKey:CategoryColorID"`
}
