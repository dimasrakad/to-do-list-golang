package dtos

type CreateTodoRequest struct {
	Title       string `json:"title" binding:"required"`
	Priority    string `json:"priority" binding:"oneof=low medium high"`
	Description string `json:"description"`
	Due         string `json:"due" binding:"required"`
	CategoryID  uint   `json:"categoryId"`
}
