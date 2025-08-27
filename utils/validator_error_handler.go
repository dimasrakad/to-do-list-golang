package utils

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func ValidateInput(c *gin.Context, input any) bool {
	if err := c.ShouldBindJSON(input); err != nil {
		var validationError validator.ValidationErrors

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
			c.JSON(http.StatusBadRequest, gin.H{"error": out})
		} else {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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
