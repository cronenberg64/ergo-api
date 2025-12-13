package config

import (
	"os"

)

type Config struct {
	Port        string
	BackendURL  string
	JWTSecret   string
}

func Load() *Config {
	return &Config{
		Port:       getEnv("PORT", "8080"),
		BackendURL: getEnv("BACKEND_URL", "http://localhost:8081"),
		JWTSecret:  getEnv("JWT_SECRET", "secret"),
	}
}

func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
