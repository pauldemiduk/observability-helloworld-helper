# Hello World Microservices

This repository contains sample **Hello World** microservices implemented in **Golang, Python, and C#**. These microservices are designed to run in **Kubernetes** across **AWS, GCP, and Azure**, with built-in **observability** using **OpenTelemetry**.

## Features
- **Multi-cloud deployment support** (AWS, GCP, Azure)
- **Logging and Metrics** with OpenTelemetry (OTEL)
- **Prometheus-compatible metrics exposure** on `/metrics` endpoint
- **Industry-standard structured logging**
- **Minimal yet scalable architecture**

## Microservices Overview

| Language | Framework/Runtime | Logging | Metrics |
|----------|------------------|---------|---------|
| Golang   | Go 1.x           | OTEL    | Prometheus |
| Python   | FastAPI          | OTEL    | Prometheus |
| C#       | .NET 7+          | OTEL    | Prometheus |

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

### Kubernetes Deployment
To deploy to Kubernetes, use Helm:

```sh
# Install dependencies
helm repo add otel https://open-telemetry.github.io/opentelemetry-helm-charts
helm repo update

# Deploy microservice
kubectl apply -f k8s/
```

## Observability

### Logs
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

### Metrics
Metrics are exposed on the `/metrics` endpoint in **Prometheus format**:
```txt
# HELP service_up 1 if the service is running
# TYPE service_up gauge
service_up{service="hello-world-golang"} 1
```

## Contributing
1. Fork the repository
2. Create a feature branch
3. Open a pull request

## License
[MIT License](LICENSE)

