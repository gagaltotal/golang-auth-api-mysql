#!/bin/bash

echo "Generating Swagger documentation..."

# Check if swag is installed
if ! command -v swag &> /dev/null; then
    echo "Error: swag is not installed"
    echo "Run: go install github.com/swaggo/swag/cmd/swag@latest"
    exit 1
fi

# Generate swagger docs
swag init -g cmd/api/main.go -o docs --parseDependency --parseInternal

if [ $? -eq 0 ]; then
    echo "Swagger documentation generated successfully!"
    echo "Docs generated in ./docs directory"
    echo ""
    echo "To view Swagger UI, run the application and visit:"
    echo "http://localhost:8080/swagger/index.html"
else
    echo "Error: Failed to generate Swagger documentation"
    exit 1
fi