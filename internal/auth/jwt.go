package auth

import (
	"fmt"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// Claims adalah data yang akan kita simpan di JWT.
type Claims struct {
	UserID   uint64 `json:"uid"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
	jwt.RegisteredClaims
}

// JWTManager menangani proses membuat dan membaca token.
type JWTManager struct {
	secret []byte
	issuer string
	ttl    time.Duration
}

func NewJWTManager(secret string, issuer string, ttl time.Duration) *JWTManager {
	return &JWTManager{
		secret: []byte(secret),
		issuer: issuer,
		ttl:    ttl,
	}
}

func (m *JWTManager) GenerateAccessToken(user User) (string, time.Time, error) {
	now := time.Now()
	expiresAt := now.Add(m.ttl)

	claims := Claims{
		UserID:   user.ID,
		TenantID: user.TenantID,
		Email:    user.Email,
		Role:     user.Role,
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    m.issuer,
			Subject:   strconv.FormatUint(user.ID, 10),
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signed, err := token.SignedString(m.secret)
	if err != nil {
		return "", time.Time{}, err
	}

	return signed, expiresAt, nil
}

func (m *JWTManager) ParseAccessToken(tokenString string) (*Claims, error) {
	claims := &Claims{}

	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (any, error) {
		// Validasi algoritma penting agar token tidak diterima dengan metode yang tidak kita harapkan.
		if token.Method == nil || token.Method.Alg() != jwt.SigningMethodHS256.Alg() {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return m.secret, nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("token tidak valid")
	}

	return claims, nil
}
