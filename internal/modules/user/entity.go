package user

import (
	"time"

	"gorm.io/gorm"
)

// Entity adalah representasi tabel users.
// Entity utama tabel users disimpan di module user agar tidak diduplikasi.
type Entity struct {
	ID           uint64         `gorm:"column:id;primaryKey" json:"id"`
	TenantID     string         `gorm:"column:tenant_id" json:"tenant_id"`
	Name         string         `gorm:"column:name" json:"name"`
	Email        string         `gorm:"column:email" json:"email"`
	PasswordHash string         `gorm:"column:password_hash" json:"-"`
	Role         string         `gorm:"column:role" json:"role"`
	IsActive     bool           `gorm:"column:is_active" json:"is_active"`
	CreatedAt    time.Time      `gorm:"column:created_at" json:"created_at"`
	UpdatedAt    time.Time      `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"column:deleted_at;index" json:"-"`
}

func (Entity) TableName() string {
	return "users"
}
