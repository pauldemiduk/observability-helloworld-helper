# Hello World - Golang Microservice

## Overview
The **hello-world-go** microservice is a **Golang-based** implementation within the **Hello World Microservices** project. It is designed with **observability in mind**, integrating **OpenTelemetry (OTEL)** for structured logging and **Prometheus** for metrics collection.

This microservice serves as the **reference implementation**, setting the standard for future enhancements in Python and C#.

### **Key Features:**
✅ **RESTful API with structured logging & tracing**  
✅ **OpenTelemetry integration for logs & traces**  
✅ **Prometheus-compatible metrics exposure on `/metrics`**  
✅ **Middleware-based logging & metrics collection**  
✅ **Graceful shutdown & error handling**  
✅ **Configurable via environment variables**  
✅ **Batch processing for outbound requests (Future Implementation)**  

## **Getting Started**

### **Prerequisites**
Ensure you have the following installed:
- **Go** (latest version recommended)
- **Docker** (optional, for containerized deployment)
- **Kubernetes (minikube/kind)** (for running in a cluster)
- **Prometheus & Grafana** (for monitoring metrics)

### **Running Locally**
```sh
# Clone the repository
git clone https://github.com/your-repo/hello-world-microservices.git
cd hello-world-microservices/hello-world-go

# Install dependencies
go mod tidy

# Run the service
go run main.go
```

### **API Endpoints**
| Method | Endpoint | Description |
|--------|---------|-------------|
| **GET** | `/hello` | Returns a Hello World response |
| **PUT** | `/inbound` | Processes inbound requests |
| **PUT** | `/outbound` | Handles outbound service calls |
| **GET** | `/status` | Returns service health check response |
| **GET** | `/metrics` | Exposes Prometheus metrics |
| **GET** | `/swagger` | Exposes API documentation |

## **Observability: Logs & Metrics**

### **Structured Logging (JSON Format)**
The microservice produces structured logs using OpenTelemetry:
```json
{
  "timestamp": "2025-02-16T14:00:00Z",
  "service": "hello-world-go",
  "level": "info",
  "message": "Service started",
  "cloud": "aws",
  "region": "us-east-1",
  "trace_id": "1234abcd5678efgh"
}
```

### **Metrics (Prometheus Format)**
Prometheus-compatible metrics are exposed on `/metrics`:
```txt
# HELP http_requests_total Total number of HTTP requests received
# TYPE http_requests_total counter
http_requests_total{operation="hello",status_code="200"} 5
```
To scrape metrics using Prometheus:
```yaml
scrape_configs:
  - job_name: 'hello-world-go'
    static_configs:
      - targets: ['localhost:8080']
```

## **Middleware**
The **logging and metrics middleware** automatically captures request information:
- Logs request start and completion with duration tracking
- Extracts **trace ID, span ID, and request ID**
- Captures request count and latency for Prometheus

## **Deployment Guide**

### **Running with Docker**
```sh
docker build -t hello-world-go .
docker run -p 8080:8080 -e CLOUD_PROVIDER=aws -e LOG_BACKEND=loki hello-world-go
```

### **Deploying to Kubernetes**
```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: hello-world-go
spec:
  replicas: 2
  selector:
    matchLabels:
      app: hello-world-go
  template:
    metadata:
      labels:
        app: hello-world-go
    spec:
      containers:
      - name: hello-world-go
        image: hello-world-go:latest
        ports:
        - containerPort: 8080
```
Apply deployment:
```sh
kubectl apply -f deployment.yaml
```

## **Grafana Observability Dashboards**
The service **exports metrics to Prometheus**, which can be visualized in **Grafana**.

**Example Grafana Dashboard Query (Request Latency Histogram):**
```txt
histogram_quantile(0.95, sum(rate(http_request_duration_seconds_bucket[5m])) by (le))
```

**Setting Up Grafana:**
```sh
helm repo add grafana https://grafana.github.io/helm-charts
helm repo update
helm install grafana grafana/grafana
```

- **Step 1:** Install Grafana
- **Step 2:** Add Prometheus as a data source in Grafana UI.

## **Error Handling & Resilience**
- **Panic Recovery:** Middleware catches unexpected panics.
- **Graceful Shutdown:** Handles SIGINT/SIGTERM for a clean shutdown.
- **Environment Variable Handling:** Uses defaults and logs warnings for missing variables.

## **Future Enhancements**
- Implement **batch processing** for outbound requests.
- Improve **log routing to different cloud providers**.
- Introduce **request delay simulation** for real-world traffic testing.
- Align Python and C# versions to match **Golang’s feature set**.

## **Contributing**
1. Fork the repository
2. Create a feature branch
3. Open a pull request

## **License**
[MIT License](LICENSE)

