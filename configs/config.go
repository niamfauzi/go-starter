package configs

import (
	"log"
	"os"
	"strconv"
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
		AppEnv:            getEnv("APP_ENV", "development"),
		AppPort:           getEnv("APP_PORT", "8080"),
		MySQLDSN:          getEnv("MYSQL_DSN", "root:root@tcp(localhost:3306)/go_starter?parseTime=true"),
		RedisAddr:         getEnv("REDIS_ADDR", "localhost:6379"),
		RedisPassword:     getEnv("REDIS_PASSWORD", ""),
		RedisDB:           getEnvInt("REDIS_DB", 0),
		JWTSecret:         getEnv("JWT_SECRET", "please-change-me"),
		JWTIssuer:         getEnv("JWT_ISSUER", "github.com/niamfauzi/go-starter"),
		JWTAccessTokenTTL: getEnvDuration("JWT_ACCESS_TOKEN_TTL", 15*time.Minute),
	}

	if cfg.JWTSecret == "" {
		log.Fatal("JWT_SECRET wajib diisi")
	}

	return cfg
}

func getEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

func getEnvInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	n, err := strconv.Atoi(val)
	if err != nil {
		log.Fatalf("nilai %s tidak valid: %v", key, err)
	}

	return n
}

func getEnvDuration(key string, fallback time.Duration) time.Duration {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	dur, err := time.ParseDuration(val)
	if err != nil {
		log.Fatalf("duration %s tidak valid: %v", key, err)
	}

	return dur
}
