// handlers.go - Logging and metrics middleware with OTEL tracing and structured JSON logging
package middleware

import (
	"fmt"
	"hello-world-go/internal/logging"
	"hello-world-go/internal/models"
	"hello-world-go/internal/util"

	"net/http"
	"time"

	"go.opentelemetry.io/contrib/instrumentation/net/http/otelhttp"
	"go.opentelemetry.io/otel"
)

// WrapHandlerWithObservability -
// Applies OTEL Tracing and Prometheus Middleware to HTTP handlers.
func WrapHandlerWithObservability(pattern string, handler http.Handler) {
	http.Handle(pattern, otelhttp.NewHandler(ObservabilityMiddleware(handler, pattern), pattern))
}

// ObservabilityMiddleware -
// Wraps all http server route resolution with proper 011y signales and telemetry
func ObservabilityMiddleware(handler http.Handler, operation string) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		startTime := time.Now()
		lrw := &LoggingResponseWriter{ResponseWriter: w, StatusCode: http.StatusOK}

		logCtx := models.LoadLogContext()

		tracer := otel.GetTracerProvider().Tracer("hello-world-service")
		ctx, span := tracer.Start(r.Context(), operation)
		defer span.End()
		r = r.WithContext(ctx)

		// ✅ Compute request queue time before extracting metadata
		queueTimeMS := time.Since(startTime).Seconds() * 1000

		// ✅ Extract structured request metadata
		requestMetadata := util.ExtractRequestMetadata(r)

		//log.Printf("DEBUG: RequestSizeBytes: %d", requestMetadata.RequestSizeBytes)

		logOpts := models.LogEntryOptions{
			LRW:              lrw,
			Request:          r,
			StartTime:        startTime,
			QueueTimeMS:      queueTimeMS,
			LogCtx:           logCtx,
			LogLevel:         "info",
			LogType:          "request",
			RequestStage:     "received",
			Message:          fmt.Sprintf("Received request %s, awaiting response..", r.URL.Path),
			RequestMetadata:  requestMetadata,
			ResponseMetadata: models.ResponseMetadata{},
		}

		logging.SendLog(BuildLogEntry(logOpts))

		handler.ServeHTTP(lrw, r)

		responseMetadata := util.ExtractResponseMetadata(lrw, time.Since(startTime))

		logOpts.LogType = "response"
		logOpts.RequestStage = "processed"
		logOpts.Message = fmt.Sprintf("Processed request %s with status %d", r.URL.Path, responseMetadata.HTTPStatusCode)
		logOpts.ResponseMetadata = responseMetadata

		logging.SendLog(BuildLogEntry(logOpts))
	})
}
