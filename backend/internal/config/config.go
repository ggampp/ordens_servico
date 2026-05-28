package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config holds all runtime configuration loaded from environment variables.
type Config struct {
	Port           string
	DatabaseURL    string
	JWTSecret      string
	JWTExpiry      time.Duration
	LogLevel       string
	SeedAdminEmail string
	SeedAdminPass  string
}

// Load reads configuration from the environment, applying sensible defaults
// so the service can run locally without an exhaustive setup.
func Load() *Config {
	expiryHours, _ := strconv.Atoi(getEnv("JWT_EXPIRY_HOURS", "24"))

	return &Config{
		Port:           getEnv("PORT", "8080"),
		DatabaseURL:    buildDatabaseURL(),
		JWTSecret:      getEnv("JWT_SECRET", "change-me-in-production"),
		JWTExpiry:      time.Duration(expiryHours) * time.Hour,
		LogLevel:       getEnv("LOG_LEVEL", "info"),
		SeedAdminEmail: getEnv("SEED_ADMIN_EMAIL", "admin@ordens.local"),
		SeedAdminPass:  getEnv("SEED_ADMIN_PASSWORD", "admin123"),
	}
}

func buildDatabaseURL() string {
	if url := os.Getenv("DATABASE_URL"); url != "" {
		return url
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		getEnv("DB_USER", "ordens"),
		getEnv("DB_PASSWORD", "ordens"),
		getEnv("DB_HOST", "localhost"),
		getEnv("DB_PORT", "5432"),
		getEnv("DB_NAME", "ordens_servico"),
		getEnv("DB_SSLMODE", "disable"),
	)
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
