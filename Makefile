# ABT Analytics Dashboard Backend Makefile

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
BINARY_NAME=abt-analytics-dashboard
BINARY_UNIX=$(BINARY_NAME)_unix

# Build the application
.PHONY: build
build:
	$(GOBUILD) -o $(BINARY_NAME) -v ./...

# Build for Linux
.PHONY: build-linux
build-linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) -v

# Run the application
.PHONY: run
run:
	$(GOCMD) run main.go

# Run with sample data
.PHONY: run-sample
run-sample:
	$(GOCMD) run main.go

# Run with custom dataset
.PHONY: run-data
run-data:
	DATA_FILE_PATH=$(DATA_FILE) $(GOCMD) run main.go

# Clean build artifacts
.PHONY: clean
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_UNIX)

# Run tests
.PHONY: test
test:
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage:
	$(GOTEST) -v -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

# Run tests with race condition detection
.PHONY: test-race
test-race:
	$(GOTEST) -v -race ./...

# Run benchmarks
.PHONY: bench
bench:
	$(GOTEST) -v -bench=. -benchmem ./...

# Format code
.PHONY: fmt
fmt:
	$(GOCMD) fmt ./...

# Run linter
.PHONY: lint
lint:
	golangci-lint run

# Vet code
.PHONY: vet
vet:
	$(GOCMD) vet ./...

# Download dependencies
.PHONY: deps
deps:
	$(GOMOD) download

# Tidy dependencies
.PHONY: tidy
tidy:
	$(GOMOD) tidy

# Install dependencies for development
.PHONY: install-deps
install-deps:
	$(GOGET) github.com/gorilla/mux@v1.8.1
	$(GOGET) github.com/gorilla/cors@v1.10.1
	$(GOGET) github.com/stretchr/testify@v1.8.4

# Generate test data
.PHONY: generate-test-data
generate-test-data:
	$(GOCMD) run scripts/generate_test_data.go

# Check for security vulnerabilities
.PHONY: security
security:
	$(GOCMD) list -json -m all | nancy sleuth

# Run all checks
.PHONY: check
check: fmt vet lint test-coverage

# Build Docker image
.PHONY: docker-build
docker-build:
	docker build -t abt-analytics-dashboard .

# Run Docker container
.PHONY: docker-run
docker-run:
	docker run -p 8080:8080 abt-analytics-dashboard

# Display help
.PHONY: help
help:
	@echo "Available targets:"
	@echo "  build          - Build the application"
	@echo "  build-linux    - Build for Linux"
	@echo "  run            - Run the application with sample data"
	@echo "  run-data       - Run with custom dataset (set DATA_FILE variable)"
	@echo "  clean          - Clean build artifacts"
	@echo "  test           - Run tests"
	@echo "  test-coverage  - Run tests with coverage report"
	@echo "  test-race      - Run tests with race condition detection"
	@echo "  bench          - Run benchmarks"
	@echo "  fmt            - Format code"
	@echo "  lint           - Run linter"
	@echo "  vet            - Vet code"
	@echo "  deps           - Download dependencies"
	@echo "  tidy           - Tidy dependencies"
	@echo "  check          - Run all checks (fmt, vet, lint, test-coverage)"
	@echo "  docker-build   - Build Docker image"
	@echo "  docker-run     - Run Docker container"
	@echo "  help           - Display this help"

.DEFAULT_GOAL := help
