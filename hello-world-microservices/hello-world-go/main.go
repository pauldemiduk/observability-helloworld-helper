package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	httpSwagger "github.com/swaggo/http-swagger"
)

// Function to retrieve cloud metadata from environment variables or defaults
func getCloudMetadata() (string, string, string) {
	provider := os.Getenv("CLOUD_PROVIDER")
	if provider == "" {
		provider = "docker"
	}
	region := os.Getenv("CLOUD_REGION")
	if region == "" {
		region = "local"
	}
	zone := os.Getenv("CLOUD_ZONE")
	if zone == "" {
		zone = "local"
	}
	return provider, region, zone
}

// HTTP handler for "/hello" endpoint
func helloHandler(w http.ResponseWriter, r *http.Request) {
	provider, region, zone := getCloudMetadata()
	runtimeLanguage := os.Getenv("RUNTIME_LANGUAGE")
	if runtimeLanguage == "" {
		runtimeLanguage = "Golang"
	}

	response := fmt.Sprintf("Hello, World! (Language: %s, Provider: %s, Region: %s, Zone: %s)", runtimeLanguage, provider, region, zone)

	// Log request with standardized format
	log.Printf("👋 Received request at /hello | Language: %s, Provider: %s, Region: %s, Zone: %s", runtimeLanguage, provider, region, zone)
	fmt.Fprintln(w, response)
}

func main() {
	// Set runtime language inside main
	runtimeLanguage := os.Getenv("RUNTIME_LANGUAGE")
	if runtimeLanguage == "" {
		runtimeLanguage = "Golang"
	}

	// Set default host and port from environment variables
	port := os.Getenv("SERVICE_PORT")
	if port == "" {
		port = "8080"
	}
	host := os.Getenv("SERVICE_HOST")
	if host == "" {
		host = "0.0.0.0"
	}

	// Retrieve cloud metadata for logging
	provider, region, zone := getCloudMetadata()
	log.Printf("🚀 Starting Hello World %s service on %s:%s | Provider: %s, Region: %s, Zone: %s", runtimeLanguage, host, port, provider, region, zone)

	http.Handle("/swagger/", httpSwagger.WrapHandler)
	// Define HTTP route
	http.HandleFunc("/hello", helloHandler)

	// Start HTTP server
	err := http.ListenAndServe(host+":"+port, nil)
	if err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
