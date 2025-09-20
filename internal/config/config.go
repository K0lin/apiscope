package config

import (
	"log"
	"os"
	"strconv"
	"strings"
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
	AllowVersionDeletion    bool
	AllowVersionDownload    bool
	AllowServerEditing      bool
	AutoAdjustServerOrigin  bool
	StripServers            bool
	AllowedOrigins          []string
	CORSAllowCredentials    bool
	CORSAllowedMethods      []string
	CORSAllowedHeaders      []string
	CORSExposeHeaders       []string
	CORSMaxAgeSeconds       int
	CORSDebug               bool
}

func Load() *Config {
	// Load .env file if it exists
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: No .env file found, using environment variables and defaults")
	}

	openAPIEnabled := getBoolEnv("OPENAPI_GENERATOR_ENABLED", false)
	openAPIServer := getEnv("OPENAPI_GENERATOR_SERVER", "https://api.openapi-generator.tech")
	allowVersionDeletion := getBoolEnv("ALLOW_VERSION_DELETION", false)
	allowVersionDownload := getBoolEnv("ALLOW_VERSION_DOWNLOAD", true)
	allowServerEditing := getBoolEnv("ALLOW_SERVER_EDITING", false)
	autoAdjustServerOrigin := getBoolEnv("AUTO_ADJUST_SERVER_ORIGIN", false)
	stripServers := getBoolEnv("STRIP_OPENAPI_SERVERS", false)
	allowedOriginsRaw := getEnv("ALLOWED_ORIGINS", "*")
	corsAllowCreds := getBoolEnv("CORS_ALLOW_CREDENTIALS", false)
	allowedMethodsRaw := getEnv("CORS_ALLOWED_METHODS", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	allowedHeadersRaw := getEnv("CORS_ALLOWED_HEADERS", "Authorization,Content-Type,Accept,Origin")
	exposeHeadersRaw := getEnv("CORS_EXPOSE_HEADERS", "Content-Length")
	corsMaxAgeStr := getEnv("CORS_MAX_AGE", "600")
	corsDebug := getBoolEnv("CORS_DEBUG", false)
	corsMaxAge := 600
	if v, err := strconv.Atoi(corsMaxAgeStr); err == nil && v >= 0 {
		corsMaxAge = v
	}

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
		AllowVersionDeletion:    allowVersionDeletion,
		AllowVersionDownload:    allowVersionDownload,
		AllowServerEditing:      allowServerEditing,
		AutoAdjustServerOrigin:  autoAdjustServerOrigin,
		StripServers:            stripServers,
		AllowedOrigins:          parseCSV(allowedOriginsRaw),
		CORSAllowCredentials:    corsAllowCreds,
		CORSAllowedMethods:      parseCSV(allowedMethodsRaw),
		CORSAllowedHeaders:      parseCSV(allowedHeadersRaw),
		CORSExposeHeaders:       parseCSV(exposeHeadersRaw),
		CORSMaxAgeSeconds:       corsMaxAge,
		CORSDebug:               corsDebug,
	}
}

func parseCSV(s string) []string {
	var res []string
	for _, part := range strings.Split(s, ",") {
		p := strings.TrimSpace(part)
		if p != "" {
			res = append(res, p)
		}
	}
	if len(res) == 0 {
		return []string{"*"}
	}
	return res
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
