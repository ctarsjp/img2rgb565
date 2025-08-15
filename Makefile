# Makefile for img2rgb565

# Binary name
BINARY_NAME=img2rgb565
BINARY_WINDOWS=$(BINARY_NAME).exe

# Go parameters
GOCMD=go
GOBUILD=$(GOCMD) build
GOCLEAN=$(GOCMD) clean
GOTEST=$(GOCMD) test
GOGET=$(GOCMD) get
GOMOD=$(GOCMD) mod

# Build flags
LDFLAGS=-ldflags "-s -w"

# Default target
all: build

# Build the binary
build:
	$(GOBUILD) -o $(BINARY_NAME) -v $(LDFLAGS)

# Build for Windows
build-windows:
	GOOS=windows GOARCH=amd64 $(GOBUILD) -o $(BINARY_WINDOWS) -v $(LDFLAGS)

# Build for Linux
build-linux:
	GOOS=linux GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v $(LDFLAGS)

# Build for macOS
build-darwin:
	GOOS=darwin GOARCH=amd64 $(GOBUILD) -o $(BINARY_NAME) -v $(LDFLAGS)

# Build for all platforms
build-all: build-windows build-linux build-darwin

# Clean build artifacts
clean:
	$(GOCLEAN)
	rm -f $(BINARY_NAME)
	rm -f $(BINARY_WINDOWS)

# Download dependencies
deps:
	$(GOGET) -v ./...
	$(GOMOD) download

# Update dependencies
update-deps:
	$(GOGET) -u ./...
	$(GOMOD) tidy

# Run with test image
test-run: build
	./$(BINARY_NAME) logo.png

# Install to GOPATH/bin
install: build
	$(GOCMD) install

# Help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary for current OS"
	@echo "  make build-windows - Build Windows executable"
	@echo "  make build-linux   - Build Linux binary"
	@echo "  make build-darwin  - Build macOS binary"
	@echo "  make build-all     - Build for all platforms"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make deps          - Download dependencies"
	@echo "  make update-deps   - Update dependencies"
	@echo "  make test-run      - Build and test with logo.png"
	@echo "  make install       - Install to GOPATH/bin"
	@echo "  make help          - Show this help message"

.PHONY: all build build-windows build-linux build-darwin build-all clean deps update-deps test-run install help