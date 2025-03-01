// http_helpers.go -
package util

import (
	"crypto/tls"
	"hello-world-go/internal/models"
	"log"

	"net/http"
	"strings"
	"time"
)

// GetTLSVersionString -
// Converts a numeric TLS version to a readable string.
func GetTLSVersionString(tlsProtocolVersion uint16) string {
	switch tlsProtocolVersion {
	case tls.VersionTLS10:
		return "TLS 1.0"
	case tls.VersionTLS11:
		return "TLS 1.1"
	case tls.VersionTLS12:
		return "TLS 1.2"
	case tls.VersionTLS13:
		return "TLS 1.3"
	default:
		return "unknown"
	}
}

// SetSecurityHeaders -
// ** VERIFY.>, is this local.?, should it be placed in middleware.go.?
func SetSecurityHeaders(w http.ResponseWriter) {
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.Header().Set("X-Frame-Options", "DENY")
	w.Header().Set("Content-Security-Policy", "default-src 'self'")
}

// ExtractRequestMetadata -
// Extracts all relevant request attributes and returns them as a structured object.
func ExtractRequestMetadata(r *http.Request) models.RequestMetadata {

	// Extract request size and url
	size := int(r.ContentLength)
	requestUrl := r.URL.String()

	tlsProtocolVersion := "unknown"
	if r.TLS != nil {
		tlsProtocolVersion = GetTLSVersionString(r.TLS.Version)
	}

	traceID, spanID := ExtractTraceContext(r)
	// ** uplift.? - If ParentTraceID extraction is redundant, just use TraceID
	parentTraceID := traceID
	clientIP, clientPort := ExtractClientIP(r)

	// Extract request headers into a normalized map
	// Be careful when testing with curl, you have to include these to see them in the logs
	requestHeaders := make(map[string]string)
	for key, values := range r.Header {
		normalizedKey := strings.ToLower(key)
		requestHeaders[normalizedKey] = strings.Join(values, ", ")
		// log.Printf("DEBUG: key: %s", normalizedKey)
	}

	// Ensure all required headers have valid defaults if missing
	requestID := requestHeaders["x-request-id"]
	if requestID == "" {
		log.Println("WARNING: X-Request-ID header is missing, request may not be traceable.")
		requestID = "unknown"
		// ** uplift
		// Some systems force and pass x-request-id in the payload
		// Most cloud native systems will generate and pass one, we can do the same if needed
	}

	userSessionID := requestHeaders["x-session-id"]
	if userSessionID == "" {
		userSessionID = "unknown"
	}

	authHeaderStatus := "missing"
	if _, exists := requestHeaders["authorization"]; exists {
		authHeaderStatus = "present"
	}

	userAgent := requestHeaders["user-agent"]
	if userAgent == "" {
		userAgent = "unknown"
	}

	origin := requestHeaders["origin"]
	if origin == "" {
		origin = "unknown"
	}

	referer := requestHeaders["referer"]
	if referer == "" {
		referer = "unknown"
	}

	protocol := requestHeaders["x-forwarded-proto"]
	if protocol == "" {
		protocol = "unknown"
	}

	clientID := requestHeaders["x-client-id"]
	if clientID == "" {
		clientID = "unknown"
	}

	correlationID := requestHeaders["x-correlation-id"]
	if correlationID == "" {
		correlationID = "unknown"
	}

	userID := requestHeaders["x-user-id"]
	if userID == "" {
		userID = "unknown"
	}

	return models.RequestMetadata{
		RequestID:          requestID,
		RequestURL:         requestUrl,
		RequestSizeBytes:   size,
		RequestHeaders:     requestHeaders,
		TraceID:            traceID,
		SpanID:             spanID,
		ParentTraceID:      parentTraceID,
		TLSProtocolVersion: tlsProtocolVersion,
		UserSessionID:      userSessionID,
		AuthHeaderStatus:   authHeaderStatus,
		ClientIP:           clientIP,
		ClientPort:         clientPort,
		UserAgent:          userAgent,
		Origin:             origin,
		Referer:            referer,
		Protocol:           protocol,
		ClientID:           clientID,
		CorrelationID:      correlationID,
		UserID:             userID,
	}
}

// ExtractResponseMetadata -
// Extracts all relevant response attributes and returns them as a structured object.
func ExtractResponseMetadata(lrw models.ResponseWriterInterface, duration time.Duration) models.ResponseMetadata {
	responseTimeMS := duration.Seconds() * 1000

	// ✅ Extract Response Headers
	responseHeaders := make(map[string]string)
	for key, values := range lrw.Header() {
		responseHeaders[key] = strings.Join(values, ", ")
	}

	return models.ResponseMetadata{
		HTTPStatusCode:       lrw.GetStatusCode(),
		ResponseOutcome:      CategorizeHTTPStatusCode(lrw.GetStatusCode()),
		ResponseSizeBytes:    lrw.GetSize(),
		ResponseHeaders:      responseHeaders,
		ResponseTimeMS:       responseTimeMS,
		ResponseTimeCategory: CategorizeResponseTime(responseTimeMS, lrw.GetStatusCode()),
	}
}
