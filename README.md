# Hello World Microservices - Observability Helper

## Overview
The **Hello World Microservices** project is designed to provide **observability-ready** service examples in multiple programming languages (**Golang, Python, C#**). These microservices are built to run in **Kubernetes** across **AWS, GCP, and Azure**, with **OpenTelemetry (OTEL)** for logs and **Prometheus** for metrics.

### **Project Scope**
This repository is part of the **Observability-Helloworld-Helper** solution, providing:
- **Multi-cloud deployment examples** (AWS, GCP, Azure)
- **Structured Logging** with OpenTelemetry (OTEL)
- **Prometheus-compatible metrics** on `/metrics` endpoint
- **Minimal yet scalable architecture**
- **Docker & Kubernetes deployment options**

## **Microservices Overview**

| Language | Framework/Runtime | Logging | Metrics | Status |
|----------|------------------|---------|---------|--------|
| Golang   | Go 1.x           | OTEL    | Prometheus | ✅ Advanced |
| Python   | FastAPI          | OTEL    | Prometheus | 🚧 Planned |
| C#       | .NET 7+          | OTEL    | Prometheus | 🚧 Planned |

The **Golang service (`hello-world-go`) is the most advanced**, with full OpenTelemetry support. The Python and C# versions are planned to reach the same feature set.

## Getting Started

### **Prerequisites**
Ensure you have the following installed:
- **Docker** (for local containerized execution)
- **Kubernetes** (minikube/kind or a cloud provider)
- **kubectl** (for Kubernetes deployments)
- **Helm** (for deploying dependencies like OpenTelemetry)
- **Prometheus & Grafana** (for observability monitoring)

### **Running Locally**
Each microservice has a `Dockerfile` and `docker-compose.yml`. To run locally:

```sh
# Clone the repository
git clone https://github.com/your-repo/hello-world-microservices.git
cd hello-world-microservices

# Start all services using Docker Compose
docker-compose up --build
```

## **Kubernetes Deployment**
### **Using Helm**
```sh
# Install dependencies (OpenTelemetry Helm charts)
helm repo add otel https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

# Deploy microservice
tkubectl apply -f k8s/
```

## **Observability: Logs & Metrics**

### **Structured Logging (JSON Format)**
Each service outputs structured logs using OpenTelemetry:
```json
{
  "timestamp": "2025-02-11T12:00:00Z",
  "service": "hello-world-golang",
  "level": "info",
  "message": "Service started",
  "cloud": "aws",
  "region": "us-east-1"
}
```

### **Metrics (Prometheus Format)**
Metrics are exposed on the `/metrics` endpoint:
```txt
# HELP service_up 1 if the service is running
# TYPE service_up gauge
service_up{service="hello-world-golang"} 1
```

To scrape metrics using Prometheus:
```yaml
scrape_configs:
  - job_name: 'hello-world-microservices'
    static_configs:
      - targets: ['localhost:8080']
```

## **API Documentation**
### **Endpoints:**
| Method | Endpoint | Description |
|--------|---------|-------------|
| **GET** | `/hello` | Returns a Hello World response |
| **PUT** | `/inbound` | Processes inbound requests |
| **PUT** | `/outbound` | Handles outbound service calls |
| **GET** | `/status` | Returns service health check response |
| **GET** | `/metrics` | Exposes Prometheus metrics |

Example request:
```sh
curl -X GET http://localhost:8080/hello
```

## **Error Handling & Resilience**
- **Panic Recovery:** Middleware automatically recovers from crashes.
- **Graceful Shutdown:** Services handle SIGINT/SIGTERM signals cleanly.
- **Missing Config Handling:** Default values and warning logs ensure smooth execution.

## **Future Enhancements**
- **Complete Python & C# implementations to match Golang**
- **Improve batch processing & outbound request simulation**
- **Enhance Kubernetes Helm charts for automated deployments**

## **Contributing**
1. Fork the repository
2. Create a feature branch
3. Open a pull request

## **License**
[MIT License](LICENSE)

---

This README serves as the **entry point for all implementations**, with **future per-language READMEs** planned as services evolve. 🚀 Let me know if any refinements are needed!
