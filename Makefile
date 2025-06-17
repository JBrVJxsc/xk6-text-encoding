.PHONY: build test clean deps

# Build the extension
build:
	go build -o xk6-text-encoding .

# Build with xk6
build-xk6:
	xk6 build --with github.com/JBrVJxsc/xk6-text-encoding@latest

# Download dependencies
deps:
	go mod download
	go mod tidy

# Run tests
test:
	go test -v ./...

# Clean build artifacts
clean:
	rm -f xk6-text-encoding
	go clean

# Run the test script with k6
run-test:
	k6 run test.js

# Format code
fmt:
	go fmt .

# Lint code
lint:
	golangci-lint run

# Install dependencies for development
dev-deps:
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

# Show help
help:
	@echo "Available targets:"
	@echo "  build      - Build the extension"
	@echo "  build-xk6  - Build with xk6"
	@echo "  deps       - Download and tidy dependencies"
	@echo "  test       - Run Go tests"
	@echo "  run-test   - Run the k6 test script"
	@echo "  clean      - Clean build artifacts"
	@echo "  fmt        - Format Go code"
	@echo "  lint       - Lint Go code"
	@echo "  dev-deps   - Install development dependencies"
	@echo "  help       - Show this help" 