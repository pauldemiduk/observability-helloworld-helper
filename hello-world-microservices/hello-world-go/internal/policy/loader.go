// ActivePolicies - Applies Filtering, Enrichment, and Routing Policies
package policy

import (
	"hello-world-go/internal/models"
	"log"
)

// ✅ Define FilteringPolicy struct
type FilteringPolicy struct {
	Enabled                 bool     // ✅ Allows enabling/disabling filtering
	DropEventsBelowSeverity int      // Minimum severity level to keep an event
	DropCategories          []string // List of categories to drop
	RejectIfMissingFields   []string // List of required fields
}

// ✅ Define RoutingPolicy struct
type RoutingPolicy struct {
	Enabled         bool              // ✅ Allows enabling/disabling filtering
	RouteByCategory map[string]string // ✅ Simple mapping is enough
}

// ✅ Define EnrichmentPolicy struct
type EnrichmentPolicy struct {
	Enabled     bool           // ✅ Allows enabling/disabling filtering
	AddMetadata []MetadataRule // ✅ Defines enrichment rules
}

// ✅ Define MetadataRule Struct (Represents each enrichment rule)
type MetadataRule struct {
	Name  string // Key name to add
	Value string // Value to assign
	Type  string // Type of metadata (e.g., "timestamp", "static")
}

// ✅ Define ActivePolicies as a Struct
type ActivePolicies struct {
	Filtering  FilteringPolicy
	Enrichment EnrichmentPolicy
	Routing    RoutingPolicy
}

// ✅ Global Variable to Hold Policies
var ActivePolicyRules = ActivePolicies{}

// ✅ Apply Active Policies to an Event
func ApplyActivePolicies(event models.EventPayload) models.EventPayload {
	conf := models.HelloInstance.Config

	// ✅ Apply Filtering Policy (If Enabled)
	if conf.PolicyFilteringEnabled {
		if !ApplyFilteringPolicy(event) {
			log.Printf("INFO: Event %s discarded by Filtering Policy.", event.EventID)
			return models.EventPayload{} // Empty struct means rejection
		}
	}

	// ✅ Apply Enrichment Policy (If Enabled)
	if conf.PolicyEnrichmentEnabled {
		event = ApplyEnrichmentPolicy(event)
		log.Printf("INFO: Event %s enriched successfully.", event.EventID)
	}

	// ✅ Apply Routing Policy (If Enabled)
	if conf.PolicyRoutingEnabled {
		action := ApplyRoutingPolicy(event)
		if action == "discard" {
			log.Printf("INFO: Event %s discarded by Routing Policy.", event.EventID)
			return models.EventPayload{} // Empty struct means rejection
		}
		log.Printf("INFO: Event %s routed successfully.", event.EventID)
	}

	// ✅ Return processed event
	return event
}
