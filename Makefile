# Project name
PROJECT_NAME := folder-diff

# Go parameters
GO := go
GOFMT := gofmt -s
GOBUILD := $(GO) build
GOCLEAN := $(GO) clean
GOTEST := $(GO) test
GOGET := $(GO) get
GOINSTALL := $(GO) install
BUILD_DIR := build
BINARY_NAME := $(BUILD_DIR)/$(PROJECT_NAME)
BINARY_UNIX := $(BUILD_DIR)/$(PROJECT_NAME)_unix
BINARY_WINDOWS := $(BUILD_DIR)/$(PROJECT_NAME).exe
BINARY_MAC := $(BUILD_DIR)/$(PROJECT_NAME)_mac
BINARY_SYNOLOGY_ARM := $(BUILD_DIR)/$(PROJECT_NAME)_synology_arm
BINARY_SYNOLOGY_X86 := $(BUILD_DIR)/$(PROJECT_NAME)_synology_x86

# Directories
CMD_DIR := ./cmd/$(PROJECT_NAME)

# Default target executed when no arguments are given to make.
all: test build

# Format the code
fmt:
	$(GOFMT) -w .

# Get the dependencies
deps:
	$(GOGET) -v ./...

# Build the project
build: fmt
	mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BINARY_NAME) $(CMD_DIR)

# Cross-compile for different platforms
build-all: fmt
	mkdir -p $(BUILD_DIR)
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_UNIX) $(CMD_DIR)
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) $(CMD_DIR)
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_MAC) $(CMD_DIR)
	GOOS=linux GOARCH=arm $(GOBUILD) -o $(BINARY_SYNOLOGY_ARM) $(CMD_DIR)
	GOOS=linux GOARCH=386 $(GOBUILD) -o $(BINARY_SYNOLOGY_X86) $(CMD_DIR)

# Run tests
test:
	$(GOTEST) -v ./...

# Clean the build directory
clean:
	$(GOCLEAN)
	rm -rf $(BUILD_DIR)

# Install the binary
install:
	$(GOINSTALL) $(CMD_DIR)

# Run the application
run: build
	./$(BINARY_NAME)

.PHONY: all fmt deps build build-all test clean install run
