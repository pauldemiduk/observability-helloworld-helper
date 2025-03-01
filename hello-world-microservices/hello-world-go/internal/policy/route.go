package policy

import "hello-world-go/internal/models"

func ApplyRoutingPolicy(event models.EventPayload) string {
	if dest, exists := ActivePolicyRules.Routing.RouteByCategory[event.Category]; exists {
		return dest // ✅ Send event to mapped destination
	}
	return "route_to_workload_queue" // ✅ Default action
}
