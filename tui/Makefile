# DotFiles TUI — Makefile
# Run `make help` to see all available targets.

BINARY_NAME := dotfiles-tui
BUILD_DIR := ./build
INSTALL_DIR := /usr/local/bin

# Go build flags
# -ldflags="-s -w" strips debug info for smaller binaries
# See: https://pkg.go.dev/cmd/go#hdr-Compile_packages_and_dependencies
GO_FLAGS := -ldflags="-s -w"

.PHONY: build run install clean lint test help

## build: Compile the application binary
build:
	@echo "Building $(BINARY_NAME)..."
	@mkdir -p $(BUILD_DIR)
	go build $(GO_FLAGS) -o $(BUILD_DIR)/$(BINARY_NAME) .

## run: Build and run the application
run:
	go run .

## install: Build and install the binary to /usr/local/bin
install: build
	@echo "Installing $(BINARY_NAME) to $(INSTALL_DIR)..."
	sudo cp $(BUILD_DIR)/$(BINARY_NAME) $(INSTALL_DIR)/$(BINARY_NAME)
	@echo "Done! Run '$(BINARY_NAME)' to start."

## clean: Remove build artifacts
clean:
	@echo "Cleaning..."
	rm -rf $(BUILD_DIR)

## lint: Run Go linter (requires golangci-lint)
lint:
	@if command -v golangci-lint >/dev/null 2>&1; then \
		golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## test: Run all tests
test:
	go test -v ./...

## tidy: Clean up go.mod and go.sum
tidy:
	go mod tidy

## help: Show this help message
help:
	@echo "DotFiles TUI — Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | sort
