# Go Boilerplate Makefile
# Usage: make init MODULE_PATH=github.com/your-org/your-project

.PHONY: help init clean build test

# Default target
help:
	@echo "Available targets:"
	@echo "  init MODULE_PATH=<new-module-path>  - Initialize boilerplate with new module path"
	@echo "  clean                               - Clean build artifacts"
	@echo "  build                               - Build the application"
	@echo "  test                                - Run tests"
	@echo ""
	@echo "Example: make init MODULE_PATH=github.com/mycompany/myapp"

# Initialize boilerplate with new module path
init:
	@if [ -z "$(MODULE_PATH)" ]; then \
		echo "Error: MODULE_PATH is required"; \
		echo "Usage: make init MODULE_PATH=github.com/your-org/your-project"; \
		exit 1; \
	fi
	@echo "Initializing boilerplate with module path: $(MODULE_PATH)"
	@echo "This will update all import paths and configuration files..."
	@echo ""
	
	# Extract project name from module path
	$(eval PROJECT_NAME := $(shell echo $(MODULE_PATH) | sed 's/.*\///'))
	@echo "Project name: $(PROJECT_NAME)"
	
	# Update go.mod file
	@echo "Updating go.mod..."
	@sed -i 's|module github.com/wasay-usmani/go-boilerplate|module $(MODULE_PATH)|' go.mod
	
	# Update all Go files with new import paths
	@echo "Updating import paths in Go files..."
	@find . -name "*.go" -type f -exec sed -i 's|github.com/wasay-usmani/go-boilerplate|$(MODULE_PATH)|g' {} \;
	
	# Rename cmd directory if needed
	@if [ -d "cmd/go-boilerplate" ]; then \
		echo "Renaming cmd/go-boilerplate to cmd/$(PROJECT_NAME)..."; \
		mv cmd/go-boilerplate cmd/$(PROJECT_NAME); \
	fi
	
	# Rename internal directory if needed
	@if [ -d "internal/go-boilerplate" ]; then \
		echo "Renaming internal/go-boilerplate to internal/$(PROJECT_NAME)..."; \
		mv internal/go-boilerplate internal/$(PROJECT_NAME); \
	fi
	
	# Update any remaining references to go-boilerplate in internal structure
	@if [ -d "internal/$(PROJECT_NAME)" ]; then \
		echo "Updating internal structure references..."; \
		find internal/$(PROJECT_NAME) -name "*.go" -type f -exec sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' {} \; 2>/dev/null || true; \
	fi
	
	# Update Dockerfile if it exists
	@if [ -f "Dockerfile" ]; then \
		echo "Updating Dockerfile..."; \
		sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' Dockerfile; \
	fi
	
	# Update docker-compose.yaml if it exists
	@if [ -f "docker-compose.yaml" ]; then \
		echo "Updating docker-compose.yaml..."; \
		sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' docker-compose.yaml; \
	fi
	
	# Update .gitignore if it contains go-boilerplate references
	@if [ -f ".gitignore" ]; then \
		echo "Updating .gitignore..."; \
		sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' .gitignore; \
	fi
	
	# Update README.md if it exists
	@if [ -f "README.md" ]; then \
		echo "Updating README.md..."; \
		sed -i 's|github.com/wasay-usmani/go-boilerplate|$(MODULE_PATH)|g' README.md; \
		sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' README.md; \
	fi
	
	# Update any configuration files in resources directory
	@if [ -d "resources" ]; then \
		echo "Updating configuration files in resources..."; \
		find resources -type f -exec sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' {} \; 2>/dev/null || true; \
		find resources -type f -exec sed -i 's|github.com/wasay-usmani/go-boilerplate|$(MODULE_PATH)|g' {} \; 2>/dev/null || true; \
	fi
	
	# Update .github workflows if they exist
	@if [ -d ".github" ]; then \
		echo "Updating GitHub workflows..."; \
		find .github -name "*.yml" -o -name "*.yaml" -exec sed -i 's|go-boilerplate|$(PROJECT_NAME)|g' {} \; 2>/dev/null || true; \
		find .github -name "*.yml" -o -name "*.yaml" -exec sed -i 's|github.com/wasay-usmani/go-boilerplate|$(MODULE_PATH)|g' {} \; 2>/dev/null || true; \
	fi
	
	# Run go mod tidy to clean up dependencies
	@echo "Running go mod tidy..."
	@go mod tidy
	
	@echo ""
	@echo "✅ Boilerplate initialization complete!"
	@echo "📁 New module path: $(MODULE_PATH)"
	@echo "📁 Project name: $(PROJECT_NAME)"
	@echo ""
	@echo "Next steps:"
	@echo "1. Review the changes made"
	@echo "2. Update any additional configuration files specific to your project"
	@echo "3. Run 'go mod tidy' if needed"
	@echo "4. Test your application with 'make build' and 'make test'"

# Clean build artifacts
clean:
	@echo "Cleaning build artifacts..."
	@go clean -cache -modcache -testcache
	@rm -rf bin/
	@rm -rf dist/
	@find . -name "*.exe" -delete
	@find . -name "*.test" -delete

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/app ./cmd/*/main.go

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Run tests with coverage
test-coverage:
	@echo "Running tests with coverage..."
	@go test -v -coverprofile=coverage.out ./...
	@go tool cover -html=coverage.out -o coverage.html
	@echo "Coverage report generated: coverage.html"

# Install dependencies
deps:
	@echo "Installing dependencies..."
	@go mod download
	@go mod verify

# Format code
fmt:
	@echo "Formatting code..."
	@go fmt ./...

# Run linter
lint:
	@echo "Running linter..."
	@golangci-lint run

# Run all checks (format, lint, test)
check: fmt lint test
	@echo "All checks passed!"
