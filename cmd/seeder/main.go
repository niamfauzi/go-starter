package main

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"

	"github.com/niamfauzi/go-starter/configs"
	"github.com/niamfauzi/go-starter/internal/auth"
	"github.com/niamfauzi/go-starter/internal/platform/db"
	"github.com/niamfauzi/go-starter/internal/platform/logger"
	"github.com/niamfauzi/go-starter/internal/product"
)

func main() {
	cfg := configs.MustLoad()

	appLogger, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}

	gormDB, err := db.NewGorm(cfg.MySQLDSN, appLogger)
	if err != nil {
		panic(err)
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)
	if err != nil {
		panic(err)
	}

	// Seed user demo.
	user := auth.User{
		TenantID:     "tenant-demo",
		Name:         "Demo User",
		Email:        "demo@example.com",
		PasswordHash: string(passwordHash),
		Role:         "admin",
		IsActive:     true,
	}

	if err := gormDB.Where("tenant_id = ? AND email = ?", user.TenantID, user.Email).
		FirstOrCreate(&user).Error; err != nil {
		panic(err)
	}

	// Seed tenant kedua agar lebih mudah mencoba tenant isolation.
	otherUser := auth.User{
		TenantID:     "tenant-other",
		Name:         "Other User",
		Email:        "other@example.com",
		PasswordHash: string(passwordHash),
		Role:         "admin",
		IsActive:     true,
	}

	if err := gormDB.Where("tenant_id = ? AND email = ?", otherUser.TenantID, otherUser.Email).
		FirstOrCreate(&otherUser).Error; err != nil {
		panic(err)
	}

	// Seed product demo.
	products := []product.Product{
		{TenantID: "tenant-demo", Name: "Keyboard Mechanical", Description: "Keyboard untuk belajar coding", Price: 450000, Stock: 12},
		{TenantID: "tenant-demo", Name: "Mouse Wireless", Description: "Mouse simpel untuk kerja harian", Price: 150000, Stock: 20},
		{TenantID: "tenant-demo", Name: "Monitor 24 Inch", Description: "Monitor untuk produktivitas", Price: 1800000, Stock: 5},
	}

	for _, item := range products {
		var existing product.Product
		err := gormDB.Where("tenant_id = ? AND name = ?", item.TenantID, item.Name).First(&existing).Error
		if err == nil {
			continue
		}

		if err := gormDB.Create(&item).Error; err != nil {
			panic(err)
		}
	}

	fmt.Println("seeder selesai")
	fmt.Println("user demo:")
	fmt.Println("tenant_id: tenant-demo")
	fmt.Println("email    : demo@example.com")
	fmt.Println("password : password123")
	fmt.Println()
	fmt.Println("user tenant lain:")
	fmt.Println("tenant_id: tenant-other")
	fmt.Println("email    : other@example.com")
	fmt.Println("password : password123")
}
