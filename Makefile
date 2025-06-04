.PHONY: build build-all clean test test-coverage test-race test-bench help version bump-major bump-minor bump-patch

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS := -ldflags "-X gaia-mcp-go/version.Version=$(VERSION) -X gaia-mcp-go/version.GitCommit=$(COMMIT) -X gaia-mcp-go/version.BuildDate=$(DATE)"

# Test configuration
COVERAGE_OUT := coverage.out
COVERAGE_HTML := coverage.html

# Default target
help: ## Show this help message
	@echo 'Usage: make [target]'
	@echo ''
	@echo 'Available targets:'
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "  %-15s %s\n", $$1, $$2}' $(MAKEFILE_LIST)

version: ## Show current version information
	@echo "Current version: $(VERSION)"
	@echo "Git commit: $(COMMIT)"
	@echo "Build date: $(DATE)"

build: ## Build for current platform
	go build $(LDFLAGS) -o bin/gaia-mcp-go ./main.go

build-all: ## Build for all platforms
	@echo "Building for multiple platforms..."
	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o bin/gaia-mcp-go-darwin-amd64 ./main.go
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o bin/gaia-mcp-go-darwin-arm64 ./main.go
	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o bin/gaia-mcp-go-windows-amd64.exe ./main.go
	GOOS=windows GOARCH=arm64 go build $(LDFLAGS) -o bin/gaia-mcp-go-windows-arm64.exe ./main.go
	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o bin/gaia-mcp-go-linux-amd64 ./main.go
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o bin/gaia-mcp-go-linux-arm64 ./main.go

# Testing targets
test: ## Run tests
	go test -v ./...

test-short: ## Run tests with short flag (skip long-running tests)
	go test -short -v ./...

test-race: ## Run tests with race detection
	go test -race -v ./...

test-coverage: ## Run tests with coverage
	go test -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -func=$(COVERAGE_OUT)

test-coverage-html: test-coverage ## Generate HTML coverage report
	go tool cover -html=$(COVERAGE_OUT) -o $(COVERAGE_HTML)
	@echo "Coverage report generated: $(COVERAGE_HTML)"

test-bench: ## Run benchmarks
	go test -bench=. -benchmem ./...

test-verbose: ## Run tests with verbose output
	go test -v -count=1 ./...

test-clean: ## Clean test cache and run tests fresh
	go clean -testcache
	go test -v ./...

# Test specific packages
test-shared: ## Test shared package only
	go test -v ./pkg/shared

test-imageutil: ## Test imageutil package only
	go test -v ./pkg/imageutil

test-version: ## Test version package only
	go test -v ./version

test-api: ## Test API package only
	go test -v ./internal/api

test-tools: ## Test tools package only
	go test -v ./internal/tools

# Advanced testing
test-integration: ## Run integration tests (when available)
	go test -tags=integration -v ./...

test-stress: ## Run stress tests (multiple iterations)
	go test -count=100 -v ./pkg/shared
	go test -count=50 -v ./pkg/imageutil

test-memory: ## Run tests with memory profiling
	go test -memprofile=mem.prof -v ./...
	@echo "Memory profile generated: mem.prof"
	@echo "View with: go tool pprof mem.prof"

test-cpu: ## Run tests with CPU profiling
	go test -cpuprofile=cpu.prof -v ./...
	@echo "CPU profile generated: cpu.prof"
	@echo "View with: go tool pprof cpu.prof"

# Quality checks
lint: ## Run linter (requires golangci-lint)
	@which golangci-lint > /dev/null || (echo "Please install golangci-lint: https://golangci-lint.run/usage/install/" && exit 1)
	golangci-lint run

fmt: ## Format code
	go fmt ./...

vet: ## Run go vet
	go vet ./...

mod-tidy: ## Tidy up go.mod
	go mod tidy

mod-verify: ## Verify dependencies
	go mod verify

# Quality gate (run before committing)
check: fmt vet lint test ## Run all quality checks

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf dist/
	rm -f $(COVERAGE_OUT)
	rm -f $(COVERAGE_HTML)
	rm -f *.prof

install: ## Install binary to GOPATH/bin
	go install $(LDFLAGS) .

# Version bumping shortcuts
bump-patch: ## Bump patch version (x.x.X)
	./scripts/bump-version.sh patch

bump-minor: ## Bump minor version (x.X.0)
	./scripts/bump-version.sh minor

bump-major: ## Bump major version (X.0.0)
	./scripts/bump-version.sh major

bump-dry: ## Dry run version bump (patch)
	./scripts/bump-version.sh patch dry_run

# Development helpers
dev-setup: ## Set up development environment
	@echo "Setting up development environment..."
	go mod tidy
	@echo "Installing development tools..."
	@which golangci-lint > /dev/null || (echo "Installing golangci-lint..." && curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin)
	@echo "Development environment ready!"

watch-test: ## Watch for changes and run tests (requires entr)
	@which entr > /dev/null || (echo "Please install entr for file watching" && exit 1)
	find . -name "*.go" | entr -c make test

# Documentation
docs-serve: ## Serve documentation locally
	@echo "Documentation available at:"
	@echo "  Testing Guide: file://$(PWD)/docs/testing-guide.md"
	@echo "  Semantic Versioning: file://$(PWD)/docs/semantic-versioning.md"

# CI/CD helpers
ci-test: ## Run tests suitable for CI environment
	go test -race -coverprofile=$(COVERAGE_OUT) ./...
	go tool cover -func=$(COVERAGE_OUT)

ci-build: ## Build for CI environment
	go build $(LDFLAGS) ./...

# Show test statistics
test-stats: test-coverage ## Show test coverage statistics
	@echo ""
	@echo "📊 Test Coverage Summary:"
	@echo "========================="
	@go tool cover -func=$(COVERAGE_OUT) | grep total | awk '{print "Total Coverage: " $$3}'
	@echo ""
	@echo "📁 Package Coverage:"
	@echo "===================="
	@go tool cover -func=$(COVERAGE_OUT) | grep -v total | head -10
	@echo ""
	@echo "🎯 Coverage Goals:"
	@echo "=================="
	@echo "  Critical packages (api, tools): 90%+"
	@echo "  Utility packages (shared, imageutil): 80%+"
	@echo "  Configuration packages: 70%+"