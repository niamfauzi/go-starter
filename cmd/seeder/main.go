package main

import (
	"context"
	"fmt"

	"github.com/niamfauzi/go-starter/configs"
	"github.com/niamfauzi/go-starter/internal/modules/auth"
	"github.com/niamfauzi/go-starter/internal/modules/product"
	"github.com/niamfauzi/go-starter/internal/modules/user"
	"github.com/niamfauzi/go-starter/internal/platform/db"
	"github.com/niamfauzi/go-starter/internal/platform/logger"
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

	passwordHash, err := auth.HashPassword("password123")
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	userRepo := user.NewRepository(gormDB)

	demoUser := user.Entity{
		TenantID:     "tenant-demo",
		Name:         "Demo User",
		Email:        "demo@example.com",
		PasswordHash: string(passwordHash),
		Role:         "admin",
		IsActive:     true,
	}

	if err := userRepo.CreateIfNotExists(ctx, &demoUser); err != nil {
		panic(err)
	}

	otherUser := user.Entity{
		TenantID:     "tenant-other",
		Name:         "Other User",
		Email:        "other@example.com",
		PasswordHash: string(passwordHash),
		Role:         "admin",
		IsActive:     true,
	}

	if err := userRepo.CreateIfNotExists(ctx, &otherUser); err != nil {
		panic(err)
	}

	products := []product.Entity{
		{TenantID: "tenant-demo", Name: "Keyboard Mechanical", Description: "Keyboard untuk belajar coding", Price: 450000, Stock: 12},
		{TenantID: "tenant-demo", Name: "Mouse Wireless", Description: "Mouse simpel untuk kerja harian", Price: 150000, Stock: 20},
		{TenantID: "tenant-demo", Name: "Monitor 24 Inch", Description: "Monitor untuk produktivitas", Price: 1800000, Stock: 5},
	}

	for _, item := range products {
		var existing product.Entity
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
