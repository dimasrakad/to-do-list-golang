package dtos

type CreateTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description *string `json:"description"`
	AssignedTo  []uint  `json:"assignedTo" binding:"required"`
	Priority    string  `json:"priority" binding:"oneof=low medium high"`
	Due         string  `json:"due" binding:"required"`
	CategoryID  uint    `json:"categoryId"`
}
