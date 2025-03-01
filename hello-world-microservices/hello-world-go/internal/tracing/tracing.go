// tracing.go - Tracing middleware with OpenTelemetry

package tracing

import (
	"fmt"
	"log"
	"sync"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

var tracerProvider *trace.TracerProvider
var traceOnce sync.Once
var tracerErr error

func SetupTracing() (*trace.TracerProvider, error) {
	traceOnce.Do(func() {
		tracerProvider = trace.NewTracerProvider(
			trace.WithResource(resource.NewWithAttributes(
				semconv.SchemaURL,
				semconv.ServiceNameKey.String("hello-world-go"),
			)),
		)

		// ✅ Check if Tracer was created successfully
		if tracerProvider == nil {
			tracerErr = fmt.Errorf("failed to initialize tracer provider")
			log.Printf("ERROR: OpenTelemetry tracing setup failed: %v", tracerErr)
		} else {
			otel.SetTracerProvider(tracerProvider)
			log.Println("INFO: OpenTelemetry tracing initialized")
		}
	})
	return tracerProvider, tracerErr
}
