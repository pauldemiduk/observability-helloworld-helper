// notfound.go -
package handlers

import (
	"fmt"
	"hello-world-go/internal/logging"
	"hello-world-go/internal/middleware"
	"hello-world-go/internal/models"
	"hello-world-go/internal/util"
	"net/http"
	"time"
)

// 404-Not Found
func NotFoundHandler(w http.ResponseWriter, r *http.Request) {
	util.SetSecurityHeaders(w) // ✅ Secure responses
	w.Header().Set("Content-Type", "application/json")

	lrw := &middleware.LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusNotFound}

	// ✅ Extract structured request metadata
	requestMetadata := util.ExtractRequestMetadata(r)

	// ✅ Create response metadata for a 404 error
	responseMetadata := models.ResponseMetadata{
		HTTPStatusCode:       http.StatusNotFound,
		ResponseOutcome:      "failure",
		ResponseSizeBytes:    len("404 Not Found"),
		ResponseHeaders:      map[string]string{"Content-Type": "application/json"},
		ResponseTimeMS:       0,
		ResponseTimeCategory: "error",
	}

	// ✅ Log "404 Not Found"
	logOpts := models.LogEntryOptions{
		LRW:              lrw,
		Request:          r,
		StartTime:        time.Now(),
		Duration:         0,
		LogCtx:           models.LoadLogContext(),
		LogLevel:         "error",
		LogType:          "response",
		RequestStage:     "processed",
		Message:          fmt.Sprintf("Unknown endpoint: %s", r.URL.Path),
		RequestMetadata:  requestMetadata,  // ✅ Pass structured request metadata
		ResponseMetadata: responseMetadata, // ✅ Pass structured response metadata
	}

	logging.SendLog(middleware.BuildLogEntry(logOpts)) // ✅ Send structured log

	// ✅ Send response
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("404 Not Found"))
}
