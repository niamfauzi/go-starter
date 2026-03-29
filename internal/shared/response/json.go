package response

import (
	"encoding/json"
	"net/http"

	"github.com/niamfauzi/go-starter/internal/shared/apperror"
)

// SuccessEnvelope adalah bentuk response sukses yang konsisten.
type SuccessEnvelope struct {
	Success bool   `json:"success"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// ErrorEnvelope adalah bentuk response error yang konsisten.
type ErrorEnvelope struct {
	Success bool       `json:"success"`
	Message string     `json:"message"`
	Error   *ErrorBody `json:"error,omitempty"`
}

type ErrorBody struct {
	Code    string `json:"code"`
	Details any    `json:"details,omitempty"`
}

func JSON(w http.ResponseWriter, status int, payload any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(payload)
}

func Success(w http.ResponseWriter, status int, message string, data any) {
	JSON(w, status, SuccessEnvelope{
		Success: true,
		Message: message,
		Data:    data,
	})
}

func Error(w http.ResponseWriter, err error) {
	appErr := apperror.As(err)

	JSON(w, appErr.Status, ErrorEnvelope{
		Success: false,
		Message: appErr.Message,
		Error: &ErrorBody{
			Code:    appErr.Code,
			Details: appErr.Details,
		},
	})
}
