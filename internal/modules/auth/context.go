package auth

import "context"

type contextKey string

const authUserKey contextKey = "auth_user"

// NewContext menyimpan user hasil autentikasi ke context.
func NewContext(ctx context.Context, user AuthUser) context.Context {
	return context.WithValue(ctx, authUserKey, user)
}

// FromContext mengambil user hasil autentikasi dari context.
func FromContext(ctx context.Context) (AuthUser, bool) {
	user, ok := ctx.Value(authUserKey).(AuthUser)
	return user, ok
}
