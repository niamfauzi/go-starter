package middleware

import (
	"net/http"
	"strings"

	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"github.com/niamfauzi/go-starter/internal/shared/response"
	"github.com/niamfauzi/go-starter/internal/shared/tenant"
)

// Tenant membaca X-Tenant-Id dari header lalu menyimpannya ke context.
// Ini membantu supaya handler tidak perlu membaca header berulang-ulang.
func Tenant(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenantID := strings.TrimSpace(r.Header.Get("X-Tenant-Id"))
		if tenantID == "" {
			response.Error(w, apperror.BadRequest("header X-Tenant-Id wajib diisi"))
			return
		}

		ctx := tenant.NewContext(r.Context(), tenantID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
