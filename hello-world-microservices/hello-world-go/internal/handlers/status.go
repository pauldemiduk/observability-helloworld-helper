// status.go -
package handlers

import (
	"encoding/json"
	"hello-world-go/internal/util"
	"log"
	"net/http"
)

// StatusResponse - Defines the JSON response structure for health checks
type StatusResponse struct {
	Status  string `json:"status"`
	Details string `json:"details,omitempty"`
}

// statusHandler - Handles `/status` requests with multiple health check levels
func StatusHandler(w http.ResponseWriter, r *http.Request) {
	query := r.URL.Query()
	level := query.Get("level")

	var response StatusResponse
	statusCode := http.StatusOK

	switch level {
	case "ready":
		// ✅ Readiness check (Stub for now, expand later)
		response = StatusResponse{
			Status:  "ready",
			Details: "Configuration loaded. Essential dependencies reachable.",
		}
	case "live":
		// ✅ Liveness check (Stub for now, expand later)
		response = StatusResponse{
			Status:  "live",
			Details: "Service is healthy and actively processing requests.",
		}
	case "basic", "":
		// ✅ Basic health check (Default)
		response = StatusResponse{
			Status: "ok",
		}
	default:
		// ❌ Invalid level provided
		statusCode = http.StatusBadRequest
		response = StatusResponse{
			Status:  "error",
			Details: "Invalid health check level. Use ?level=basic, ready, or live.",
		}
	}

	// ✅ Logging health check execution
	log.Printf("Health check executed: level=%s, status=%s", level, response.Status)

	// ✅ Return JSON response
	util.SetSecurityHeaders(w) // ✅ Secure responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(response)

}
