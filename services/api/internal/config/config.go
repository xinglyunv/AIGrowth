package config

import (
	"fmt"
	"os"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	DatabaseURL string
	RedisURL    string
	JWTSecret   string
	ServerPort  string
	AppEnv      string
	Pool        *pgxpool.Pool
}

// Load reads configuration from environment variables with sensible defaults.
func Load() (*Config, error) {
	cfg := &Config{
		DatabaseURL: envOrDefault("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/aige?sslmode=disable"),
		RedisURL:    envOrDefault("REDIS_URL", "redis://localhost:6379/0"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		ServerPort:  envOrDefault("SERVER_PORT", "8080"),
		AppEnv:      envOrDefault("APP_ENV", "development"),
	}

	if cfg.JWTSecret == "" {
		return nil, fmt.Errorf("JWT_SECRET environment variable is required")
	}

	return cfg, nil
}

func envOrDefault(key, defaultVal string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return defaultVal
}
