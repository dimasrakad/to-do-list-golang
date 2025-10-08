package dtos

type SuccessResponse struct {
	Data    any    `json:"data"`
	Message string `json:"message"`
}
