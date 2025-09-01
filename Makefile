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

# Database operations
migrate:
	@echo "Running database migrations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh

migrate-seed:
	@echo "Running database migrations with seed data..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh --seed

# Performance testing
test-performance:
	@echo "Running performance tests..."
	@chmod +x scripts/performance_test.sh
	@./scripts/performance_test.sh

# Phase 2 testing
test-phase2:
	@echo "Running Phase 2 tests..."
	@chmod +x scripts/phase2_test.sh
	@./scripts/phase2_test.sh

# Phase 3 testing
test-phase3:
	@echo "Running Phase 3 tests..."
	@chmod +x scripts/phase3_test.sh
	@./scripts/phase3_test.sh

# Apply performance optimizations
optimize:
	@echo "Applying performance optimizations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh
	@echo "✅ Performance optimizations applied successfully!"

# Apply Phase 2 optimizations
optimize-phase2:
	@echo "Applying Phase 2 optimizations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh
	@echo "✅ Phase 2 optimizations applied successfully!"

# Apply Phase 3 optimizations
optimize-phase3:
	@echo "Applying Phase 3 optimizations..."
	@chmod +x scripts/migrate.sh
	@./scripts/migrate.sh
	@echo "✅ Phase 3 optimizations applied successfully!"
