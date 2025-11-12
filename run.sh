#!/bin/bash

# River Queue Demo - Main Runner Script
# This script starts the River Queue worker

echo "🚀 Starting River Queue Demo..."
echo "================================"

# Check if we're in the correct directory
if [ ! -f "main.go" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Check if config file exists
if [ ! -f "setting/config_DEV.jsonc" ]; then
    echo "❌ Error: Config file not found: setting/config_DEV.jsonc"
    echo "� Please create the config file with your database credentials"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    echo "💡 Please install Go from https://go.dev/"
    exit 1
fi

echo "✅ Go version: $(go version)"
echo ""

# Download dependencies if needed
if [ ! -d "vendor" ] && [ ! -f "go.sum" ]; then
    echo "📦 Downloading dependencies..."
    go mod download
fi

echo "🏃 Starting River Queue worker..."
echo "Press Ctrl+C to stop"
echo ""

# Run the application
go run main.go
