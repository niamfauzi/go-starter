package middleware

import (
	"net/http"

	chimiddleware "github.com/go-chi/chi/v5/middleware"
)

func Recovery(next http.Handler) http.Handler {
	return chimiddleware.Recoverer(next)
}
