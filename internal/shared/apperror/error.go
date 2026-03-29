package apperror

import "net/http"

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
