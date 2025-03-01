# Hello World Microservices

This repository contains sample **Hello World** microservices implemented in **Golang, Python, and C#**. These microservices are designed to run in **Kubernetes** across **AWS, GCP, and Azure**, with built-in **observability** using **OpenTelemetry (OTEL)**.

## Features
- **Multi-cloud deployment support** (AWS, GCP, Azure)
- **Logging and Metrics** with OpenTelemetry (OTEL)
- **Prometheus-compatible metrics exposure** on `/metrics` endpoint
- **Industry-standard structured logging** (Loki, CloudWatch, GCP, Azure)
- **Minimal yet scalable architecture**
- **Built-in Swagger API documentation**

## Microservices Overview

| Language | Framework/Runtime | Logging | Metrics | Status |
|----------|------------------|---------|---------|--------|
| Golang   | Go 1.x           | OTEL    | Prometheus | ✅ Fully Implemented |
| Python   | FastAPI          | OTEL    | Prometheus | 🛠️ In Progress |
| C#       | .NET 7+          | OTEL    | Prometheus | 🚀 Planned |

## Getting Started

### Prerequisites
Ensure you have the following installed:
- **Docker**
- **Kubernetes (minikube/kind or a cloud provider)**
- **kubectl**
- **Helm** (for deployment)
- **Prometheus & Grafana** (for observability)

### Running Locally

Each microservice has a `Dockerfile` and `docker-compose.yml`. To run locally:

```sh
# Clone the repository
git clone https://github.com/your-repo/hello-world-microservices.git
cd hello-world-microservices

# Build and run (example for Golang)
docker-compose up --build
```

## API Standardization

All microservices follow a **standardized API structure** for consistency.

| Method | Endpoint | Description |
|--------|---------|-------------|
| **GET** | `/hello` | Returns a Hello World response |
| **PUT** | `/inbound` | Processes inbound requests |
| **PUT** | `/outbound` | Handles outbound service calls |
| **GET** | `/status` | Returns service health check response |
| **GET** | `/metrics` | Exposes Prometheus metrics |
| **GET** | `/swagger` | Exposes API documentation |

## Observability

### Logs
Each service outputs structured logs using OpenTelemetry:
```json
{
  "timestamp": "2025-02-16T14:00:00Z",
  "service": "hello-world-golang",
  "level": "info",
  "message": "Service started",
  "cloud": "aws",
  "region": "us-east-1",
  "trace_id": "1234abcd5678efgh"
}
```

### Metrics
Metrics are exposed on the `/metrics` endpoint in **Prometheus format**:
```txt
# HELP service_up 1 if the service is running
# TYPE service_up gauge
service_up{service="hello-world-golang"} 1
```

## Observability Dashboards (Grafana)

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

## Deployment Guide

### Running in Docker
```sh
docker build -t hello-world-go .
docker run -p 8080:8080 -e CLOUD_PROVIDER=aws -e LOG_BACKEND=loki hello-world-go
```

### Deploying to Kubernetes
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

## Multi-Cloud Observability
These microservices are **cloud-agnostic** and can integrate with:
- **AWS CloudWatch**
- **GCP Cloud Logging**
- **Azure Monitor**
- **Grafana Loki** (for log aggregation)

### Pushing Logs to Cloud Providers
Logs can be forwarded to **various cloud observability services** via environment variables:
```sh
export LOG_BACKEND=cloudwatch  # Options: stdout, loki, cloudwatch, gcp, azure
```

## Debugging & Troubleshooting
Common issues when running locally or in Kubernetes:
```sh
# Check logs for errors
kubectl logs -l app=hello-world-go

# Verify service is running
kubectl get pods -l app=hello-world-go

# Test API responses
curl -X GET http://localhost:8080/hello
```

## Contributing
1. Fork the repository
2. Create a feature branch
3. Open a pull request

## License
[MIT License](LICENSE)

