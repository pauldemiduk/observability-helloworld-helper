// service.go - Centralized models for Service -> Config, Instance, Context
package models

import (
	"log"
)

// var ServiceVersion = "unknown" // ✅ Default value, overridden at build time
var HelloInstance *ServiceInstance // ✅ Now it's a pointer

// ServiceInstance -
// Encapsulates both application configuration (`ServiceConfig`) and runtime service metadata (`ServiceContext`).
type ServiceInstance struct {
	Config        ServiceConfig
	Context       ServiceContext
	IngestQueue   *IngestQueue
	WorkloadQueue *WorkloadQueue
}

// NewServiceInstance - Initializes ServiceInstance with queues
func NewServiceInstance() *ServiceInstance {
	instance := &ServiceInstance{
		Config:        ServiceConfig{},        // Config will be loaded later in LoadConfig()
		Context:       ServiceContext{},       // Context will be populated dynamically
		IngestQueue:   NewIngestQueue(100, 3), // Default: Max 100 events, Retry 3 times
		WorkloadQueue: NewWorkloadQueue(50),   // Default: Max 50 workload items
	}

	log.Println("INFO: ServiceInstance initialized with IngestQueue & WorkloadQueue.")
	return instance
}

// ServiceConfig -
// Stores static configuration values loaded from environment variables.
// Defines structured settings for logging, observability, networking, and cloud deployments.
type ServiceConfig struct {
	// ✅ Service Metadata
	ServiceName     string // Logical name of the service (e.g., "hello-world-service")
	ServiceVersion  string // Version of the deployed service (e.g., "1.0.3")
	ServiceHostname string // Hostname or IP address where the service is bound
	ServiceIP       string // IP address used by the service for internal communication
	ServicePort     int    // Port number the service listens on

	// ✅ Application Environment
	Environment string // Deployment environment (e.g., "dev", "staging", "prod")

	// ✅ Logging & Observability Config
	LogLevel      string  // Log level (e.g., "debug", "info", "warn", "error")
	LogBackend    string  // Log destination (e.g., "stdout", "file", "loki", "cloudwatch")
	LogFile       string  // File path for local log storage (if applicable)
	LogSampleRate float64 // Percentage of logs sampled for analysis (0.0 - 1.0)

	// ✅ Log Rotation Settings
	LogRotationEnabled bool // Enable/disable log rotation
	LogRotationSizeMB  int  // ✅ Rotate logs when exceeding this size (MB)
	LogRetentionCount  int  // ✅ Maximum number of log backups to keep
	LogMaxSizeMB       int  // ✅ When rotation is disabled, clear logs at this size

	// ✅ Event payload policy for for ingets to workload flow
	PolicyFilteringEnabled  bool // Enables Filtering Policy
	PolicyEnrichmentEnabled bool // Enables Enrichment Policy
	PolicyRoutingEnabled    bool // Enables Routing Policy

	// ✅ Hello, world schema
	SchemaFilePath string // File path for HelloInstance schema

	// ✅ Observability: Log Aggregation Endpoints
	LokiEndpoint         string // URL for Loki log aggregation service
	CloudWatchEndpoint   string // AWS CloudWatch logging endpoint
	GCPLoggingEndpoint   string // Google Cloud Logging endpoint
	AzureMonitorEndpoint string // Azure Monitor Logs endpoint

	// ✅ Response Time Buckets (Used for Categorization of API Latency)
	ResponseTimeFast     float64 // Threshold (in ms) for a fast response (e.g., <200ms)
	ResponseTimeModerate float64 // Threshold (in ms) for a moderate response (e.g., <500ms)
	ResponseTimeSlow     float64 // Threshold (in ms) for a slow response (e.g., >=1000ms)

	// ✅ Response Latency Levels (Used for Alerting on Critical Latencies)
	ResponseLatencyNormal   float64 // Upper bound for "normal" response latency
	ResponseLatencyWarning  float64 // Upper bound for "warning" level response latency
	ResponseLatencyCritical float64 // Upper bound for "critical" latency requiring attention

}

// ServiceContext -
// Stores runtime metadata about the running service instance.
// Derived from `ServiceConfig`, it provides contextual details for logs, monitoring, and tracing.
type ServiceContext struct {
	// ✅ Service Context
	ServiceName     string // Name of the running service (e.g., "hello-world-service")
	ServiceVersion  string // Version of the deployed service instance
	ServiceHostname string // Hostname where the service is running
	ServiceIP       string // IP address of the service (internal or external)
	ServicePort     int    // Port on which the service is exposed
	ProcessID       int    // Process ID of the running service instance
	ThreadID        int    // Thread or Goroutine ID executing the request

	// ✅ Cloud & Deployment Context
	NodeIP           string // IP address of the host (VM, EC2, Kubernetes Node) where the service runs
	NodePort         int    // Port on the host machine mapped to the service
	CloudProvider    string // Cloud provider name (e.g., "aws", "gcp", "azure", "on-prem")
	CloudRegion      string // Cloud region where the service is deployed (e.g., "us-west-2")
	AvailabilityZone string // Availability Zone within the region (e.g., "us-west-2a")
	LoadBalancerIP   string // IP address of the load balancer handling external traffic
	InstanceID       string // Unique identifier for the cloud VM/container instance

	// ✅ Kubernetes Context
	K8sNamespace   string // Kubernetes namespace where the service is deployed
	K8sPodName     string // Name of the Kubernetes pod running the service
	K8sContainerID string // Unique container ID assigned by the runtime (Docker, containerd)
}
