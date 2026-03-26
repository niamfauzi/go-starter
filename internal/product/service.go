package product

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/niamfauzi/go-starter/internal/shared/apperror"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

// Service menyimpan business logic product.
// Di sini kita juga meletakkan cache Redis agar handler tetap tipis.
type Service struct {
	repo     *Repository
	cache    *redis.Client
	logger   *zap.Logger
	cacheTTL time.Duration
}

func NewService(repo *Repository, cache *redis.Client, logger *zap.Logger) *Service {
	return &Service{
		repo:     repo,
		cache:    cache,
		logger:   logger,
		cacheTTL: 30 * time.Second,
	}
}

func (s *Service) List(ctx context.Context, tenantID string) ([]Product, error) {
	cacheKey := fmt.Sprintf("products:list:%s", tenantID)

	// Coba ambil dari Redis dulu.
	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var products []Product
			if err := json.Unmarshal([]byte(cached), &products); err == nil {
				return products, nil
			}
		}
	}

	products, err := s.repo.List(ctx, tenantID)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if b, err := json.Marshal(products); err == nil {
			_ = s.cache.Set(ctx, cacheKey, string(b), s.cacheTTL).Err()
		}
	}

	return products, nil
}

func (s *Service) GetByID(ctx context.Context, tenantID string, productID uint64) (*Product, error) {
	cacheKey := fmt.Sprintf("products:get:%s:%d", tenantID, productID)

	if s.cache != nil {
		cached, err := s.cache.Get(ctx, cacheKey).Result()
		if err == nil && cached != "" {
			var product Product
			if err := json.Unmarshal([]byte(cached), &product); err == nil {
				return &product, nil
			}
		}
	}

	product, err := s.repo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}

	if s.cache != nil {
		if b, err := json.Marshal(product); err == nil {
			_ = s.cache.Set(ctx, cacheKey, string(b), s.cacheTTL).Err()
		}
	}

	return product, nil
}

func (s *Service) Create(ctx context.Context, tenantID string, req CreateRequest) (*Product, error) {
	product := &Product{
		TenantID:    tenantID,
		Name:        req.Name,
		Description: req.Description,
		Price:       req.Price,
		Stock:       req.Stock,
	}

	if err := s.repo.Create(ctx, product); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, tenantID, product.ID)

	return product, nil
}

func (s *Service) Update(ctx context.Context, tenantID string, productID uint64, req UpdateRequest) (*Product, error) {
	product, err := s.repo.FindByID(ctx, tenantID, productID)
	if err != nil {
		return nil, err
	}

	product.Name = req.Name
	product.Description = req.Description
	product.Price = req.Price
	product.Stock = req.Stock

	if err := s.repo.Update(ctx, product); err != nil {
		return nil, err
	}

	s.invalidateCache(ctx, tenantID, product.ID)

	return product, nil
}

func (s *Service) Delete(ctx context.Context, tenantID string, productID uint64) error {
	// Pastikan product benar-benar milik tenant tersebut.
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
	if s.cache == nil {
		return
	}

	listKey := fmt.Sprintf("products:list:%s", tenantID)
	getKey := fmt.Sprintf("products:get:%s:%d", tenantID, productID)

	if err := s.cache.Del(ctx, listKey, getKey).Err(); err != nil {
		s.logger.Warn("gagal invalidasi cache product", zap.Error(err))
	}
}
