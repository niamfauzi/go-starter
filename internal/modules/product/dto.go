package product

import "time"

// CreateRequest adalah body request saat membuat product.
type CreateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"max=500"`
	Price       int64  `json:"price" validate:"required,gte=0"`
	Stock       int    `json:"stock" validate:"required,gte=0"`
}

// UpdateRequest adalah body request saat mengubah product.
type UpdateRequest struct {
	Name        string `json:"name" validate:"required,min=3,max=150"`
	Description string `json:"description" validate:"max=500"`
	Price       int64  `json:"price" validate:"required,gte=0"`
	Stock       int    `json:"stock" validate:"required,gte=0"`
}

// Response adalah bentuk data product yang dikirim ke client.
type Response struct {
	ID          uint64    `json:"id"`
	TenantID    string    `json:"tenant_id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       int64     `json:"price"`
	Stock       int       `json:"stock"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}
