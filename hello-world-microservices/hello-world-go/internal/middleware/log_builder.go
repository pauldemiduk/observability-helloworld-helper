// log_builder.go - Logging and metrics middleware with OTEL tracing and structured JSON logging
package middleware

import (
	"fmt"
	"log"
	"time"

	"hello-world-go/internal/models"
	"hello-world-go/internal/util"
)

// BuildLogEntry -
// Constructs a structured log entry using LogEntryOptions
func BuildLogEntry(opts models.LogEntryOptions) models.LogEntry {

	// start a log entry
	var logEntry models.LogEntry

	// Populate using structured helper functions based on log type
	switch opts.LogType {
	case "request":
		logEntry = buildRequestLogEntry(opts)
	case "response":
		logEntry = buildResponseLogEntry(opts)
	case "execution":
		logEntry = buildExecutionLogEntry(opts)
	case "error":
		logEntry = buildErrorLogEntry(opts)
	default:
		log.Printf("WARNING: Unknown log type '%s' in buildLogEntry", opts.LogType)
	}

	// ✅ Core Metadata
	logEntry.Timestamp = time.Now().UTC().Format(time.RFC3339Nano)
	logEntry.Environment = models.HelloInstance.Config.Environment
	logEntry.LogLevel = opts.LogLevel

	// ✅ Service Metadata (Service Details)
	logEntry.ServiceName = models.HelloInstance.Config.ServiceName
	logEntry.ServiceVersion = models.HelloInstance.Config.ServiceVersion
	logEntry.ServiceHostname = models.HelloInstance.Config.ServiceHostname
	logEntry.ServiceIP = models.HelloInstance.Config.ServiceIP
	logEntry.ServicePort = models.HelloInstance.Config.ServicePort
	logEntry.InstanceID = models.HelloInstance.Context.InstanceID
	logEntry.Language = "go"

	// ✅ Host & Deployment Context (Where the service runs)
	logEntry.ProcessID = models.HelloInstance.Context.ProcessID
	logEntry.ThreadID = models.HelloInstance.Context.ThreadID
	logEntry.NodeIP = models.HelloInstance.Context.NodeIP
	logEntry.NodePort = models.HelloInstance.Context.NodePort
	logEntry.LoadBalancerIP = models.HelloInstance.Context.LoadBalancerIP

	// ✅ Cloud & Infrastructure Metadata
	logEntry.CloudProvider = models.HelloInstance.Context.CloudProvider
	logEntry.CloudRegion = models.HelloInstance.Context.CloudRegion
	logEntry.AvailabilityZone = models.HelloInstance.Context.AvailabilityZone
	logEntry.K8sNamespace = models.HelloInstance.Context.K8sNamespace
	logEntry.K8sPodName = models.HelloInstance.Context.K8sPodName
	logEntry.K8sContainerID = models.HelloInstance.Context.K8sContainerID

	// ✅ Debugging & Observability Data
	logEntry.CallerFunction = util.ExtractCallerFunction()
	logEntry.ExecutionUnitID = util.GetExecutionUnitID()
	logEntry.CallerLocation = util.ExtractCallerLocation()
	logEntry.ComputeTimeMS = opts.Duration.Seconds() * 1000 // ✅ Standardized naming
	logEntry.QueueTimeMS = opts.QueueTimeMS
	// ** VERIFY.>, should this be on ResponseMetadata.?
	logEntry.ResponseLatencyLevel = util.CategorizeResponseLatency(opts.ResponseMetadata.ResponseTimeMS)

	return logEntry
}

// buildRequestLogEntry -
// Populate LogEntry fields based on http.Request - GET, PUT, POST
func buildRequestLogEntry(opts models.LogEntryOptions) models.LogEntry {
	requestMetadata := opts.RequestMetadata

	// ** RequestSizeBytes
	return models.LogEntry{
		// ✅ Request Lifecycle
		RequestStage: "received",

		// ✅ Request Identification
		RequestID:     requestMetadata.RequestID,
		Operation:     requestMetadata.Operation,
		RequestMethod: requestMetadata.RequestMethod,
		RequestURL:    opts.Request.URL.String(),

		// ✅ Request Headers & Client Data
		RequestHeaders:     requestMetadata.RequestHeaders,
		ClientIP:           requestMetadata.ClientIP,
		ClientPort:         requestMetadata.ClientPort,
		UserAgent:          requestMetadata.UserAgent,
		Origin:             requestMetadata.Origin,
		Referer:            requestMetadata.Referer,
		Protocol:           requestMetadata.Protocol,
		TLSProtocolVersion: requestMetadata.TLSProtocolVersion,
		AuthHeaderStatus:   requestMetadata.AuthHeaderStatus,

		// ✅ Tracing & Correlation
		TraceID:       requestMetadata.TraceID,
		ParentTraceID: requestMetadata.ParentTraceID,
		SpanID:        requestMetadata.SpanID,

		// ✅ Request Metadata
		UserSessionID:    requestMetadata.UserSessionID,
		RequestSizeBytes: requestMetadata.RequestSizeBytes,
		ClientID:         requestMetadata.ClientID,
		CorrelationID:    requestMetadata.CorrelationID,
		UserID:           requestMetadata.UserID,

		// ✅ Performance Metrics
		QueueTimeMS: opts.QueueTimeMS,

		// ✅ Log Message
		Message: fmt.Sprintf("Received request %s, awaiting response..", opts.Request.URL.Path),
	}
}

// buildResponseLogEntry -
// Populate LogEntry fields based on http.Response - Route handlers
func buildResponseLogEntry(opts models.LogEntryOptions) models.LogEntry {
	responseMetadata := opts.ResponseMetadata
	requestMetadata := opts.RequestMetadata

	return models.LogEntry{
		// ✅ Request Lifecycle
		RequestStage: "processed",

		// ✅ Request Identification
		RequestID:     requestMetadata.RequestID,
		Operation:     requestMetadata.Operation,
		RequestMethod: requestMetadata.RequestMethod,
		RequestURL:    opts.Request.URL.String(),

		// ✅ Request Headers & Client Data
		RequestHeaders:     requestMetadata.RequestHeaders,
		ClientIP:           requestMetadata.ClientIP,
		ClientPort:         requestMetadata.ClientPort,
		UserAgent:          requestMetadata.UserAgent,
		Origin:             requestMetadata.Origin,
		Referer:            requestMetadata.Referer,
		Protocol:           requestMetadata.Protocol,
		TLSProtocolVersion: requestMetadata.TLSProtocolVersion,
		AuthHeaderStatus:   requestMetadata.AuthHeaderStatus,

		// ✅ Tracing & Correlation
		TraceID:       requestMetadata.TraceID,
		ParentTraceID: requestMetadata.ParentTraceID,
		SpanID:        requestMetadata.SpanID,

		// ✅ Request Metadata
		UserSessionID:    requestMetadata.UserSessionID,
		RequestSizeBytes: requestMetadata.RequestSizeBytes,
		ClientID:         requestMetadata.ClientID,
		CorrelationID:    requestMetadata.CorrelationID,
		UserID:           requestMetadata.UserID,

		// ✅ Response Metadata
		HTTPStatusCode:       responseMetadata.HTTPStatusCode,
		ResponseOutcome:      responseMetadata.ResponseOutcome,
		ResponseSizeBytes:    responseMetadata.ResponseSizeBytes,
		ResponseHeaders:      responseMetadata.ResponseHeaders,
		ResponseTimeMS:       responseMetadata.ResponseTimeMS,
		ResponseTimeCategory: responseMetadata.ResponseTimeCategory,

		// ✅ Log Message
		Message: fmt.Sprintf("Processed request %s with status %d", opts.Request.URL.Path, responseMetadata.HTTPStatusCode),
	}
}

// buildExecutionLogEntry -
// Populate LogEntry fields for simple logging use cases, code flow and execution not http lifecycle
// ** VERIFY.>, more attributes needed on code execution logs.?
func buildExecutionLogEntry(opts models.LogEntryOptions) models.LogEntry {
	cpuUsage, memoryUsage := util.CaptureResourceUsage()

	return models.LogEntry{
		CallerFunction:  util.ExtractCallerFunction(),
		ExecutionUnitID: util.GetExecutionUnitID(),
		CallerLocation:  util.ExtractCallerLocation(),
		CPUUsage:        cpuUsage,
		MemoryUsage:     memoryUsage,
		Message:         opts.Message,
	}
}

// buildErrorLogEntry -
// Populate LogEntry fields for simple logging use cases, code flow and execution errors
// ** VERIFY.>, more attributes needed on errors.?
func buildErrorLogEntry(opts models.LogEntryOptions) models.LogEntry {
	return models.LogEntry{
		ErrorMessage: opts.ErrorMessage,
		StackTrace:   opts.StackTrace,
	}
}
