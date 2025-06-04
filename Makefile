.PHONY: build build-all clean test help version bump-major bump-minor bump-patch

# Version information
VERSION ?= $(shell git describe --tags --always --dirty)
COMMIT := $(shell git rev-parse --short HEAD)
DATE := $(shell date -u '+%Y-%m-%d_%H:%M:%S')

# Build flags
LDFLAGS := -ldflags "-X gaia-mcp-go/version.Version=$(VERSION) -X gaia-mcp-go/version.GitCommit=$(COMMIT) -X gaia-mcp-go/version.BuildDate=$(DATE)"

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

test: ## Run tests
	go test -v ./...

clean: ## Clean build artifacts
	rm -rf bin/
	rm -rf dist/

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