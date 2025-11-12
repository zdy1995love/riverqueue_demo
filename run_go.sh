#!/bin/bash

# River Queue Demo - Go Task Inserter
# This script inserts tasks into River Queue using Go

echo "📝 Inserting tasks using Go..."
echo "=============================="

# Check if we're in the correct directory
if [ ! -f "examples/insert_tasks.go" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Check if config file exists
if [ ! -f "setting/config_DEV.jsonc" ]; then
    echo "❌ Error: Config file not found: setting/config_DEV.jsonc"
    exit 1
fi

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "❌ Error: Go is not installed"
    exit 1
fi

echo "🏃 Running Go task inserter..."
echo ""

# Run the Go task inserter
go run examples/insert_tasks.go
