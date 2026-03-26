package auth

// LoginRequest adalah body request untuk login.
type LoginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

// LoginResponse adalah response ketika login berhasil.
type LoginResponse struct {
	AccessToken string   `json:"access_token"`
	TokenType   string   `json:"token_type"`
	ExpiresIn   int64    `json:"expires_in"`
	User        UserInfo `json:"user"`
}

// UserInfo adalah bentuk data user yang aman untuk dikirim ke client.
type UserInfo struct {
	ID       uint64 `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
