// *
// ==> graph_filter.go - [blueprint:include]
// *
// =>  graph_filter.go explained:
//
//	Description: GraphQL routines to filter in-memory datastore.
//
// # GraphQL composition
//
//	Schema
//	- Entity
//	-- Interface
//	--- Query
//	---- Filter  <= (graph_filter.go)
package models

import (
	"log"
	"time"
)

// ResponseWriterInterface -
// Defines an interface that LoggingResponseWriter must implement.
/* type ResponseWriterInterface interface {
	http.ResponseWriter
	WriteHeader(statusCode int)
	Write(b []byte) (int, error)
	GetStatusCode() int
	GetSize() int
} */

// LogFilterInput -
// Define LogFilterInput Struct - Supports GraphQL UI query filters
// ** ROADMAP.>, more filters on LogEntry fields
type LogFilterInput struct {
	Since                string `json:"since"`
	Level                string `json:"level"`
	HTTPStatusCode       int    `json:"httpStatusCode"`
	ResponseTimeCategory string `json:"responseTimeCategory"`
	RequestMethod        string `json:"requestMethod"`
	ResponseLatencyLevel string `json:"responseLatencyLevel"`
	ServiceName          string `json:"serviceName"` // not-implemented
	TraceID              string `json:"traceId"`     // not-implemented
}

// ParseLogFilter -
// Parse the GraphQL log input filter
func ParseLogFilter(filterMap map[string]interface{}) LogFilterInput {
	var filter LogFilterInput

	if level, ok := filterMap["level"].(string); ok {
		filter.Level = level
	}
	if serviceName, ok := filterMap["serviceName"].(string); ok {
		filter.ServiceName = serviceName
	}
	if traceID, ok := filterMap["traceId"].(string); ok {
		filter.TraceID = traceID
	}
	if since, ok := filterMap["since"].(string); ok {
		filter.Since = since
	}
	if httpStatusCode, ok := filterMap["httpStatusCode"].(int); ok {
		filter.HTTPStatusCode = httpStatusCode
	}
	if responseTimeCategory, ok := filterMap["responseTimeCategory"].(string); ok {
		filter.ResponseTimeCategory = responseTimeCategory
	}
	if requestMethod, ok := filterMap["requestMethod"].(string); ok {
		filter.RequestMethod = requestMethod
	}
	if responseLatencyLevel, ok := filterMap["responseLatencyLevel"].(string); ok {
		filter.ResponseLatencyLevel = responseLatencyLevel
	}

	return filter
}

// filterLogsByTime -
// Filters logs by timestamp
func filterLogsByTime(logs []LogEntry, since string) []LogEntry {
	var filtered []LogEntry

	// Parse the "since" timestamp
	sinceTime, err := time.Parse(time.RFC3339Nano, since) // Ensures correct time format
	if err != nil {
		log.Printf("WARNING: Invalid time format for 'since': %s", since)
		return logs // Return all logs if parsing fails
	}

	// Filter logs
	for _, logEntry := range logs {
		entryTime, err := time.Parse(time.RFC3339Nano, logEntry.Timestamp) // Re-added missing line
		if err == nil && entryTime.After(sinceTime) {
			filtered = append(filtered, logEntry)
		}
	}
	return filtered
}

// filterLogsByLevel -
// Filters logs by level (info, debug, warn, .. )
func filterLogsByLevel(logs []LogEntry, level string) []LogEntry {
	var filtered []LogEntry
	for _, log := range logs {
		if log.LogLevel == level {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// filterLogsByStatusCode -
// Filters logs by HTTP status code
func filterLogsByStatusCode(logs []LogEntry, statusCode int) []LogEntry {
	var filtered []LogEntry
	log.Printf("DEBUG: Filtering for HTTP Status Code: %d", statusCode)

	for _, entry := range logs { // Renamed loop variable from `log` to `entry`
		log.Printf("DEBUG: Checking Log - HTTPStatusCode: %d", entry.HTTPStatusCode)
		if entry.HTTPStatusCode == statusCode {
			log.Printf("✅ MATCH FOUND: %d", entry.HTTPStatusCode)
			filtered = append(filtered, entry)
		}
	}

	log.Printf("DEBUG: Total Matches Found: %d", len(filtered))
	return filtered
}

// filterLogsByResponseTimeCategory -
// Filters logs by response time category (fast, normal, slow)
func filterLogsByResponseTimeCategory(logs []LogEntry, category string) []LogEntry {
	var filtered []LogEntry
	for _, log := range logs {
		if log.ResponseTimeCategory == category {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// filterLogsByRequestMethod -
// Filters logs by HTTP request method (GET, POST, etc.)
func filterLogsByRequestMethod(logs []LogEntry, method string) []LogEntry {
	var filtered []LogEntry
	for _, log := range logs {
		if log.RequestMethod == method {
			filtered = append(filtered, log)
		}
	}
	return filtered
}

// filterLogsByResponseLatencyLevel -
// Filters logs based on latency severity (normal, warning, critical)
func filterLogsByResponseLatencyLevel(logs []LogEntry, level string) []LogEntry {
	var filtered []LogEntry
	for _, entry := range logs { // Change loop variable to 'entry'
		log.Printf("DEBUG: Checking log [%s] with response_latency_level: %s", entry.Timestamp, entry.ResponseLatencyLevel)
		if entry.ResponseLatencyLevel == level {
			log.Printf("DEBUG: ✅ MATCH FOUND for response_latency_level = %s", level)
			filtered = append(filtered, entry)
		}
	}
	log.Printf("DEBUG: Filtered %d logs matching response_latency_level = %s", len(filtered), level)
	return filtered
}
