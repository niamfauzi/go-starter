package configs

import (
	"log"
	"os"
	"strconv"
	"time"
)

// GetEnv membaca string dari environment variable.
func GetEnv(key string, fallback string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return fallback
}

// GetEnvInt membaca integer dari environment variable.
func GetEnvInt(key string, fallback int) int {
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

// GetEnvDuration membaca duration dari environment variable.
func GetEnvDuration(key string, fallback time.Duration) time.Duration {
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
