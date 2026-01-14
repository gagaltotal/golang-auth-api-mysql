#!/bin/bash

echo "Installing Swag CLI for Swagger documentation generation..."

# Check if Go is installed
if ! command -v go &> /dev/null; then
    echo "Error: Go is not installed"
    exit 1
fi

# Install swag
go install github.com/swaggo/swag/cmd/swag@latest

# Verify installation
if command -v swag &> /dev/null; then
    echo "Swag installed successfully!"
    swag --version
else
    echo "Error: Swag installation failed"
    echo "Make sure your GOPATH/bin is in your PATH"
    exit 1
fi

echo ""
echo "To generate Swagger docs, run: swag init -g cmd/api/main.go -o docs"