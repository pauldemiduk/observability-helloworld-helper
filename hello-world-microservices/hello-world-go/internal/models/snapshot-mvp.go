// *
// ==> snapshot.go - [blueprint:include]
// *
// =>  snapshot.go explained:
//
//	Description: Structures and Data routines for Snapshot/Context foundation
package models

import (
	"fmt"
	"hello-world-go/internal/metrics"
	"log"
	"os"
	"strconv"

	"gopkg.in/yaml.v2"
)

//
// Define Data Structures Matching the YAML Schema
//

// Snapshot -
// Struct for Snapshot views.
type SnapshotMvp struct {
	ID        string       `yaml:"id" json:"id"`
	Name      string       `yaml:"name" json:"name"`
	CreatedAt string       `yaml:"createdAt" json:"createdAt"`
	Duration  string       `yaml:"duration" json:"duration"`
	Contexts  []ContextMvp `yaml:"contexts" json:"contexts"`
}

// Context -
// Struct for Context stats and details.
type ContextMvp struct {
	ID          string           `yaml:"id" json:"id"`
	Name        string           `yaml:"name" json:"name"`
	Description string           `yaml:"description" json:"description"`
	Queries     []QueryExecution `yaml:"queries" json:"queries"`
	UseCase     UseCase          `yaml:"useCase" json:"useCase"`
}

type UseCase struct {
	ID          string                 `yaml:"id" json:"id"`
	Name        string                 `yaml:"name" json:"name"`
	Description string                 `yaml:"description" json:"description"`
	Metrics     map[string]interface{} `yaml:"metrics,omitempty" json:"metrics,omitempty"`
}

//
// Context hydration content helpers
//

// ContextSlowAPICalls -
// Struct for Context use case hydration.
type ContextSlowAPICalls struct {
	ID           string         `yaml:"id" json:"id"`
	Name         string         `yaml:"name" json:"name"`
	Duration     int            `yaml:"duration" json:"duration"`
	SlowRequests int            `yaml:"slowRequests" json:"slowRequests"`
	P95Latency   float64        `yaml:"p95Latency" json:"p95Latency"`
	P99Latency   float64        `yaml:"p99Latency" json:"p99Latency"`
	TopEndpoints []SlowEndpoint `yaml:"topEndpoints" json:"topEndpoints"`
}

// ErrorRate -
// Struct for Context use case hydration.
type ContextErrorRates struct {
	ID              string           `yaml:"id"`
	Name            string           `yaml:"name"`
	Duration        int              `yaml:"duration" json:"duration"`
	TotalRequests   int              `yaml:"totalRequests"`
	FailedRequests  int              `yaml:"failedRequests"`
	ErrorRate       float64          `yaml:"errorRate"`
	TopErrorCodes   []ErrorCode      `yaml:"topErrorCodes"`
	FailedEndpoints []FailedEndpoint `yaml:"failedEndpoints"`
}

// ContextSLAViolation -
// Defines Go struct for SLA violation data
type ContextSLAViolation struct {
	ID              string `json:"id"`
	AffectedService string `json:"affectedService"`
	BreachTime      string `json:"breachTime"`
	Reason          string `json:"reason"`
	Impact          string `json:"impact"`
}

// ContextSLOCompliance -
// Defines Go struct for SLO compliance data
type ContextSLOCompliance struct {
	ID     string `json:"id"`
	Target string `json:"target"`
	Actual string `json:"actual"`
	Status string `json:"status"`
}

// ContextErrorBudget -
// Defines Go struct for error budget tracking
type ContextErrorBudget struct {
	ID                string `json:"id"`
	TotalAllocated    string `json:"totalAllocated"`
	Used              string `json:"used"`
	Remaining         string `json:"remaining"`
	DepletionForecast string `json:"depletionForecast"`
}

//
// Snapshot hydration content helpers
//

// GetSnapshots -
// Load the Snapshot and Context templates from datastore-snapshots.yaml
func GetSnapshots() []SnapshotMvp {
	log.Println("INFO: Fetching snapshots...") // Add log before calling the loader

	filename := "configs/datastore-snapshots.yaml" // Correct file path
	snapshots, err := LoadSnapshotsFromYAML(filename)

	if err != nil {
		log.Println("WARNING: Falling back to empty snapshot list")
		return []SnapshotMvp{}
	}

	log.Printf("INFO: Returning %d snapshots", len(snapshots)) // Confirm how many snapshots are returned
	return snapshots
}

func GetSlowAPICallSnapshot(duration int) (*ContextSlowAPICalls, error) {
	log.Printf("DEBUG: Fetching slow API call snapshot for duration: %d minutes", duration)

	// 🚀 Fetch slow request count
	slowRequests, err := metrics.GetMetricStats(duration, "slowRequests")
	if err != nil {
		log.Printf("ERROR: Failed to fetch slowRequests: %v", err)
		return nil, fmt.Errorf("failed to fetch slowRequests: %w", err)
	}
	log.Printf("DEBUG: Slow Requests: %v", slowRequests)

	// 🚀 Fetch P95 latency
	p95Latency, err := metrics.GetMetricStats(duration, "p95Latency")
	if err != nil {
		log.Printf("ERROR: Failed to fetch p95Latency: %v", err)
		return nil, fmt.Errorf("failed to fetch p95Latency: %w", err)
	}
	log.Printf("DEBUG: P95 Latency: %v", p95Latency)

	// 🚀 Fetch P99 latency
	p99Latency, err := metrics.GetMetricStats(duration, "p99Latency")
	if err != nil {
		log.Printf("ERROR: Failed to fetch p99Latency: %v", err)
		return nil, fmt.Errorf("failed to fetch p99Latency: %w", err)
	}
	log.Printf("DEBUG: P99 Latency: %v", p99Latency)

	// 🚀 Fetch top slow API endpoints dynamically
	log.Printf("DEBUG: Querying Prometheus for topEndpoints")
	nameValuePairs, err := metrics.ConvertVectorToNameValuePairs(duration, "topEndpoints")
	if err != nil {
		log.Printf("ERROR: Failed to fetch topEndpoints: %v", err)
		return nil, fmt.Errorf("failed to fetch topEndpoints: %w", err)
	}
	log.Printf("DEBUG: Retrieved %d topEndpoints", len(nameValuePairs))

	// ✅ Convert NameValuePair to SlowEndpoint
	var topEndpoints []SlowEndpoint
	for _, pair := range nameValuePairs {
		log.Printf("DEBUG: Processing topEndpoint: Name=%s, Value=%f", pair.Name, pair.Value)
		topEndpoints = append(topEndpoints, SlowEndpoint{
			Endpoint: pair.Name,
			Latency:  pair.Value,
		})
	}

	log.Printf("DEBUG: Constructing ContextSlowAPICalls struct with %d topEndpoints", len(topEndpoints))

	// 🚀 Construct ContextSlowAPICalls struct
	return &ContextSlowAPICalls{
		ID:           "slow-api-1",
		Name:         "Slow API Calls",
		Duration:     duration,
		SlowRequests: int(slowRequests.(float64) * 30 * 60), // Convert rate to approximate count
		P95Latency:   p95Latency.(float64),
		P99Latency:   p99Latency.(float64),
		TopEndpoints: topEndpoints,
	}, nil
}

// GetErrorRateSnapshot -
func GetErrorRateSnapshot(duration int) (*ContextErrorRates, error) {
	log.Printf("DEBUG: Fetching error rate snapshot for duration: %d minutes", duration)

	// 🚀 Fetch total requests
	totalRequests, err := metrics.GetMetricStats(duration, "totalRequests")
	if err != nil {
		log.Printf("ERROR: Failed to fetch totalRequests: %v", err)
		return nil, fmt.Errorf("failed to fetch totalRequests: %w", err)
	}
	log.Printf("DEBUG: Total Requests: %v", totalRequests)

	// 🚀 Fetch failed requests
	failedRequests, err := metrics.GetMetricStats(duration, "failedRequests")
	if err != nil {
		log.Printf("ERROR: Failed to fetch failedRequests: %v", err)
		return nil, fmt.Errorf("failed to fetch failedRequests: %w", err)
	}
	log.Printf("DEBUG: Failed Requests: %v", failedRequests)

	// 🚀 Fetch error rate
	errorRate, err := metrics.GetMetricStats(duration, "errorRate")
	if err != nil {
		log.Printf("ERROR: Failed to fetch errorRate: %v", err)
		return nil, fmt.Errorf("failed to fetch errorRate: %w", err)
	}
	log.Printf("DEBUG: Error Rate: %v", errorRate)

	// 🚀 Fetch top error codes dynamically
	log.Printf("DEBUG: Querying Prometheus for topErrorCodes")
	nameValueErrorCodes, err := metrics.ConvertVectorToNameValuePairs(duration, "topErrorCodes")
	if err != nil {
		log.Printf("ERROR: Failed to fetch topErrorCodes: %v", err)
		return nil, fmt.Errorf("failed to fetch topErrorCodes: %w", err)
	}
	log.Printf("DEBUG: Retrieved %d topErrorCodes", len(nameValueErrorCodes))

	// ✅ Convert NameValuePair to ErrorCode
	var topErrorCodes []ErrorCode
	for _, pair := range nameValueErrorCodes {
		code, err := strconv.Atoi(pair.Name)
		if err != nil {
			log.Printf("WARNING: Skipping invalid status code: %s", pair.Name)
			continue
		}
		log.Printf("DEBUG: Processing error code: %d with count: %f", code, pair.Value)

		topErrorCodes = append(topErrorCodes, ErrorCode{
			Code:  code,
			Count: int(pair.Value),
		})
	}

	// 🚀 Fetch failed endpoints dynamically
	log.Printf("DEBUG: Querying Prometheus for failedEndpoints")
	nameValueFailedEndpoints, err := metrics.ConvertVectorToNameValuePairs(duration, "failedEndpoints")
	if err != nil {
		log.Printf("ERROR: Failed to fetch failedEndpoints: %v", err)
		return nil, fmt.Errorf("failed to fetch failedEndpoints: %w", err)
	}
	log.Printf("DEBUG: Retrieved %d failedEndpoints", len(nameValueFailedEndpoints))

	// ✅ Convert NameValuePair to FailedEndpoint
	var failedEndpoints []FailedEndpoint

	for _, pair := range nameValueFailedEndpoints {
		failureCount := int(pair.Value) // ✅ Get the real failure count
		failureRate := 0.0

		if failureCount > 0 {
			failureRate = (float64(failureCount) / float64(totalRequests.(float64))) * 100

			if failureRate > 100 {
				log.Printf("WARNING: Adjusted failure rate for %s - Original: %.2f%%, Capped: 100%%",
					pair.Name, failureRate)
				failureRate = 100.0
			}
		}

		log.Printf("DEBUG: Failed Endpoint: %s | Failures: %d | Failure Rate: %.2f%%",
			pair.Name, failureCount, failureRate)

		failedEndpoints = append(failedEndpoints, FailedEndpoint{
			Endpoint:    pair.Name,
			Failures:    failureCount,
			FailureRate: failureRate, // ✅ This should now be between 0-100%
		})
	}

	log.Printf("DEBUG: Constructing ContextErrorRates struct with %d topErrorCodes and %d failedEndpoints", len(topErrorCodes), len(failedEndpoints))

	// 🚀 Construct ContextErrorRates struct
	return &ContextErrorRates{
		ID:              "error-rate-1",
		Name:            "API Failure Rate",
		Duration:        duration,
		TotalRequests:   int(totalRequests.(float64)),
		FailedRequests:  int(failedRequests.(float64)),
		ErrorRate:       errorRate.(float64),
		TopErrorCodes:   topErrorCodes,
		FailedEndpoints: failedEndpoints,
	}, nil
}

// GetSLAViolationSnapshot -
// Hydrates SLA Violation data from multiple Prometheus queries
func GetSLAViolationSnapshot(duration int) ([]ContextSLAViolation, error) {
	log.Printf("DEBUG: Fetching SLA violation snapshot for duration: %d minutes", duration)

	// Fetch affected services
	affectedServices, err := metrics.ConvertVectorToNameValuePairs(duration, "slaAffectedServices")
	if err != nil {
		log.Printf("ERROR: Failed to fetch affected services: %v", err)
		return nil, fmt.Errorf("failed to fetch affected services: %w", err)
	}

	// Construct SLA Violations response
	var violations []ContextSLAViolation
	for _, service := range affectedServices {
		violations = append(violations, ContextSLAViolation{
			ID:              fmt.Sprintf("sla-%s", service.Name),
			AffectedService: service.Name,
			BreachTime:      "2025-03-08T08:30:00Z",
			Reason:          "Latency exceeded SLO threshold",
			Impact:          "High",
		})
	}

	return violations, nil
}

// GetSLOComplianceSnapshot - Hydrates SLO compliance data
func GetSLOComplianceSnapshot(duration int) ([]ContextSLOCompliance, error) {
	log.Printf("DEBUG: Fetching SLO compliance snapshot for duration: %d minutes", duration)

	// Fetch compliance percentages
	uptimeCompliance, err := metrics.GetMetricStats(duration, "sloCompliance")
	if err != nil {
		log.Printf("ERROR: Failed to fetch SLO Compliance: %v", err)
		return nil, fmt.Errorf("failed to fetch SLO Compliance: %w", err)
	}

	// Construct SLO Compliance response
	compliance := []ContextSLOCompliance{
		{
			ID:     "slo-uptime",
			Target: "99.9% uptime",
			Actual: fmt.Sprintf("%.2f%%", uptimeCompliance.(float64)),
			Status: "Near Violation",
		},
	}

	return compliance, nil
}

// GetErrorBudgetSnapshot - Hydrates Error Budget data
func GetErrorBudgetSnapshot(duration int) ([]ContextErrorBudget, error) {
	log.Printf("DEBUG: Fetching Error Budget snapshot for duration: %d minutes", duration)

	// Fetch error budget remaining percentage
	remainingBudget, err := metrics.GetMetricStats(duration, "errorBudget")
	if err != nil {
		log.Printf("ERROR: Failed to fetch Error Budget: %v", err)
		return nil, fmt.Errorf("failed to fetch Error Budget: %w", err)
	}

	// Construct Error Budget response
	budget := []ContextErrorBudget{
		{
			ID:                "budget-remaining",
			TotalAllocated:    "100 hours",
			Used:              "76 hours",
			Remaining:         fmt.Sprintf("%.2f%%", remainingBudget.(float64)),
			DepletionForecast: "5 days",
		},
	}

	return budget, nil
}

// LoadSnapshotsFromYAML -
func LoadSnapshotsFromYAML(filePath string) ([]SnapshotMvp, error) {
	// Read the YAML file
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read snapshot datastore file: %w", err)
	}

	// Unmarshal YAML into a slice of Snapshots
	var snapshots []SnapshotMvp
	err = yaml.Unmarshal(yamlFile, &snapshots)
	if err != nil {
		return nil, fmt.Errorf("failed to parse snapshot YAML: %w", err)
	}

	return snapshots, nil
}
