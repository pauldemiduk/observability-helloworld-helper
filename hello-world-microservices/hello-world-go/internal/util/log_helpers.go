// http_helpers.go -
// Routines to help resolve LogEntry,
package util

import (
	"bufio"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"hello-world-go/internal/models"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"time"

	"go.opentelemetry.io/otel/trace"
)

// Cached container ID
var cachedContainerID string

// ValidateCloudProvider -
func ValidateCloudProvider(provider string) string {
	validProviders := map[string]bool{"aws": true, "gcp": true, "azure": true, "localhost": true}
	if !validProviders[provider] {
		log.Printf("WARNING: Invalid CLOUD_PROVIDER '%s', defaulting to 'localhost'", provider)
		return "localhost"
	}
	return provider
}

// ResolveHostName -
// Retrieves the host machine IP address
func ResolveHostName() string {
	host, err := os.Hostname()
	if err != nil {
		return "unknown"
	}
	return host
}

// GetProcessID -
// Retrieves service process id
func GetProcessID() int {
	return os.Getpid()
}

// GetThreadID -
// Placeholder, Go doesn’t expose true thread IDs
func GetThreadID() int {
	return os.Getpid()
}

// MeasureRequestQueueTime -
// Measures request queue time
func MeasureRequestQueueTime(startTime time.Time) float64 {
	return time.Since(startTime).Seconds()
}

// CaptureResourceUsage -
// Retrieves current CPU and memory usage in MB
func CaptureResourceUsage() (float64, float64) {
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)
	return float64(memStats.Sys) / (1024 * 1024), float64(memStats.Alloc) / (1024 * 1024)
}

// GetContainerID -
// Retrieves and caches the container ID from /proc/self/cgroup
func GetContainerID() string {
	if cachedContainerID != "" {
		return cachedContainerID
	}

	file, err := os.Open("/proc/self/cgroup")
	if err != nil {
		//log.Println("WARNING: Unable to read container ID, defaulting to 'unknown'")
		cachedContainerID = "unknown"
		return cachedContainerID
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(line, "docker") || strings.Contains(line, "kubepods") {
			parts := strings.Split(line, "/")
			for _, part := range parts {
				if len(part) == 64 {
					cachedContainerID = part[:12] // Cache first 12 characters (Docker standard)
					return cachedContainerID
				}
			}
		}
	}

	log.Println("DEBUG: Could not determine container ID, defaulting to 'unknown'")
	cachedContainerID = "unknown"
	return cachedContainerID
}

// ExtractClientIP -
// Extracts the IP and port from a request
func ExtractClientIP(r *http.Request) (string, int) {
	ip, portStr, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr, -1 // ✅ Return -1 instead of "unknown" for invalid ports
	}

	// ✅ Convert port from string to int
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return ip, -1 // ✅ Return -1 if port conversion fails
	}

	return ip, port
}

// ExtractParentTraceID -
// Retrieves the Parent Trace ID from request headers or OTEL context.
func ExtractParentTraceID(r *http.Request) string {
	traceParent := r.Header.Get("traceparent")
	if traceParent != "" {
		parts := strings.Split(traceParent, "-")
		if len(parts) >= 2 {
			return parts[1] // ✅ Extracts Parent Trace ID from traceparent format
		}
	}

	traceB3 := r.Header.Get("b3")
	if traceB3 != "" {
		parts := strings.Split(traceB3, "-")
		if len(parts) >= 2 {
			return parts[0] // ✅ Extracts Parent Trace ID from b3 format
		}
	}

	log.Println("WARN: No Parent Trace ID found in request headers")
	return "unknown"
}

// ExtractTraceContext -
// Retrieves Trace ID and Span ID from OpenTelemetry context
func ExtractTraceContext(r *http.Request) (string, string) {
	span := trace.SpanFromContext(r.Context())
	traceID := span.SpanContext().TraceID().String()
	spanID := span.SpanContext().SpanID().String()
	if traceID == "00000000000000000000000000000000" || traceID == "" {
		traceID = "unknown-trace"
	}
	if spanID == "0000000000000000" || spanID == "" {
		spanID = "unknown-span"
	}
	return traceID, spanID
}

// IsAuthHeaderPresent -
// Checks if an Authorization header is present
func IsAuthHeaderPresent(r *http.Request) string {
	if r.Header.Get("Authorization") == "" {
		return "missing"
	}
	return "present"
}

// CategorizeResponseTime -
// Maps response time to human readable buckets for easy SRE work, writes to LogEvent.ResponseTimeCategory
// UPLIFT.>, bucket levels belong in config
func CategorizeResponseTime(responseTimeMS float64, httpStatusCode int) string {
	// ✅ If HTTPStatusCode is -1 (request received but not responded), return "pending"
	if httpStatusCode == -1 {
		return "pending"
	}

	switch {
	case responseTimeMS <= models.HelloInstance.Config.ResponseTimeFast:
		return "fast"
	case responseTimeMS <= models.HelloInstance.Config.ResponseTimeModerate:
		return "moderate"
	case responseTimeMS <= models.HelloInstance.Config.ResponseTimeSlow:
		return "slow"
	default:
		return "very slow"
	}
}

// CategorizeResponseLatency -
// Maps response time to latency buckets for easy SRE work, writes to LogEvent.ResponseLatencyLevel
func CategorizeResponseLatency(responseTimeMS float64) string {
	switch {
	case responseTimeMS <= models.HelloInstance.Config.ResponseLatencyNormal:
		return "normal"
	case responseTimeMS <= models.HelloInstance.Config.ResponseLatencyWarning:
		return "warning"
	default:
		return "critical"
	}
}

// CategorizeHTTPStatusCode -
// Maps response code int values to human readable buckets for easy SRE work, writes to LogEvent.HTTPStatusCode
func CategorizeHTTPStatusCode(statusCode int) string {
	//log.Printf("DEBUG: http status code: %d", statusCode)
	switch {
	case statusCode >= 200 && statusCode < 300:
		//log.Printf("DEBUG: http status code: %d, success", statusCode)
		return "success"
	case statusCode >= 400 && statusCode < 500:
		return "client_error"
	case statusCode >= 500:
		return "server_error"
	default:
		return "unknown"
	}
}

// ExtractCallerFunction -
// Returns the name of the calling function, writes to LogEvent.CallerFunction
func ExtractCallerFunction() string {
	pc := make([]uintptr, 1)
	n := runtime.Callers(3, pc) // Get only one caller frame
	if n == 0 {
		return "unknown-function"
	}
	fn := runtime.FuncForPC(pc[0])
	if fn != nil {
		return fn.Name()[strings.LastIndex(fn.Name(), "/")+1:] // Extract function name
	}
	return "unknown-function"
}

// GetExecutionUnitID -
// (Go → Goroutines, Python → Threads/Async, C# → Tasks), writes to LogEvent.ExecutionUnitID
func GetExecutionUnitID() int {
	var buf [64]byte
	n := runtime.Stack(buf[:], false)
	stack := string(buf[:n])

	// Ensure brackets exist in the stack trace output
	idStart := strings.Index(stack, "goroutine ") + 10
	idEnd := strings.Index(stack[idStart:], " ")
	if idStart > 9 && idEnd > 0 { // Ensure valid indices
		id, err := strconv.Atoi(stack[idStart : idStart+idEnd])
		if err == nil {
			return id
		}
	}
	return -1 // Default to -1 if extraction fails
}

// ExtractCallerLocation -
func ExtractCallerLocation() string {
	_, file, line, ok := runtime.Caller(2)
	if !ok {
		return "unknown:0"
	}
	shortFile := file[strings.LastIndex(file, "/")+1:] // Extract only file name
	return fmt.Sprintf("%s:%d", shortFile, line)
}

// GetInstanceID - Determines a unique instance identifier dynamically
func GetInstanceID() string {
	// ✅ 1. Check Kubernetes Pod Name
	if podName := os.Getenv("KUBERNETES_POD_NAME"); podName != "" {
		return podName
	}

	// ✅ 2. Check Docker Container ID
	if containerID := getDockerContainerID(); containerID != "" {
		return containerID
	}

	// ✅ 3. Check AWS EC2 Instance Metadata
	if awsID := getAWSInstanceID(); awsID != "" {
		return awsID
	}

	// ✅ 4. Check GCP Instance Metadata
	if gcpID := getGCPInstanceID(); gcpID != "" {
		return gcpID
	}

	// ✅ 5. Generate a random UUID as fallback
	return generateRandomInstanceID()
}

// getDockerContainerID - Extracts the Docker container ID
func getDockerContainerID() string {
	data, err := os.ReadFile("/proc/self/cgroup")
	if err == nil {
		lines := strings.Split(string(data), "\n")
		for _, line := range lines {
			parts := strings.Split(line, "/")
			if len(parts) > 2 {
				return parts[len(parts)-1]
			}
		}
	}
	return ""
}

// getAWSInstanceID - Queries AWS EC2 Metadata API for Instance ID
func getAWSInstanceID() string {
	data, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_uuid")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

// getGCPInstanceID - Queries GCP Metadata API for Instance ID
func getGCPInstanceID() string {
	data, err := os.ReadFile("/sys/devices/virtual/dmi/id/product_serial")
	if err == nil {
		return strings.TrimSpace(string(data))
	}
	return ""
}

// generateRandomInstanceID - Generates a random UUID as a fallback
func generateRandomInstanceID() string {
	bytes := make([]byte, 16)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
