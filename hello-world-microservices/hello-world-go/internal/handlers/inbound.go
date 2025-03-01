// inbound.go -
package handlers

import (
	"encoding/json"
	"fmt"
	"hello-world-go/internal/models"
	"hello-world-go/internal/routing"
	"hello-world-go/internal/util"
	"log"
	"net/http"
)

// inbound route
//
// router.HandleFunc("/inbound", handlers.InboundHandler).Methods("POST")
//
// func inboundHandler(w http.ResponseWriter, r *http.Request) {
// 	util.SetSecurityHeaders(w) // ✅ Secure responses
// 	w.Header().Set("Content-Type", "application/json")
// 	w.WriteHeader(http.StatusOK)
// 	json.NewEncoder(w).Encode(map[string]string{"message": "Inbound request processed successfully."})
// }

// InboundHandler - Handles incoming event payloads at `/inbound`
func InboundHandler(w http.ResponseWriter, r *http.Request) {
	var event models.EventPayload

	// ✅ Decode JSON payload
	err := json.NewDecoder(r.Body).Decode(&event)
	if err != nil {
		log.Printf("ERROR: Failed to parse event payload: %v", err)
		http.Error(w, "Invalid event payload", http.StatusBadRequest)
		return
	}

	// ✅ Validate required fields
	if event.EventID == "" || event.Timestamp.IsZero() || event.Source == "" || event.Category == "" || event.Severity == "" || event.Message == "" {
		log.Println("ERROR: Missing required event fields")
		http.Error(w, "Missing required event fields", http.StatusBadRequest)
		return
	}

	// ✅ Assign Default Values if Missing
	if event.MetricName == "" {
		event.MetricName = "unknown"
	}

	// ✅ Log event receipt
	log.Printf("INFO: Received event: %+v", event)

	// ✅ Route the event
	routing.RouteEvent(event) // ✅ New routing logic

	// ✅ Send JSON response
	util.SetSecurityHeaders(w) // ✅ Secure responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)

	response := map[string]string{
		"message": fmt.Sprintf("Event processed successfully, event_id: %s", event.EventID),
	}
	json.NewEncoder(w).Encode(response)
}
