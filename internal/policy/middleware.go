package policy

import (
	"net/http"
	"strings"
)

func EnforcePolicy(engine *Engine) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Extract user role from context or token (mocked for now)
			// In a real scenario, this would come from the JWT claims in the context
			// For Phase 2 verification, we'll parse the "Authorization" header again or assume it's passed down.
			// Let's re-parse for simplicity or assume the Auth middleware put it in context.
			// Since we didn't implement context passing in Auth middleware yet, let's just grab it from header here for now.
			
			authHeader := r.Header.Get("Authorization")
			role := "guest"
			if authHeader != "" {
				parts := strings.Split(authHeader, " ")
				if len(parts) == 2 {
					token := parts[1]
					// Mock token decoding: "admin-token" -> role: admin, "user-token" -> role: user
					if token == "admin-token" {
						role = "admin"
					} else if token == "user-token" {
						role = "user"
					}
				}
			}

			input := map[string]interface{}{
				"method": r.Method,
				"path":   r.URL.Path,
				"role":   role,
			}

			allowed, err := engine.Evaluate(r.Context(), input)
			if err != nil {
				http.Error(w, "Policy Evaluation Error", http.StatusInternalServerError)
				return
			}

			if !allowed {
				http.Error(w, "Forbidden", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
