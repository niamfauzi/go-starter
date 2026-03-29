package user

// Summary adalah bentuk data user yang aman untuk dikirim ke layer lain / client.
type Summary struct {
	ID       uint64 `json:"id"`
	TenantID string `json:"tenant_id"`
	Name     string `json:"name"`
	Email    string `json:"email"`
	Role     string `json:"role"`
}

func ToSummary(entity *Entity) Summary {
	return Summary{
		ID:       entity.ID,
		TenantID: entity.TenantID,
		Name:     entity.Name,
		Email:    entity.Email,
		Role:     entity.Role,
	}
}
