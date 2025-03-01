// metrics.go - Enhanced Prometheus metrics setup

package metrics

import (
	"log"
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

// Prometheus metrics with labels for better observability
var (
	// ✅ NEW: Service uptime counter
	ServiceUptime = prometheus.NewCounter(
		prometheus.CounterOpts{
			Name: "service_uptime_seconds",
			Help: "Total uptime of the service in seconds",
		},
	)

	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "Total number of HTTP requests received",
		},
		[]string{"operation", "method", "status_code"},
	)

	httpRequestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_duration_seconds",
			Help:    "Histogram of response time for requests",
			Buckets: []float64{0.001, 0.005, 0.01, 0.025, 0.05, 0.1, 0.25, 0.5, 1, 2, 5},
		},
		[]string{"operation", "method"},
	)

	httpRequestSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_request_size_bytes",
			Help:    "Histogram of request sizes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"operation", "method"},
	)

	httpResponseSizeBytes = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "http_response_size_bytes",
			Help:    "Histogram of response sizes",
			Buckets: prometheus.ExponentialBuckets(100, 2, 10),
		},
		[]string{"operation", "method", "status_code"},
	)

	metricsOnce sync.Once
)

// SetupMetrics ensures Prometheus metrics are initialized only once
func SetupMetrics() {
	metricsOnce.Do(func() {
		prometheus.MustRegister(ServiceUptime, httpRequestsTotal, httpRequestDuration, httpRequestSizeBytes, httpResponseSizeBytes)
		ServiceUptime.Add(1) // ✅ Start tracking uptime immediately
		log.Println("INFO: Prometheus metrics registered")
	})
}
