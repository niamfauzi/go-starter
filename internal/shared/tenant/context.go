package tenant

import "context"

type contextKey string

const tenantKey contextKey = "tenant_id"

// NewContext menyimpan tenant_id ke context.
func NewContext(ctx context.Context, tenantID string) context.Context {
	return context.WithValue(ctx, tenantKey, tenantID)
}

// FromContext mengambil tenant_id dari context.
func FromContext(ctx context.Context) (string, bool) {
	tenantID, ok := ctx.Value(tenantKey).(string)
	return tenantID, ok
}
