// *
// ==> observability.go - [blueprint:include]
// *
// =>  observability.go explained:
//
//	Description: Structures and Data routines for observability data flow
package models

// ObservabilityDataFlow - Mapped to Scorecard Fields
type ObservabilityDataFlow struct {
	AlertVolume  int `json:"alert_volume"`
	LogVolume    int `json:"log_volume"`
	EventVolume  int `json:"event_volume"`
	MetricVolume int `json:"metric_volume"`
	TraceVolume  int `json:"trace_volume"` // Placeholder for future trace integration
}
