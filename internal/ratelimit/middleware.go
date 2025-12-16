package ratelimit

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

func Limit(l *Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Limit by IP
			ip := r.RemoteAddr
			host, _, err := net.SplitHostPort(r.RemoteAddr)
			if err == nil {
				ip = host
			}
			if forwarded := r.Header.Get("X-Forwarded-For"); forwarded != "" {
				ip = strings.Split(forwarded, ",")[0]
			}

			allowed, err := l.Allow(r.Context(), ip)
			if err != nil {
				// Fail open or closed? Let's fail open for now but log error
				// In production, might want to fail open to avoid outage if Redis is down
				// But for security gateway, maybe fail closed?
				// Let's log and allow for now to avoid blocking everyone if Redis dies
				fmt.Printf("Rate limit error for IP %s: %v\n", ip, err)
				next.ServeHTTP(w, r)
				return
			}

			if !allowed {
				http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
