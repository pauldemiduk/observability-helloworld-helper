// *
// ==> graph_interface.go - [blueprint:include]
// *
// =>  graph_interface.go explained:
//
//	Description:
//	Description: GraphQL schema interface variables.
//
// # GraphQL composition
//
//	Schema
//	- Entity
//	-- Interface 	<= (graph_interface.go)
//	--- Query
//	---- Filter
package models

import "github.com/graphql-go/graphql"

//*
// ==> GraphQL UI filter option interace

// LogFilterInputType -
// Filter the in-mem log datastore
var LogFilterInputType = graphql.NewInputObject(graphql.InputObjectConfig{
	Name: "LogFilterInput",
	Fields: graphql.InputObjectConfigFieldMap{
		"level": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"since": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"httpStatusCode": &graphql.InputObjectFieldConfig{
			Type: graphql.Int,
		},
		"responseTimeCategory": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"requestMethod": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
		"responseLatencyLevel": &graphql.InputObjectFieldConfig{
			Type: graphql.String,
		},
	},
})

//*
// ==> GraphQL variables that map to schema entities.

// queryExecutionType -
// Defines interface for Query Execution Type
var QueryExecutionType = graphql.NewObject(graphql.ObjectConfig{
	Name: "QueryExecution",
	Fields: graphql.Fields{
		"id":   &graphql.Field{Type: graphql.String},
		"name": &graphql.Field{Type: graphql.String},
		"data": &graphql.Field{Type: graphql.String},
	},
})

// serviceConfigType -
// Defines interface for ServiceConfig entity.
// ROADMAP.>, add full attributes.
var serviceConfigType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ServiceConfig",
	Fields: graphql.Fields{
		"serviceName":    &graphql.Field{Type: graphql.String},
		"environment":    &graphql.Field{Type: graphql.String},
		"serviceVersion": &graphql.Field{Type: graphql.String},
		"logLevel":       &graphql.Field{Type: graphql.String},
	},
})

// requestMetadataType -
// Defines interface for RequestMetadata entity.
// ROADMAP.>, add full attributes.
var requestMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "RequestMetadata",
	Fields: graphql.Fields{
		"requestID": &graphql.Field{Type: graphql.String},
		"clientIP":  &graphql.Field{Type: graphql.String},
		"userAgent": &graphql.Field{Type: graphql.String},
	},
})

// responseMetadataType -
// Defines interface for ResponseMetadata entity.
// ROADMAP.>, add full attributes.
var responseMetadataType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ResponseMetadata",
	Fields: graphql.Fields{
		"statusCode":     &graphql.Field{Type: graphql.Int},
		"responseTimeMS": &graphql.Field{Type: graphql.Float},
	},
})

// logContextType -
// Defines interface for LogContext entity.
// ROADMAP.>, add full attributes.
var logContextType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogContext",
	Fields: graphql.Fields{
		"instanceID": &graphql.Field{Type: graphql.String},
		"processID":  &graphql.Field{Type: graphql.Int},
		"serviceIP":  &graphql.Field{Type: graphql.String},
	},
})

// logEntryType -
// Defines interface for LogEntry entity.
var logEntryType = graphql.NewObject(graphql.ObjectConfig{
	Name: "LogEntry",
	Fields: graphql.Fields{
		"timestamp":              &graphql.Field{Type: graphql.String},
		"environment":            &graphql.Field{Type: graphql.String},
		"log_level":              &graphql.Field{Type: graphql.String},
		"service_name":           &graphql.Field{Type: graphql.String},
		"service_version":        &graphql.Field{Type: graphql.String},
		"service_hostname":       &graphql.Field{Type: graphql.String},
		"service_ip":             &graphql.Field{Type: graphql.String},
		"service_port":           &graphql.Field{Type: graphql.Int},
		"process_id":             &graphql.Field{Type: graphql.Int},
		"thread_id":              &graphql.Field{Type: graphql.Int},
		"instance_id":            &graphql.Field{Type: graphql.String},
		"language":               &graphql.Field{Type: graphql.String},
		"node_ip":                &graphql.Field{Type: graphql.String},
		"node_port":              &graphql.Field{Type: graphql.Int},
		"cloud_provider":         &graphql.Field{Type: graphql.String},
		"cloud_region":           &graphql.Field{Type: graphql.String},
		"availability_zone":      &graphql.Field{Type: graphql.String},
		"k8s_namespace":          &graphql.Field{Type: graphql.String},
		"k8s_pod_name":           &graphql.Field{Type: graphql.String},
		"k8s_container_id":       &graphql.Field{Type: graphql.String},
		"load_balancer_ip":       &graphql.Field{Type: graphql.String},
		"request_id":             &graphql.Field{Type: graphql.String},
		"operation":              &graphql.Field{Type: graphql.String},
		"request_method":         &graphql.Field{Type: graphql.String},
		"request_size_bytes":     &graphql.Field{Type: graphql.Int},
		"request_url":            &graphql.Field{Type: graphql.String},
		"request_headers":        &graphql.Field{Type: graphql.NewList(graphql.String)},
		"client_ip":              &graphql.Field{Type: graphql.String},
		"client_port":            &graphql.Field{Type: graphql.Int},
		"user_agent":             &graphql.Field{Type: graphql.String},
		"origin":                 &graphql.Field{Type: graphql.String},
		"referer":                &graphql.Field{Type: graphql.String},
		"protocol":               &graphql.Field{Type: graphql.String},
		"tls_protocol_version":   &graphql.Field{Type: graphql.String},
		"auth_header_status":     &graphql.Field{Type: graphql.String},
		"request_stage":          &graphql.Field{Type: graphql.String},
		"client_id":              &graphql.Field{Type: graphql.String},
		"user_session_id":        &graphql.Field{Type: graphql.String},
		"trace_id":               &graphql.Field{Type: graphql.String},
		"parent_trace_id":        &graphql.Field{Type: graphql.String},
		"span_id":                &graphql.Field{Type: graphql.String},
		"correlation_id":         &graphql.Field{Type: graphql.String},
		"user_id":                &graphql.Field{Type: graphql.String},
		"http_status_code":       &graphql.Field{Type: graphql.Int},
		"response_outcome":       &graphql.Field{Type: graphql.String},
		"response_size_bytes":    &graphql.Field{Type: graphql.Int},
		"response_time_ms":       &graphql.Field{Type: graphql.Float},
		"response_time_category": &graphql.Field{Type: graphql.String},
		"response_headers":       &graphql.Field{Type: graphql.NewList(graphql.String)},
		"caller_function":        &graphql.Field{Type: graphql.String},
		"execution_unit_id":      &graphql.Field{Type: graphql.Int},
		"caller_location":        &graphql.Field{Type: graphql.String},
		"compute_time_ms":        &graphql.Field{Type: graphql.Float},
		"queue_time_ms":          &graphql.Field{Type: graphql.Float},
		"response_latency_level": &graphql.Field{Type: graphql.String},
		"cpu_usage":              &graphql.Field{Type: graphql.Float},
		"memory_usage":           &graphql.Field{Type: graphql.Float},
		"error_message":          &graphql.Field{Type: graphql.String},
		"stack_trace":            &graphql.Field{Type: graphql.String},
		"message":                &graphql.Field{Type: graphql.String},
	},
})

// contextType -
// Defines interface for Context entity.
var contextType = graphql.NewObject(graphql.ObjectConfig{
	Name: "Context",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"queries":     &graphql.Field{Type: graphql.NewList(queryExecutionType)},
		"useCase":     &graphql.Field{Type: useCaseType},
	},
})

// snapshotType -
// Defines interface for Snapshot entity.
// Note: ==> load dependency is important, this goes after Context.
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

var useCaseType = graphql.NewObject(graphql.ObjectConfig{
	Name: "UseCase",
	Fields: graphql.Fields{
		"id":          &graphql.Field{Type: graphql.String},
		"name":        &graphql.Field{Type: graphql.String},
		"description": &graphql.Field{Type: graphql.String},
		"slowAPICall": &graphql.Field{Type: slowAPICallType},
	},
})

//
// ==> Define Context Detail Helpers
//

// slowAPICallType -
// Defines interface for slowAPICall Type (name/value helper)
var slowAPICallType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlowAPICall",
	Fields: graphql.Fields{
		"id":           &graphql.Field{Type: graphql.String},
		"name":         &graphql.Field{Type: graphql.String},
		"duration":     &graphql.Field{Type: graphql.Int},
		"slowRequests": &graphql.Field{Type: graphql.Int},
		"p95Latency":   &graphql.Field{Type: graphql.Float},
		"p99Latency":   &graphql.Field{Type: graphql.Float},
		"topEndpoints": &graphql.Field{Type: graphql.NewList(slowEndpointType)},
	},
})

// slowEndpointType -
// Defines interface for SlowEndpoint Type (name/value helper)
var slowEndpointType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SlowEndpoint",
	Fields: graphql.Fields{
		"endpoint": &graphql.Field{Type: graphql.String},
		"latency":  &graphql.Field{Type: graphql.Float},
	},
})

var errorRateType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ErrorRate",
		Fields: graphql.Fields{
			"id": &graphql.Field{
				Type: graphql.String,
			},
			"name": &graphql.Field{
				Type: graphql.String,
			},
			"duration": &graphql.Field{
				Type: graphql.Int,
			},
			"totalRequests": &graphql.Field{
				Type: graphql.Int,
			},
			"failedRequests": &graphql.Field{
				Type: graphql.Int,
			},
			"errorRate": &graphql.Field{
				Type: graphql.Float,
			},
			"topErrorCodes": &graphql.Field{
				Type: graphql.NewList(errorCodeType), // ✅ Reference error codes
			},
			"failedEndpoints": &graphql.Field{
				Type: graphql.NewList(failedEndpointType), // ✅ Reference failed endpoints
			},
		},
	},
)

var errorCodeType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "ErrorCode",
		Fields: graphql.Fields{
			"code": &graphql.Field{
				Type: graphql.Int,
			},
			"count": &graphql.Field{
				Type: graphql.Int,
			},
		},
	},
)

var failedEndpointType = graphql.NewObject(
	graphql.ObjectConfig{
		Name: "FailedEndpoint",
		Fields: graphql.Fields{
			"endpoint": &graphql.Field{
				Type: graphql.String,
			},
			"failures": &graphql.Field{
				Type: graphql.Int,
			},
			"failureRate": &graphql.Field{
				Type: graphql.Float,
			},
		},
	},
)

// slaViolationType - Defines interface for SLA violations
var slaViolationType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SLAViolation",
	Fields: graphql.Fields{
		"id":              &graphql.Field{Type: graphql.String},
		"affectedService": &graphql.Field{Type: graphql.String},
		"breachTime":      &graphql.Field{Type: graphql.String},
		"reason":          &graphql.Field{Type: graphql.String},
		"impact":          &graphql.Field{Type: graphql.String},
	},
})

// sloComplianceType - Defines interface for SLO tracking
var sloComplianceType = graphql.NewObject(graphql.ObjectConfig{
	Name: "SLOCompliance",
	Fields: graphql.Fields{
		"id":     &graphql.Field{Type: graphql.String},
		"target": &graphql.Field{Type: graphql.String},
		"actual": &graphql.Field{Type: graphql.String},
		"status": &graphql.Field{Type: graphql.String},
	},
})

// errorBudgetType - Defines interface for error budget monitoring
var errorBudgetType = graphql.NewObject(graphql.ObjectConfig{
	Name: "ErrorBudget",
	Fields: graphql.Fields{
		"id":                &graphql.Field{Type: graphql.String},
		"totalAllocated":    &graphql.Field{Type: graphql.String},
		"used":              &graphql.Field{Type: graphql.String},
		"remaining":         &graphql.Field{Type: graphql.String},
		"depletionForecast": &graphql.Field{Type: graphql.String},
	},
})
