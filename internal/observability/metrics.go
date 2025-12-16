package observability

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var (
	RequestsTotal = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "ergo_requests_total",
			Help: "The total number of processed requests",
		},
		[]string{"method", "path", "status"},
	)

	RequestDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "ergo_request_duration_seconds",
			Help: "The duration of requests in seconds",
		},
		[]string{"method", "path"},
	)
)
