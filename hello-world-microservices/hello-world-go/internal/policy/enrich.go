package policy

import (
	"hello-world-go/internal/models"
	"time"
)

func ApplyEnrichmentPolicy(event models.EventPayload) models.EventPayload {
	for _, meta := range ActivePolicyRules.Enrichment.AddMetadata { // ✅ Use `ActivePolicyRules`
		switch meta.Type {
		case "timestamp":
			event.Labels[meta.Name] = time.Now().Format(time.RFC3339) // ✅ Add timestamp
		default:
			event.Labels[meta.Name] = meta.Value // ✅ Assign static metadata
		}
	}
	return event
}
