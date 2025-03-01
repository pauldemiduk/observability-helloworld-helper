package handlers

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/graphql-go/graphql"
)

// ✅ Inject Schema at Runtime
var schema graphql.Schema

func SetGraphQLSchema(s graphql.Schema) {
	schema = s
}

// ✅ Handles GraphQL Requests
func GraphQLHandler(w http.ResponseWriter, r *http.Request) {
	log.Println("📡 Received GraphQL request at /graph")

	// Read the request body
	body, err := io.ReadAll(r.Body)
	if err != nil {
		log.Printf("ERROR: Failed to read GraphQL request body: %v", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Log the raw request body
	log.Printf("DEBUG: Raw GraphQL Request Body -> %s", string(body))

	// Decode JSON
	var requestBody struct {
		Query string `json:"query"`
	}
	err = json.Unmarshal(body, &requestBody)
	if err != nil {
		log.Printf("ERROR: Failed to decode GraphQL request: %v", err)
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	// Log parsed query
	log.Printf("DEBUG: Parsed GraphQL Query -> %s", requestBody.Query)

	if requestBody.Query == "" {
		log.Println("ERROR: GraphQL query is empty!")
		http.Error(w, "GraphQL query cannot be empty", http.StatusBadRequest)
		return
	}

	// Ensure schema is set
	if schema.QueryType() == nil {
		log.Println("ERROR: GraphQL schema is NOT initialized!")
		http.Error(w, "Internal Server Error: GraphQL schema not initialized", http.StatusInternalServerError)
		return
	}

	// Execute GraphQL query
	result := graphql.Do(graphql.Params{
		Schema:        schema,
		RequestString: requestBody.Query,
	})

	// Log the GraphQL response
	log.Printf("DEBUG: GraphQL Response -> %+v", result)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(result)
}
