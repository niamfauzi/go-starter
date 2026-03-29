package apperror

import (
	"errors"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// As mengubah error biasa menjadi AppError jika memungkinkan.
func As(err error) *AppError {
	if err == nil {
		return nil
	}

	var appErr *AppError
	if errors.As(err, &appErr) {
		return appErr
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return NotFound("data tidak ditemukan")
	}

	return Internal(err)
}

// FromValidator mengubah hasil validator menjadi detail yang ramah untuk client.
func FromValidator(err error) []FieldError {
	var validationErrors validator.ValidationErrors
	if !errors.As(err, &validationErrors) {
		return nil
	}

	result := make([]FieldError, 0, len(validationErrors))
	for _, fieldErr := range validationErrors {
		result = append(result, FieldError{
			Field:   fieldErr.Field(),
			Message: validatorMessage(fieldErr),
		})
	}

	return result
}

func validatorMessage(err validator.FieldError) string {
	switch err.Tag() {
	case "required":
		return "field ini wajib diisi"
	case "email":
		return "format email tidak valid"
	case "min":
		return "nilai terlalu pendek / kecil"
	case "max":
		return "nilai terlalu panjang / besar"
	case "gte":
		return "nilai harus lebih besar atau sama dengan batas minimum"
	default:
		return "nilai tidak valid"
	}
}
