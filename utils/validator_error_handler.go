package utils

import (
	"errors"
	"fmt"
	"net/http"
	"to-do-list-golang/models/dtos"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateInput(c *gin.Context, input any) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		var validationError validator.ValidationErrors
		res := dtos.ErrorResponse{}

		if errors.As(err, &validationError) {
			out := ""
			for i, fieldError := range validationError {
				if i > 0 {
					out += "\n"
				}
				if fieldError.Tag() == "password" {
					out += "Password must be at least 8 characters, contain uppercase, lowercase, number, and symbol"
				} else {
					out += ValidationErrorToText(fieldError)
				}
			}

			res.Error = out
			c.JSON(http.StatusBadRequest, res)
		} else {
			res.Error = err.Error()
			c.JSON(http.StatusBadRequest, res)
		}
		return false
	}
	return true
}

func ValidationErrorToText(fieldError validator.FieldError) string {
	switch fieldError.Tag() {
	case "required":
		return fmt.Sprintf("%s is required", fieldError.Field())
	case "email":
		return fmt.Sprintf("%s must be a valid email", fieldError.Field())
	case "min":
		return fmt.Sprintf("%s must be at least %s characters", fieldError.Field(), fieldError.Param())
	case "max":
		return fmt.Sprintf("%s must be at most %s characters", fieldError.Field(), fieldError.Param())
	default:
		return fmt.Sprintf("%s is not valid", fieldError.Field())
	}
}
