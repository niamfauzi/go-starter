package product

import (
	"time"

	"gorm.io/gorm"
)

// Product adalah model database untuk tabel products.
type Product struct {
	ID          uint64         `gorm:"column:id;primaryKey" json:"id"`
	TenantID    string         `gorm:"column:tenant_id" json:"tenant_id"`
	Name        string         `gorm:"column:name" json:"name"`
	Description string         `gorm:"column:description" json:"description"`
	Price       int64          `gorm:"column:price" json:"price"`
	Stock       int            `gorm:"column:stock" json:"stock"`
	CreatedAt   time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt   time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Product) TableName() string {
	return "products"
}
