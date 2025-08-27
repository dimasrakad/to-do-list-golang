package validators

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

var (
	uppercase = regexp.MustCompile(`[A-Z]`)
	lowercase = regexp.MustCompile(`[a-z]`)
	number    = regexp.MustCompile(`[0-9]`)
	symbol    = regexp.MustCompile(`[\W_]`)
)

func PasswordValidator(fieldLevel validator.FieldLevel) bool {
	password := fieldLevel.Field().String()

	if len(password) < 8 {
		return false
	}
	if !uppercase.MatchString(password) {
		return false
	}
	if !lowercase.MatchString(password) {
		return false
	}
	if !number.MatchString(password) {
		return false
	}
	if !symbol.MatchString(password) {
		return false
	}
	return true
}
