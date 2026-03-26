package product

import (
	"context"

	"gorm.io/gorm"
)

// Repository menangani semua query product ke database.
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) Create(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Create(product).Error
}

func (r *Repository) List(ctx context.Context, tenantID string) ([]Product, error) {
	var products []Product
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Order("id DESC").
		Find(&products).Error
	return products, err
}

func (r *Repository) FindByID(ctx context.Context, tenantID string, productID uint64) (*Product, error) {
	var product Product
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("id = ?", productID).
		First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *Repository) Update(ctx context.Context, product *Product) error {
	return r.db.WithContext(ctx).Save(product).Error
}

func (r *Repository) Delete(ctx context.Context, tenantID string, productID uint64) error {
	return r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Delete(&Product{}, productID).Error
}
