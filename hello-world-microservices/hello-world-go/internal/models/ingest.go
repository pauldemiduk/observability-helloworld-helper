package models

import (
	"log"
	"sync"
)

// IngestQueue - Holds raw inbound events before processing
type IngestQueue struct {
	mu         sync.Mutex
	queue      []EventPayload
	maxSize    int
	retryMap   map[string]int
	retryLimit int
}

// NewIngestQueue - Initializes a new ingest queue with limits
func NewIngestQueue(maxSize, retryLimit int) *IngestQueue {

	log.Printf("INFO: IngestQueue initialized with max size: %d, retry limit: %d", maxSize, retryLimit)

	return &IngestQueue{
		queue:      make([]EventPayload, 0),
		maxSize:    maxSize,
		retryMap:   make(map[string]int),
		retryLimit: retryLimit,
	}
}

// Enqueue - Adds an event to the IngestQueue (drops if full)
func (q *IngestQueue) Enqueue(event EventPayload) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// ✅ Enforce max queue size
	if q.maxSize > 0 && len(q.queue) >= q.maxSize {
		log.Printf("WARNING: IngestQueue full, dropping event: %s", event.EventID)
		return
	}

	q.queue = append(q.queue, event)
	log.Printf("INFO: Event %s enqueued to IngestQueue. Current size: %d", event.EventID, len(q.queue))

}

// Dequeue - Removes and returns an event from the IngestQueue
func (q *IngestQueue) Dequeue() *EventPayload {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	// ✅ Retrieve and remove the first event
	event := q.queue[0]
	q.queue = q.queue[1:]

	// ✅ Check retry count
	if q.retryMap[event.EventID] >= q.retryLimit {
		log.Printf("ERROR: Event %s exceeded retry limit (%d), dropping.", event.EventID, q.retryLimit)
		delete(q.retryMap, event.EventID) // Remove from retry tracking
		return nil
	}

	// ✅ Track retries
	q.retryMap[event.EventID]++
	return &event
}

// Length - Returns the current size of the queue
func (q *IngestQueue) Length() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}

// ✅ GetMaxSize - Returns the max queue size (needed for logging)
func (q *IngestQueue) GetMaxSize() int {
	return q.maxSize
}
