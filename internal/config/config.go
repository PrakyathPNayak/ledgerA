package config

import (
	"fmt"
	"os"

	"github.com/joho/godotenv"
)

// Config represents the application configuration.
type Config struct {
	DatabaseURL         string
	Port                string
	GinMode             string
	FirebaseCredentials string
	FirebaseProjectID   string
	CorsAllowedOrigins  string
}

// Load reads config from environment variables.
func Load() (*Config, error) {
	_ = godotenv.Load() // ignore error, as .env might not exist in prod

	cfg := &Config{
		DatabaseURL:         os.Getenv("DATABASE_URL"),
		Port:                os.Getenv("PORT"),
		GinMode:             os.Getenv("GIN_MODE"),
		FirebaseCredentials: os.Getenv("FIREBASE_CREDENTIALS"),
		FirebaseProjectID:   os.Getenv("FIREBASE_PROJECT_ID"),
		CorsAllowedOrigins:  os.Getenv("CORS_ALLOWED_ORIGINS"),
	}

	if cfg.DatabaseURL == "" {
		return nil, fmt.Errorf("DATABASE_URL is required")
	}
	if cfg.Port == "" {
		cfg.Port = "8080"
	}

	return cfg, nil
}
