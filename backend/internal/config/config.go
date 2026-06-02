package config

import (
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Host           string
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	LogLevel       string
	SeedAdminEmail string
	SeedAdminPass  string
	// StaticDir, when set and containing an index.html, makes the backend
	// also serve the built SPA (single-port monolith). Empty disables it.
	StaticDir string
}

// Load reads configuration from the environment. DATABASE_URL is intentionally
// read without a fallback because the database must be provided by the runtime.
func Load() *Config {
	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		Host: getEnv("HOST", "0.0.0.0"),
		Port: getEnv("PORT", "8080"),
		// The database is configured exclusively through the connection URL,
		// which already carries user, password, host and database name.
		DatabaseURL:    os.Getenv("DATABASE_URL"),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiry:      time.Duration(expiryHours) * time.Hour,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		SeedAdminEmail: getEnv("SEED_ADMIN_EMAIL", "admin@ordens.local"),
		SeedAdminPass:  getEnv("SEED_ADMIN_PASSWORD", "admin123"),
		StaticDir:      getEnv("STATIC_DIR", ""),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
