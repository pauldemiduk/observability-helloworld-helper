package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"gopkg.in/yaml.v2"

	"hello-world-go/internal/config"
	"hello-world-go/internal/graph"
	"hello-world-go/internal/handlers"
	"hello-world-go/internal/logging"
	"hello-world-go/internal/metrics"
	"hello-world-go/internal/middleware"
	"hello-world-go/internal/models"
	"hello-world-go/internal/tracing"
)

// ✅ Structs to Parse `_instance.yaml`
type InstanceConfig struct {
	StartupFlow map[string]StepConfig `yaml:"Startup-Flow"`
}

type StepConfig struct {
	ExecuteStep string `yaml:"execute_step"` // ✅ Controls execution
	LogStep     string `yaml:"log_step"`     // ✅ Logs the attempt
	LogOutcome  string `yaml:"log_outcome"`  // ✅ Logs the result
	LogValues   string `yaml:"log_values"`   // ✅ New field to log detailed values
}

const instanceFilePath = "configs/_instance.yaml"

func main() {
	log.Println("DEBUG: Forced test log - If this appears, debug logs are working!")

	startTime := time.Now()

	// ✅ Load `_instance.yaml` with validation
	instance, err := loadInstanceConfig(instanceFilePath)
	if err != nil {
		log.Printf("WARNING: Failed to load instance configuration: %v", err)
		log.Println("⚠️  Falling back to default startup flow...")
		instance = getDefaultInstanceConfig()
	}

	log.Println("INFO: Initializing service...")

	// ✅ Initialize Router
	logStep(instance, "new-router", "Attempting to initialize router")
	router := mux.NewRouter()
	logStepOutcome(instance, "new-router", "Router initialized")

	registerRoutes(router, instance)

	// ✅ Set GraphQL Schema
	logStep(instance, "set-graph-schema", "Attempting to set GraphQL schema")
	handlers.SetGraphQLSchema(graph.GetGraphQLSchema())
	logStepOutcome(instance, "set-graph-schema", "GraphQL schema set")

	// ✅ Initialize Instance
	logStep(instance, "new-instance", "Attempting to create service instance")
	models.HelloInstance = models.NewServiceInstance()
	logStepOutcome(instance, "new-instance", "Service instance created")

	//

	// ✅ Load Configuration (Ensuring Required Defaults)
	logStep(instance, "load-config", "Attempting to load service config")
	models.HelloInstance.Config = config.LoadConfig()

	// ✅ Log Config Details if `log_values: "yes"` (BEFORE marking step complete)
	if instance.StartupFlow["load-config"].LogValues == "yes" {
		logServiceConfig(models.HelloInstance.Config)
	}

	// ✅ Now Mark the Step as Complete
	logStepOutcome(instance, "load-config", "Service config loaded")

	//

	// ✅ Build Service Context (New Step from _instance.yaml)
	logStep(instance, "load-service-context", "Attempting to load service context")
	models.HelloInstance.Context = config.LoadServiceContext(models.HelloInstance.Config)
	// ✅ Log Service Context Details if `log_values: "yes"` (BEFORE marking step complete)
	if instance.StartupFlow["load-service-context"].LogValues == "yes" {
		logServiceContext(models.HelloInstance.Context)
	}
	// ✅ Now Mark the Step as Complete
	logStepOutcome(instance, "load-service-context", "Service context loaded")

	// ✅ Force GetSnapshots() (Ensure Default Behavior)
	logStep(instance, "get-snapshots", "Attempting to load snapshots")
	snapshots := models.GetSnapshots()
	if len(snapshots) == 0 {
		log.Println("WARNING: No snapshots loaded at startup.")
	}
	logStepOutcome(instance, "get-snapshots", fmt.Sprintf("Loaded %d snapshots", len(snapshots)))

	// ✅ Initialize Prometheus Metrics
	logStep(instance, "setup-metrics", "Attempting to setup metrics")
	metrics.SetupMetrics()
	logStepOutcome(instance, "setup-metrics", "Metrics initialized")

	// ✅ Initialize OpenTelemetry Tracing
	logStep(instance, "setup-tracing", "Attempting to setup tracing")
	tp, err := tracing.SetupTracing()
	if err != nil {
		logStepOutcome(instance, "setup-tracing", "Tracing initialization failed")
		log.Fatalf("ERROR: Failed to initialize OTEL tracer: %v", err)
	}
	defer tp.Shutdown(context.Background())
	logStepOutcome(instance, "setup-tracing", "Tracing initialized")

	// ✅ Initialize Log Context (Ensure Required Values Exist)
	logStep(instance, "load-log-context", "Attempting to load log context")
	logCtx := models.LoadLogContext()
	// ✅ Log Log Context Details if `log_values: "yes"` (BEFORE marking step complete)
	if instance.StartupFlow["load-log-context"].LogValues == "yes" {
		logLogContext(logCtx)
	}

	logStepOutcome(instance, "load-log-context", fmt.Sprintf("Log context initialized: %+v", logCtx))

	// ✅ Register Middleware & Handlers
	logStep(instance, "register-handlers", "Attempting to register handlers")
	logStepOutcome(instance, "register-handlers", "Handlers registered")

	// ✅ Start HTTP Server
	logStep(instance, "start-http-server", "Attempting to start HTTP server")
	server := &http.Server{Addr: models.HelloInstance.Config.ServiceHostname + ":" + strconv.Itoa(models.HelloInstance.Config.ServicePort)}
	logStepOutcome(instance, "start-http-server", "HTTP server started")
	startHTTPServer(server, instance, startTime)

	// ✅ Handle Graceful Shutdown
	handleShutdown(instance, server)

}

// ✅ Load `_instance.yaml` with validation
func loadInstanceConfig(filepath string) (InstanceConfig, error) {
	file, err := ioutil.ReadFile(filepath)
	if err != nil {
		return InstanceConfig{}, err
	}

	var instance InstanceConfig
	err = yaml.Unmarshal(file, &instance)
	if err != nil {
		return InstanceConfig{}, err
	}

	// ✅ Validate Required Fields in `_instance.yaml`
	for step, config := range instance.StartupFlow {
		if config.LogStep != "yes" && config.LogStep != "no" {
			log.Printf("WARNING: Invalid log_step value in %s, defaulting to 'no'", step)
			updatedConfig := config
			updatedConfig.LogStep = "no"
			instance.StartupFlow[step] = updatedConfig // ✅ Reassign struct back into the map
		}

		if config.LogOutcome != "yes" && config.LogOutcome != "no" {
			log.Printf("WARNING: Invalid log_outcome value in %s, defaulting to 'no'", step)
			updatedConfig := config
			updatedConfig.LogOutcome = "no"
			instance.StartupFlow[step] = updatedConfig // ✅ Reassign struct back into the map
		}
	}

	return instance, nil
}

// ✅ Get Default Startup Flow if `_instance.yaml` is missing
func getDefaultInstanceConfig() InstanceConfig {
	return InstanceConfig{
		StartupFlow: map[string]StepConfig{
			"log-start":            {"yes", "yes", "yes", "no"},
			"new-router":           {"yes", "yes", "yes", "no"},
			"register-routes":      {"yes", "yes", "yes", "no"},
			"set-graph-schema":     {"yes", "yes", "yes", "no"},
			"new-instance":         {"yes", "yes", "yes", "yes"},
			"load-config":          {"yes", "yes", "yes", "yes"},
			"load-service-context": {"yes", "yes", "yes", "yes"},
			"get-snapshots":        {"yes", "yes", "yes", "no"},
			"setup-metrics":        {"yes", "yes", "yes", "no"},
			"uptime-ticker":        {"yes", "yes", "yes", "no"},
			"setup-tracing":        {"yes", "yes", "yes", "no"},
			"load-log-context":     {"yes", "yes", "yes", "yes"},
			"register-handlers":    {"yes", "yes", "yes", "no"},
			"start-http-server":    {"yes", "yes", "yes", "no"},
			"channel-for-errors":   {"yes", "yes", "yes", "no"},
			"listen-and-serve":     {"yes", "yes", "yes", "no"},
			"confirm-startup":      {"yes", "yes", "yes", "no"},
		},
	}
}

// ✅ Log Step Attempt
func logStep(instance InstanceConfig, step string, message string) {
	if instance.StartupFlow[step].LogStep == "yes" {
		log.Printf("INFO: %s - %s", step, message)
	}
}

// ✅ Log Step Outcome with `log-values` toggle = yes
func logStepOutcome(instance InstanceConfig, step string, message string) {
	if instance.StartupFlow[step].LogOutcome == "yes" {
		log.Printf("INFO: %s - %s", step, message)
	}

	// ✅ Log additional values if `log-values` is enabled for this step
	if instance.StartupFlow[step].LogValues == "yes" {
		switch step {
		case "log-start":
			log.Println("INFO: Log system initialized.")

		case "new-router":
			log.Println("INFO: Router instance created and ready.")

		case "register-routes":
			log.Println("INFO: Routes registered with observability middleware.")

		case "set-graph-schema":
			log.Println("INFO: GraphQL schema applied.")

		case "new-instance":
			logServiceInstance(models.HelloInstance) // ✅ Uses the new function
			//log.Printf("INFO: ServiceInstance Details: %+v", models.HelloInstance)

		case "load-config":
			log.Printf("INFO: Service config initialized.")

		case "load-service-context":
			log.Printf("INFO: Service context initialized.")

		case "get-snapshots":
			snapshots := models.GetSnapshots()
			log.Printf("INFO: Loaded %d snapshots.", len(snapshots))

		case "setup-metrics":
			log.Println("INFO: Metrics registered with Prometheus.")

		case "uptime-ticker":
			log.Println("INFO: Uptime ticker initialized.")

		case "setup-tracing":
			log.Println("INFO: OpenTelemetry tracing enabled.")

		case "load-log-context":
			log.Printf("INFO: Log context initialized.")

		case "register-handlers":
			log.Println("INFO: All HTTP handlers registered.")

		case "start-http-server":
			log.Printf("INFO: HTTP server running at %s:%d", models.HelloInstance.Config.ServiceHostname, models.HelloInstance.Config.ServicePort)

		case "channel-for-errors":
			log.Println("INFO: Error channel for tracking service failures is active.")

		case "listen-and-serve":
			log.Println("INFO: Listening for incoming HTTP requests.")

		case "confirm-startup":
			log.Println("INFO: Service startup confirmed, all systems operational.")

		default:
			log.Printf("INFO: Step '%s' executed, but no specific log-values message is defined.", step)
		}
	}
}

// ✅ Log ServiceConfig as Compact JSON
func logServiceConfig(config models.ServiceConfig) {
	configJSON, err := json.Marshal(config)
	if err != nil {
		log.Printf("ERROR: Failed to serialize ServiceConfig: %v", err)
		return
	}
	log.Printf("INFO: ServiceConfig details: %s", configJSON)
}

// ✅ Log ServiceContext as Compact JSON
func logServiceContext(ctx models.ServiceContext) {
	ctxJSON, err := json.Marshal(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to serialize ServiceContext: %v", err)
		return
	}
	log.Printf("INFO: ServiceContext Details: %s", ctxJSON)
}

// ✅ Log LogContext as Compact JSON
func logLogContext(ctx models.LogContext) {
	ctxJSON, err := json.Marshal(ctx)
	if err != nil {
		log.Printf("ERROR: Failed to serialize LogContext: %v", err)
		return
	}
	log.Printf("INFO: LogContext Details: %s", ctxJSON)
}

// ✅ Log Full ServiceInstance (Includes Config + Context)
func logServiceInstance(instance *models.ServiceInstance) {
	if instance == nil {
		log.Println("WARNING: ServiceInstance is nil—unable to log details.")
		return
	}
	// ✅ Include Queue Details in ServiceInstance Logging
	instanceData := map[string]interface{}{
		"Config":  instance.Config,
		"Context": instance.Context,
		"Queues":  json.RawMessage(formatQueueStatus(instance)),
	}

	instanceJSON, err := json.Marshal(instanceData)
	if err != nil {
		log.Printf("ERROR: Failed to serialize ServiceInstance: %v", err)
		return
	}
	log.Printf("INFO: ServiceInstance Initialized: %s", instanceJSON)
}

// ✅ Format Queue Details
func formatQueueStatus(instance *models.ServiceInstance) string {
	ingestQueue := "  - IngestQueue: Not initialized"
	workloadQueue := "  - WorkloadQueue: Not initialized"

	if instance.IngestQueue != nil {
		ingestQueue = fmt.Sprintf("  - IngestQueue: %d items / Max: %d", instance.IngestQueue.Length(), instance.IngestQueue.GetMaxSize())
	}
	if instance.WorkloadQueue != nil {
		workloadQueue = fmt.Sprintf("  - WorkloadQueue: %d items / Max: %d", instance.WorkloadQueue.Length(), instance.WorkloadQueue.GetMaxSize())
	}

	return fmt.Sprintf("\n%s\n%s\n", ingestQueue, workloadQueue)
}

//
//

// ✅ Standardized Route Registration - All Routes Are Now Consistent
func registerRoutes(router *mux.Router, instance InstanceConfig) {
	if instance.StartupFlow["register-routes"].ExecuteStep == "yes" {
		// ✅ Register Core API Handlers in `mux.Router`
		router.HandleFunc("/hello", handlers.HelloHandler).Methods("GET")
		router.HandleFunc("/status", handlers.StatusHandler).Methods("GET")
		router.HandleFunc("/inbound", handlers.InboundHandler).Methods("POST")
		router.Handle("/metrics", promhttp.Handler()).Methods("GET")
		router.HandleFunc("/queue/push", handlers.QueuePushHandler).Methods("POST")
		router.HandleFunc("/queue/pop", handlers.QueuePopHandler).Methods("GET")
		router.HandleFunc("/schema", handlers.SchemaHandler).Methods("GET")
		router.HandleFunc("/graph", handlers.GraphQLHandler).Methods("GET", "POST")
		router.HandleFunc("/outbound", handlers.OutboundHandler).Methods("POST")

		// ✅ Apply Observability Middleware for All Routes
		middleware.WrapHandlerWithObservability("/hello", http.HandlerFunc(handlers.HelloHandler))
		middleware.WrapHandlerWithObservability("/status", http.HandlerFunc(handlers.StatusHandler))
		middleware.WrapHandlerWithObservability("/inbound", http.HandlerFunc(handlers.InboundHandler))
		middleware.WrapHandlerWithObservability("/metrics", promhttp.Handler())
		middleware.WrapHandlerWithObservability("/queue/push", http.HandlerFunc(handlers.QueuePushHandler))
		middleware.WrapHandlerWithObservability("/queue/pop", http.HandlerFunc(handlers.QueuePopHandler))
		middleware.WrapHandlerWithObservability("/schema", http.HandlerFunc(handlers.SchemaHandler))
		middleware.WrapHandlerWithObservability("/graph", http.HandlerFunc(handlers.GraphQLHandler))
		middleware.WrapHandlerWithObservability("/outbound", http.HandlerFunc(handlers.OutboundHandler))

		// ✅ Resolve Swagger Documentation Path
		swaggerPath := config.ResolveDocsPath()
		swaggerFileServer := http.FileServer(http.Dir(swaggerPath))
		swaggerHandler := http.StripPrefix("/swagger/", swaggerFileServer)
		router.PathPrefix("/swagger/").Handler(swaggerHandler)
		middleware.WrapHandlerWithObservability("/swagger/", swaggerHandler)

		// ✅ 404 - Not Found Handler
		router.NotFoundHandler = http.HandlerFunc(handlers.NotFoundHandler)

		logStepOutcome(instance, "register-routes", "All routes registered with observability middleware, including Swagger & 404 handler")
	}
}

// ✅ Start HTTP Server
func startHTTPServer(server *http.Server, instance InstanceConfig, startTime time.Time) {
	serverErrChan := make(chan error, 1)
	go func() {
		err := server.ListenAndServe()
		serverErrChan <- err
	}()

	select {
	case err := <-serverErrChan:
		if err != nil && err != http.ErrServerClosed {
			log.Fatalf("ERROR: Service failed to start: %v", err)
		}
	case <-time.After(500 * time.Millisecond):
		startupDuration := time.Since(startTime).Seconds() * 1000
		log.Printf("INFO: Service started successfully (startup time: %.6fms)", startupDuration)
	}
}

// ✅ Handle Graceful Shutdown
func handleShutdown(instance InstanceConfig, server *http.Server) {
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	sig := <-sigChan
	log.Printf("INFO: Received shutdown signal: %v", sig)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	logStepOutcome(instance, "shutdown", "Service shutting down")

	// ✅ Fix: Use `server.Shutdown(ctx)` Instead of `http.DefaultServer`
	if err := server.Shutdown(ctx); err != nil {
		log.Printf("ERROR: Service shutdown failed: %v", err)
		os.Exit(1)
	}

	logging.FlushPendingLogs()
	log.Println("INFO: Service exited cleanly.")
}
