# 8086 Simulator Makefile

# Variables
BINARY_NAME=8086-simulator
MAIN_PACKAGE=./main.go
BUILD_DIR=build
COVERAGE_FILE=coverage.out
LINT_TIMEOUT=5m

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod
GOLINT=golangci-lint

# Default target
.PHONY: all
all: clean deps test build

# Help target
.PHONY: help
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

# Clean build artifacts
.PHONY: clean
clean: ## Clean build artifacts
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)
	rm -f $(COVERAGE_FILE)
	rm -f $(BINARY_NAME)

# Download dependencies
.PHONY: deps
deps: ## Download dependencies
	$(GOMOD) download
	$(GOMOD) verify

# Tidy dependencies
.PHONY: tidy
tidy: ## Tidy dependencies
	$(GOMOD) tidy

# Run tests
.PHONY: test
test: ## Run tests
	$(GOTEST) -v ./...

# Run tests with coverage
.PHONY: test-coverage
test-coverage: ## Run tests with coverage
	$(GOTEST) -v -coverprofile=$(COVERAGE_FILE) ./...
	$(GOCMD) tool cover -html=$(COVERAGE_FILE) -o coverage.html

# Show test coverage
.PHONY: coverage
coverage: test-coverage ## Show test coverage
	$(GOCMD) tool cover -func=$(COVERAGE_FILE)

# Build the application
.PHONY: build
build: ## Build the application
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME) $(MAIN_PACKAGE)

# Build for multiple platforms
.PHONY: build-all
build-all: ## Build for multiple platforms
	@mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PACKAGE)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PACKAGE)
	GOOS=darwin GOARCH=arm64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PACKAGE)

# Run the application
.PHONY: run
run: build ## Build and run the application
	./$(BUILD_DIR)/$(BINARY_NAME)

# Run the application directly (without building)
.PHONY: run-direct
run-direct: ## Run the application directly
	$(GOCMD) run $(MAIN_PACKAGE)

# Install golangci-lint
.PHONY: install-lint
install-lint: ## Install golangci-lint
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.54.2

# Run linter
.PHONY: lint
lint: ## Run linter
	$(GOLINT) run --timeout=$(LINT_TIMEOUT)

# Run linter with fix
.PHONY: lint-fix
lint-fix: ## Run linter with auto-fix
	$(GOLINT) run --timeout=$(LINT_TIMEOUT) --fix

# Format code
.PHONY: fmt
fmt: ## Format code
	$(GOCMD) fmt ./...

# Vet code
.PHONY: vet
vet: ## Vet code
	$(GOCMD) vet ./...

# Check code quality (fmt, vet, lint)
.PHONY: check
check: fmt vet lint ## Check code quality

# Install development tools
.PHONY: install-tools
install-tools: ## Install development tools
	$(GOGET) -u github.com/golangci/golangci-lint/cmd/golangci-lint

# Benchmark tests
.PHONY: benchmark
benchmark: ## Run benchmark tests
	$(GOTEST) -bench=. -benchmem ./...

# Race condition tests
.PHONY: test-race
test-race: ## Run tests with race detection
	$(GOTEST) -race ./...

# Generate mocks (if using mockgen)
.PHONY: generate
generate: ## Generate code
	$(GOCMD) generate ./...

# Docker build
.PHONY: docker-build
docker-build: ## Build Docker image
	docker build -t $(BINARY_NAME) .

# Docker run
.PHONY: docker-run
docker-run: ## Run Docker container
	docker run --rm $(BINARY_NAME)

# CI target (for GitHub Actions)
.PHONY: ci
ci: deps test build lint ## Run CI pipeline

# Development target
.PHONY: dev
dev: deps test build ## Development setup
