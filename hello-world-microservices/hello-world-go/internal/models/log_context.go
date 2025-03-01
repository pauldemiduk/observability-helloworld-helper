// log_context.go - Defines structured log models.

package models

// LogContext -
// Stores reusable metadata for logging, ensuring efficiency by avoiding repeated function calls.
type LogContext struct {
	// ✅ Service Metadata (Core Service Information)
	ServiceName     string // Name of the microservice (e.g., "hello-world-go")
	ServiceVersion  string // Injected build-time version
	ServiceHostname string // Hostname where the service is running
	ServiceIP       string // IP address of the service (internal or external)
	ServicePort     int    // Port on which the service is exposed
	InstanceID      string // Unique identifier for the running instance of the service

	// ✅ Host & Process Details (Where the service is running)
	ProcessID      int    // OS Process ID (PID) of the running service
	ThreadID       int    // Thread ID, usually the same as ProcessID for single-threaded processes
	NodeIP         string // IP address or hostname of the node the service is running on
	NodePort       int    // Port number the service is listening on
	LoadBalancerIP string // IP address of the load balancer (if applicable)

	// ✅ Cloud & Infrastructure Metadata
	CloudProvider    string // Cloud provider identifier (e.g., "AWS", "GCP", "Azure", "localhost")
	CloudRegion      string // Cloud region where the service is deployed
	AvailabilityZone string // Cloud zone (availability zone) within the region

	// ✅ Kubernetes Metadata (Containerized Deployment)
	K8sNamespace   string // Kubernetes namespace of the service
	K8sPodName     string // Name of the Kubernetes pod the service is running in
	K8sContainerID string // Docker or Kubernetes container ID, if applicable

	// ✅ User Session Metadata
	UserSessionID string // Unique session identifier (if applicable)
}
