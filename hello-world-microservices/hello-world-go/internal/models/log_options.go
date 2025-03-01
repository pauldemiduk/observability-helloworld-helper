// log_options.go - Defines structured log models.

package models

import (
	"net/http"
	"time"
)

// ResponseWriterInterface -
// Defines an interface that LoggingResponseWriter must implement.
type ResponseWriterInterface interface {
	http.ResponseWriter
	WriteHeader(statusCode int)
	Write(b []byte) (int, error)
	GetStatusCode() int // ✅ Allows retrieving HTTP status code
	GetSize() int       // ✅ Allows retrieving response size
}

// LogEntryOptions -
// Defines the attributes required for constructing a structured log entry.
// Provides metadata related to request lifecycle, tracing, and security.
type LogEntryOptions struct {
	LRW              ResponseWriterInterface // Response writer used for capturing status and response details
	Request          *http.Request           // Incoming HTTP request object
	StartTime        time.Time               // Timestamp when the request was received
	Duration         time.Duration           // Total time taken to process the request
	QueueTimeMS      float64                 //
	LogCtx           LogContext              // Context containing service and cloud metadata
	LogLevel         string                  // Log level (e.g., "info", "warn", "error")
	LogType          string                  // ✅ Defines whether this is a "request", "response", "execution", or "error" log
	RequestMetadata  RequestMetadata         // ✅ Structured request attributes
	ResponseMetadata ResponseMetadata        // ✅ Structured response attributes
	RequestStage     string                  // Stage of request processing ("received", "processed")
	Message          string                  // Human-readable log message
	ErrorMessage     string                  // ✅ Now added for error logs
	StackTrace       string                  // ✅ Now added for error logs
}
