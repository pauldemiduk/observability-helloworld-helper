// config_helper.go -
package config

import (
	"hello-world-go/internal/models"
	"log"
	"os"
	"path/filepath"
)

// GetLogEndpoint -
// Returns the appropriate log aggregation endpoint based on the log backend.
func GetLogEndpoint(config models.ServiceConfig) string {
	switch config.LogBackend {
	case "loki":
		return config.LokiEndpoint
	case "aws":
		return config.CloudWatchEndpoint
	case "gcp":
		return config.GCPLoggingEndpoint
	case "azure":
		return config.AzureMonitorEndpoint
	default:
		return "unknown"
	}
}

// ResolveDocsPath -
// Dynamically resolves the correct project documentation path.
func ResolveDocsPath() string {
	wd, err := os.Getwd()
	if err != nil {
		log.Fatalf("ERROR: Failed to get working directory: %v", err)
	}

	// Adjust path resolution based on project structure
	docsPath := filepath.Join(wd, "..", "docs")

	// Check if the resolved docs path exists
	if _, err := os.Stat(docsPath); os.IsNotExist(err) {
		log.Printf("WARNING: Documentation directory does not exist: %s", docsPath)
	}

	return docsPath
}
