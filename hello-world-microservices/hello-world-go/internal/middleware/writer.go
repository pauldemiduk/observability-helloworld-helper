// writer.go -
package middleware

import "net/http"

// LoggingResponseWriter - Captures response status codes and response size.
type LoggingResponseWriter struct {
	http.ResponseWriter
	StatusCode int
	Size       int
}

// GetStatusCode returns the captured HTTP status code.
func (lrw *LoggingResponseWriter) GetStatusCode() int {
	return lrw.StatusCode
}

// GetSize returns the captured response size.
func (lrw *LoggingResponseWriter) GetSize() int {
	return lrw.Size
}

// WriteHeader - Captures response status codes.
func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	if lrw.StatusCode == 0 { // Prevent duplicate WriteHeader calls
		lrw.StatusCode = code
		lrw.ResponseWriter.WriteHeader(code)
	}
}

// Write - Captures response size.
func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.Size += size
	return size, err
}
