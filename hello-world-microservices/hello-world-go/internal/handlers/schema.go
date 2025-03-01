package handlers

import (
	"fmt"
	"hello-world-go/internal/models"
	"io/ioutil"
	"log"
	"net/http"
)

// SchemaHandler - Serves the YAML schema file
func SchemaHandler(w http.ResponseWriter, r *http.Request) {
	// Get schema file path from config
	schemaFile := models.HelloInstance.Config.SchemaFilePath

	if schemaFile == "" {
		http.Error(w, "Schema file path not configured", http.StatusInternalServerError)
		log.Println("ERROR: Schema file path is not set in configuration.")
		return
	}

	// Read YAML file
	data, err := ioutil.ReadFile(schemaFile)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to read schema file: %v", err), http.StatusInternalServerError)
		log.Printf("ERROR: Unable to read schema file: %v", err)
		return
	}

	// Set response headers and serve file content
	w.Header().Set("Content-Type", "text/yaml")
	w.WriteHeader(http.StatusOK)
	_, err = w.Write(data)
	if err != nil {
		log.Printf("ERROR: Failed to write schema response: %v", err)
	}
}
