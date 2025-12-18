package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	BackendURL  string
	JWTSecret   string
	PolicyPath  string
	MaxRiskScore float64
	RedisAddr    string
}

func Load() *Config {
	// Load .env file if it exists
	_ = godotenv.Load()

	cfg := &Config{
		Port:         getEnv("PORT", "8080"),
		BackendURL:   getEnv("BACKEND_URL", "http://localhost:8081"),
		JWTSecret:    getEnv("JWT_SECRET", ""), // No default secret
		PolicyPath:   getEnv("POLICY_PATH", "policies/rbac.rego"),
		MaxRiskScore: 0.8,
		RedisAddr:    getEnv("REDIS_ADDR", "localhost:6379"),
	}

	if cfg.JWTSecret == "" {
		// Fail secure: do not start without a secret
		panic("JWT_SECRET environment variable is required")
	}

	return cfg
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
