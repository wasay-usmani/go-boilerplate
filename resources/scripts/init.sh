#!/bin/bash

# Initialize boilerplate with new module path
# Usage: ./init.sh <MODULE_PATH>
# Example: ./init.sh github.com/your-org/your-project

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Function to print colored output
print_info() {
    echo -e "${BLUE}ℹ️  $1${NC}"
}

print_success() {
    echo -e "${GREEN}✅ $1${NC}"
}

print_warning() {
    echo -e "${YELLOW}⚠️  $1${NC}"
}

print_error() {
    echo -e "${RED}❌ $1${NC}"
}

# Check if MODULE_PATH is provided
if [ $# -eq 0 ]; then
    print_error "MODULE_PATH is required"
    echo "Usage: $0 <MODULE_PATH>"
    echo "Example: $0 github.com/your-org/your-project"
    exit 1
fi

MODULE_PATH="$1"

# Validate module path format
if [[ ! "$MODULE_PATH" =~ ^[a-zA-Z0-9._-]+/[a-zA-Z0-9._-]+(/[a-zA-Z0-9._-]+)*$ ]]; then
    print_error "Invalid module path format: $MODULE_PATH"
    echo "Module path should be in format: domain.com/org/project"
    exit 1
fi

print_info "Initializing boilerplate with module path: $MODULE_PATH"
print_info "This will update all import paths and configuration files..."
echo ""

# Extract project name from module path
PROJECT_NAME=$(echo "$MODULE_PATH" | sed 's/.*\///')
# Convert to valid Go package name (lowercase, no hyphens)
PACKAGE_NAME=$(echo "$PROJECT_NAME" | tr '[:upper:]' '[:lower:]' | sed 's/-//g')
print_info "Project name: $PROJECT_NAME"
print_info "Package name: $PACKAGE_NAME"

# Update go.mod file
print_info "Updating go.mod..."
if [ -f "go.mod" ]; then
    sed -i "s|module github.com/wasay-usmani/go-boilerplate|module $MODULE_PATH|" go.mod
    print_success "Updated go.mod"
else
    print_warning "go.mod file not found"
fi

# Update all Go files with new import paths
print_info "Updating import paths in Go files..."
find . -name "*.go" -type f -exec sed -i "s|github.com/wasay-usmani/go-boilerplate|$MODULE_PATH|g" {} \;
print_success "Updated import paths in Go files"

# Rename cmd directory to match project name
if [ -d "cmd/go-boilerplate" ]; then
    print_info "Renaming cmd/go-boilerplate to cmd/$PROJECT_NAME..."
    mv cmd/go-boilerplate "cmd/$PROJECT_NAME"
    print_success "Renamed cmd directory"
fi

# Rename internal directory to match project name
if [ -d "internal/go-boilerplate" ]; then
    print_info "Renaming internal/go-boilerplate to internal/$PROJECT_NAME..."
    mv internal/go-boilerplate "internal/$PROJECT_NAME"
    print_success "Renamed internal directory"
fi

# Rename app subdirectory to match project name
if [ -d "internal/$PROJECT_NAME/app/go-boilerplate" ]; then
    print_info "Renaming internal/$PROJECT_NAME/app/go-boilerplate to internal/$PROJECT_NAME/app/$PROJECT_NAME..."
    mv "internal/$PROJECT_NAME/app/go-boilerplate" "internal/$PROJECT_NAME/app/$PROJECT_NAME"
    print_success "Renamed app subdirectory"
fi

# Update any remaining references to go-boilerplate in internal structure
if [ -d "internal/$PROJECT_NAME" ]; then
    print_info "Updating internal structure references..."
    
    # Update go-boilerplate references to project name
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|go-boilerplate|$PROJECT_NAME|g" {} \; 2>/dev/null || true
    
    # Update internal package names in import paths
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|internal/go-boilerplate|internal/$PROJECT_NAME|g" {} \; 2>/dev/null || true
    
    # Update app package names in import paths
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|app/go-boilerplate|app/$PROJECT_NAME|g" {} \; 2>/dev/null || true
    
    # Update package declarations (goboilerplate -> packagename)
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|^package goboilerplate|package $PACKAGE_NAME|g" {} \; 2>/dev/null || true
    
    # Update import aliases (goboilerplate -> packagename)
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|goboilerplate \"|$PACKAGE_NAME \"|g" {} \; 2>/dev/null || true
    
    # Update struct field names (Boilerplate -> ProjectName)
    PROJECT_NAME_CAPITALIZED=$(echo "$PROJECT_NAME" | sed 's/-\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|Boilerplate|$PROJECT_NAME_CAPITALIZED|g" {} \; 2>/dev/null || true
    
    # Update type references (goboilerplate.App -> packagename.App)
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|goboilerplate\.App|$PACKAGE_NAME.App|g" {} \; 2>/dev/null || true
    
    # Update function calls (goboilerplate.New -> packagename.New)
    find "internal/$PROJECT_NAME" -name "*.go" -type f -exec sed -i "s|goboilerplate\.New|$PACKAGE_NAME.New|g" {} \; 2>/dev/null || true
    
    print_success "Updated internal structure references"
fi

# Update import paths in all Go files to use the new internal structure
print_info "Updating import paths to use new internal structure..."
find . -name "*.go" -type f -exec sed -i "s|$MODULE_PATH/internal/go-boilerplate|$MODULE_PATH/internal/$PROJECT_NAME|g" {} \;
print_success "Updated import paths for internal structure"

# Update Dockerfile if it exists
if [ -f "Dockerfile" ]; then
    print_info "Updating Dockerfile..."
    sed -i "s|go-boilerplate|$PROJECT_NAME|g" Dockerfile
    print_success "Updated Dockerfile"
fi

# Update .gitignore if it contains go-boilerplate references
if [ -f ".gitignore" ]; then
    print_info "Updating .gitignore..."
    sed -i "s|go-boilerplate|$PROJECT_NAME|g" .gitignore
    print_success "Updated .gitignore"
fi

# Update README.md if it exists
if [ -f "README.md" ]; then
    print_info "Updating README.md..."
    sed -i "s|github.com/wasay-usmani/go-boilerplate|$MODULE_PATH|g" README.md
    sed -i "s|go-boilerplate|$PROJECT_NAME|g" README.md
    print_success "Updated README.md"
fi

# Update .github workflows if they exist
if [ -d ".github" ]; then
    print_info "Updating GitHub workflows..."
    find .github -name "*.yml" -o -name "*.yaml" -exec sed -i "s|go-boilerplate|$PROJECT_NAME|g" {} \; 2>/dev/null || true
    find .github -name "*.yml" -o -name "*.yaml" -exec sed -i "s|github.com/wasay-usmani/go-boilerplate|$MODULE_PATH|g" {} \; 2>/dev/null || true
    print_success "Updated GitHub workflows"
fi

# Update resources directory
if [ -d "resources" ]; then
    print_info "Updating resources directory..."
    
    # Rename the go-boilerplate directory
    if [ -d "resources/go-boilerplate" ]; then
        print_info "Renaming resources/go-boilerplate to resources/$PROJECT_NAME..."
        mv resources/go-boilerplate "resources/$PROJECT_NAME"
        print_success "Renamed resources directory"
    fi
    
    # Only update the Dockerfile in resources/${PROJECT_NAME}/
    if [ -f "resources/$PROJECT_NAME/Dockerfile" ]; then
        sed -i "s|go-boilerplate|$PROJECT_NAME|g" "resources/$PROJECT_NAME/Dockerfile"
        sed -i "s|github.com/wasay-usmani/go-boilerplate|$MODULE_PATH|g" "resources/$PROJECT_NAME/Dockerfile"
        print_success "Updated Dockerfile in resources/$PROJECT_NAME/"
    fi
    
    print_success "Updated resources directory"
fi

# Run go mod tidy to clean up dependencies
print_info "Running go mod tidy..."
if command -v go >/dev/null 2>&1; then
    go mod tidy
    print_success "Ran go mod tidy"
else
    print_warning "Go command not found, skipping go mod tidy"
fi

echo ""
print_success "Boilerplate initialization complete!"
print_info "📁 New module path: $MODULE_PATH"
print_info "📁 Project name: $PROJECT_NAME"
echo ""
print_info "Next steps:"
echo "1. Review the changes made"
echo "2. Update any additional configuration files specific to your project"
echo "3. Run 'go mod tidy' if needed"
echo "4. Test your application with 'make build' and 'make test'" 