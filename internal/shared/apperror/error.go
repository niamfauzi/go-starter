package apperror

import (
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// FieldError dipakai untuk memberi tahu field mana yang salah saat validasi.
type FieldError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// AppError adalah error aplikasi yang sudah memiliki:
// - code untuk kode error internal
// - status HTTP
// - message yang aman ditampilkan ke client
// - details opsional
type AppError struct {
	Code    string `json:"code"`
	Message string `json:"message"`
	Status  int    `json:"-"`
	Details any    `json:"details,omitempty"`
	Err     error  `json:"-"`
}

func (e *AppError) Error() string {
	if e.Err != nil {
		return e.Err.Error()
	}
	return e.Message
}

func (e *AppError) Unwrap() error {
	return e.Err
}

func New(code string, message string, status int, details any, err error) *AppError {
	return &AppError{
		Code:    code,
		Message: message,
		Status:  status,
		Details: details,
		Err:     err,
	}
}

func Validation(message string, details []FieldError) *AppError {
	return New("VALIDATION_ERROR", message, http.StatusBadRequest, details, nil)
}

func BadRequest(message string) *AppError {
	return New("BAD_REQUEST", message, http.StatusBadRequest, nil, nil)
}

func Unauthorized(message string) *AppError {
	return New("UNAUTHORIZED", message, http.StatusUnauthorized, nil, nil)
}

func Forbidden(message string) *AppError {
	return New("FORBIDDEN", message, http.StatusForbidden, nil, nil)
}

func NotFound(message string) *AppError {
	return New("NOT_FOUND", message, http.StatusNotFound, nil, nil)
}

func Conflict(message string) *AppError {
	return New("CONFLICT", message, http.StatusConflict, nil, nil)
}

func Internal(err error) *AppError {
	return New("INTERNAL_SERVER_ERROR", "terjadi kesalahan pada server", http.StatusInternalServerError, nil, err)
}

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
