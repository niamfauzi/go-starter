package validator

import "github.com/go-playground/validator/v10"

// New membuat validator instance untuk dipakai handler.
func New() *validator.Validate {
	return validator.New()
}
