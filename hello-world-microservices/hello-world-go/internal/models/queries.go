package models

import (
	"log"
	"strconv"
)

// ✅ Function to Fetch Error Rates from Real Logs
func GetErrorRatesFromLogs() []ErrorRate {
	log.Println("DEBUG: Querying logs for error rates...")

	// ✅ Fetch Events (Ensure real data is retrieved)
	events := FetchRecentErrorEvents()
	log.Printf("INFO: Retrieved %d events", len(events)) // ✅ New log statement

	if len(events) == 0 {
		log.Println("WARNING: No error events found in logs!")
		return nil
	}

	// ✅ Process Events into Error Rate Data
	errorRates := processErrorRates(events)

	return errorRates
}

// ✅ Function to Process Events into Error Rate Data
// ✅ Function to Process Events into Error Rate Data
func processErrorRates(events []EventPayload) []ErrorRate {
	log.Println("DEBUG: Processing error rates from events...")

	// ✅ Create a map to count errors per code
	errorCodeCounts := make(map[int]int)
	totalRequests := 0
	failedRequests := 0

	for _, event := range events {
		totalRequests++

		// ✅ Extract error code from Labels (if exists)
		errorCode := 0
		if codeStr, exists := event.Labels["error_code"]; exists {
			code, err := strconv.Atoi(codeStr)
			if err == nil {
				errorCode = code
			}
		}

		// ✅ If there's no explicit error code, infer from severity
		if errorCode == 0 {
			switch event.Severity {
			case "critical":
				errorCode = 500 // Assume 500 for critical failures
			case "warning":
				errorCode = 400 // Assume 400 for warnings
			}
		}

		// ✅ Count valid errors
		if errorCode != 0 {
			errorCodeCounts[errorCode]++
			failedRequests++
		}
	}

	// ✅ Convert map into a list of ErrorCode objects
	var topErrorCodes []ErrorCode
	for code, count := range errorCodeCounts {
		topErrorCodes = append(topErrorCodes, ErrorCode{Code: code, Count: count})
	}

	// ✅ Compute error rate percentage (avoid division by zero)
	errorRate := 0.0
	if totalRequests > 0 {
		errorRate = (float64(failedRequests) / float64(totalRequests)) * 100
	}

	// ✅ Construct and return the ErrorRate struct
	return []ErrorRate{
		{
			ID:             "error-rate-realtime",
			Name:           "API Failure Rate (Real Logs)",
			TotalRequests:  totalRequests,
			FailedRequests: failedRequests,
			ErrorRate:      errorRate,
			TopErrorCodes:  topErrorCodes,
		},
	}
}
