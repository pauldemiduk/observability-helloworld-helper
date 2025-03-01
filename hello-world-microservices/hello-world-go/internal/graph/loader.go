package graph

import (
	"hello-world-go/internal/models"
	"log"

	"github.com/graphql-go/graphql"
)

// ✅ Declare Schema Variable
var schema graphql.Schema

// ✅ Initialize Schema in init() After Query Definition
func init() {
	// ✅ Ensure queryType is initialized before schema creation
	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"serviceConfig":    serviceConfigQuery,
			"requestMetadata":  requestMetadataQuery,
			"responseMetadata": responseMetadataQuery,
			"logContext":       logContextQuery,
			"logEntry":         logEntryQuery,
			"snapshots":        snapshotQuery,
			"contexts":         contextQuery,
		},
	})

	// ✅ Ensure queryType is NOT nil before using it
	if queryType == nil {
		log.Fatalf("CRITICAL ERROR: queryType is still nil before initializing GraphQL schema")
	}

	// ✅ Now create the GraphQL schema
	var err error
	schema, err = graphql.NewSchema(graphql.SchemaConfig{
		Query: queryType, // ✅ queryType is guaranteed to be initialized now
	})
	if err != nil {
		log.Fatalf("ERROR: Failed to create GraphQL schema: %v", err)
	}
}

// ✅ Provide Schema Without Direct Import
func GetGraphQLSchema() graphql.Schema {
	return schema
}

// ✅ Helper function to create GraphQL field with logging
func createField(fieldType graphql.Output, resolver func() (interface{}, error), debugMessage string) *graphql.Field {
	return &graphql.Field{
		Type: fieldType,
		Resolve: func(p graphql.ResolveParams) (interface{}, error) {
			log.Println(debugMessage)
			return resolver()
		},
	}
}

// ✅ Define ServiceConfig Type
var serviceConfigType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceConfig",
	Fields: graphql.Fields{
		"serviceName":    &graphql.Field{Type: graphql.String},
		"environment":    &graphql.Field{Type: graphql.String},
		"serviceVersion": &graphql.Field{Type: graphql.String},
		"logLevel":       &graphql.Field{Type: graphql.String},
	},
})

// ✅ Define RequestMetadata Type
var requestMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RequestMetadata",
	Fields: graphql.Fields{
		"requestID": &graphql.Field{Type: graphql.String},
		"clientIP":  &graphql.Field{Type: graphql.String},
		"userAgent": &graphql.Field{Type: graphql.String},
	},
})

// ✅ Define ResponseMetadata Type
var responseMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ResponseMetadata",
	Fields: graphql.Fields{
		"statusCode":     &graphql.Field{Type: graphql.Int},
		"responseTimeMS": &graphql.Field{Type: graphql.Float},
	},
})

// ✅ Define LogContext Type
var logContextType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogContext",
	Fields: graphql.Fields{
		"instanceID": &graphql.Field{Type: graphql.String},
		"processID":  &graphql.Field{Type: graphql.Int},
		"serviceIP":  &graphql.Field{Type: graphql.String},
	},
})

// ✅ Define LogEntry Type
var logEntryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogEntry",
	Fields: graphql.Fields{
		"message":     &graphql.Field{Type: graphql.String},
		"timestamp":   &graphql.Field{Type: graphql.String},
		"logLevel":    &graphql.Field{Type: graphql.String},
		"serviceName": &graphql.Field{Type: graphql.String},
	},
})

// ✅ Define Query Fields
var serviceConfigQuery = createField(serviceConfigType, func() (interface{}, error) {
	log.Printf("DEBUG: Resolving ServiceConfig -> %+v", models.HelloInstance.Config)
	return models.HelloInstance.Config, nil
}, "DEBUG: Resolving ServiceConfig")

var requestMetadataQuery = createField(requestMetadataType, func() (interface{}, error) {
	return "RequestMetadata not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving RequestMetadata - Placeholder Implementation")

var responseMetadataQuery = createField(responseMetadataType, func() (interface{}, error) {
	return "ResponseMetadata not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving ResponseMetadata - Placeholder Implementation")

var logContextQuery = createField(logContextType, func() (interface{}, error) {
	return "LogContext not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving LogContext - Placeholder Implementation")

var logEntryQuery = createField(logEntryType, func() (interface{}, error) {
	return "LogEntry not yet integrated into GraphQL snapshots.", nil
}, "DEBUG: Resolving LogEntry - Placeholder Implementation")

//

// ✅ Declare queryType as a nil pointer (avoids cycle)
var queryType *graphql.Object

// ✅ Define Query Execution Type
var queryExecutionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "QueryExecution",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
		"data": &graphql.Field{Type: graphql.String},
	},
})

// ✅ Define Snapshot Type (AFTER Context Type)
var snapshotType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Snapshot",
	Fields: graphql.Fields{
		"id":        &graphql.Field{Type: graphql.String},
		"name":      &graphql.Field{Type: graphql.String},
		"createdAt": &graphql.Field{Type: graphql.String},
		"duration":  &graphql.Field{Type: graphql.String},
		"contexts":  &graphql.Field{Type: graphql.NewList(contextType)},
	},
})

// ✅ Define GraphQL Queries
var snapshotQuery = &graphql.Field{
	Type: graphql.NewList(snapshotType),
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		log.Println("DEBUG: Resolving GraphQL Snapshot Query")
		return models.GetSnapshots(), nil
	},
}

// ✅ Define ErrorCode Type for GraphQL
var errorCodeType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ErrorCode",
	Fields: graphql.Fields{
		"code":  &graphql.Field{Type: graphql.Int},
		"count": &graphql.Field{Type: graphql.Int},
	},
})

// ✅ Define SlowEndpoint Type for GraphQL
var slowEndpointType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlowEndpoint",
	Fields: graphql.Fields{
		"endpoint": &graphql.Field{Type: graphql.String},
		"latency":  &graphql.Field{Type: graphql.Float},
	},
})

// ✅ Define SlowAPICall Type for GraphQL
var slowAPICallType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlowAPICall",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.String},
		"name":         &graphql.Field{Type: graphql.String},
		"slowRequests": &graphql.Field{Type: graphql.Int},
		"p95Latency":   &graphql.Field{Type: graphql.Float},
		"p99Latency":   &graphql.Field{Type: graphql.Float},
		"topEndpoints": &graphql.Field{Type: graphql.NewList(slowEndpointType)}, // ✅ Expect list
	},
})

var slaViolationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlaViolation",
	Fields: graphql.Fields{
		"incident":     &graphql.Field{Type: graphql.String},
		"impact":       &graphql.Field{Type: graphql.String},
		"compensation": &graphql.Field{Type: graphql.String},
	},
})

// ✅ Modify Context Type to Include `Error Rates`
var contextType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Context",
	Fields: graphql.Fields{
		"id":                &graphql.Field{Type: graphql.String},
		"name":              &graphql.Field{Type: graphql.String},
		"description":       &graphql.Field{Type: graphql.String},
		"queries":           &graphql.Field{Type: graphql.NewList(queryExecutionType)},   // ✅ Exists in YAML
		"errorRates":        &graphql.Field{Type: graphql.NewList(errorRateType)},        // ✅ Exists in YAML
		"slowestAPICalls":   &graphql.Field{Type: graphql.NewList(slowAPICallType)},      // ✅ Exists in YAML
		"failedEndpoints":   &graphql.Field{Type: graphql.NewList(failedEndpointType)},   // ✅ Exists in YAML
		"slaViolations":     &graphql.Field{Type: graphql.NewList(slaViolationType)},     // ✅ Exists in YAML
		"slowQueries":       &graphql.Field{Type: graphql.NewList(slowQueryType)},        // ✅ Exists in YAML
		"autoscalingEvents": &graphql.Field{Type: graphql.NewList(autoscalingEventType)}, // ✅ Exists in YAML
		"errorBudget":       &graphql.Field{Type: errorBudgetType},                       // ✅ Exists in YAML
		"sloThresholds":     &graphql.Field{Type: graphql.NewList(sloThresholdType)},     // ✅ Exists in YAML
	},
})

// ✅ Ensure Context Query Includes `Error Rates`
var contextQuery = &graphql.Field{
	Type: graphql.NewList(contextType),
	Args: graphql.FieldConfigArgument{
		"snapshotID": &graphql.ArgumentConfig{Type: graphql.String},
	},
	Resolve: func(p graphql.ResolveParams) (interface{}, error) {
		snapshotID, _ := p.Args["snapshotID"].(string)
		log.Printf("DEBUG: Resolving GraphQL Context Query for snapshot %s", snapshotID)

		// ✅ Replace Mock Data with Real Event Queries
		contexts := models.GetSnapshotData(snapshotID)

		// ✅ For each context, replace mock error rates with real data
		// ** FIX.>,
		for i, _ := range contexts {
			log.Println("DEBUG: Fetching error rates from logs...")
			contexts[i].ErrorRates = models.GetErrorRatesFromLogs() // ✅ Use real logs!
		}

		return contexts, nil
	},
}

// ✅ Define FailedEndpoint Type for GraphQL
var failedEndpointType = graphql.NewObject(graphql.ObjectConfig{
	Name: "FailedEndpoint",
	Fields: graphql.Fields{
		"endpoint":    &graphql.Field{Type: graphql.String},
		"failures":    &graphql.Field{Type: graphql.Int},
		"failureRate": &graphql.Field{Type: graphql.Float},
	},
})

// ✅ Modify ErrorRate Type to Include `FailedEndpoints`
var errorRateType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ErrorRate",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.String},
		"name":            &graphql.Field{Type: graphql.String},
		"totalRequests":   &graphql.Field{Type: graphql.Int},
		"failedRequests":  &graphql.Field{Type: graphql.Int},
		"errorRate":       &graphql.Field{Type: graphql.Float},
		"topErrorCodes":   &graphql.Field{Type: graphql.NewList(errorCodeType)},
		"failedEndpoints": &graphql.Field{Type: graphql.NewList(failedEndpointType)},
		"slowestAPICalls": &graphql.Field{Type: graphql.NewList(slowAPICallType)}, // ✅ Re-added slow API call tracking
	},
})

var slowQueryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlowQuery",
	Fields: graphql.Fields{
		"query":         &graphql.Field{Type: graphql.String},
		"executionTime": &graphql.Field{Type: graphql.String},
		"impact":        &graphql.Field{Type: graphql.String},
	},
})

var autoscalingEventType = graphql.NewObject(graphql.ObjectConfig{
	Name: "AutoscalingEvent",
	Fields: graphql.Fields{
		"timestamp":       &graphql.Field{Type: graphql.String},
		"event":           &graphql.Field{Type: graphql.String},
		"instancesBefore": &graphql.Field{Type: graphql.Int},
		"instancesAfter":  &graphql.Field{Type: graphql.Int},
	},
})

var errorBudgetType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ErrorBudget",
	Fields: graphql.Fields{
		"remaining":     &graphql.Field{Type: graphql.String},
		"burnRate":      &graphql.Field{Type: graphql.String},
		"depletionDate": &graphql.Field{Type: graphql.String},
	},
})

var sloThresholdType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SloThreshold",
	Fields: graphql.Fields{
		"metric": &graphql.Field{Type: graphql.String},
		"actual": &graphql.Field{Type: graphql.String},
		"target": &graphql.Field{Type: graphql.String},
		"status": &graphql.Field{Type: graphql.String},
	},
})

// ✅ Now, initialize queryType inside init() (AFTER all dependencies are resolved)
func init() {
	queryType = graphql.NewObject(graphql.ObjectConfig{
		Name: "Query",
		Fields: graphql.Fields{
			"serviceConfig":    serviceConfigQuery,
			"requestMetadata":  requestMetadataQuery,
			"responseMetadata": responseMetadataQuery,
			"logContext":       logContextQuery,
			"logEntry":         logEntryQuery,
			"snapshots":        snapshotQuery,
			"contexts":         contextQuery,
		},
	})

	// ✅ Now replace placeholder in contextType (ensuring queryType is fully built)
	contextType.AddFieldConfig("queries", &graphql.Field{
		Type: graphql.NewList(queryExecutionType),
	})
}

// ✅ Extend GraphQL Query Type
func extendQueryType(query *graphql.Object) {
	query.AddFieldConfig("snapshots", snapshotQuery)
	query.AddFieldConfig("contexts", contextQuery)
}
