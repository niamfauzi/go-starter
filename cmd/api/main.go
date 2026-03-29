package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/niamfauzi/go-starter/configs"
	"github.com/niamfauzi/go-starter/internal/modules/auth"
	"github.com/niamfauzi/go-starter/internal/modules/product"
	"github.com/niamfauzi/go-starter/internal/modules/user"
	"github.com/niamfauzi/go-starter/internal/platform/cache"
	"github.com/niamfauzi/go-starter/internal/platform/db"
	platformhttp "github.com/niamfauzi/go-starter/internal/platform/http"
	"github.com/niamfauzi/go-starter/internal/platform/logger"
	sharedvalidator "github.com/niamfauzi/go-starter/internal/shared/validator"
)

func main() {
	cfg := configs.MustLoad()

	appLogger, err := logger.New(cfg.AppEnv)
	if err != nil {
		panic(err)
	}
	defer func(appLogger *zap.Logger) {
		_ = appLogger.Sync()
	}(appLogger)

	gormDB, err := db.NewGorm(cfg.MySQLDSN, appLogger)
	if err != nil {
		appLogger.Fatal("gagal konek ke mysql", zap.Error(err))
	}

	redisClient := cache.NewRedis(cfg)
	validate := sharedvalidator.New()

	userRepo := user.NewRepository(gormDB)

	authRepo := auth.NewRepository(userRepo)
	jwtManager := auth.NewJWTManager(cfg.JWTSecret, cfg.JWTIssuer, cfg.JWTAccessTokenTTL)
	authService := auth.NewService(authRepo, jwtManager)
	authHandler := auth.NewHTTPDelivery(authService, validate, appLogger)

	productRepo := product.NewRepository(gormDB)
	productCache := product.NewCache(redisClient, 30*time.Second)
	productService := product.NewService(productRepo, productCache, appLogger)
	productHandler := product.NewHTTPDelivery(productService, validate, appLogger)

	router := platformhttp.NewRouter(appLogger, authHandler, productHandler, jwtManager)
	server := platformhttp.NewServer(cfg.AppPort, router)

	go func() {
		appLogger.Info("server berjalan", zap.String("port", cfg.AppPort))
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			appLogger.Fatal("server gagal berjalan", zap.Error(err))
		}
	}()

	// Graceful shutdown membuat server punya waktu untuk menutup koneksi dengan baik.
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	appLogger.Info("mematikan server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		appLogger.Error("graceful shutdown gagal", zap.Error(err))
	}

	appLogger.Info("server berhenti")
}
