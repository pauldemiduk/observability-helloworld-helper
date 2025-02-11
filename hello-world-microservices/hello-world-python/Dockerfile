# Use official Python image
FROM python:3.9-slim

# Set the working directory inside the container
WORKDIR /app

# Copy application code
COPY . .

# Install dependencies
RUN pip install --no-cache-dir flask requests flasgger prometheus_client

# Set default environment variables
ENV SERVICE_HOST="0.0.0.0"
ENV SERVICE_PORT="8081"
ENV CLOUD_PROVIDER="docker"
ENV CLOUD_REGION="local"
ENV CLOUD_ZONE="local"
ENV RUNTIME_LANGUAGE="Python"

# Expose the port dynamically
EXPOSE ${SERVICE_PORT}

# Command to run the application
CMD ["python3", "app.py"]
