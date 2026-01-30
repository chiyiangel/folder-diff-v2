# folder-diff Makefile
# Cross-platform build automation

# Metadata
PROJECT_NAME := folder-diff
VERSION ?= $(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
BUILD_TIME := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
GO_VERSION := $(shell go version | awk '{print $$3}')

# Directories
BUILD_DIR := dist
CMD_DIR := ./cmd/$(PROJECT_NAME)
COVERAGE_DIR := coverage

# Build flags
LDFLAGS := -ldflags="-s -w \
	-X 'main.Version=$(VERSION)' \
	-X 'main.Commit=$(COMMIT)' \
	-X 'main.BuildTime=$(BUILD_TIME)' \
	-X 'main.GoVersion=$(GO_VERSION)'"

# Go commands
GO := go
GOFMT := gofmt
GOLINT := golangci-lint
GOTEST := $(GO) test
GOBUILD := $(GO) build $(LDFLAGS)

# Detect OS for platform-specific commands
ifeq ($(OS),Windows_NT)
    SHELL := powershell.exe
    RM := Remove-Item -Recurse -Force
    MKDIR := New-Item -ItemType Directory -Force
    EXE_EXT := .exe
else
    RM := rm -rf
    MKDIR := mkdir -p
    EXE_EXT :=
endif

# Build targets for different platforms
PLATFORMS := \
	linux/amd64 \
	linux/arm64 \
	linux/386 \
	linux/arm \
	darwin/amd64 \
	darwin/arm64 \
	windows/amd64 \
	windows/386

# Default target
.DEFAULT_GOAL := help

###################
# Primary Targets
###################

## help: Show this help message
.PHONY: help
help:
	@echo "folder-diff - Cross-platform Makefile"
	@echo ""
	@echo "Usage: make [target]"
	@echo ""
	@echo "Build Targets:"
	@echo "  build           - Build for current platform"
	@echo "  build-all       - Build for all platforms"
	@echo "  build-linux     - Build for Linux (amd64, arm64, 386, arm)"
	@echo "  build-darwin    - Build for macOS (amd64, arm64)"
	@echo "  build-windows   - Build for Windows (amd64, 386)"
	@echo ""
	@echo "Development Targets:"
	@echo "  dev             - Build and run locally"
	@echo "  run             - Run without building"
	@echo "  test            - Run tests"
	@echo "  test-coverage   - Run tests with coverage report"
	@echo "  lint            - Run linters"
	@echo "  fmt             - Format code"
	@echo ""
	@echo "Maintenance:"
	@echo "  clean           - Remove build artifacts"
	@echo "  deps            - Download dependencies"
	@echo "  tidy            - Tidy and verify dependencies"
	@echo "  version         - Show version info"
	@echo ""
	@echo "Current version: $(VERSION)"

## version: Show version information
.PHONY: version
version:
	@echo "Version:    $(VERSION)"
	@echo "Commit:     $(COMMIT)"
	@echo "Build Time: $(BUILD_TIME)"
	@echo "Go Version: $(GO_VERSION)"

###################
# Build Targets
###################

## build: Build binary for current platform
.PHONY: build
build:
	@echo "Building $(PROJECT_NAME) for current platform..."
	@$(MKDIR) $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(PROJECT_NAME)$(EXE_EXT) $(CMD_DIR)
	@echo "✓ Built: $(BUILD_DIR)/$(PROJECT_NAME)$(EXE_EXT)"

## build-all: Build binaries for all platforms
.PHONY: build-all
build-all: $(PLATFORMS)

.PHONY: $(PLATFORMS)
$(PLATFORMS):
	@echo "Building for $@..."
	@$(MKDIR) $(BUILD_DIR)
	$(eval GOOS=$(word 1,$(subst /, ,$@)))
	$(eval GOARCH=$(word 2,$(subst /, ,$@)))
	$(eval OUTPUT=$(BUILD_DIR)/$(PROJECT_NAME)_$(GOOS)_$(GOARCH)$(if $(filter windows,$(GOOS)),.exe,))
	GOOS=$(GOOS) GOARCH=$(GOARCH) $(GOBUILD) -o $(OUTPUT) $(CMD_DIR)
	@echo "✓ Built: $(OUTPUT)"

## build-linux: Build for Linux platforms
.PHONY: build-linux
build-linux:
	@$(MAKE) linux/amd64 linux/arm64 linux/386 linux/arm

## build-darwin: Build for macOS platforms
.PHONY: build-darwin
build-darwin:
	@$(MAKE) darwin/amd64 darwin/arm64

## build-windows: Build for Windows platforms
.PHONY: build-windows
build-windows:
	@$(MAKE) windows/amd64 windows/386

###################
# Development
###################

## dev: Build and run for development
.PHONY: dev
dev: fmt lint test build
	@echo "Running $(PROJECT_NAME)..."
	./$(BUILD_DIR)/$(PROJECT_NAME)$(EXE_EXT) --help

## run: Run the application (requires prior build)
.PHONY: run
run:
	$(GO) run $(CMD_DIR) $(ARGS)

## fmt: Format code
.PHONY: fmt
fmt:
	@echo "Formatting code..."
	$(GOFMT) -s -w .
	$(GO) mod tidy
	@echo "✓ Code formatted"

## lint: Run linters
.PHONY: lint
lint:
	@echo "Running linters..."
	@if command -v $(GOLINT) >/dev/null 2>&1; then \
		$(GOLINT) run ./...; \
		echo "✓ Linting complete"; \
	else \
		echo "⚠ golangci-lint not installed, skipping..."; \
		echo "  Install: curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin"; \
	fi

## test: Run tests
.PHONY: test
test:
	@echo "Running tests..."
	$(GOTEST) -v -race -timeout 30s ./...
	@echo "✓ Tests passed"

## test-coverage: Run tests with coverage
.PHONY: test-coverage
test-coverage:
	@echo "Running tests with coverage..."
	@$(MKDIR) $(COVERAGE_DIR)
	$(GOTEST) -v -race -coverprofile=$(COVERAGE_DIR)/coverage.out -covermode=atomic ./...
	$(GO) tool cover -html=$(COVERAGE_DIR)/coverage.out -o $(COVERAGE_DIR)/coverage.html
	@echo "✓ Coverage report: $(COVERAGE_DIR)/coverage.html"

## bench: Run benchmarks
.PHONY: bench
bench:
	@echo "Running benchmarks..."
	$(GOTEST) -bench=. -benchmem ./...

###################
# Maintenance
###################

## clean: Remove build artifacts
.PHONY: clean
clean:
	@echo "Cleaning build artifacts..."
	@$(RM) $(BUILD_DIR) 2>/dev/null || true
	@$(RM) $(COVERAGE_DIR) 2>/dev/null || true
	$(GO) clean -cache -testcache
	@echo "✓ Clean complete"

## deps: Download dependencies
.PHONY: deps
deps:
	@echo "Downloading dependencies..."
	$(GO) mod download
	@echo "✓ Dependencies downloaded"

## tidy: Tidy and verify dependencies
.PHONY: tidy
tidy:
	@echo "Tidying dependencies..."
	$(GO) mod tidy
	$(GO) mod verify
	@echo "✓ Dependencies verified"

## install: Install binary to GOPATH/bin
.PHONY: install
install:
	@echo "Installing $(PROJECT_NAME)..."
	$(GO) install $(LDFLAGS) $(CMD_DIR)
	@echo "✓ Installed to $(shell go env GOPATH)/bin/$(PROJECT_NAME)"

## uninstall: Remove installed binary
.PHONY: uninstall
uninstall:
	@echo "Uninstalling $(PROJECT_NAME)..."
	@$(RM) $(shell go env GOPATH)/bin/$(PROJECT_NAME)$(EXE_EXT)
	@echo "✓ Uninstalled"

###################
# Release (for local testing)
###################

## release: Create release archives (requires build-all)
.PHONY: release
release: build-all
	@echo "Creating release archives..."
	@for file in $(BUILD_DIR)/$(PROJECT_NAME)_*; do \
		base=$$(basename $$file); \
		if [ "$${base##*.}" = "exe" ]; then \
			base=$${base%.exe}; \
			zip -j $(BUILD_DIR)/$$base.zip $$file README.md LICENSE 2>/dev/null || true; \
		else \
			tar -czf $(BUILD_DIR)/$$base.tar.gz -C $(BUILD_DIR) $$(basename $$file) -C .. README.md LICENSE 2>/dev/null || true; \
		fi; \
		echo "✓ Created archive for $$base"; \
	done
	@echo "✓ Release archives created in $(BUILD_DIR)/"

# Ensure build directory exists
$(BUILD_DIR):
	@$(MKDIR) $(BUILD_DIR)
