from flask import Flask, request, jsonify
import os
import logging
import time
import random
import datetime
import uuid
import json
from prometheus_client import Counter, Histogram, generate_latest
import waitress

app = Flask(__name__)

logging.basicConfig(
    level=logging.INFO,
    format='[%(levelname)s] %(asctime)s | Service: %(message)s'
)

# Define Prometheus metrics
http_requests_total = Counter('http_requests_total', 'Total number of HTTP requests received')
http_request_duration = Histogram('http_request_duration_seconds', 'Histogram of response time for requests')

# Function to retrieve cloud metadata from environment variables
def get_cloud_metadata():
    return (
        os.getenv("CLOUD_PROVIDER", "docker"),
        os.getenv("CLOUD_REGION", "local"),
        os.getenv("CLOUD_ZONE", "local")
    )

# Retrieve additional metadata
RUNTIME_LANGUAGE = os.getenv("RUNTIME_LANGUAGE", "Python")
IP_ADDRESS = os.getenv("SERVICE_IP", "127.0.0.1")
SERVICE_PID = os.getpid()
CONTAINER_ID = os.getenv("CONTAINER_ID", "unknown")

def log_event(operation, message):
    cloud_provider, region, zone = get_cloud_metadata()
    log_data = {
        "timestamp": datetime.datetime.utcnow().isoformat() + "Z",
        "log_level": "INFO",
        "service_name": "hello-world-python",
        "service_operation": operation,
        "runtime_language": RUNTIME_LANGUAGE,
        "cloud_provider": cloud_provider,
        "region": region,
        "zone": zone,
        "ip_address": IP_ADDRESS,
        "service_pid": SERVICE_PID,
        "container_id": CONTAINER_ID,
        "message": message
    }
    app.logger.info(json.dumps(log_data))

@app.route('/hello', methods=['GET'])
def hello():
    http_requests_total.inc()
    log_event("hello", "Received request at /hello")
    response = {
        "message": "Hello, World!",
        "language": RUNTIME_LANGUAGE,
        "cloud_provider": get_cloud_metadata()[0],
        "region": get_cloud_metadata()[1],
        "zone": get_cloud_metadata()[2]
    }
    return jsonify(response)

@app.route('/inbound', methods=['POST'])
def inbound():
    http_requests_total.inc()
    log_event("inbound", "Processing inbound request")
    return jsonify({"status": "inbound request processed"}), 200

@app.route('/outbound', methods=['POST'])
def outbound():
    http_requests_total.inc()
    log_event("outbound", "Processing outbound request")
    return jsonify({"status": "outbound request processed"}), 200

@app.route('/metrics')
def metrics():
    return generate_latest()

if __name__ == '__main__':
    HOST = os.getenv("SERVICE_HOST", "0.0.0.0")
    PORT = int(os.getenv("SERVICE_PORT", 8081))
    log_event("service_start", "Service started successfully")
    waitress.serve(app, host=HOST, port=PORT)
