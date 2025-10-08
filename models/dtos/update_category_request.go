package dtos

type UpdateCategoryRequest struct {
	Name            *string `json:"name"`
	CategoryColorID *uint   `json:"categoryColorId"`
}
