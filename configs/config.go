package configs

import (
	"log"
	"time"
)

// Config menyimpan semua konfigurasi aplikasi di satu tempat.
// Ini memudahkan kita saat ingin mengganti environment development / production.
type Config struct {
	AppEnv            string
	AppPort           string
	MySQLDSN          string
	RedisAddr         string
	RedisPassword     string
	RedisDB           int
	JWTSecret         string
	JWTIssuer         string
	JWTAccessTokenTTL time.Duration
}

// MustLoad membaca environment variable.
// Jika ada konfigurasi penting yang salah, aplikasi akan dihentikan lebih awal.
func MustLoad() Config {
	cfg := Config{
		AppEnv:            GetEnv("APP_ENV", "development"),
		AppPort:           GetEnv("APP_PORT", "8080"),
		MySQLDSN:          GetEnv("MYSQL_DSN", "root:root@tcp(localhost:3306)/go_starter?parseTime=true"),
		RedisAddr:         GetEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     GetEnv("REDIS_PASSWORD", ""),
		RedisDB:           GetEnvInt("REDIS_DB", 0),
		JWTSecret:         GetEnv("JWT_SECRET", "please-change-me"),
		JWTIssuer:         GetEnv("JWT_ISSUER", "github.com/niamfauzi/go-starter"),
		JWTAccessTokenTTL: GetEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET wajib diisi")
	}

	return cfg
}
