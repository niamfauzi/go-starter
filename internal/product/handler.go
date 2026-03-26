package product

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-playground/validator/v10"
	"go.uber.org/zap"

	"github.com/niamfauzi/go-starter/internal/auth"
	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"github.com/niamfauzi/go-starter/internal/shared/response"
)

// Handler menangani HTTP request untuk fitur product.
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

func (h *Handler) List(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	products, err := h.service.List(r.Context(), authUser.TenantID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "list product berhasil diambil", products)
}

func (h *Handler) GetByID(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	productID, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	product, err := h.service.GetByID(r.Context(), authUser.TenantID, productID)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "detail product berhasil diambil", product)
}

func (h *Handler) Create(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	var req CreateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadRequest("body request tidak valid"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apperror.Validation("validation failed", apperror.FromValidator(err)))
		return
	}

	product, err := h.service.Create(r.Context(), authUser.TenantID, req)
	if err != nil {
		h.logger.Error("gagal create product", zap.Error(err))
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusCreated, "product berhasil dibuat", product)
}

func (h *Handler) Update(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	productID, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	var req UpdateRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		response.Error(w, apperror.BadRequest("body request tidak valid"))
		return
	}

	if err := h.validate.Struct(req); err != nil {
		response.Error(w, apperror.Validation("validation failed", apperror.FromValidator(err)))
		return
	}

	product, err := h.service.Update(r.Context(), authUser.TenantID, productID, req)
	if err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "product berhasil diubah", product)
}

func (h *Handler) Delete(w http.ResponseWriter, r *http.Request) {
	authUser, ok := auth.FromContext(r.Context())
	if !ok {
		response.Error(w, apperror.Unauthorized("unauthorized"))
		return
	}

	productID, err := parseIDParam(r, "id")
	if err != nil {
		response.Error(w, err)
		return
	}

	if err := h.service.Delete(r.Context(), authUser.TenantID, productID); err != nil {
		response.Error(w, err)
		return
	}

	response.Success(w, http.StatusOK, "product berhasil dihapus", nil)
}

func parseIDParam(r *http.Request, key string) (uint64, error) {
	raw := chi.URLParam(r, key)
	id, err := strconv.ParseUint(raw, 10, 64)
	if err != nil {
		return 0, apperror.BadRequest("id tidak valid")
	}
	return id, nil
}
