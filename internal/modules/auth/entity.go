package auth

// AuthUser adalah data user yang sudah lolos autentikasi dan aman disimpan ke context.
type AuthUser struct {
	UserID   uint64 `json:"user_id"`
	TenantID string `json:"tenant_id"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}
