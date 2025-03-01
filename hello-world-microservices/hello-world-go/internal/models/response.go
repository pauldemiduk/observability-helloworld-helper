// response.go - Defines request and response metadata models.

package models

// ResponseMetadata -
// Captures key details about an HTTP response, ensuring structured logging and observability.
type ResponseMetadata struct {
	// ✅ Response Status
	HTTPStatusCode  int    `json:"http_status_code,omitempty"` // HTTP status code for the response (e.g., 200, 404, 500).
	ResponseOutcome string `json:"response_outcome,omitempty"` // High-level response outcome (e.g., success, failure, retry).

	// ✅ Response Timing
	ResponseTimeMS       float64 `json:"response_time_ms,omitempty"`       // Time taken to generate a response in milliseconds.
	ResponseTimeCategory string  `json:"response_time_category,omitempty"` // Categorization of response time (e.g., fast, normal, slow).

	// ✅ Response Content
	ResponseSizeBytes int               `json:"response_size_bytes,omitempty"` // Outgoing response size in bytes.
	ResponseHeaders   map[string]string `json:"response_headers,omitempty"`    // Captures response headers for debugging and observability.
}
