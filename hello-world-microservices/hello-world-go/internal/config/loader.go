// config.go - Enhanced centralized configuration management with log verbosity
package config

import (
	"hello-world-go/internal/models"
	"hello-world-go/internal/util"

	"log"
	"sync"
)

func init() {
	//log.Printf("DEBUG: ServiceVersion at startup: %s", ServiceVersion) // ✅ Confirm value before `LoadLogContext`
}

var configOnce sync.Once

// LoadConfig -
// Reads environment variables and runtime config settings.
// It ensures values have proper defaults and validates critical parameters.
func LoadConfig() models.ServiceConfig {
	configOnce.Do(func() {
		// ✅ Initialize `ServiceConfig` first
		serviceConfig := models.ServiceConfig{
			// ✅ Service Metadata
			ServiceName:     util.GetEnvOrDefault("SERVICE_NAME", "hello-world-go"),
			ServiceVersion:  "unknown",
			ServiceHostname: util.GetEnvOrDefault("SERVICE_HOST", "0.0.0.0"),
			ServiceIP:       util.GetEnvOrDefault("SERVICE_IP", "127.0.0.1"),
			ServicePort:     util.GetEnvIntValidated("SERVICE_PORT", 8080, 1024, 49151),

			// ✅ Application Environment
			Environment: util.GetEnvOrDefault("ENV", "dev"),

			// ✅ Logging & Observability Config
			LogLevel:      util.GetEnvOrDefault("LOG_LEVEL", "debug"),
			LogBackend:    util.GetEnvOrDefault("LOG_BACKEND", "stdout"),
			LogFile:       util.GetEnvOrDefault("LOG_FILE", "service.log"),
			LogSampleRate: util.GetEnvFloatValidated("LOG_SAMPLE_RATE", 0.3, 0.0, 1.0),

			// ✅ Log Rotation Settings with Validated Limits
			LogRotationEnabled: util.GetEnvOrDefaultBool("LOG_ROTATION_ENABLED", true),
			LogRotationSizeMB:  util.GetEnvIntValidated("LOG_ROTATION_SIZE_MB", 10, 1, 100),
			LogRetentionCount:  util.GetEnvIntValidated("LOG_RETENTION_COUNT", 5, 1, 50),
			LogMaxSizeMB:       util.GetEnvIntValidated("LOG_MAX_SIZE_MB", 100, 10, 1000),

			// ✅ Schema File Path
			SchemaFilePath: util.GetEnvOrDefault("SCHEMA_FILE_PATH", "configs/hello.yaml"),

			// ✅ Observability Endpoints
			LokiEndpoint:         util.GetEnvOrDefault("LOKI_ENDPOINT", ""),
			CloudWatchEndpoint:   util.GetEnvOrDefault("CLOUDWATCH_ENDPOINT", ""),
			GCPLoggingEndpoint:   util.GetEnvOrDefault("GCP_LOGGING_ENDPOINT", ""),
			AzureMonitorEndpoint: util.GetEnvOrDefault("AZURE_MONITOR_ENDPOINT", ""),

			// ✅ Response Time Buckets
			ResponseTimeFast:     util.GetEnvFloatValidated("RESPONSE_TIME_FAST", 200.0, 1, 5000),
			ResponseTimeModerate: util.GetEnvFloatValidated("RESPONSE_TIME_MODERATE", 500.0, 1, 5000),
			ResponseTimeSlow:     util.GetEnvFloatValidated("RESPONSE_TIME_SLOW", 1000.0, 1, 10000),

			ResponseLatencyNormal:   util.GetEnvFloatValidated("RESPONSE_LATENCY_NORMAL", 500.0, 1, 5000),
			ResponseLatencyWarning:  util.GetEnvFloatValidated("RESPONSE_LATENCY_WARNING", 1000.0, 1, 10000),
			ResponseLatencyCritical: util.GetEnvFloatValidated("RESPONSE_LATENCY_CRITICAL", 2000.0, 1, 20000),
		}

		// ✅ Assign `ServiceConfig` to `HelloInstance` (but NOT `ServiceContext`)
		models.HelloInstance = &models.ServiceInstance{
			Config:        serviceConfig,
			IngestQueue:   models.NewIngestQueue(100, 3),
			WorkloadQueue: models.NewWorkloadQueue(50),
		}
	})

	// ✅ Return `ServiceConfig`
	return models.HelloInstance.Config
}

// LoadServiceContext - Initializes ServiceContext from existing ServiceConfig
func LoadServiceContext(config models.ServiceConfig) models.ServiceContext {
	return models.ServiceContext{
		ServiceName:      config.ServiceName,
		ServiceVersion:   config.ServiceVersion,
		ServiceHostname:  config.ServiceHostname,
		ServiceIP:        config.ServiceIP,
		ServicePort:      config.ServicePort,
		ProcessID:        util.GetProcessID(),
		ThreadID:         util.GetThreadID(),
		NodeIP:           util.ResolveHostName(),
		NodePort:         config.ServicePort,
		CloudProvider:    util.GetEnvOrDefault("CLOUD_PROVIDER", "aws"),
		CloudRegion:      util.GetEnvOrDefault("CLOUD_REGION", "us-west-2"),
		AvailabilityZone: util.GetEnvOrDefault("CLOUD_ZONE", "us-west-2a"),
		LoadBalancerIP:   util.GetEnvOrDefault("LOAD_BALANCER_IP", "unknown"),
		InstanceID:       util.GetInstanceID(),
		K8sNamespace:     util.GetEnvOrDefault("K8S_NAMESPACE", "unknown"),
		K8sPodName:       util.GetEnvOrDefault("K8S_POD_NAME", "unknown"),
		K8sContainerID:   util.GetContainerID(),
	}
}

// ConfigManager - Manages dynamic config reloads
type ConfigManager struct {
	sync.RWMutex
	Thresholds map[string]float64 // Example: Adjustable thresholds for response times
}

// Global Config Instance
var RuntimeConfig = &ConfigManager{
	Thresholds: map[string]float64{
		"fast":   100,
		"normal": 500,
		"slow":   1000,
	},
}

// ReloadConfig - Dynamically reloads config without restarting service
func ReloadConfig() {
	RuntimeConfig.Lock() // Ensure thread safety
	defer RuntimeConfig.Unlock()

	// Simulate reloading config (future: load from file, API, DB, etc.)
	RuntimeConfig.Thresholds["fast"] = 150 // Example: Adjust threshold dynamically
	RuntimeConfig.Thresholds["normal"] = 550
	RuntimeConfig.Thresholds["slow"] = 1100

	log.Println("INFO: Config reloaded dynamically.")
}
