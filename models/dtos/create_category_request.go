package dtos

type CreateCategoryRequest struct {
	Name            string `json:"name" binding:"required"`
	CategoryColorID uint   `json:"categoryColorId" binding:"required"`
}
