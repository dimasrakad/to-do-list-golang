package models

type Category struct {
	ID              uint          `json:"id" gorm:"primaryKey"`
	Name            string        `json:"name" gorm:"unique;not null"`
	Todos           []Todo        `json:"-" gorm:"foreignKey:CategoryID"`
	CategoryColorID uint          `json:"categoryColorId"`
	CategoryColor   CategoryColor `json:"categoryColor" gorm:"constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
