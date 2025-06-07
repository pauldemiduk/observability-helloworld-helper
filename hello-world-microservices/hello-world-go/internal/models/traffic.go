// *
// ==> traffic.go - [blueprint:include]
// *
// =>  traffic.go explained:
//
//	Description: Structures and Data routines for test data generation
package models

// ServiceTraffic -
type ServiceTraffic struct {
	TotalRequests int           `json:"total_requests"` // Total request count
	RequestRate   float64       `json:"request_rate"`   // Requests per second (RPS)
	TopEndpoints  []TopEndpoint `json:"top_endpoints"`  // Most accessed API paths
	TopConsumers  []TopConsumer `json:"top_consumers"`  // Services generating most requests
}

// TopConsumer - Tracks API consumers with highest request volume
type TopConsumer struct {
	ClientID string `yaml:"client_id" json:"client_id"` // Unique API consumer ID
	Requests int    `yaml:"requests" json:"requests"`   // Total number of requests made
}

// TopEndpoint - Tracks API endpoints with highest request volume
type TopEndpoint struct {
	Endpoint string `yaml:"endpoint" json:"endpoint"`
	Requests int    `yaml:"requests" json:"requests"` // Total number of requests made
}
