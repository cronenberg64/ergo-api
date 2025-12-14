package threat

import (
	"fmt"
	"net/http"
)

func AnalyzeRequest(detector *Detector) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			features := ExtractFeatures(r)
			riskScore := detector.CalculateRisk(features)

			// Add Risk Score to Header for visibility
			w.Header().Set("X-Risk-Score", fmt.Sprintf("%.2f", riskScore))

			if riskScore > detector.maxRiskScore {
				http.Error(w, "Blocked by Threat Detection", http.StatusForbidden)
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}
