# Use official Golang base image
FROM golang:1.22-alpine

# Set the working directory inside the container
WORKDIR /app

# Copy go.mod and go.sum first (for dependency caching)
COPY go.mod go.sum ./

# Ensure build dependencies are installed
RUN apk add --no-cache git

# Download dependencies before copying source code
RUN go mod download

# Copy the application source code
COPY . .

# Build the application
RUN go build -o hello-world .

# Set default environment variables
ENV SERVICE_HOST="0.0.0.0"
ENV SERVICE_PORT="8080"
ENV CLOUD_PROVIDER="docker"
ENV CLOUD_REGION="local"
ENV CLOUD_ZONE="local"
ENV RUNTIME_LANGUAGE="Golang"

# Expose the port dynamically
EXPOSE ${SERVICE_PORT}

# Command to run the application
CMD ["/app/hello-world"]
