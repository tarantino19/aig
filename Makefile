# AIG (AI Git CLI) Makefile

BINARY_NAME=aig
MAIN_PATH=cmd/aig/main.go
INSTALL_PATH=/usr/local/bin
BUILD_DIR=build

# Build info
VERSION=$(shell git describe --tags --always --dirty 2>/dev/null || echo "dev")
BUILD_TIME=$(shell date +%Y-%m-%dT%H:%M:%S)
LDFLAGS=-ldflags "-X main.version=$(VERSION) -X main.buildTime=$(BUILD_TIME)"

.PHONY: all build clean install uninstall test help

# Default target
all: build

# Build the binary
build:
	@echo "🔨 Building $(BINARY_NAME)..."
	@go build $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)
	@echo "✅ Build complete: $(BINARY_NAME)"

# Build for multiple platforms
build-all:
	@echo "🔨 Building for multiple platforms..."
	@mkdir -p $(BUILD_DIR)
	@echo "Building for Linux..."
	@GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-linux-amd64 $(MAIN_PATH)
	@echo "Building for macOS..."
	@GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-amd64 $(MAIN_PATH)
	@GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-darwin-arm64 $(MAIN_PATH)
	@echo "Building for Windows..."
	@GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/$(BINARY_NAME)-windows-amd64.exe $(MAIN_PATH)
	@echo "✅ Multi-platform build complete in $(BUILD_DIR)/"

# Install system-wide
install: build
	@echo "📦 Installing $(BINARY_NAME) to $(INSTALL_PATH)..."
	@if [ -w "$(INSTALL_PATH)" ]; then \
		cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
		chmod +x $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		sudo cp $(BINARY_NAME) $(INSTALL_PATH)/$(BINARY_NAME); \
		sudo chmod +x $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@mkdir -p ~/.config/aig
	@echo "✅ Installation complete! Run 'aig --help' to get started."

# Uninstall
uninstall:
	@echo "🗑️  Uninstalling $(BINARY_NAME)..."
	@if [ -w "$(INSTALL_PATH)" ]; then \
		rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	else \
		sudo rm -f $(INSTALL_PATH)/$(BINARY_NAME); \
	fi
	@echo "✅ Uninstallation complete."

# Run tests
test:
	@echo "🧪 Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "🧹 Cleaning..."
	@rm -f $(BINARY_NAME)
	@rm -rf $(BUILD_DIR)
	@echo "✅ Clean complete."

# Development build with race detection
dev:
	@echo "🔨 Building development version..."
	@go build -race $(LDFLAGS) -o $(BINARY_NAME) $(MAIN_PATH)

# Format code
fmt:
	@echo "🎨 Formatting code..."
	@go fmt ./...

# Lint code
lint:
	@echo "🔍 Linting code..."
	@golangci-lint run

# Update dependencies
deps:
	@echo "📦 Updating dependencies..."
	@go mod tidy
	@go mod download

# Show help
help:
	@echo "AIG (AI Git CLI) Build Commands:"
	@echo ""
	@echo "🔨 Build Commands:"
	@echo "  make build      - Build the binary"
	@echo "  make build-all  - Build for multiple platforms"
	@echo "  make dev        - Build with race detection"
	@echo ""
	@echo "📦 Installation:"
	@echo "  make install    - Install system-wide"
	@echo "  make uninstall  - Remove installation"
	@echo ""
	@echo "🧪 Development:"
	@echo "  make test       - Run tests"
	@echo "  make fmt        - Format code"
	@echo "  make lint       - Lint code"
	@echo "  make deps       - Update dependencies"
	@echo ""
	@echo "🧹 Cleanup:"
	@echo "  make clean      - Clean build artifacts"
	@echo ""
	@echo "❓ Help:"
	@echo "  make help       - Show this help" 