# DotFiles TUI — Makefile
# Run from tui/ or repo root: `make -C tui build`
# The Go module lives at the repository root (parent of this directory).

BINARY_NAME := dotfiles-tui
BUILD_DIR := ./build
INSTALL_DIR := /usr/local/bin
APPIMAGE_DIR := $(BUILD_DIR)/AppDir
VERSION ?= dev
REPO_ROOT := $(abspath ..)
PKG := ./tui

GO_FLAGS := -ldflags="-s -w -X main.version=$(VERSION)"

.PHONY: build run install clean lint test tidy appimage help generate-modules

## build: Compile the application binary
build:
	@echo "Building $(BINARY_NAME) (version=$(VERSION))..."
	@mkdir -p $(BUILD_DIR)
	cd $(REPO_ROOT) && go build $(GO_FLAGS) -o $(CURDIR)/$(BUILD_DIR)/$(BINARY_NAME) $(PKG)

## run: Build and run the application
run:
	cd $(REPO_ROOT) && go run $(PKG)

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
		cd $(REPO_ROOT) && golangci-lint run ./...; \
	else \
		echo "golangci-lint not installed. Run: go install github.com/golangci-lint/cmd/golangci-lint@latest"; \
	fi

## test: Run all tests
test:
	cd $(REPO_ROOT) && go test -v ./...

## tidy: Clean up go.mod and go.sum
tidy:
	cd $(REPO_ROOT) && go mod tidy

## generate-modules: Refresh tui/modules_register.go blank-imports from modules/*/module.go
generate-modules:
	cd $(REPO_ROOT) && go run ./tui/tools/genregister.go

## appimage: Build a Linux AppImage (requires appimagetool on PATH)
appimage: build
	@echo "Packaging AppImage..."
	@mkdir -p $(APPIMAGE_DIR)/usr/bin
	cp $(BUILD_DIR)/$(BINARY_NAME) $(APPIMAGE_DIR)/usr/bin/$(BINARY_NAME)
	cp appimage/AppRun $(APPIMAGE_DIR)/AppRun
	chmod +x $(APPIMAGE_DIR)/AppRun
	cp appimage/dotfiles-tui.desktop $(APPIMAGE_DIR)/dotfiles-tui.desktop
	cp appimage/dotfiles-tui.svg $(APPIMAGE_DIR)/dotfiles-tui.svg
	cd $(BUILD_DIR) && ARCH=x86_64 appimagetool AppDir $(BINARY_NAME)-$(VERSION)-x86_64.AppImage
	@echo "AppImage created: $(BUILD_DIR)/$(BINARY_NAME)-$(VERSION)-x86_64.AppImage"

## help: Show this help message
help:
	@echo "DotFiles TUI — Available targets:"
	@echo ""
	@grep -E '^## ' $(MAKEFILE_LIST) | sed 's/## /  /' | sort
