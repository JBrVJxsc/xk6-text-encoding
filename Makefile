.PHONY: build test clean deps build-xk6 run-test fmt lint dev-deps help test-all test-go test-k6 install-xk6 check-xk6

# Variables
EXTENSION_NAME = xk6-text-encoding
MODULE_PATH = github.com/JBrVJxsc/xk6-text-encoding
TEST_FILE = text_encoding_k6_test.js

# Default target
all: deps build test-all

# Build the extension (standalone)
build: deps
	go build -o $(EXTENSION_NAME) .

# Build with xk6
build-xk6: check-xk6
	xk6 build --with $(MODULE_PATH)@latest

# Build with xk6 using local code for development
build-xk6-local: check-xk6
	xk6 build --with $(MODULE_PATH)=.

# Download and tidy dependencies
deps:
	go mod download
	go mod tidy

# Install xk6 if not present
install-xk6:
	@which xk6 >/dev/null 2>&1 || { \
		echo "Installing xk6..."; \
		go install go.k6.io/xk6/cmd/xk6@latest; \
	}

# Check if xk6 is installed
check-xk6:
	@which xk6 >/dev/null 2>&1 || { \
		echo "Error: xk6 not found. Run 'make install-xk6' to install it."; \
		exit 1; \
	}

# Run Go tests only
test-go:
	go test -v ./...

# Run Go tests with coverage
test-go-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Run K6 tests only (requires xk6 build)
test-k6: build-xk6
	./k6 run $(TEST_FILE)

# Run both Go and K6 tests
test-all: test-go test-k6
	@echo ""
	@echo "✅ All tests completed successfully!"

# Run tests and generate coverage
test-full: test-go-coverage test-k6
	@echo ""
	@echo "✅ Full test suite completed with coverage!"

# Clean build artifacts
clean:
	rm -f $(EXTENSION_NAME)
	rm -f k6
	rm -f coverage.out
	rm -f coverage.html
	go clean

# Format Go code
fmt:
	go fmt ./...

# Lint Go code
lint:
	golangci-lint run

# Lint and fix issues automatically
lint-fix:
	golangci-lint run --fix

# Install development dependencies
dev-deps: install-xk6
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	@echo "Development dependencies installed."

# Run the K6 test script (alias for test-k6)
run-test: test-k6

# Verify the module
verify:
	go mod verify

# Security check
security:
	@which gosec >/dev/null 2>&1 || { \
		echo "Installing gosec..."; \
		go install github.com/securecodewarrior/gosec/v2/cmd/gosec@latest; \
	}
	gosec ./...

# Performance test with K6
perf-test: build-xk6
	./k6 run --vus 10 --duration 30s $(TEST_FILE)

# Quick development cycle: format, lint, test
dev: fmt lint test-all

# Prepare for release: full checks
release-check: deps fmt lint security verify test-full
	@echo "✅ Release checks completed successfully!"

# Initialize project (run once after clone)
init: dev-deps deps
	@echo "Project initialized successfully!"

# Show detailed help
help:
	@echo "=== $(EXTENSION_NAME) Build System ==="
	@echo ""
	@echo "Quick Commands:"
	@echo "  make dev          - Quick development cycle (fmt + lint + test)"
	@echo "  make test-all     - Run both Go and K6 tests"
	@echo "  make build-xk6    - Build with xk6"
	@echo ""
	@echo "Build Commands:"
	@echo "  make build        - Build the extension (standalone)"
	@echo "  make build-xk6    - Build with xk6 (latest from GitHub)"
	@echo "  make build-xk6-local - Build with xk6 (local code)"
	@echo ""
	@echo "Test Commands:"
	@echo "  make test-go      - Run Go tests only"
	@echo "  make test-k6      - Run K6 tests only"
	@echo "  make test-all     - Run both Go and K6 tests"
	@echo "  make test-full    - Run tests with coverage report"
	@echo "  make perf-test    - Run performance test with K6"
	@echo ""
	@echo "Development Commands:"
	@echo "  make deps         - Download and tidy dependencies"
	@echo "  make fmt          - Format Go code"
	@echo "  make lint         - Lint Go code"
	@echo "  make lint-fix     - Lint and auto-fix issues"
	@echo "  make security     - Run security checks"
	@echo ""
	@echo "Utility Commands:"
	@echo "  make clean        - Clean build artifacts"
	@echo "  make verify       - Verify module dependencies"
	@echo "  make dev-deps     - Install development dependencies"
	@echo "  make install-xk6  - Install xk6 tool"
	@echo "  make init         - Initialize project (run once)"
	@echo "  make release-check - Full pre-release validation"
	@echo "  make help         - Show this help"
	@echo ""
	@echo "Files:"
	@echo "  Test file: $(TEST_FILE)"
	@echo "  Module:    $(MODULE_PATH)"