package config

import (
	"os"
	"time"
)

type Config struct {
	Port             string
	DatabasePath     string
	DatabasePassword string
	StoragePath      string
	LinkExpiration   time.Duration
	MaxFileSize      int64
	MaxVersions      int
}

func Load() *Config {
	return &Config{
		Port:             getEnv("PORT", "8080"),
		DatabasePath:     getEnv("REDIS_ADDR", "localhost:6379"),
		DatabasePassword: getEnv("REDIS_PASSWORD", ""),
		StoragePath:      getEnv("STORAGE_PATH", "./storage/documents"),
		LinkExpiration:   time.Hour * 24 * 30, // 30 days
		MaxFileSize:      50 * 1024 * 1024,    // 50 MB
		MaxVersions:      20,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
