package auth

import (
	"context"
	"errors"
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"

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
	user, err := s.repo.FindByEmail(ctx, tenantID, req.Email)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, apperror.Unauthorized("email atau password salah")
		}
		return nil, err
	}

	if !user.IsActive {
		return nil, apperror.Unauthorized("akun tidak aktif")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		return nil, apperror.Unauthorized("email atau password salah")
	}

	token, expiresAt, err := s.jwt.GenerateAccessToken(*user)
	if err != nil {
		return nil, err
	}

	return &LoginResponse{
		AccessToken: token,
		TokenType:   "Bearer",
		ExpiresIn:   int64(time.Until(expiresAt).Seconds()),
		User: UserInfo{
			ID:       user.ID,
			TenantID: user.TenantID,
			Name:     user.Name,
			Email:    user.Email,
			Role:     user.Role,
		},
	}, nil
}

func (s *Service) Me(ctx context.Context, tenantID string, userID uint64) (*UserInfo, error) {
	user, err := s.repo.FindByID(ctx, tenantID, userID)
	if err != nil {
		return nil, err
	}

	return &UserInfo{
		ID:       user.ID,
		TenantID: user.TenantID,
		Name:     user.Name,
		Email:    user.Email,
		Role:     user.Role,
	}, nil
}
