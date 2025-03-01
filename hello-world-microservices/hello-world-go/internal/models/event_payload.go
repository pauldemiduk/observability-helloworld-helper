package models

import "time"

// EventPayload - Defines the structure of inbound events
type EventPayload struct {
	// ✅ Core Event Metadata
	EventID   string    `json:"event_id"`  // Unique identifier for the event (UUID)
	Timestamp time.Time `json:"timestamp"` // Event creation timestamp
	Source    string    `json:"source"`    // Service/system generating the event
	Category  string    `json:"category"`  // Event type (alert, incident, log, metric)
	Severity  string    `json:"severity"`  // Severity level (info, warning, critical)
	Message   string    `json:"message"`   // Human-readable event description

	// ✅ Optional Metric Details
	MetricName  string  `json:"metric_name,omitempty"`  // Metric name (if applicable)
	MetricValue float64 `json:"metric_value,omitempty"` // Metric value (if applicable)
	Threshold   float64 `json:"threshold,omitempty"`    // Alert threshold
	Unit        string  `json:"unit,omitempty"`         // Measurement unit

	// ✅ Service Context (Runtime Metadata from the Source Service)
	ServiceVersion   string `json:"service_version"`   // Service version generating the event
	InstanceID       string `json:"instance_id"`       // Unique runtime instance ID
	ServiceIP        string `json:"service_ip"`        // IP address of the service
	CloudProvider    string `json:"cloud_provider"`    // AWS, GCP, Azure, localhost, etc.
	Region           string `json:"region"`            // Cloud region where event was generated
	AvailabilityZone string `json:"availability_zone"` // Cloud availability zone
	NodeIP           string `json:"node_ip"`           // Host machine or node IP
	K8sPodName       string `json:"k8s_pod_name"`      // Kubernetes pod name (if applicable)
	K8sContainerID   string `json:"k8s_container_id"`  // Kubernetes/Docker container ID

	// ✅ Labels & Custom Attributes
	Labels map[string]string `json:"labels,omitempty"` // Key-value metadata for custom event attributes
}
