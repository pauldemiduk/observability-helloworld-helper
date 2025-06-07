// *
// ==> graph_query.go - [blueprint:include]
// *
// =>  graph_filter.go explained:
//
//	Description: GraphQL routines to filter log data in the in-memory datastore.
//
// # GraphQL composition
//
//	Schema
//	- Entity
//	-- Interface
//	--- Query 		<= (graph_query.go)
//	---- Filter
package models

import (
	"hello-world-go/internal/metrics"
	"log"

	"github.com/graphql-go/graphql"
)

// createField -
// Helper function to create GraphQL field with logging
func createField(fieldType graphql.Output, resolver func() (interface{}, error), debugMessage string) *graphql.Field {
	return &graphql.Field{
		Type: fieldType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			log.Println(debugMessage)
			return resolver()
		},
	}
}

// isPrometheusEnabled -
// Helper function to toggle Prometheus query execution
func isPrometheusEnabled() bool {
	return true // Set to true once Prometheus is installed
}

//
// ==> Define Query Fields -
//

// serviceConfigQuery -
// Resolve query => not-implemented
var serviceConfigQuery = createField(serviceConfigType, func() (interface{}, error) {
	log.Printf("DEBUG: Resolving ServiceConfig -> %+v", HelloInstance.Config)
	return HelloInstance.Config, nil
}, "DEBUG: Resolving ServiceConfig")

// requestMetadataQuery -
// Resolve query => not-implemented
var requestMetadataQuery = createField(requestMetadataType, func() (interface{}, error) {
	return "RequestMetadata not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving RequestMetadata - Placeholder Implementation")

// responseMetadataQuery -
// Resolve query => not-implemented
var responseMetadataQuery = createField(responseMetadataType, func() (interface{}, error) {
	return "ResponseMetadata not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving ResponseMetadata - Placeholder Implementation")

// logContextQuery -
// Resolve query => not-implemented
var logContextQuery = createField(logContextType, func() (interface{}, error) {
	return "LogContext not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving LogContext - Placeholder Implementation")

// logEntryQuery -
// Resolve query => not-implemented
var logEntryQuery = createField(logEntryType, func() (interface{}, error) {
	return "LogEntry not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving LogEntry - Placeholder Implementation")

// contextQuery -
// Resolve query => contexts := LoadMovingData(snapshotID), contextType
var contextQuery = &graphql.Field{
	Type: graphql.NewList(contextType),
	Args: graphql.FieldConfigArgument{
		"snapshotID": &graphql.ArgumentConfig{Type: graphql.String},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		snapshotID, _ := p.Args["snapshotID"].(string)
		snapshots, err := LoadSnapshotsFromYAML("configs/datastore-snapshots.yaml")
		if err != nil {
			return nil, err
		}

		// Find the snapshot by ID and return its contexts
		for _, snapshot := range snapshots {
			if snapshot.ID == snapshotID {
				return snapshot.Contexts, nil
			}
		}

		return nil, nil
	},
}

// snapshotQuery -
// Resolve query => GetSnapshots(), snapshotType
var snapshotQuery = &graphql.Field{
	Type: graphql.NewList(snapshotType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		return LoadSnapshotsFromYAML("configs/datastore-snapshots.yaml")
	},
}

// slowAPICallQuery -
// Resolve query => GetSlowAPICallSnapshot(), hydrate SlowAPICall <- SlowAPICallType <- ** slowAPICallQuery
var slowAPICallQuery = &graphql.Field{
	Type:        slowAPICallType,
	Description: "Fetch slow API calls with optional duration filter",
	Args: graphql.FieldConfigArgument{
		"duration": &graphql.ArgumentConfig{
			Type:         graphql.Int, // Ensure it's defined as an integer
			DefaultValue: 5,           // Default value if not provided
			Description:  "Time window in minutes for Prometheus queries",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		// Extract the duration from query arguments
		duration, ok := params.Args["duration"].(int)
		if !ok || duration <= 0 {
			duration = 5 // Enforce a default value of 5 minutes
		}

		// Fetch the full snapshot
		return GetSlowAPICallSnapshot(duration)
	},
}

var errorRateQuery = &graphql.Field{
	Type:        errorRateType,
	Description: "Fetch API error rates with optional duration filter",
	Args: graphql.FieldConfigArgument{
		"duration": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 5,
			Description:  "Time window in minutes for Prometheus queries",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		duration, ok := params.Args["duration"].(int)
		if !ok || duration <= 0 {
			duration = 5
		}
		return GetErrorRateSnapshot(duration)
	},
}

// slaViolationQuery - Fetches SLA violation data from snapshot
var slaViolationQuery = &graphql.Field{
	Type:        graphql.NewList(slaViolationType),
	Description: "Fetch structured SLA violations from snapshots",
	Args: graphql.FieldConfigArgument{
		"duration": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 5, // Default to 5 minutes if not provided
			Description:  "Time window in minutes for Prometheus queries",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		duration, ok := params.Args["duration"].(int)
		if !ok || duration <= 0 {
			duration = 5
		}
		log.Println("DEBUG: Resolving SLA Violations Query")
		return GetSLAViolationSnapshot(duration)
	},
}

// sloComplianceQuery - Fetches SLO compliance structured data from snapshot
var sloComplianceQuery = &graphql.Field{
	Type:        graphql.NewList(sloComplianceType),
	Description: "Fetch structured SLO compliance data from snapshots",
	Args: graphql.FieldConfigArgument{
		"duration": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 5,
			Description:  "Time window in minutes for Prometheus queries",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		duration, ok := params.Args["duration"].(int)
		if !ok || duration <= 0 {
			duration = 5
		}
		log.Println("DEBUG: Resolving SLO Compliance Query")
		return GetSLOComplianceSnapshot(duration)
	},
}

// errorBudgetQuery - Fetches structured error budget data from snapshot
var errorBudgetQuery = &graphql.Field{
	Type:        graphql.NewList(errorBudgetType),
	Description: "Fetch structured error budget consumption data",
	Args: graphql.FieldConfigArgument{
		"duration": &graphql.ArgumentConfig{
			Type:         graphql.Int,
			DefaultValue: 5,
			Description:  "Time window in minutes for Prometheus queries",
		},
	},
	Resolve: func(params graphql.ResolveParams) (interface{}, error) {
		duration, ok := params.Args["duration"].(int)
		if !ok || duration <= 0 {
			duration = 5
		}
		log.Println("DEBUG: Resolving Error Budget Query")
		return GetErrorBudgetSnapshot(duration)
	},
}

// metricQueries -
// Resolve query => GetMetricData(), uses metric names and maps them to Prometheus queries
var metricQueries = graphql.Fields{
	"slowRequests": &graphql.Field{
		Type:        graphql.Float,
		Description: "Total slow API requests",
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 5},
		},
		Resolve: metrics.GetMetricData,
	},
	"p95Latency": &graphql.Field{
		Type:        graphql.Float,
		Description: "P95 Latency for API calls",
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 5},
		},
		Resolve: metrics.GetMetricData,
	},
	"errorRate": &graphql.Field{
		Type:        graphql.Float,
		Description: "API error rate over time",
		Args: graphql.FieldConfigArgument{
			"duration": &graphql.ArgumentConfig{Type: graphql.Int, DefaultValue: 5},
		},
		Resolve: metrics.GetMetricData,
	},
}

// recentLogsQuery -
// Resolve query => logs := GetRecentLogs(limit, filter), logEntryType
var recentLogsQuery = &graphql.Field{
	Type: graphql.NewList(logEntryType),
	Args: graphql.FieldConfigArgument{
		"limit": &graphql.ArgumentConfig{
			Type: graphql.Int,
		},
		"filter": &graphql.ArgumentConfig{
			Type: LogFilterInputType,
		},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		limit, _ := p.Args["limit"].(int)

		// Ensure proper filter parsing
		var filter LogFilterInput
		if filterArg, ok := p.Args["filter"].(map[string]interface{}); ok {
			log.Printf("DEBUG: Raw GraphQL Filter Data -> %+v", filterArg) // 🔍 Print raw filter data
			filter = ParseLogFilter(filterArg)
		}

		log.Printf("DEBUG: Parsed LogFilterInput -> %+v", filter) // 🔍 Print parsed filter

		// httpStatusCode == 0 is a valid value for GraphUI filter
		// Only parse `filter` if it exists
		if filterArg, ok := p.Args["filter"]; ok {
			if filterMap, isMap := filterArg.(map[string]interface{}); isMap {
				filter = ParseLogFilter(filterMap)

				// Check if `httpStatusCode` was explicitly provided
				if _, exists := filterMap["httpStatusCode"]; !exists {
					log.Printf("DEBUG: No httpStatusCode filter was provided in UI")
					filter.HTTPStatusCode = -1 // Assign special value to indicate "no filter"
				}
			}
		} else {
			log.Printf("DEBUG: No filter object provided, using default LogFilterInput")
			filter.HTTPStatusCode = -1 // Ensure it's set correctly when no filter is used
		}

		// Resolve query => GetRecentLogs(limit, filter)
		logs := GetRecentLogs(limit, filter)

		// Debugging: Print first log to verify fields
		if len(logs) > 0 {
			log.Printf("DEBUG: First Log Entry -> %+v", logs[0])
		}

		return logs, nil
	},
}

//
// ==> Routines to get LogEntry data
//

// GetRecentLogs -
// Query GraphQL for recent logs using simple in-memory datastore
// Resolve query => var logs []LogEntry = fetchLogsFromStore(limit)
// Apply LogFilterInput arguments => filter.Since, .Level, .HttpStatusCode, .ResponseTimeCategory, .ResponseLatencyLevel, .RequestMethod
func GetRecentLogs(limit int, filter LogFilterInput) []LogEntry {

	// log.Printf("DEBUG: Fetching logs from GraphQLDataStore -> Count: %d", len(GraphQLDataStore))
	// log.Printf("DEBUG: Filter Received -> %+v", filter)

	log.Printf("DEBUG: Received GraphQL filter - Level: %s, HTTPStatusCode: %d, ResponseLatencyLevel: %s",
		filter.Level, filter.HTTPStatusCode, filter.ResponseLatencyLevel)

	var logs []LogEntry = fetchLogsFromStore(limit)

	log.Printf("DEBUG: Logs Before Filtering -> %d", len(logs))

	if filter.Since != "" {
		log.Printf("DEBUG: Applying Time Filter -> %s", filter.Since)
		logs = filterLogsByTime(logs, filter.Since)
	}

	if filter.Level != "" {
		log.Printf("DEBUG: Applying Level Filter -> %s", filter.Level)
		logs = filterLogsByLevel(logs, filter.Level)
	}

	if filter.HTTPStatusCode >= 0 { // Ignore 0 (default) unless explicitly requested
		log.Printf("DEBUG: Applying HTTP Status Code filter -> %d", filter.HTTPStatusCode)
		logs = filterLogsByStatusCode(logs, filter.HTTPStatusCode)
	} else {
		log.Printf("DEBUG: No HTTP Status Code filter applied")
	}

	if filter.ResponseTimeCategory != "" {
		log.Printf("DEBUG: Applying Response Time Filter -> %s", filter.ResponseTimeCategory)
		logs = filterLogsByResponseTimeCategory(logs, filter.ResponseTimeCategory)
	}

	if filter.ResponseLatencyLevel != "" {
		log.Printf("DEBUG: Applying Response Latency Filter -> %s", filter.ResponseLatencyLevel)
		logs = filterLogsByResponseLatencyLevel(logs, filter.ResponseLatencyLevel)
	}

	if filter.RequestMethod != "" {
		log.Printf("DEBUG: Applying Request Method Filter -> %s", filter.RequestMethod)
		logs = filterLogsByRequestMethod(logs, filter.RequestMethod)
	}

	log.Printf("DEBUG: Logs After Filtering -> %d", len(logs))
	return logs
}
