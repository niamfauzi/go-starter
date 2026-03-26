package http

import (
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	chimiddleware "github.com/go-chi/chi/v5/middleware"
	"go.uber.org/zap"

	"github.com/niamfauzi/go-starter/internal/auth"
	httpmiddleware "github.com/niamfauzi/go-starter/internal/platform/http/middleware"
	"github.com/niamfauzi/go-starter/internal/product"
)

// NewRouter menyusun semua route + middleware aplikasi.
func NewRouter(
	logger *zap.Logger,
	authHandler *auth.Handler,
	productHandler *product.Handler,
	jwtManager *auth.JWTManager,
) http.Handler {
	r := chi.NewRouter()

	// Middleware umum yang berlaku ke semua endpoint.
	r.Use(chimiddleware.RequestID)
	r.Use(chimiddleware.RealIP)
	r.Use(httpmiddleware.Logger(logger))
	r.Use(chimiddleware.Recoverer)
	r.Use(chimiddleware.Timeout(30 * time.Second))
	r.Use(httpmiddleware.Tenant)

	r.Get("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Menyajikan file dokumentasi Swagger/OpenAPI sederhana.
	r.Handle("/docs/*", http.StripPrefix("/docs/", http.FileServer(http.Dir("./docs"))))

	r.Route("/api/v1", func(r chi.Router) {
		// Auth
		r.Post("/auth/login", authHandler.Login)

		// Semua endpoint product dibuat private.
		// Alasan: project ini multi-tenant, jadi lebih aman bila semua akses product
		// mensyaratkan token agar konsisten sejak awal.
		r.Group(func(r chi.Router) {
			r.Use(httpmiddleware.AuthJWT(jwtManager))

			r.Get("/auth/me", authHandler.Me)

			r.Route("/products", func(r chi.Router) {
				r.Get("/", productHandler.List)
				r.Post("/", productHandler.Create)
				r.Get("/{id}", productHandler.GetByID)
				r.Put("/{id}", productHandler.Update)
				r.Delete("/{id}", productHandler.Delete)
			})
		})
	})

	return r
}
