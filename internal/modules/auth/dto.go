package auth

import "github.com/niamfauzi/go-starter/internal/modules/user"

// LoginRequest adalah body request untuk login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse adalah response ketika login berhasil.
type LoginResponse struct {
	AccessToken string       `json:"access_token"`
	TokenType   string       `json:"token_type"`
	ExpiresIn   int64        `json:"expires_in"`
	User        user.Summary `json:"user"`
}
