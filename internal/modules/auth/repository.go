package auth

import (
	"context"

	"github.com/niamfauzi/go-starter/internal/modules/user"
)

// Repository adalah auth-specific data access.
// Kita tetap punya module user sebagai pemilik entity utama tabel users.
type Repository struct {
	userRepo *user.Repository
}

func NewRepository(userRepo *user.Repository) *Repository {
	return &Repository{userRepo: userRepo}
}

func (r *Repository) FindByEmail(ctx context.Context, tenantID string, email string) (*user.Entity, error) {
	return r.userRepo.FindByEmail(ctx, tenantID, email)
}

func (r *Repository) FindByID(ctx context.Context, tenantID string, userID uint64) (*user.Entity, error) {
	return r.userRepo.FindByID(ctx, tenantID, userID)
}
