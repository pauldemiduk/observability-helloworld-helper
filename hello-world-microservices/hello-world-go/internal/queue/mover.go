package queue

import (
	"hello-world-go/internal/models"
	"hello-world-go/internal/policy"
	"log"
)

// MoveFromIngestToWorkload - Transfers validated events from IngestQueue to WorkloadQueue
func MoveFromIngestToWorkload() {

	// ✅ Extract policy settings
	filteringEnabled := policy.ActivePolicyRules.Filtering.Enabled
	enrichmentEnabled := policy.ActivePolicyRules.Enrichment.Enabled
	routingEnabled := policy.ActivePolicyRules.Routing.Enabled

	// ** verify.> policy counts go here.?

	// ✅ Process events from IngestQueue
	for models.HelloInstance.IngestQueue.Length() > 0 {
		eventPtr := models.HelloInstance.IngestQueue.Dequeue() // ✅ Now returns *models.EventPayload
		if eventPtr == nil {
			continue
		}
		event := *eventPtr // ✅ Dereference the pointer

		// ✅ Apply Filtering Policy
		if filteringEnabled && !policy.ApplyFilteringPolicy(event) {
			continue
		}

		// ✅ Apply Enrichment Policy
		if enrichmentEnabled {
			event = policy.ApplyEnrichmentPolicy(event)
		}

		// ✅ Apply Routing Policy
		action := "default"
		if routingEnabled {
			action = policy.ApplyRoutingPolicy(event)
		}

		// ✅ Handle Routing Actions
		if action == "discard" {
			continue
		}

		// ✅ Move to WorkloadQueue
		models.HelloInstance.WorkloadQueue.Enqueue(event)

		log.Printf("INFO: Moving event %s from IngestQueue to WorkloadQueue", event.EventID)
		log.Printf("INFO: Event %s passed filtering, proceeding to enrichment", event.EventID)
		log.Printf("INFO: Event %s enriched, proceeding to routing", event.EventID)
		log.Printf("INFO: Event %s routed successfully. Final action: %s", event.EventID, action)

	}

}

// validateEvent - Stub function for event validation (to be implemented later)
func validateEvent(event models.EventPayload) bool {
	// 🚨 TODO: Implement proper validation logic
	// Example: Check if required fields are non-empty

	return true // Placeholder: Assume valid for now
}
