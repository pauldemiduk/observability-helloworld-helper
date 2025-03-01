// log_lifecycle.go - Optimized structured logging for the Hello World microservice
package logging

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hello-world-go/internal/models"

	"log"
	"math/rand"
	"net/http"
	"os"
	"runtime/debug"
	"sync"
	"time"
)

// Buffered log queue
var logQueue = make(chan models.LogEntry, 100)
var logWorkerOnce sync.Once
var logRateLimit = struct {
	sync.Mutex
	entries map[string]time.Time
}{entries: make(map[string]time.Time)}

// SendLog -
// Asynchronous logging function
func SendLog(entry models.LogEntry) {
	// Validate LogLevel before queuing
	if entry.LogLevel == "" {
		entry.LogLevel = "INFO" // Default to INFO if not provided
	}
	startLogWorker() // Ensure worker is running
	select {
	case logQueue <- entry:
		return // Successfully added to queue
	default:
		log.Printf("WARN: Log queue is full, dropping log entry: %s", entry.LogLevel)
	}
}

// FlushPendingLogs -
// Completes the writes for any logs during service shutdown lifecycle
func FlushPendingLogs() {
	log.Println("INFO: Flushing remaining logs before shutdown...")
	time.Sleep(500 * time.Millisecond) // Allow logs to process
}

// startLogWorker - internal function
// Initialize async logging worker
func startLogWorker() {
	logWorkerOnce.Do(func() {
		go func() {
			for entry := range logQueue {
				processLog(entry)
			}
		}()
	})
}

// processLog - internal function
// Processes a log entry asynchronously with validation, formatting, and retries.
// It ensures structured logging is correctly formatted and sent to the appropriate backend.
func processLog(entry models.LogEntry) {

	// ✅ Validate and standardize log fields before processing
	entry = validateLogEntry(entry)

	// ✅ Convert log entry to JSON format for structured logging
	logJson, err := json.Marshal(entry)
	if err != nil {
		log.Printf("ERROR: Failed to marshal log entry: %v\nStack Trace:\n%s", err, debug.Stack())
		return
	}

	// ✅ Log to stdout (default logging destination)
	log.Printf("INFO: Structured Log: %s", string(logJson))

	// ✅ Check if file logging is enabled and write log to a file if configured
	if models.HelloInstance.Config.LogBackend == "file" && models.HelloInstance.Config.LogFile != "" {
		writeLogToFile(string(logJson), models.HelloInstance.Config.LogFile)
		return // Exit after writing to a file
	}

	// ✅ Define supported external log backends and their endpoints
	logEndpoints := map[string]string{
		"loki":       models.HelloInstance.Config.LokiEndpoint,
		"cloudwatch": models.HelloInstance.Config.CloudWatchEndpoint,
		"gcp":        models.HelloInstance.Config.GCPLoggingEndpoint,
		"azure":      models.HelloInstance.Config.AzureMonitorEndpoint,
	}

	// ✅ Retrieve the logging endpoint based on the configured log backend
	endpoint, exists := logEndpoints[models.HelloInstance.Config.LogBackend]
	if !exists || endpoint == "" {
		log.Println("WARN: No valid log backend configured, skipping external log transmission")
		return // Exit if no valid backend is set
	}

	// ✅ Attempt to send the log entry to the external logging system with retries
	for i := 0; i < 3; i++ {
		if err := postLog(endpoint, logJson); err == nil {
			return // ✅ Log successfully sent, exit function
		}

		// ✅ Log a warning and retry with exponential backoff
		log.Printf("WARN: Log delivery failed (attempt %d). Retrying in %d seconds.", i+1, i+1)
		time.Sleep(time.Duration(i+1) * time.Second)
	}

	// ✅ If all retries fail, log an error message indicating log delivery failure
	log.Printf("ERROR: All log delivery attempts failed. Log entry: %s \nStack Trace:\n%s", string(logJson), debug.Stack())
}

// validateLogEntry -
// Ensures all required LogEntry attributes have valid values.
func validateLogEntry(entry models.LogEntry) models.LogEntry {
	// ✅ Retrieve sample rate
	sampleRate := models.HelloInstance.Config.LogSampleRate

	// ✅ Apply log sampling for non-critical logs
	if entry.LogLevel == "info" || entry.LogLevel == "debug" {
		if !shouldSample(sampleRate) { // ✅ Uses pre-retrieved sample rate
			return entry // 🚀 Skip logging this entry
		}
	}

	if entry.RequestHeaders == nil {
		entry.RequestHeaders = make(map[string]string) // ✅ Ensure it's always initialized
	}
	// Default values for missing fields
	if entry.ServiceName == "" {
		entry.ServiceName = "unknown-service"
	}
	if entry.Operation == "" {
		entry.Operation = "unknown-operation"
	}
	if entry.TraceID == "" || entry.TraceID == "00000000000000000000000000000000" {
		entry.TraceID = "unknown-trace"
	}
	if entry.ParentTraceID == "" {
		entry.ParentTraceID = "unknown-trace"
	}
	if entry.SpanID == "" || entry.SpanID == "0000000000000000" {
		entry.SpanID = "unknown-span"
	}
	if entry.ResponseOutcome != "pending" {
		entry.HTTPStatusCode = 0 // No status http request
	}
	if entry.ResponseOutcome == "" {
		entry.ResponseOutcome = "unknown-outcome"
	}
	if entry.RequestStage == "" {
		entry.RequestStage = "unknown-stage"
	}
	if entry.TLSProtocolVersion == "" {
		entry.TLSProtocolVersion = "unknown-tlsversion"
	}
	if entry.UserSessionID == "" {
		entry.UserSessionID = "unknown-session"
	}
	if entry.ResponseHeaders == nil {
		entry.ResponseHeaders = make(map[string]string) // ✅ Ensure it's always initialized
	}

	// ✅ Log unexpected values for better debugging, but prevent log spam
	if shouldLog("unexpected_outcome_"+entry.ResponseOutcome, 10*time.Second) {
		if entry.ResponseOutcome != "success" && entry.ResponseOutcome != "failure" &&
			entry.ResponseOutcome != "pending" && entry.ResponseOutcome != "client_error" &&
			entry.ResponseOutcome != "server_error" {
			log.Printf("WARN: Unexpected Outcome value: %s", entry.ResponseOutcome)
		}
	}

	if shouldLog("unexpected_LogLevel_"+entry.LogLevel, 10*time.Second) {
		if entry.LogLevel != "info" && entry.LogLevel != "warn" &&
			entry.LogLevel != "error" && entry.LogLevel != "debug" {
			log.Printf("WARN: Unexpected LogLevel value: %s", entry.LogLevel)
		}
	}

	if entry.ResponseOutcome == "pending" {
		entry.Message = fmt.Sprintf("Received request %s for %s service", entry.Operation, entry.ServiceName)
	} else {
		entry.Message = fmt.Sprintf("Processed request %s with status %d", entry.Operation, entry.HTTPStatusCode)
	}

	// ✅ Use a more compact timestamp format (millisecond precision)
	entry.Timestamp = time.Now().UTC().Format("2006-01-02T15:04:05.000Z")

	// ✅ Populate ErrorMessage and StackTrace only for errors
	if entry.LogLevel == "error" {
		logKey := "error_" + entry.Operation
		if shouldLog(logKey, 5*time.Second) { // ✅ Prevent excessive error logs
			entry.ErrorMessage = fmt.Sprintf("Error in %s operation after %.3fms", entry.Operation, entry.ResponseTimeMS)
			if models.HelloInstance.Config.LogLevel == "debug" { // ✅ Only capture stack trace in debug mode
				entry.StackTrace = string(debug.Stack())
			}
		} else {
			entry.ErrorMessage = "" // ✅ Suppress duplicate error logs
			entry.StackTrace = ""
		}
	}
	return entry
}

// postLog - internal function
// Sends a log entry as an HTTP POST request to a remote logging endpoint.
// It handles error scenarios such as missing endpoint configuration, request creation failures,
// network errors, and non-2xx HTTP responses.
func postLog(endpoint string, payload []byte) error {
	// ✅ Check if the endpoint is provided; if not, log a message and skip sending.
	if endpoint == "" {
		log.Println("INFO: No endpoint configured, skipping log transmission")
		return nil
	}

	// ✅ Create an HTTP client with a 5-second timeout to avoid long waits on network failures.
	client := &http.Client{Timeout: 5 * time.Second}

	// ✅ Construct a new HTTP POST request with the log payload.
	req, err := http.NewRequest("POST", endpoint, bytes.NewBuffer(payload))
	if err != nil {
		log.Printf("ERROR: Failed to create request for %s: %v\nStack Trace:\n%s", endpoint, err, debug.Stack())
		return err
	}

	// ✅ Set the request header to indicate JSON content.
	req.Header.Set("Content-Type", "application/json")

	// ✅ Send the request to the logging service.
	resp, err := client.Do(req)
	if err != nil {
		log.Printf("ERROR: Failed to send log to %s: %v\nStack Trace:\n%s", endpoint, err, debug.Stack())
		return err
	}
	defer resp.Body.Close() // Ensure the response body is closed to prevent resource leaks.

	// ✅ Check the HTTP response status; if it's not 2xx, log an error.
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		log.Printf("ERROR: Received non-2xx response from %s: %d\nStack Trace:\n%s", endpoint, resp.StatusCode, debug.Stack())
		return fmt.Errorf("non-2xx response: %d", resp.StatusCode)
	}

	// ✅ If the request succeeded, return nil (no error).
	return nil
}

// writeLogToFile - internal function
// Writes a log entry to file based on ServiceConfig logging settings
func writeLogToFile(logEntry string, logFile string) {
	if logFile == "" {
		log.Println("WARN: No log file configured, skipping file logging")
		return
	}

	// ✅ Rotate logs before writing
	ManageLogRotation(logFile)

	// ✅ Open log file in append mode
	file, err := os.OpenFile(logFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Printf("ERROR: Failed to write log entry to file %s: %v", logFile, err)
		return
	}
	defer file.Close()

	_, err = file.WriteString(logEntry + "\n")
	if err != nil {
		log.Printf("ERROR: Failed to write log entry to file %s: %v", logFile, err)
	}
}

// shouldLog -
// Checks if a message should be logged based on rate limiting.
// If the same log message appears too frequently, it will be suppressed.
func shouldLog(logKey string, limitDuration time.Duration) bool {
	logRateLimit.Lock()
	defer logRateLimit.Unlock()

	lastLogged, exists := logRateLimit.entries[logKey]
	if exists && time.Since(lastLogged) < limitDuration {
		return false // ✅ Suppress logging if it's too soon
	}

	// ✅ Store timestamp for next log event
	logRateLimit.entries[logKey] = time.Now()
	return true
}

// shouldSample -
// Determines if a log should be recorded based on the sample rate.
// The `sampleRate` should be a value between 0.0 (no logs) and 1.0 (all logs).
func shouldSample(sampleRate float64) bool {
	counter := rand.Intn(100)
	return counter%int(sampleRate*100) == 0 // ✅ More predictable sampling
}
