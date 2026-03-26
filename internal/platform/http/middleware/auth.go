package middleware

import (
	"net/http"
	"strings"

	"github.com/niamfauzi/go-starter/internal/auth"
	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"github.com/niamfauzi/go-starter/internal/shared/response"
	"github.com/niamfauzi/go-starter/internal/shared/tenant"
)

// AuthJWT memvalidasi access token JWT.
// Jika valid, middleware akan menyimpan user ke context.
func AuthJWT(jwtManager *auth.JWTManager) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			tokenString := extractBearerToken(r.Header.Get("Authorization"))
			if tokenString == "" {
				response.Error(w, apperror.Unauthorized("bearer token wajib diisi"))
				return
			}

			claims, err := jwtManager.ParseAccessToken(tokenString)
			if err != nil {
				response.Error(w, apperror.Unauthorized("token tidak valid atau sudah expired"))
				return
			}

			tenantID, _ := tenant.FromContext(r.Context())
			if tenantID != "" && claims.TenantID != tenantID {
				response.Error(w, apperror.Forbidden("tenant pada token tidak cocok dengan header"))
				return
			}

			authUser := auth.AuthUser{
				UserID:   claims.UserID,
				TenantID: claims.TenantID,
				Email:    claims.Email,
				Role:     claims.Role,
			}

			ctx := auth.NewContext(r.Context(), authUser)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func extractBearerToken(header string) string {
	header = strings.TrimSpace(header)
	if header == "" {
		return ""
	}

	parts := strings.SplitN(header, " ", 2)
	if len(parts) != 2 {
		return ""
	}

	if !strings.EqualFold(parts[0], "Bearer") {
		return ""
	}

	return strings.TrimSpace(parts[1])
}
