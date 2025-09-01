.PHONY: help build run test clean docker-build docker-up docker-down dev

# Default target
help:
	@echo "Available commands:"
	@echo "  build        - Build the Go application"
	@echo "  run          - Run the application"
	@echo "  test         - Run tests"
	@echo "  clean        - Clean build artifacts"
	@echo "  docker-build - Build Docker image"
	@echo "  docker-up    - Start all services with Docker Compose"
	@echo "  docker-down  - Stop all services"
	@echo "  dev          - Run with hot reload (requires air)"

# Build the application
build:
	go build -o bin/ethosview-backend ./cmd/server

# Run the application
run:
	go run ./cmd/server

# Run tests
test:
	go test ./...

# Clean build artifacts
clean:
	rm -rf bin/
	go clean

# Build Docker image
docker-build:
	docker build -t ethosview-backend .

# Start all services with Docker Compose
docker-up:
	docker-compose up -d

# Stop all services
docker-down:
	docker-compose down

# Run with hot reload (requires air: go install github.com/cosmtrek/air@latest)
dev:
	@if ! command -v air &> /dev/null; then \
		echo "Air not found. Installing..."; \
		go install github.com/cosmtrek/air@latest; \
	fi
	air

# Install dependencies
deps:
	go mod download
	go mod tidy

# Generate go.sum
sum:
	go mod tidy
	go mod verify
