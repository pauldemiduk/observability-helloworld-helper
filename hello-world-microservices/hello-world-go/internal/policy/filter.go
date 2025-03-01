package policy

import (
	"hello-world-go/internal/models"
	"log"
	"strconv"
)

// ApplyFilteringPolicy - Evaluates event against filtering rules
func ApplyFilteringPolicy(event models.EventPayload) bool {
	policy := ActivePolicyRules.Filtering // ✅ Reference FilteringPolicy

	// ✅ Convert event.Severity from string to int
	eventSeverity, err := strconv.Atoi(event.Severity)
	if err != nil {
		log.Printf("WARNING: Event %s has invalid severity '%s', assuming lowest priority (0).", event.EventID, event.Severity)
		eventSeverity = 0 // Default to lowest severity if parsing fails
	}

	// ✅ Drop events below the configured severity threshold
	if eventSeverity < policy.DropEventsBelowSeverity {
		log.Printf("INFO: Dropping event %s due to severity filter. (Severity: %d < %d)", event.EventID, eventSeverity, policy.DropEventsBelowSeverity)
		return false
	}

	// ✅ Drop events if the category matches a filtered category
	for _, category := range policy.DropCategories {
		if event.Category == category {
			log.Printf("INFO: Dropping event %s due to category filter: %s", event.EventID, category)
			return false
		}
	}

	// ✅ Drop events missing required fields
	for _, requiredField := range policy.RejectIfMissingFields {
		if isEmpty(event, requiredField) {
			log.Printf("INFO: Dropping event %s due to missing field: %s", event.EventID, requiredField)
			return false
		}
	}

	return true
}

// isEmpty - Checks if a required field is missing
func isEmpty(event models.EventPayload, field string) bool {
	switch field {
	case "EventID":
		return event.EventID == ""
	case "Timestamp":
		return event.Timestamp.IsZero()
	case "Category":
		return event.Category == ""
	case "Severity":
		return event.Severity == "" // ✅ Since Severity is a string, check for empty string
	case "Message":
		return event.Message == ""
	}
	return false
}
