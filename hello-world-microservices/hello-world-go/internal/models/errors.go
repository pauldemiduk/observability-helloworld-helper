// *
// ==> errors.go - [blueprint:include]
// *
// =>  errors.go explained:
//
//	Description: Structures and Data routines for failed requests and errors
package models

// ServiceErrors -
type ServiceErrors struct {
	TotalRequests   int              `json:"total_requests"`   // Total request count
	FailedRequests  int              `json:"failed_requests"`  // Total number of errors (HTTP 4xx & 5xx)
	FailureRate     float64          `json:"failure_rate"`     // Requests per second (RPS)
	TopErrorCodes   []ErrorCode      `json:"top_errors"`       // Top error codes with counts
	FailedEndpoints []FailedEndpoint `json:"failed_endpoints"` // Top failed endpoints with latency and failure rate
}

// ErrorCode -
// Struct for Context case hydration.
type ErrorCode struct {
	Code  int `yaml:"code"`
	Count int `yaml:"count"`
}

// FailedEndpoint -
// Struct for Context case hydration.
type FailedEndpoint struct {
	Endpoint    string  `yaml:"endpoint"`
	Failures    int     `yaml:"failures"`
	FailureRate float64 `yaml:"failureRate"`
}
