// request.go - Defines request and response metadata models.

package models

// RequestMetadata -
// Captures key metadata from an incoming request, ensuring structured logging and observability.
type RequestMetadata struct {
	// ✅ Core Metadata
	Operation  string `json:"operation,omitempty"`   // The type of operation being performed (e.g., "hello", "inbound", "outbound").
	RequestID  string `json:"request_id,omitempty"`  // Unique identifier for the request (from "X-Request-ID" header).
	RequestURL string `json:"request_url,omitempty"` // Full request URL (path + query parameters)

	// ✅ Request Details
	RequestMethod    string            `json:"request_method,omitempty"`     // HTTP method used in the request (e.g., GET, POST, PUT).
	RequestSizeBytes int               `json:"request_size_bytes,omitempty"` // Size of the request payload in bytes.
	RequestHeaders   map[string]string `json:"request_headers,omitempty"`    // Captures all HTTP request headers for debugging and observability.

	// ✅ Client Information
	ClientIP           string `json:"client_ip,omitempty"`            // IP address from request header.
	ClientPort         int    `json:"client_port,omitempty"`          // Port number from request header.
	UserAgent          string `json:"user_agent,omitempty"`           // User-Agent string from the client.
	Origin             string `json:"origin,omitempty"`               // Origin header value, indicating where the request originated.
	Referer            string `json:"referer,omitempty"`              // Referrer URL from the request headers.
	Protocol           string `json:"protocol,omitempty"`             // Protocol used (e.g., "http", "https").
	TLSProtocolVersion string `json:"tls_protocol_version,omitempty"` // TLS version used in the request (if applicable, e.g., "1.2", "1.3").

	// ✅ Authentication & Identity
	AuthHeaderStatus string `json:"auth_header_status,omitempty"` // Indicates if an Authorization header is present ("present" or "missing").
	UserSessionID    string `json:"user_session_id,omitempty"`    // Unique session identifier (if available, from "X-Session-ID" header).
	ClientID         string `json:"client_id,omitempty"`          // Unique client identifier (from "X-Client-ID" header).
	UserID           string `json:"user_id,omitempty"`            // Authenticated user ID (if available, from "X-User-ID" header).

	// ✅ Distributed Tracing & Correlation
	TraceID       string `json:"trace_id,omitempty"`        // OpenTelemetry Trace ID for distributed tracing.
	ParentTraceID string `json:"parent_trace_id,omitempty"` // Parent Trace ID (if part of a larger distributed trace).
	SpanID        string `json:"span_id,omitempty"`         // OpenTelemetry Span ID for tracking individual operations.
	CorrelationID string `json:"correlation_id,omitempty"`  // Correlation ID used to track related requests across services.
}
