package models

type CategoryColor struct {
	ID         uint       `json:"id" gorm:"primaryKey"`
	Name       string     `json:"name" gorm:"unique"`
	Code       string     `json:"code" gorm:"unique"`
	Categories []Category `json:"-" gorm:"foreignKey:CategoryColorID"`
}
