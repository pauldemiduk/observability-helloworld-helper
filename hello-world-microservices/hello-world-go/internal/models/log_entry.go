// log_entry.go - Defines structured log models.

package models

// LogEntry -
// Represents a structured log entry with categorized attributes.
// ** UPLIFT.>, refined error capture?
type LogEntry struct {
	// Core metadata
	Timestamp   string `json:"timestamp"`             // Value in microseconds, includes timezone
	Environment string `json:"environment,omitempty"` // Application environment (dev, staging, prod)
	LogLevel    string `json:"log_level,omitempty"`   // Log level: info, warn, error, critical, fatal

	// ✅ Service Metadata
	ServiceName     string `json:"service_name,omitempty"`     // Name of the service writing the log
	ServiceVersion  string `json:"service_version,omitempty"`  // Deployed service version
	ServiceHostname string `json:"service_hostname,omitempty"` // Hostname where the service is running

	// ✅ Service-level details (Where the service is accessible)
	ServiceIP   string `json:"service_ip,omitempty"`   // IP address of the service (internal or external)
	ServicePort int    `json:"service_port,omitempty"` // Port on which the service is exposed
	ProcessID   int    `json:"process_id,omitempty"`   // Process ID of the running service
	ThreadID    int    `json:"thread_id,omitempty"`    // Identifier for the processing thread
	InstanceID  string `json:"instance_id,omitempty"`  // Unique ID of the running service instance
	Language    string `json:"language,omitempty"`     // Runtime language: Golang, Python, C#

	// ✅ Host-level details (Where the service is running)
	NodeIP   string `json:"node_ip,omitempty"`   // IP address of the host (VM, EC2, Kubernetes Node) where the service runs
	NodePort int    `json:"node_port,omitempty"` // Port on the host machine mapped to the service

	// Cloud & Infrastructure details
	CloudProvider    string `json:"cloud_provider,omitempty"`    // Cloud provider: localhost, AWS, GCP, Azure
	CloudRegion      string `json:"cloud_region,omitempty"`      // Cloud network region
	AvailabilityZone string `json:"availability_zone,omitempty"` // Cloud network zone
	K8sNamespace     string `json:"k8s_namespace,omitempty"`     // Kubernetes namespace of the service
	K8sPodName       string `json:"k8s_pod_name,omitempty"`      // Kubernetes pod name where the service runs
	K8sContainerID   string `json:"k8s_container_id,omitempty"`  // Docker/Kubernetes runtime container ID
	LoadBalancerIP   string `json:"load_balancer_ip,omitempty"`  // IP of the load balancer handling the request

	// Request details
	RequestID          string            `json:"request_id,omitempty"`           // Unique client request identifier
	Operation          string            `json:"operation,omitempty"`            // Operation type: start, hello, inbound, outbound
	RequestMethod      string            `json:"request_method,omitempty"`       // HTTP method used in the request
	RequestURL         string            `json:"request_url,omitempty"`          // Full request URL (path + query parameters)
	RequestSizeBytes   int               `json:"request_size_bytes,omitempty"`   // Incoming message size in bytes
	RequestHeaders     map[string]string `json:"request_headers,omitempty"`      // Captures request headers
	ClientIP           string            `json:"client_ip,omitempty"`            // IP address from request header
	ClientPort         int               `json:"client_port,omitempty"`          // Port number from request header
	UserAgent          string            `json:"user_agent,omitempty"`           // Client user agent string
	Origin             string            `json:"origin,omitempty"`               // Origin of the request
	Referer            string            `json:"referer,omitempty"`              // HTTP referer header
	Protocol           string            `json:"protocol,omitempty"`             // Protocol used in request (e.g., HTTP, HTTPS)
	TLSProtocolVersion string            `json:"tls_protocol_version,omitempty"` // Captures TLS version (1.2, 1.3) if HTTPS
	AuthHeaderStatus   string            `json:"auth_header_status,omitempty"`   // Authorization header presence status
	RequestStage       string            `json:"request_stage,omitempty"`        // Indicates processing stage (received, in-progress, processed)
	ClientID           string            `json:"client_id,omitempty"`            // Unique client identifier
	UserSessionID      string            `json:"user_session_id,omitempty"`      // Helps track session state in distributed requests
	TraceID            string            `json:"trace_id,omitempty"`             // OpenTelemetry Trace ID for distributed tracing
	ParentTraceID      string            `json:"parent_trace_id,omitempty"`      // OTEL Parent Trace ID (if applicable)
	SpanID             string            `json:"span_id,omitempty"`              // OpenTelemetry Span ID for distributed tracing
	CorrelationID      string            `json:"correlation_id,omitempty"`       // Links logs across non-tracing-aware systems
	UserID             string            `json:"user_id,omitempty"`              // Unique user identifier (if available)

	// Response details
	HTTPStatusCode       int               `json:"http_status_code,omitempty"`       // HTTP status code for client response
	ResponseOutcome      string            `json:"response_outcome,omitempty"`       // High-level response success/failure indicator (e.g., success, failure, retry).
	ResponseSizeBytes    int               `json:"response_size_bytes,omitempty"`    // Outgoing response size in bytes
	ResponseTimeMS       float64           `json:"response_time_ms,omitempty"`       // Time taken to generate a response in milliseconds
	ResponseTimeCategory string            `json:"response_time_category,omitempty"` //
	ResponseHeaders      map[string]string `json:"response_headers,omitempty"`       // Captures response headers

	// Debugging & Observability
	CallerFunction       string  `json:"caller_function,omitempty"`        // Fully qualified function name where log was captured (e.g., main.WrapHandlerWithObservability)
	ExecutionUnitID      int     `json:"execution_unit_id,omitempty"`      // (Go → Goroutines, Python → Threads/Async, C# → Tasks)
	CallerLocation       string  `json:"caller_location,omitempty"`        // Source file and line number
	ComputeTimeMS        float64 `json:"compute_time_ms,omitempty"`        // Measures only the time spent processing (excluding queue time).
	QueueTimeMS          float64 `json:"queue_time_ms,omitempty"`          // Captures time spent waiting in a queue (e.g., load balancer, message broker) before being processed.
	ResponseLatencyLevel string  `json:"response_latency_level,omitempty"` // More granular latency category (critical, warning, normal).

	// System performance metrics
	CPUUsage    float64 `json:"cpu_usage,omitempty"`    // CPU usage at the time of request processing
	MemoryUsage float64 `json:"memory_usage,omitempty"` // Memory usage at the time of request processing

	// Error details
	ErrorMessage string `json:"error_message,omitempty"` // Human-readable error message
	StackTrace   string `json:"stack_trace,omitempty"`   // Captured stack trace on failures

	// Log message
	Message string `json:"message,omitempty"` // Log message content
}
