package threat

import (
	"net/http"
	"strings"
)

type RequestFeatures struct {
	UserID   string
	IP       string
	Location string
}

func ExtractFeatures(r *http.Request) *RequestFeatures {
	// Extract UserID from header (mocked, assuming Auth middleware put it there or we parse token again)
	// For Phase 3, let's assume the token is in the header and we map it to a user.
	// Or simpler: use the token string itself as UserID for now.
	
	userID := "anonymous"
	authHeader := r.Header.Get("Authorization")
	if len(strings.Split(authHeader, " ")) == 2 {
		userID = strings.Split(authHeader, " ")[1]
	}

	ip := r.RemoteAddr
	// Allow overriding IP for testing
	if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
		ip = forwarded
	}

	return &RequestFeatures{
		UserID:   userID,
		IP:       ip,
		Location: mockGeoIP(ip),
	}
}

func mockGeoIP(ip string) string {
	// Mock GeoIP Database
	switch ip {
	case "1.1.1.1":
		return "NY"
	case "2.2.2.2":
		return "Tokyo"
	case "3.3.3.3":
		return "London"
	default:
		return "Unknown"
	}
}
