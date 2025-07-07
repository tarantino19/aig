#!/bin/bash

# AIG (AI Git CLI) Installation Script
# This script installs aig system-wide

set -e

BINARY_NAME="aig"
INSTALL_DIR="/usr/local/bin"
CONFIG_DIR="$HOME/.config/aig"

echo "🚀 Installing AIG (AI Git CLI Tool)..."

# Check if binary exists
if [ ! -f "$BINARY_NAME" ]; then
    echo "❌ Error: $BINARY_NAME binary not found in current directory"
    echo "Please run 'go build -o aig cmd/aig/main.go' first"
    exit 1
fi

# Create config directory if it doesn't exist
if [ ! -d "$CONFIG_DIR" ]; then
    echo "📁 Creating config directory: $CONFIG_DIR"
    mkdir -p "$CONFIG_DIR"
fi

# Install binary
echo "📦 Installing $BINARY_NAME to $INSTALL_DIR..."
if [ -w "$INSTALL_DIR" ]; then
    cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    chmod +x "$INSTALL_DIR/$BINARY_NAME"
else
    sudo cp "$BINARY_NAME" "$INSTALL_DIR/$BINARY_NAME"
    sudo chmod +x "$INSTALL_DIR/$BINARY_NAME"
fi

# Create .env template if it doesn't exist
if [ ! -f "$CONFIG_DIR/.env" ]; then
    echo "📝 Creating .env template..."
    cat > "$CONFIG_DIR/.env" << 'EOF'
# AI Git CLI Configuration
# Get your Gemini API key from: https://makersuite.google.com/app/apikey

# Gemini AI API Key
AIG_GEMINI_API_KEY=your-gemini-api-key-here

# Optional: Override default model
# AIG_AI_MODEL=gemini-1.5-flash

# Optional: Override temperature (0.0 to 1.0)
# AIG_AI_TEMPERATURE=0.7
EOF
fi

# Verify installation
if command -v aig >/dev/null 2>&1; then
    echo "✅ Installation successful!"
    echo ""
    echo "📋 Next steps:"
    echo "1. Get your Gemini API key: https://makersuite.google.com/app/apikey"
    echo "2. Set your API key:"
    echo "   aig config set ai.api_key YOUR_API_KEY"
    echo "   # OR edit: $CONFIG_DIR/.env"
    echo ""
    echo "🎉 You can now use 'aig' from anywhere!"
    echo ""
    echo "📖 Quick start:"
    echo "   aig --help           # Show help"
    echo "   aig commit           # Generate commit message"
    echo "   aig review --staged  # Review staged changes"
    echo "   aig summary          # Summarize commits"
    echo ""
    aig --version
else
    echo "❌ Installation failed. Please check your PATH."
    exit 1
fi 