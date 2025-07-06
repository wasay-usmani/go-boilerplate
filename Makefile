# Go Boilerplate Makefile
# Usage: make init MODULE_PATH=github.com/your-org/your-project
# Usage: make create-service SERVICE_NAME=my-service

.PHONY: help init create-service create-app clean build test

# Default target
help:
	@echo "Available targets:"
	@echo "  init MODULE_PATH=<new-module-path>  - Initialize boilerplate with new module path"
	@echo "  create-service SERVICE_NAME=<name>  - Create a new microservice from boilerplate"
	@echo "  create-app SERVICE_NAME=<name> APP_NAME=<name> - Add a new app module to existing service"
	@echo "  clean                               - Clean build artifacts"
	@echo "  build                               - Build all services"
	@echo "  build-service SERVICE_NAME=<name>   - Build a specific service"
	@echo "  run-service SERVICE_NAME=<name>     - Run a specific service"
	@echo "  test                                - Run tests"
	@echo ""
	@echo "Examples:"
	@echo "  make init MODULE_PATH=github.com/mycompany/myapp"
	@echo "  make create-service SERVICE_NAME=user-service"
	@echo "  make create-app SERVICE_NAME=user-service APP_NAME=auth"
	@echo "  make build-service SERVICE_NAME=user-service"
	@echo "  make run-service SERVICE_NAME=user-service"

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
	
	# Run the initialization script
	@./resources/scripts/init.sh $(MODULE_PATH)

# Create a new microservice from boilerplate templates
create-service:
	@if [ -z "$(SERVICE_NAME)" ]; then \
		echo "Error: SERVICE_NAME is required"; \
		echo "Usage: make create-service SERVICE_NAME=my-service"; \
		exit 1; \
	fi
	@echo "Creating new microservice: $(SERVICE_NAME)"
	@echo "This will generate a new microservice from templates..."
	@echo ""
	
	# Check if service already exists
	@if [ -d "cmd/$(SERVICE_NAME)" ]; then \
		echo "Error: Service '$(SERVICE_NAME)' already exists in cmd/"; \
		exit 1; \
	fi
	@if [ -d "internal/$(SERVICE_NAME)" ]; then \
		echo "Error: Service '$(SERVICE_NAME)' already exists in internal/"; \
		exit 1; \
	fi
	@if [ -d "resources/$(SERVICE_NAME)" ]; then \
		echo "Error: Service '$(SERVICE_NAME)' already exists in resources/"; \
		exit 1; \
	fi
	
	# Get current module path from go.mod
	$(eval MODULE_PATH := $(shell grep '^module ' go.mod | cut -d' ' -f2))
	@echo "Current module path: $(MODULE_PATH)"
	
	# Generate microservice from templates
	@echo "Generating microservice from templates..."
	@./resources/scripts/create-service.sh $(SERVICE_NAME) $(MODULE_PATH) .
	
	@echo ""
	@echo "✅ Microservice '$(SERVICE_NAME)' created successfully!"
	@echo "📁 Created directories:"
	@echo "  - cmd/$(SERVICE_NAME)"
	@echo "  - internal/$(SERVICE_NAME)"
	@echo "  - resources/$(SERVICE_NAME)"
	@echo ""
	@echo "Next steps:"
	@echo "1. Review the generated code in the new directories"
	@echo "2. Update service-specific configuration and business logic"
	@echo "3. Update any service-specific dependencies in go.mod"
	@echo "4. Test the new service with 'go build ./cmd/$(SERVICE_NAME)'"
	@echo "5. Add any service-specific environment variables or config"

# Add a new app module to an existing service
create-app:
	@if [ -z "$(SERVICE_NAME)" ]; then \
		echo "Error: SERVICE_NAME is required"; \
		echo "Usage: make create-app SERVICE_NAME=my-service APP_NAME=my-app"; \
		exit 1; \
	fi
	@if [ -z "$(APP_NAME)" ]; then \
		echo "Error: APP_NAME is required"; \
		echo "Usage: make create-app SERVICE_NAME=my-service APP_NAME=my-app"; \
		exit 1; \
	fi
	@echo "Adding new app module: $(APP_NAME)"
	@echo "To service: $(SERVICE_NAME)"
	@echo "This will create a new app module with cmd, module, and qrys files..."
	@echo ""
	
	# Generate app module from templates
	@echo "Generating app module from templates..."
	@./resources/scripts/create-app.sh $(SERVICE_NAME) $(APP_NAME)
	
	@echo ""
	@echo "✅ App module '$(APP_NAME)' added successfully to service '$(SERVICE_NAME)'!"
	@echo "📁 Created files:"
	@echo "  - internal/$(SERVICE_NAME)/app/$(APP_NAME)/app.go"
	@echo "  - internal/$(SERVICE_NAME)/app/$(APP_NAME)/cmds.go"
	@echo "  - internal/$(SERVICE_NAME)/app/$(APP_NAME)/qrys.go"
	@echo ""
	@echo "Next steps:"
	@echo "1. Update internal/$(SERVICE_NAME)/app/module.go to include the new app module"
	@echo "2. Add your business logic to the generated files"
	@echo "3. Update any dependencies as needed"

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
	$(eval BUILD_NAME := $(shell grep '^module ' go.mod | cut -d' ' -f2 | sed 's/.*\///'))
	@go build -o bin/$(BUILD_NAME) ./cmd/*/main.go

# Build a specific service
build-service:
	@if [ -z "$(SERVICE_NAME)" ]; then \
		echo "Error: SERVICE_NAME is required"; \
		echo "Usage: make build-service SERVICE_NAME=my-service"; \
		echo "Available services:"; \
		ls -1 cmd/ 2>/dev/null | sed 's/^/  - /' || echo "  No services found"; \
		exit 1; \
	fi
	@if [ ! -d "cmd/$(SERVICE_NAME)" ]; then \
		echo "Error: Service '$(SERVICE_NAME)' not found in cmd/"; \
		echo "Available services:"; \
		ls -1 cmd/ 2>/dev/null | sed 's/^/  - /' || echo "  No services found"; \
		exit 1; \
	fi
	@echo "Building service: $(SERVICE_NAME)"
	@mkdir -p bin
	@go build -o bin/$(SERVICE_NAME) ./cmd/$(SERVICE_NAME)/main.go
	@echo "✅ Service built successfully: bin/$(SERVICE_NAME)"

# Run a specific service
run-service:
	@if [ -z "$(SERVICE_NAME)" ]; then \
		echo "Error: SERVICE_NAME is required"; \
		echo "Usage: make run-service SERVICE_NAME=my-service"; \
		echo "Available services:"; \
		ls -1 cmd/ 2>/dev/null | sed 's/^/  - /' || echo "  No services found"; \
		exit 1; \
	fi
	@if [ ! -d "cmd/$(SERVICE_NAME)" ]; then \
		echo "Error: Service '$(SERVICE_NAME)' not found in cmd/"; \
		echo "Available services:"; \
		ls -1 cmd/ 2>/dev/null | sed 's/^/  - /' || echo "  No services found"; \
		exit 1; \
	fi
	@echo "Running service: $(SERVICE_NAME)"
	@go run ./cmd/$(SERVICE_NAME)/main.go

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
