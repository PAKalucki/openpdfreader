.PHONY: all build run test clean install-deps lint fmt

# Binary name
BINARY=openpdfreader
BINARY_WINDOWS=$(BINARY).exe

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GORUN=$(GOCMD) run
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOFMT=$(GOCMD) fmt
GOVET=$(GOCMD) vet

# Build directory
BUILD_DIR=build

all: build

build:
	@mkdir -p $(BUILD_DIR)
	$(GOBUILD) -o $(BUILD_DIR)/$(BINARY) ./cmd/openpdfreader

build-windows:
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=1 GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BUILD_DIR)/$(BINARY_WINDOWS) ./cmd/openpdfreader

run:
	$(GORUN) ./cmd/openpdfreader

test:
	$(GOTEST) -v ./...

test-coverage:
	$(GOTEST) -coverprofile=coverage.out ./...
	$(GOCMD) tool cover -html=coverage.out -o coverage.html

clean:
	rm -rf $(BUILD_DIR)
	rm -f coverage.out coverage.html

install-deps:
	$(GOGET) -u fyne.io/fyne/v2
	$(GOGET) -u github.com/pdfcpu/pdfcpu

lint:
	golangci-lint run ./...

fmt:
	$(GOFMT) ./...

vet:
	$(GOVET) ./...

# Install fyne-cross for cross-compilation
install-fyne-cross:
	go install github.com/fyne-io/fyne-cross@latest

# Cross-compile for all platforms using fyne-cross
cross-compile: install-fyne-cross
	fyne-cross linux -arch=amd64 ./cmd/openpdfreader
	fyne-cross windows -arch=amd64 ./cmd/openpdfreader

# Development helpers
dev: fmt vet run

check: fmt vet lint test
