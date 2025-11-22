# Makefile for mdns-reflector-go

BINARY_NAME=mdns-reflector-go
VERSION?=dev
BUILD_DIR=build

# Build flags
LDFLAGS=-ldflags "-X main.version=$(VERSION)"

# Go build environment
CGO_ENABLED=0
GOOS=$(shell go env GOOS)
GOARCH=$(shell go env GOARCH)

# Force disable CGO for all builds to avoid version conflicts
export CGO_ENABLED

.PHONY: all clean build build-all release test

all: build

clean:
	rm -rf $(BUILD_DIR)

build: clean
	@echo "Building $(BINARY_NAME) for $(GOOS)/$(GOARCH)..."
	@mkdir -p $(BUILD_DIR)
	CGO_ENABLED=$(CGO_ENABLED) go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

build-all: clean
	@echo "Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)

	# macOS
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 .
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 .

	# Linux
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 .
	GOOS=linux GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-arm64 .

	# Windows
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe .

	@echo "Build complete. Binaries are in $(BUILD_DIR)/"

release: build-all
	@echo "Creating release archives..."
	@cd $(BUILD_DIR) && \
	for file in $(BINARY_NAME)-*; do \
		if [[ $$file == *.exe ]]; then \
			zip $${file%.exe}-$(VERSION).zip $$file; \
		else \
			tar -czf $${file}-$(VERSION).tar.gz $$file; \
		fi; \
	done
	@echo "Release archives created in $(BUILD_DIR)/"

test:
	go test -v ./...

install: build
	@echo "Installing $(BINARY_NAME) to /usr/local/bin..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) /usr/local/bin/

# Development helpers
fmt:
	go fmt ./...

vet:
	go vet ./...

lint: fmt vet
	golint ./...

deps:
	go mod download
	go mod tidy
