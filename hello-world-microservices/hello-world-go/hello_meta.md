# HelloInstance - Service Hierarchy & Configuration

## 🌐 HelloInstance - Running Service Instance

---

### 🔹 ServiceConfig - Environment or Default

#### ✅ Service Metadata
- **ServiceName:** `util.GetEnvOrDefault("SERVICE_NAME", "hello-world-go")`
- **ServiceVersion:** `"unknown"`
- **ServiceHostname:** `util.GetEnvOrDefault("SERVICE_HOST", "0.0.0.0")`
- **ServiceIP:** `util.GetEnvOrDefault("SERVICE_IP", "127.0.0.1")`
- **ServicePort:** `util.GetEnvIntValidated("SERVICE_PORT", 8080, 1024, 49151)`

#### ✅ Application Environment
- **Environment:** `util.GetEnvOrDefault("ENV", "dev")`

#### ✅ Logging Configuration
- **LogLevel:** `util.GetEnvOrDefault("LOG_LEVEL", "info")`
- **LogBackend:** `util.GetEnvOrDefault("LOG_BACKEND", "stdout")`
- **LogFile:** `util.GetEnvOrDefault("LOG_FILE", "")`
- **LogSampleRate:** `util.GetEnvFloatValidated("LOG_SAMPLE_RATE", 0.3, 0.0, 1.0)`

#### ✅ Observability & Cloud Config
- **LokiEndpoint:** `util.GetEnvOrDefault("LOKI_ENDPOINT", "")`
- **CloudWatchEndpoint:** `util.GetEnvOrDefault("CLOUDWATCH_ENDPOINT", "")`
- **GCPLoggingEndpoint:** `util.GetEnvOrDefault("GCP_LOGGING_ENDPOINT", "")`
- **AzureMonitorEndpoint:** `util.GetEnvOrDefault("AZURE_MONITOR_ENDPOINT", "")`

#### ✅ Response Time Buckets
- **ResponseTimeFast:** `util.GetEnvFloatValidated("RESPONSE_TIME_FAST", 200.0, 1, 5000)`
- **ResponseTimeModerate:** `util.GetEnvFloatValidated("RESPONSE_TIME_MODERATE", 500.0, 1, 5000)`
- **ResponseTimeSlow:** `util.GetEnvFloatValidated("RESPONSE_TIME_SLOW", 1000.0, 1, 10000)`

---

### 🔹 ServiceContext - Runtime Context

#### ✅ RequestMetadata (HTTP Lifecycle)
- **Operation:** `"operation,omitempty"` (e.g., "hello", "inbound", "outbound")
- **RequestID:** `"request_id,omitempty"` (Unique client request identifier)
- **RequestURL:** `"request_url,omitempty"` (Full request URL)

#### ✅ ResponseMetadata (HTTP Response)
- **HTTPStatusCode:** `"http_status_code,omitempty"`
- **ResponseOutcome:** `"response_outcome,omitempty"`
- **ResponseTimeMS:** `"response_time_ms,omitempty"`

---

### 🔹 LogContext - Logging Metadata

#### ✅ Service Metadata
- **ServiceName, ServiceVersion, ServiceHostname**
- **ServiceIP, ServicePort, InstanceID**

#### ✅ Cloud & Infrastructure Metadata
- **CloudProvider, CloudRegion, AvailabilityZone**
- **K8sNamespace, K8sPodName, K8sContainerID**

#### ✅ Logging Options
- **LogLevel, LogBackend, LogSampleRate**

---

### 🔹 LogEntryOptions - Dynamic Log Entry Configuration

- **LogType:** `"request", "response", "execution", "error"`
- **RequestStage:** `"received", "processed"`
- **LogSeverity:** `"info", "warn", "error", "critical"`
- **QueueTimeMS:** `Captured queue wait time before processing`

---

## 🚀 Future Enhancements
- **Dynamic config reload (`/reload` route)**
- **Queueing improvements (store-and-forward, retries)**
- **Enhanced event routing (prioritization, filtering)**
- **Distributed tracing across multi-cloud environments**
