#!/bin/bash

# Create Service Script
# Usage: ./scripts/create-service.sh <service-name> <module-path> <output-dir>

set -e

if [ $# -ne 3 ]; then
    echo "Usage: $0 <service-name> <module-path> <output-dir>"
    echo "Example: $0 user-service github.com/mycompany/myapp ."
    exit 1
fi

SERVICE_NAME="$1"
MODULE_PATH="$2"
OUTPUT_DIR="$3"

# Convert service name to different formats
SERVICE_NAME_TITLE=$(echo "$SERVICE_NAME" | sed 's/-\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')
SERVICE_NAME_CAMEL=$(echo "$SERVICE_NAME" | sed 's/-\([a-z]\)/\U\1/g')
SERVICE_NAME_PKG=$(echo "$SERVICE_NAME" | tr '[:upper:]' '[:lower:]' | sed 's/-//g' | sed 's/_//g')

echo "Creating microservice: $SERVICE_NAME"
echo "Module path: $MODULE_PATH"
echo "Output directory: $OUTPUT_DIR"
echo "Package name: $SERVICE_NAME_PKG"
echo ""

# Function to process a template file
process_template() {
    local template_file="$1"
    local output_file="$2"
    
    if [ ! -f "$template_file" ]; then
        echo "Warning: Template file $template_file not found, skipping..."
        return
    fi
    
    # Create output directory
    mkdir -p "$(dirname "$output_file")"
    
    # Process template with sed (escape dots in module path)
    sed -e "s/{{\.ServiceName}}/$SERVICE_NAME/g" \
        -e "s/{{\.ServiceNameTitle}}/$SERVICE_NAME_TITLE/g" \
        -e "s/{{\.ServiceNameCamel}}/$SERVICE_NAME_CAMEL/g" \
        -e "s/{{\.ServiceNamePkg}}/$SERVICE_NAME_PKG/g" \
        -e "s|{{\.ModulePath}}|$MODULE_PATH|g" \
        "$template_file" > "$output_file"
    
    echo "Generated: $output_file"
}

# Define template mappings
declare -A templates=(
    ["templates/cmd/main.go.tmpl"]="$OUTPUT_DIR/cmd/$SERVICE_NAME/main.go"
    ["templates/internal/app/module.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/app/module.go"
    ["templates/internal/app/service/app.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/app/$SERVICE_NAME/app.go"
    ["templates/internal/app/service/cmds.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/app/$SERVICE_NAME/cmds.go"
    ["templates/internal/app/service/qrys.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/app/$SERVICE_NAME/qrys.go"
    ["templates/internal/config/config.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/config/config.go"
    ["templates/internal/config/config.toml.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/config/config.toml"
    ["templates/internal/repository/module.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/repository/module.go"
    ["templates/internal/server/http/api.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/server/http/api.go"
    ["templates/internal/server/http/routes.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/server/http/routes.go"
    ["templates/internal/server/http/health.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/server/http/health.go"
    ["templates/internal/server/rpc/rpc.go.tmpl"]="$OUTPUT_DIR/internal/$SERVICE_NAME/server/rpc/rpc.go"
    ["templates/resources/Dockerfile.tmpl"]="$OUTPUT_DIR/resources/$SERVICE_NAME/Dockerfile"
)

# Process each template
for template_file in "${!templates[@]}"; do
    output_file="${templates[$template_file]}"
    process_template "$template_file" "$output_file"
done

echo ""
echo "✅ Successfully generated microservice '$SERVICE_NAME' in $OUTPUT_DIR"
echo "📁 Created directories:"
echo "  - cmd/$SERVICE_NAME"
echo "  - internal/$SERVICE_NAME"
echo "  - resources/$SERVICE_NAME"
echo ""
echo "Next steps:"
echo "1. Review the generated code in the new directories"
echo "2. Update service-specific configuration and business logic"
echo "3. Update any service-specific dependencies in go.mod"
echo "4. Test the new service with 'go build ./cmd/$SERVICE_NAME'"
echo "5. Add any service-specific environment variables or config" 