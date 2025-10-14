package dtos

type UpdateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	AssignedTo  *[]uint `json:"assignedTo"`
	Priority    *string `json:"priority" binding:"omitempty,oneof=low medium high"`
	Status      *string `json:"status" binding:"omitempty,oneof='pending' 'in progress' 'done'"`
	Due         *string `json:"due"`
	CategoryID  *uint   `json:"categoryId"`
}
