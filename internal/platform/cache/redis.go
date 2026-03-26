package cache

import (
	"context"
	"time"

	"github.com/niamfauzi/go-starter/configs"
	"github.com/redis/go-redis/v9"
)

// NewRedis membuat koneksi Redis.
// Redis di project ini dipakai sederhana saja: sebagai cache.
func NewRedis(cfg configs.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     cfg.RedisAddr,
		Password: cfg.RedisPassword,
		DB:       cfg.RedisDB,
	})

	// Ping tidak wajib, tapi bagus untuk fail-fast saat startup.
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	_ = client.Ping(ctx).Err()

	return client
}
