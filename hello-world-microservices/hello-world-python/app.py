from flask import Flask, request, jsonify
import os
import logging
import time
import random
import requests
from prometheus_client import Counter, Histogram, generate_latest
from flasgger import Swagger

app = Flask(__name__)
Swagger(app)
logging.basicConfig(level=logging.INFO, format='[%(levelname)s] %(asctime)s | %(message)s')

# Define Prometheus metrics
http_requests_total = Counter('http_requests_total', 'Total number of HTTP requests received')
http_request_duration = Histogram('http_request_duration_seconds', 'Histogram of response time for requests')
http_errors_total = Counter('http_errors_total', 'Total number of HTTP errors encountered')
http_outbound_requests_total = Counter('http_outbound_requests_total', 'Total number of outbound HTTP requests made')
http_outbound_request_duration = Histogram('http_outbound_request_duration_seconds', 'Histogram of response time for outbound requests')

# Function to retrieve cloud metadata from environment variables or defaults
def get_cloud_metadata():
    provider = os.getenv("CLOUD_PROVIDER", "docker")
    region = os.getenv("CLOUD_REGION", "local")
    zone = os.getenv("CLOUD_ZONE", "local")
    return provider, region, zone

# Get runtime language from environment variables
RUNTIME_LANGUAGE = os.getenv("RUNTIME_LANGUAGE", "Python")

@app.route('/hello', methods=['GET'])
def hello():
    """Returns a Hello World response"""
    start_time = time.time()
    http_requests_total.inc()
    cloud_provider, cloud_region, cloud_zone = get_cloud_metadata()
    log_message = f"👋 Received request at /hello | Language: {RUNTIME_LANGUAGE}, Provider: {cloud_provider}, Region: {cloud_region}, Zone: {cloud_zone}"
    app.logger.info(log_message)
    
    response = {
        "message": "Hello, World!",
        "language": RUNTIME_LANGUAGE,
        "provider": cloud_provider,
        "region": cloud_region,
        "zone": cloud_zone
    }
    http_request_duration.observe(time.time() - start_time)
    return jsonify(response)

@app.route('/metrics')
def metrics():
    return generate_latest()

if __name__ == '__main__':
    # Get host and port from environment variables
    PORT = int(os.getenv("SERVICE_PORT", 8081))
    HOST = os.getenv("SERVICE_HOST", "0.0.0.0")
    
    # Retrieve cloud metadata for logging
    cloud_provider, cloud_region, cloud_zone = get_cloud_metadata()
    app.logger.info(f"🚀 Starting Hello World {RUNTIME_LANGUAGE} service on {HOST}:{PORT} | Provider: {cloud_provider}, Region: {cloud_region}, Zone: {cloud_zone}")
    
    # Start the Flask application
    app.run(host=HOST, port=PORT)
