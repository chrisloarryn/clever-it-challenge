# Beer API Makefile

.PHONY: build run test clean docker-build docker-run help

# Variables
APP_NAME=beer-api
DOCKER_IMAGE=beer-api:latest
GO_VERSION=1.17

# Default target
help: ## Show this help message
	@echo 'Usage: make <target>'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Development
build: ## Build the application
	@echo "Building $(APP_NAME)..."
	go build -o $(APP_NAME) cmd/main.go

run: ## Run the application
	@echo "Running $(APP_NAME)..."
	go run cmd/main.go

dev: ## Run in development mode with in-memory database
	@echo "Running in development mode..."
	export DB_TYPE=inmemory && \
	export LOG_LEVEL=debug && \
	export ENVIRONMENT=development && \
	go run cmd/main.go

# Testing
test: ## Run all tests
	@echo "Running tests..."
	go test ./...

test-verbose: ## Run tests with verbose output
	@echo "Running tests (verbose)..."
	go test -v ./...

test-coverage: ## Run tests with coverage
	@echo "Running tests with coverage..."
	go test -cover ./...

test-coverage-html: ## Generate HTML coverage report
	@echo "Generating coverage report..."
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Code quality
lint: ## Run linter
	@echo "Running linter..."
	golangci-lint run

fmt: ## Format code
	@echo "Formatting code..."
	go fmt ./...

vet: ## Run go vet
	@echo "Running go vet..."
	go vet ./...

# Dependencies
deps: ## Download dependencies
	@echo "Downloading dependencies..."
	go mod download

tidy: ## Tidy up dependencies
	@echo "Tidying dependencies..."
	go mod tidy

# Docker
docker-build: ## Build Docker image
	@echo "Building Docker image..."
	docker build -t $(DOCKER_IMAGE) .

docker-run: ## Run application in Docker
	@echo "Running Docker container..."
	docker run -p 8080:8080 --rm $(DOCKER_IMAGE)

docker-run-postgres: ## Run with PostgreSQL using docker-compose
	@echo "Starting with PostgreSQL..."
	docker-compose up --build

# Database
db-up: ## Start PostgreSQL database
	@echo "Starting PostgreSQL..."
	docker-compose up postgres -d

db-down: ## Stop PostgreSQL database
	@echo "Stopping PostgreSQL..."
	docker-compose down

# Clean
clean: ## Clean build artifacts
	@echo "Cleaning..."
	rm -f $(APP_NAME)
	rm -f coverage.out coverage.html
	docker system prune -f

# Production builds
build-linux: ## Build for Linux
	@echo "Building for Linux..."
	GOOS=linux GOARCH=amd64 go build -o $(APP_NAME)-linux cmd/main.go

build-windows: ## Build for Windows
	@echo "Building for Windows..."
	GOOS=windows GOARCH=amd64 go build -o $(APP_NAME)-windows.exe cmd/main.go

build-mac: ## Build for macOS
	@echo "Building for macOS..."
	GOOS=darwin GOARCH=amd64 go build -o $(APP_NAME)-mac cmd/main.go

build-all: build-linux build-windows build-mac ## Build for all platforms

# API testing
api-test: ## Test API endpoints (requires running server)
	@echo "Testing API endpoints..."
	@echo "Health check:"
	curl -s http://localhost:8080/ping | jq .
	@echo "\nGet all beers:"
	curl -s http://localhost:8080/api/v1/beers | jq .
	@echo "\nCreate a test beer:"
	curl -s -X POST http://localhost:8080/api/v1/beers \
		-H "Content-Type: application/json" \
		-d '{"id":999,"name":"Test Beer","brewery":"Test Brewery","country":"Chile","price":1500,"currency":"CLP"}' | jq .

api-test-full: ## Run comprehensive API tests (requires running server)
	@echo "Running comprehensive API tests..."
	@./scripts/test-api.sh

api-test-manual: ## Show manual testing commands
	@echo "Manual API testing commands:"
	@echo "See docs/API_TESTING.md for complete cURL examples"
	@echo "Or use test.REST file in VS Code with REST Client extension"

api-load-test: ## Create multiple test beers for load testing
	@echo "Creating multiple test beers..."
	@for i in $$(seq 2000 2010); do \
		echo "Creating beer $$i"; \
		curl -s -X POST http://localhost:8080/api/v1/beers \
			-H "Content-Type: application/json" \
			-d "{\"id\":$$i,\"name\":\"Beer $$i\",\"brewery\":\"Brewery $$i\",\"country\":\"Chile\",\"price\":$$(($$RANDOM % 2000 + 500)),\"currency\":\"CLP\"}" > /dev/null; \
	done
	@echo "Load test complete. Created 11 beers."

# Environment setup
setup: deps ## Setup development environment
	@echo "Setting up development environment..."
	@echo "Installing development tools..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Setup complete!"

# Help target should be first for default
.DEFAULT_GOAL := help
