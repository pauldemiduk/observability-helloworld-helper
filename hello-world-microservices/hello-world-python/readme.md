# Python Hello World Microservice

This is the **Python** implementation of the Hello World microservice. It is built using **FastAPI** and supports **OpenTelemetry** for logging and metrics.

## Features
- Built with **Python 3.x** and **FastAPI**
- Uses **OpenTelemetry for structured logging**
- Exposes **Prometheus-compatible metrics** on `/metrics`
- Cloud deployment-ready for **AWS, GCP, and Azure**

## Getting Started

### Prerequisites
- **Python 3.x**
- **Docker**
- **Kubernetes (minikube/kind or cloud provider)**
- **Prometheus & Grafana** (for metrics visualization)

### Running Locally

```sh
# Install dependencies
pip install -r requirements.txt

# Run the service
python main.py
```

Or run using Docker:

```sh
docker build -t hello-world-python .
docker run -p 8080:8080 hello-world-python
```

### Observability

#### Logs
This service outputs structured logs using OpenTelemetry:
```json
{
  "timestamp": "2025-02-11T12:00:00Z",
  "service": "hello-world-python",
  "level": "info",
  "message": "Service started",
  "cloud": "gcp",
  "region": "us-west1"
}
```

#### Metrics
Exposed on `/metrics` in **Prometheus format**:
```txt
# HELP service_up 1 if the service is running
# TYPE service_up gauge
service_up{service="hello-world-python"} 1
```

