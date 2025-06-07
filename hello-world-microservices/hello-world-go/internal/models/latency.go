// *
// ==> latency.go - [blueprint:include]
// *
// =>  latency.go explained:
//
//	Description: Structures and Data routines for latency calculations & slow endpoints
package models

// ServiceLatency -
type ServiceLatency struct {
	P95Duration   float64        `json:"p95_duration_metric"` // 95th percentile latency
	P99Duration   float64        `json:"p99_duration_metric"` // 99th percentile latency
	TotalRequests int            `json:"total_requests"`      // Total request count
	SlowRequests  int            `json:"slow_requests"`       // Slow request count
	SlowEndpoints []SlowEndpoint `json:"slow_endpoints"`      // Top slow endpoints
}

//
// Context use case helpers
//

// SlowEndpoint -
// Struct for Context use case hydration.
type SlowEndpoint struct {
	Endpoint string  `yaml:"endpoint" json:"endpoint"`
	Latency  float64 `yaml:"latency" json:"latency"`
}
