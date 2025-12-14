package main

import (
	"context"
	"log"
	"net/http"

	"github.com/cronenberg64/ergo-api/internal/auth"
	"github.com/cronenberg64/ergo-api/internal/config"
	"github.com/cronenberg64/ergo-api/internal/policy"
	"github.com/cronenberg64/ergo-api/internal/proxy"
	"github.com/cronenberg64/ergo-api/internal/threat"
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

	// Initialize Threat Detector
	threatDetector := threat.NewDetector(cfg.MaxRiskScore)

	// Chain Middleware: Auth -> Threat -> Policy -> Proxy
	// handler = Auth(Threat(Policy(Proxy)))
	
	policyHandler := policy.EnforcePolicy(policyEngine)(p)
	threatHandler := threat.AnalyzeRequest(threatDetector)(policyHandler)
	handler := auth.ValidateJWT(threatHandler, cfg.JWTSecret)

	log.Printf("Starting Ergo API on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, handler); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
