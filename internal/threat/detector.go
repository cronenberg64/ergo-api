package threat

import (
	"sync"
	"time"
)

type Detector struct {
	mu           sync.Mutex
	userHistory  map[string]*UserHistory // Key: UserID
	maxRiskScore float64
}

type UserHistory struct {
	LastSeen     time.Time
	LastLocation string // Mocked: "NY", "Tokyo", etc.
	LastIP       string
}

func NewDetector(maxRiskScore float64) *Detector {
	return &Detector{
		userHistory:  make(map[string]*UserHistory),
		maxRiskScore: maxRiskScore,
	}
}

func (d *Detector) CalculateRisk(features *RequestFeatures) float64 {
	d.mu.Lock()
	defer d.mu.Unlock()

	risk := 0.0

	// 1. Impossible Travel Check
	if history, exists := d.userHistory[features.UserID]; exists {
		// Simple mock logic: If location changed and time diff is small
		if history.LastLocation != features.Location {
			// Assume any location change requires at least 1 hour
			if time.Since(history.LastSeen) < 1*time.Hour {
				risk += 1.0 // Instant block
			}
		}
	}

	// Update History
	d.userHistory[features.UserID] = &UserHistory{
		LastSeen:     time.Now(),
		LastLocation: features.Location,
		LastIP:       features.IP,
	}

	return risk
}
