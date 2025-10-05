.PHONY: help build run test test-verbose test-coverage test-race clean dev install-deps fmt vet lint benchmark

# Variabel
APP_NAME=htmx-app
BUILD_DIR=bin
PORT=8080

# Default target
help:
	@echo "Available commands:"
	@echo "  make build          - Build the application"
	@echo "  make run            - Run the application"
	@echo "  make dev            - Run with auto-reload (requires air)"
	@echo "  make test           - Run tests"
	@echo "  make test-verbose   - Run tests with verbose output"
	@echo "  make test-coverage  - Run tests with coverage report"
	@echo "  make test-race      - Run tests with race detection"
	@echo "  make benchmark      - Run benchmarks"
	@echo "  make fmt            - Format code"
	@echo "  make vet            - Run go vet"
	@echo "  make lint           - Run golangci-lint (requires golangci-lint)"
	@echo "  make clean          - Remove build artifacts"
	@echo "  make install-deps   - Install development dependencies"
	@echo "  make all            - Format, vet, test, and build"

# Build app
build:
	@echo "Building $(APP_NAME)..."
	@mkdir -p $(BUILD_DIR)
	@go build -o $(BUILD_DIR)/$(APP_NAME) .
	@echo "Build complete: $(BUILD_DIR)/$(APP_NAME)"

# Run app
run: build
	@echo "Starting $(APP_NAME) on port $(PORT)..."
	@./$(BUILD_DIR)/$(APP_NAME)


# Run tests
test:
	@echo "Running tests..."
	@go test ./...

# Run tests with verbose output
test-verbose:
	@echo "Running tests (verbose)..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run tests with race detection
test-race:
	@echo "Running tests with race detection..."
	@go test -race -v ./...

# Run benchmarks
benchmark:
	@echo "Running benchmarks..."
	@go test -bench=. -benchmem ./...

# Format kode
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run go vet
vet:
	@echo "Running go vet..."
	@go vet ./...

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@rm -rf $(BUILD_DIR)
	@rm -f coverage.out coverage.html
	@echo "Clean complete"

# Install development dependencies
install-deps:
	@echo "Installing development dependencies..."
	@go install github.com/air-verse/air@latest
	@echo "Dependencies installed"

# Run all checks
all: fmt vet test build
	@echo "All checks passed!"
