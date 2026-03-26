package auth

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"github.com/niamfauzi/go-starter/internal/shared/response"
	"github.com/niamfauzi/go-starter/internal/shared/tenant"
)

// Handler menangani request/response HTTP untuk fitur auth.
type Handler struct {
	service  *Service
	validate *validator.Validate
	logger   *zap.Logger
}

func NewHandler(service *Service, validate *validator.Validate, logger *zap.Logger) *Handler {
	return &Handler{
		service:  service,
		validate: validate,
		logger:   logger,
	}
}

// Login menangani proses login.
func (h *Handler) Login(w http.ResponseWriter, r *http.Request) {
	tenantID, ok := tenant.FromContext(r.Context())
	if !ok || tenantID == "" {
		response.Error(w, apperror.BadRequest("header X-Tenant-Id wajib diisi"))
		return
	}

	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadRequest("body request tidak valid"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apperror.Validation("validation failed", apperror.FromValidator(err)))
		return
	}

	res, err := h.service.Login(r.Context(), tenantID, req)
	if err != nil {
		var appErr *apperror.AppError
		if !errors.As(err, &appErr) {
			h.logger.Error("login gagal", zap.Error(err), zap.String("tenant_id", tenantID))
		}
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "login berhasil", res)
}

// Me mengambil profile user yang sedang login.
func (h *Handler) Me(w http.ResponseWriter, r *http.Request) {
	authUser, ok := FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	res, err := h.service.Me(r.Context(), authUser.TenantID, authUser.UserID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "profile berhasil diambil", res)
}
