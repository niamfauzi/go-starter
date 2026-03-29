package auth

import (
	"context"
	"errors"
	"time"

	"gorm.io/gorm"

	"github.com/niamfauzi/go-starter/internal/modules/user"
	"github.com/niamfauzi/go-starter/internal/shared/apperror"
)

// Service berisi business logic untuk fitur auth.
type Service struct {
	repo *Repository
	jwt  *JWTManager
}

func NewService(repo *Repository, jwtManager *JWTManager) *Service {
	return &Service{
		repo: repo,
		jwt:  jwtManager,
	}
}

func (s *Service) Login(ctx context.Context, tenantID string, req LoginRequest) (*LoginResponse, error) {
	foundUser, err := s.repo.FindByEmail(ctx, tenantID, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("email atau password salah")
		}
		return nil, err
	}

	if !foundUser.IsActive {
		return nil, apperror.Unauthorized("akun tidak aktif")
	}

	if err := ComparePassword(foundUser.PasswordHash, req.Password); err != nil {
		return nil, apperror.Unauthorized("email atau password salah")
	}

	token, expiresAt, err := s.jwt.GenerateAccessToken(*foundUser)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
		User:        user.ToSummary(foundUser),
	}, nil
}

func (s *Service) Me(ctx context.Context, tenantID string, userID uint64) (*user.Summary, error) {
	foundUser, err := s.repo.FindByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	summary := user.ToSummary(foundUser)
	return &summary, nil
}
