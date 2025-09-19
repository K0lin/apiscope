package config

import (
	"log"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

type Config struct {
	Port                    string
	DatabasePath            string
	DatabasePassword        string
	StoragePath             string
	LinkExpiration          time.Duration
	MaxFileSize             int64
	MaxVersions             int
	OpenAPIGeneratorEnabled bool
	OpenAPIGeneratorServer  string
}

func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using environment variables and defaults")
	}

	openAPIEnabled := getBoolEnv("OPENAPI_GENERATOR_ENABLED", false)
	openAPIServer := getEnv("OPENAPI_GENERATOR_SERVER", "https://api.openapi-generator.tech")

	return &Config{
		Port:                    getEnv("PORT", "8080"),
		DatabasePath:            getEnv("REDIS_ADDR", "localhost:6379"),
		DatabasePassword:        getEnv("REDIS_PASSWORD", ""),
		StoragePath:             getEnv("STORAGE_PATH", "./storage/documents"),
		LinkExpiration:          time.Hour * 24 * 30, // 30 days
		MaxFileSize:             50 * 1024 * 1024,    // 50 MB
		MaxVersions:             20,
		OpenAPIGeneratorEnabled: openAPIEnabled,
		OpenAPIGeneratorServer:  openAPIServer,
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getBoolEnv(key string, defaultValue bool) bool {
	if value := os.Getenv(key); value != "" {
		if parsed, err := strconv.ParseBool(value); err == nil {
			return parsed
		}
	}
	return defaultValue
}
