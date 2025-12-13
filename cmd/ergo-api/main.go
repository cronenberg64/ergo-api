package main

import (
	"log"
	"net/http"

	"github.com/cronenberg64/ergo-api/internal/auth"
	"github.com/cronenberg64/ergo-api/internal/config"
	"github.com/cronenberg64/ergo-api/internal/proxy"
)

func main() {
	cfg := config.Load()

	p, err := proxy.NewProxy(cfg.BackendURL)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	handler := auth.ValidateJWT(p, cfg.JWTSecret)

	log.Printf("Starting Ergo API on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
