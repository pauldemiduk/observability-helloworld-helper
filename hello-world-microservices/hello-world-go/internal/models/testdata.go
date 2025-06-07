// *
// ==> testdata.go - [blueprint:include]
// *
// =>  testdata.go explained:
//
//	Description: Structures and Data routines for test data generation
package models

import (
	"fmt"
	"hello-world-go/internal/metrics"
	"log"
	"math/rand"
	"time"
)

func GenerateTestLogs(duration time.Duration, logRatePerSec int) {
	interval := time.Second / time.Duration(logRatePerSec)
	jitter := time.Duration(rand.Intn(10)-5) * time.Millisecond
	time.Sleep(interval + jitter)

	endTime := time.Now().Add(duration)

	cloudProviders := []string{"aws", "gcp", "azure", "on-prem"}
	cloudRegions := []string{"us-east-1", "us-west-2", "europe-west1", "asia-southeast1"}
	cloudZones := []string{"availability-zone-1", "availability-zone-2", "availability-zone-3"}
	latencyLevels := []string{"fast", "normal", "warning", "critical"}

	//operations := []string{"op-1", "op-2", "op-3", "op-4", "op-5", "op-6"}
	sources := []string{"checkout-service", "payment-service", "user-service", "inventory-service", "auth-service"}
	categories := []string{"api_error", "db_error", "timeout", "security_violation", "validation_failure"}
	severities := []string{"info", "warning", "error", "critical"}
	endpoints := []string{"/checkout", "/cart", "/payment", "/order", "/inventory", "/user/profile", "/user/login", "/auth/token"}
	statusCodes := []int{200, 201, 400, 401, 403, 404, 408, 429, 500, 502, 503, 504}
	environments := []string{"dev", "staging", "prod"}
	requestMethods := []string{"GET", "POST", "PUT", "DELETE"}
	clientIDs := []string{"client-123", "client-456", "client-789", "client-999"}
	userAgents := []string{
		"Mozilla/5.0 (Windows NT 10.0; Win64; x64)",
		"Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)",
		"PostmanRuntime/7.26.8",
		"curl/7.68.0",
	}

	requestCounter := 0

	for time.Now().Before(endTime) {
		// 🚀 Keep existing request data randomization
		//operation := operations[rand.Intn(len(operations))]
		serviceName := sources[rand.Intn(len(sources))]
		requestURL := endpoints[rand.Intn(len(endpoints))]
		httpStatusCode := statusCodes[rand.Intn(len(statusCodes))]
		requestMethod := requestMethods[rand.Intn(len(requestMethods))]
		environment := environments[rand.Intn(len(environments))]
		cloudProvider := cloudProviders[rand.Intn(len(cloudProviders))]
		cloudRegion := cloudRegions[rand.Intn(len(cloudRegions))]
		cloudZone := cloudZones[rand.Intn(len(cloudZones))]
		latencyLevel := latencyLevels[rand.Intn(len(latencyLevels))]
		clientID := clientIDs[rand.Intn(len(clientIDs))]
		userAgent := userAgents[rand.Intn(len(userAgents))]

		// 🚀 Use updated simulateLatency() logic
		responseTimeMS := simulateLatency(latencyLevel, requestCounter)

		// Assign Log Level Based on HTTP Status Code
		var logLevel string
		switch {
		case httpStatusCode >= 500:
			logLevel = "critical"
		case httpStatusCode >= 400:
			logLevel = "error"
		default:
			logLevel = severities[rand.Intn(len(severities))]
		}

		// Assign ResponseOutcome based on HTTP Status Code
		var responseOutcome string
		switch {
		case httpStatusCode >= 200 && httpStatusCode < 300:
			responseOutcome = "success"
		case httpStatusCode >= 400 && httpStatusCode < 500:
			responseOutcome = "client_error"
		default:
			responseOutcome = "server_error"
		}

		// Assign ResponseLatencyLevel Based on ResponseTimeMS
		var responseLatencyLevel, responseTimeCategory string
		switch {
		case responseTimeMS <= 500.0:
			responseLatencyLevel = "normal"
			responseTimeCategory = "fast"
		case responseTimeMS <= 1000.0:
			responseLatencyLevel = "warning"
			responseTimeCategory = "normal"
		default:
			responseLatencyLevel = "critical"
			responseTimeCategory = "slow"
		}

		// Randomize AuthHeaderStatus
		authHeaderStatus := []string{"present", "missing"}[rand.Intn(2)]

		// Generate message
		message := fmt.Sprintf("Generated %s event from %s (%s)",
			logLevel, serviceName, categories[rand.Intn(len(categories))])

		// Store in GraphQL Data Store
		entry := LogEntry{
			Timestamp:            time.Now().UTC().Format(time.RFC3339Nano),
			Environment:          environment,
			LogLevel:             logLevel,
			ServiceName:          serviceName,
			RequestURL:           requestURL,
			ClientID:             clientID,
			UserAgent:            userAgent,
			CloudProvider:        cloudProvider,
			CloudRegion:          cloudRegion,
			AvailabilityZone:     cloudZone,
			HTTPStatusCode:       httpStatusCode,
			RequestMethod:        requestMethod,
			ResponseTimeMS:       responseTimeMS, // ✅ Uses float64
			ResponseOutcome:      responseOutcome,
			ResponseLatencyLevel: responseLatencyLevel,
			ResponseTimeCategory: responseTimeCategory,
			AuthHeaderStatus:     authHeaderStatus,
			Message:              message,
		}

		log.Printf("DEBUG: TestGen request_url=%s, request_method=%s", requestURL, requestMethod)

		// 🚀 Keep all existing logging and telemetry (unchanged)
		StoreLogForGraphQL(entry)
		metrics.IncrementRequestCount(serviceName, requestMethod, requestURL, environment, logLevel, cloudProvider, cloudRegion, cloudZone, clientID)
		metrics.IncrementResponseCount(serviceName, requestMethod, requestURL, environment, logLevel, cloudProvider, cloudRegion, cloudZone, clientID, httpStatusCode, responseTimeCategory, responseLatencyLevel)

		metrics.ObserveResponseTime(
			serviceName,
			requestMethod,
			requestURL,
			environment,
			logLevel,
			cloudProvider,
			cloudRegion,
			cloudZone,
			clientID,
			responseTimeCategory,
			responseLatencyLevel,
			responseOutcome,
			httpStatusCode,
			responseTimeMS)

		// Sleep to control log generation rate
		time.Sleep(interval)

		// Increment request counter for p99 tracking
		requestCounter++

		// 🚀 Trigger MDM Hydration after 500 requests
		if requestCounter == 500 {
			log.Println("INFO: Reached 500 requests, triggering MDM hydration...")
			// HydrateMDM()
		}
	}
}

// simulateLatency -
func simulateLatency(latencyLevel string, requestCounter int) float64 {
	// Define base latencies
	normalLatency := rand.Float64()*450 + 50 // 50.0 - 500.0 ms (fast response)
	p95Latency := rand.Float64()*400 + 800   // 800.0 - 1200.0 ms (p95 region)
	p99Latency := rand.Float64()*3000 + 2000 // 2000.0 - 5000.0 ms (p99 region)

	// 🚀 Smarter p99 injection logic (reduces over-weighting)
	if requestCounter >= 200 && requestCounter%50 == 0 {
		return p99Latency
	}

	// Normal distribution logic (80/15/5 rule)
	roll := rand.Float64()
	if roll < 0.85 {
		return normalLatency // 85% fast requests
	} else if roll < 0.98 {
		return p95Latency // 13% p95 range
	}
	return p99Latency // 2% p99 range
}

// weightedRandomChoice - (not-used)
// Picks a value based on weighted probability
func weightedRandomChoice(values []float64, weights []float64) float64 {
	r := rand.Float64()
	total := 0.0
	for i, w := range weights {
		total += w
		if r <= total {
			return values[i]
		}
	}
	return values[len(values)-1] // Return the slowest latency if nothing matches
}
