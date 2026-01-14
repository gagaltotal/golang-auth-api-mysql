.PHONY: swagger run test build docker-up docker-down clean

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

# Run application
run:
	@echo "Running application..."
	go run cmd/api/main.go

# Run with hot reload (requires air)
dev:
	@echo "Running with hot reload..."
	air

# Run tests
test:
	@echo "Running tests..."
	go test ./tests/... -v

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	go test ./tests/... -coverprofile=coverage.out
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Build binary
build:
	@echo "Building application..."
	go build -o bin/api cmd/api/main.go

# Build for production
build-prod:
	@echo "Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/api cmd/api/main.go

# Docker commands
docker-up:
	@echo "Starting Docker containers..."
	docker-compose -f docker/docker-compose.yml up -d

docker-down:
	@echo "Stopping Docker containers..."
	docker-compose -f docker/docker-compose.yml down

docker-logs:
	@echo "Showing Docker logs..."
	docker-compose -f docker/docker-compose.yml logs -f

# Clean build artifacts
clean:
	@echo "Cleaning..."
	rm -rf bin/
	rm -f coverage.out coverage.html
	go clean

# Install development tools
install-tools:
	@echo "Installing development tools..."
	go install github.com/swaggo/swag/cmd/swag@latest
	go install github.com/cosmtrek/air@latest

# Database migrations
migrate:
	@echo "Running migrations..."
	go run cmd/api/main.go migrate

# Format code
fmt:
	@echo "Formatting code..."
	go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	golangci-lint run

# Tidy dependencies
tidy:
	@echo "Tidying dependencies..."
	go mod tidy

# Help command
help:
	@echo "Available commands:"
	@echo "  make swagger         - Generate Swagger documentation"
	@echo "  make run            - Run application"
	@echo "  make dev            - Run with hot reload"
	@echo "  make test           - Run tests"
	@echo "  make test-coverage  - Run tests with coverage"
	@echo "  make build          - Build binary"
	@echo "  make build-prod     - Build for production"
	@echo "  make docker-up      - Start Docker containers"
	@echo "  make docker-down    - Stop Docker containers"
	@echo "  make docker-logs    - Show Docker logs"
	@echo "  make clean          - Clean build artifacts"
	@echo "  make install-tools  - Install development tools"
	@echo "  make fmt            - Format code"
	@echo "  make lint           - Run linter"
	@echo "  make tidy           - Tidy dependencies"