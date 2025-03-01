package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"hello-world-go/internal/models"
	"hello-world-go/internal/queue"
	"hello-world-go/internal/util"
)

// QueuePushHandler - Pushes event into the IngestQueue
func QueuePushHandler(w http.ResponseWriter, r *http.Request) {
	var event models.EventPayload
	if err := json.NewDecoder(r.Body).Decode(&event); err != nil {
		http.Error(w, "Invalid event payload", http.StatusBadRequest)
		return
	}

	// ✅ Ensure IngestQueue has Enqueue()
	models.HelloInstance.IngestQueue.Enqueue(event)
	log.Printf("INFO: Event %s added to IngestQueue", event.EventID)

	util.SetSecurityHeaders(w) // ✅ Secure responses
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusAccepted)
	json.NewEncoder(w).Encode(map[string]string{"message": "Event added to queue!"})

}

// QueuePopHandler - Fetches an event from the WorkloadQueue
func QueuePopHandler(w http.ResponseWriter, r *http.Request) {
	eventPtr := models.HelloInstance.WorkloadQueue.Dequeue()
	if eventPtr == nil {
		http.Error(w, "No events in workload queue", http.StatusNotFound)
		return
	}

	event := *eventPtr // ✅ Dereference

	log.Printf("INFO: Event %s dequeued from WorkloadQueue", event.EventID)
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(event)
}

// QueueProcessHandler - Moves events from IngestQueue to WorkloadQueue
func QueueProcessHandler(w http.ResponseWriter, r *http.Request) {
	//queue.MoveFromIngestToWorkload()
	go queue.MoveFromIngestToWorkload() // ✅ Process in background
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Queue processing triggered!"})
}
