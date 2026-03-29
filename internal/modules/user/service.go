package user

import "context"

// Service menampung business logic sederhana untuk user.
// Saat ini tetap tipis karena kebutuhan user belum kompleks.
type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{repo: repo}
}

func (s *Service) FindByEmail(ctx context.Context, tenantID string, email string) (*Entity, error) {
	return s.repo.FindByEmail(ctx, tenantID, email)
}

func (s *Service) FindByID(ctx context.Context, tenantID string, userID uint64) (*Entity, error) {
	return s.repo.FindByID(ctx, tenantID, userID)
}
