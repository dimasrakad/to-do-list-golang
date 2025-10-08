package dtos

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	Status      *string `json:"status" binding:"omitempty,oneof='pending' 'in progress' 'done'"`
	Description *string `json:"description"`
	Due         *string `json:"due"`
}
