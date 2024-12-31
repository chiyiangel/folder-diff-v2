# Project configuration
PROJECT_NAME := folder-diff
VERSION := $(shell git describe --tags --always --dirty)
BUILD_TIME := $(shell date +%Y%m%d.%H%M%S)
LDFLAGS := -ldflags="-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

# Go parameters
GO := go
GOFMT := gofmt -s
GOBUILD := $(GO) build $(LDFLAGS)
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOGET := $(GO) get

# Directories
BUILD_DIR := build
CMD_DIR := ./cmd/$(PROJECT_NAME)

# Build targets
TARGETS := \
	$(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64 \
	$(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe \
	$(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64 \
	$(BUILD_DIR)/$(PROJECT_NAME)_linux_arm \
	$(BUILD_DIR)/$(PROJECT_NAME)_linux_386 \
	$(BUILD_DIR)/$(PROJECT_NAME)_synology_arm \
	$(BUILD_DIR)/$(PROJECT_NAME)_synology_x86

# Default target
all: build

# Build all targets
build: $(TARGETS)

# Individual platform builds
$(BUILD_DIR)/$(PROJECT_NAME)_linux_amd64:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_windows_amd64.exe:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_darwin_amd64:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_linux_arm:
	GOOS=linux GOARCH=arm $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_linux_386:
	GOOS=linux GOARCH=386 $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_synology_arm:
	GOOS=linux GOARCH=arm $(GOBUILD) -o $@ $(CMD_DIR)

$(BUILD_DIR)/$(PROJECT_NAME)_synology_x86:
	GOOS=linux GOARCH=386 $(GOBUILD) -o $@ $(CMD_DIR)

# Development targets
dev: fmt test
	$(GOBUILD) -o $(BUILD_DIR)/$(PROJECT_NAME) $(CMD_DIR)

# Format code
fmt:
	$(GOFMT) -w .

# Run tests
test:
	$(GOTEST) -v -cover ./...

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Install dependencies
deps:
	$(GOGET) -v ./...

# Run the application
run: dev
	./$(BUILD_DIR)/$(PROJECT_NAME)

# Show help
help:
	@echo "Available targets:"
	@echo "  all       - Build all platform binaries (default)"
	@echo "  build     - Build all platform binaries"
	@echo "  dev       - Build for development (with tests and formatting)"
	@echo "  fmt       - Format source code"
	@echo "  test      - Run tests"
	@echo "  clean     - Remove build artifacts"
	@echo "  deps      - Install dependencies"
	@echo "  run       - Build and run the application"
	@echo "  help      - Show this help message"

.PHONY: all build dev fmt test clean deps run help
