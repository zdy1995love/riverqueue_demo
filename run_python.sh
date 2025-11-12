#!/bin/bash

# River Queue Demo - Python Task Inserter
# This script inserts tasks into River Queue using Python

echo "📝 Inserting tasks using Python..."
echo "=================================="

# Check if we're in the correct directory
if [ ! -f "examples/insert_tasks.py" ]; then
    echo "❌ Error: Please run this script from the project root directory"
    exit 1
fi

# Check if config file exists
if [ ! -f "setting/config_DEV.jsonc" ]; then
    echo "❌ Error: Config file not found: setting/config_DEV.jsonc"
    exit 1
fi

# Check if Python is installed
if ! command -v python3 &> /dev/null; then
    echo "❌ Error: Python 3 is not installed"
    exit 1
fi

echo "✅ Python version: $(python3 --version)"
echo ""

# Check if riverqueue is installed
if ! python3 -c "import riverqueue" 2>/dev/null; then
    echo "⚠️  Warning: riverqueue-python is not installed"
    echo "💡 To install, run:"
    echo "   conda create -n riverqueue python=3.12"
    echo "   conda activate riverqueue"
    echo "   cd /path/to/riverqueue-python"
    echo "   pip install ."
    echo ""
    echo "Or install from PyPI (if available):"
    echo "   pip install riverqueue"
    echo ""
    read -p "Do you want to continue anyway? (y/N): " -n 1 -r
    echo
    if [[ ! $REPLY =~ ^[Yy]$ ]]; then
        exit 1
    fi
fi

echo "🏃 Running Python task inserter..."
echo ""

# Run the Python task inserter
python3 examples/insert_tasks.py
