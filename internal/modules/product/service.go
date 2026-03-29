package product

import (
	"context"

	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"go.uber.org/zap"
)

// Service menyimpan business logic product.
type Service struct {
	repo   *Repository
	cache  *Cache
	logger *zap.Logger
}

func NewService(repo *Repository, cache *Cache, logger *zap.Logger) *Service {
	return &Service{
		repo:   repo,
		cache:  cache,
		logger: logger,
	}
}

func (s *Service) List(ctx context.Context, tenantID string) ([]Response, error) {
	if cached, ok := s.cache.GetList(ctx, tenantID); ok {
		return cached, nil
	}

	products, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	responses := toResponseList(products)
	s.cache.SetList(ctx, tenantID, responses)

	return responses, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID string, productID uint64) (*Response, error) {
	if cached, ok := s.cache.GetDetail(ctx, tenantID, productID); ok {
		return cached, nil
	}

	entity, err := s.repo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}

	response := toResponse(*entity)
	s.cache.SetDetail(ctx, tenantID, productID, response)

	return &response, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (*Response, error) {
	entity := &Entity{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.Create(ctx, entity); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, tenantID, entity.ID)

	response := toResponse(*entity)
	return &response, nil
}

func (s *Service) Update(ctx context.Context, tenantID string, productID uint64, req UpdateRequest) (*Response, error) {
	entity, err := s.repo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}

	entity.Name = req.Name
	entity.Description = req.Description
	entity.Price = req.Price
	entity.Stock = req.Stock

	if err := s.repo.Update(ctx, entity); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, tenantID, entity.ID)

	response := toResponse(*entity)
	return &response, nil
}

func (s *Service) Delete(ctx context.Context, tenantID string, productID uint64) error {
	if _, err := s.repo.FindByID(ctx, tenantID, productID); err != nil {
		return apperror.NotFound("product tidak ditemukan")
	}

	if err := s.repo.Delete(ctx, tenantID, productID); err != nil {
		return err
	}

	s.invalidateCache(ctx, tenantID, productID)
	return nil
}

func (s *Service) invalidateCache(ctx context.Context, tenantID string, productID uint64) {
	if err := s.cache.Invalidate(ctx, tenantID, productID); err != nil {
		s.logger.Warn("gagal invalidasi cache product", zap.Error(err))
	}
}
