// *
// ==> graph.go - [blueprint:include]
// *
// =>  graph.go explained:
//
//	Description: Structures and Data routines for GraphQL foundation.
//
// # GraphQL composition
//
//	Schema			<= (graph.go)
//	- Entity
//	-- Interface
//	--- Query
//	---- Filter
package models

import (
	"log"

	"github.com/graphql-go/graphql"
)

// Declare Schema Variable
var schema graphql.Schema

// NewGraphSchema -
// Initialize GraphQL Schema - Invoked by => NewServiceInstance, LoadConfig
func NewGraphSchema() graphql.Schema {

	// Declare queryType as a nil pointer (avoids cycle)
	var queryType *graphql.Object

	log.Printf("INFO: GraphSchema initialization -> started..")

	// Ensure queryType is initialized before schema creation
	// ** ROADMAP.>, (entity, interface, queryType..) -- new query types plug in here..
	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"serviceConfig":    serviceConfigQuery,
			"requestMetadata":  requestMetadataQuery,
			"responseMetadata": responseMetadataQuery,
			//"logContext":       logContextQuery,
			"logEntry":   logEntryQuery,
			"recentLogs": recentLogsQuery,
			//"snapshots":        snapshotQuery,
			//"contexts":         contextQuery,
			//"slowAPICalls":     slowAPICallQuery,
			//"errorRates":       errorRateQuery,
			//"slaViolations":    slaViolationQuery,
			//"sloCompliance":    sloComplianceQuery,
			//"errorBudget":      errorBudgetQuery,
		},
	})

	// Ensure queryType is NOT nil before using it
	if queryType == nil {
		log.Fatalf("CRITICAL ERROR: queryType is still nil before initializing GraphQL schema")
	}

	// queryType okay, now create the GraphQL schema
	var err error
	schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType, // ==> queryType is guaranteed to be initialized now
	})
	if err != nil {
		log.Fatalf("ERROR: Failed to create GraphQL schema: %v", err)
	}

	// Add missing field modifications after queryType is created
	contextType.AddFieldConfig("queries", &graphql.Field{
		Type: graphql.NewList(queryExecutionType),
	})

	log.Printf("INFO: GraphSchema initialization -> complete.")

	return schema
}

// GetGraphQLSchema -
// Helper access to Schema without direct file import.
func GetGraphQLSchema() graphql.Schema {
	return schema
}

// SetGraphQLSchema -
// Helper access to Schema without direct file import.
func SetGraphQLSchema(s graphql.Schema) {
	schema = s
}

// QueryExecution -
// Define Query Execution Entity
// (entity, interface, query) => (*QueryExecution, queryExecutionType, query)
type QueryExecution struct {
	ID   string `yaml:"id" json:"id"`
	Name string `yaml:"name" json:"name"`
	Data string `yaml:"data" json:"data"`
}

//
// ==> Define Helper Types
//

// queryExecutionType -
// Define Query Execution Interface
// (entity, interface, query) => (QueryExecution, *queryExecutionType, query)
var queryExecutionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "QueryExecution",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
		"data": &graphql.Field{Type: graphql.String},
	},
})
