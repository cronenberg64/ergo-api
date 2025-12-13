package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cronenberg64/ergo-api/internal/auth"
	"github.com/cronenberg64/ergo-api/internal/config"
	"github.com/cronenberg64/ergo-api/internal/policy"
	"github.com/cronenberg64/ergo-api/internal/proxy"
)

func main() {
	cfg := config.Load()

	p, err := proxy.NewProxy(cfg.BackendURL)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	// Initialize Policy Engine
	policyEngine, err := policy.NewEngine(context.Background(), cfg.PolicyPath)
	if err != nil {
		log.Fatalf("Failed to initialize policy engine: %v", err)
	}

	// Chain Middleware: Auth -> Policy -> Proxy
	// Note: We want Auth to run first to validate identity, then Policy to validate access.
	// Since we wrap in reverse order:
	// handler = Auth(Policy(Proxy))
	
	policyHandler := policy.EnforcePolicy(policyEngine)(p)
	handler := auth.ValidateJWT(policyHandler, cfg.JWTSecret)

	log.Printf("Starting Ergo API on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
