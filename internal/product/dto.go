package product

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
