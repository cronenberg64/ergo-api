package main

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/cronenberg64/ergo-api/internal/auth"
	"github.com/cronenberg64/ergo-api/internal/config"
	"github.com/cronenberg64/ergo-api/internal/observability"
	"github.com/cronenberg64/ergo-api/internal/policy"
	"github.com/cronenberg64/ergo-api/internal/proxy"
	"github.com/cronenberg64/ergo-api/internal/ratelimit"
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

	// Initialize Rate Limiter (100 req/min)
	rateLimiter := ratelimit.NewLimiter(cfg.RedisAddr, 100, 1*time.Minute)


	// Chain Middleware: Observability -> RateLimit -> Auth -> Threat -> Policy -> Proxy
	// handler = Obs(Rate(Auth(Threat(Policy(Proxy)))))
	
	policyHandler := policy.EnforcePolicy(policyEngine)(p)
	threatHandler := threat.AnalyzeRequest(threatDetector)(policyHandler)
	authHandler := auth.ValidateJWT(threatHandler, cfg.JWTSecret)
	rateHandler := ratelimit.Limit(rateLimiter)(authHandler)
	mainHandler := observability.Middleware(rateHandler)

	// Use a Mux to route /metrics separately (bypass middleware chain)
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	mux.Handle("/", mainHandler)

	log.Printf("Starting Ergo API on port %s", cfg.Port)
	if err := http.ListenAndServe(":"+cfg.Port, mux); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
