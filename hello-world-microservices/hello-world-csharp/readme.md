# C# Hello World Microservice

This is the **C#** implementation of the Hello World microservice. It is built using **.NET 7+** and supports **OpenTelemetry** for logging and metrics.

## Features
- Built with **.NET 7+**
- Uses **OpenTelemetry for structured logging**
- Exposes **Prometheus-compatible metrics** on `/metrics`
- Cloud deployment-ready for **AWS, GCP, and Azure**

## Getting Started

### Prerequisites
- **.NET SDK 7+**
- **Docker**
- **Kubernetes (minikube/kind or cloud provider)**
- **Prometheus & Grafana** (for metrics visualization)

### Running Locally

```sh
# Build and run locally
dotnet build
dotnet run
```

Or run using Docker:

```sh
docker build -t hello-world-csharp .
docker run -p 8080:8080 hello-world-csharp
```

### Observability

#### Logs
This service outputs structured logs using OpenTelemetry:
```json
{
  "timestamp": "2025-02-11T12:00:00Z",
  "service": "hello-world-csharp",
  "level": "info",
  "message": "Service started",
  "cloud": "azure",
  "region": "east-us"
}
```

#### Metrics
Exposed on `/metrics` in **Prometheus format**:
```txt
# HELP service_up 1 if the service is running
# TYPE service_up gauge
service_up{service="hello-world-csharp"} 1
```