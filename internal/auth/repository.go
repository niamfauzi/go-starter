package auth

import (
	"context"

	"gorm.io/gorm"
)

// Repository menangani akses data user.
type Repository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) *Repository {
	return &Repository{db: db}
}

func (r *Repository) FindByEmail(ctx context.Context, tenantID string, email string) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("email = ?", email).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}

func (r *Repository) FindByID(ctx context.Context, tenantID string, userID uint64) (*User, error) {
	var user User
	err := r.db.WithContext(ctx).
		Where("tenant_id = ?", tenantID).
		Where("id = ?", userID).
		First(&user).Error
	if err != nil {
		return nil, err
	}
	return &user, nil
}
