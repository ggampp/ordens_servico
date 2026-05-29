package config

import (
	"os"
	"strconv"
	"time"
)

// defaultDatabaseURL is used for local development when DATABASE_URL is unset.
const defaultDatabaseURL = "postgres://ordens_servico:ordens_servico@localhost:5432/ordens_servico?sslmode=disable"

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
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

// Load reads configuration from the environment, applying sensible defaults
// so the service can run locally without an exhaustive setup.
func Load() *Config {
	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		Port: getEnv("PORT", "8080"),
		// The database is configured exclusively through the connection URL,
		// which already carries user, password, host and database name.
		DatabaseURL:    getEnv("DATABASE_URL", defaultDatabaseURL),
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
