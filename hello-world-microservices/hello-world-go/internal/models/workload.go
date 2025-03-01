package models

import (
	"log"
	"sync"
)

// WorkloadQueue - Holds validated & processed events that require action
type WorkloadQueue struct {
	mu      sync.Mutex
	queue   []EventPayload
	maxSize int
}

// NewWorkloadQueue - Initializes a new workload queue with a max size
func NewWorkloadQueue(maxSize int) *WorkloadQueue {

	log.Printf("INFO: WorkloadQueue initialized with max size: %d", maxSize)

	return &WorkloadQueue{
		queue:   make([]EventPayload, 0),
		maxSize: maxSize,
	}
}

// Enqueue - Adds an event to the WorkloadQueue (drops if full)
func (q *WorkloadQueue) Enqueue(event EventPayload) {
	q.mu.Lock()
	defer q.mu.Unlock()

	// ✅ Enforce max queue size
	if q.maxSize > 0 && len(q.queue) >= q.maxSize {
		log.Printf("WARNING: WorkloadQueue full, dropping event: %s", event.EventID)
		return
	}

	q.queue = append(q.queue, event)
	log.Printf("INFO: Event %s enqueued from WorkloadQueue. Remaining size: %d", event.EventID, len(q.queue))
}

// Dequeue - Removes and returns an event from the WorkloadQueue
func (q *WorkloadQueue) Dequeue() *EventPayload {
	q.mu.Lock()
	defer q.mu.Unlock()

	if len(q.queue) == 0 {
		return nil
	}

	// ✅ Retrieve and remove the first event
	event := q.queue[0]
	q.queue = q.queue[1:]

	log.Printf("INFO: Event %s dequeued from WorkloadQueue. Remaining size: %d", event.EventID, len(q.queue))

	return &event
}

// ✅ GetMaxSize - Returns the max queue size
func (q *WorkloadQueue) GetMaxSize() int {
	return q.maxSize
}

// ✅ Length - Returns the current size of the queue
func (q *WorkloadQueue) Length() int {
	q.mu.Lock()
	defer q.mu.Unlock()
	return len(q.queue)
}
