#!/bin/bash

# Create App Script
# Usage: ./resources/scripts/create-app.sh <service-name> <app-name>

set -e

if [ $# -ne 2 ]; then
    echo "Usage: $0 <service-name> <app-name>"
    echo "Example: $0 user-service auth"
    exit 1
fi

SERVICE_NAME="$1"
APP_NAME="$2"

# Convert app name to different formats
APP_NAME_TITLE=$(echo "$APP_NAME" | sed 's/-\([a-z]\)/\U\1/g' | sed 's/^\([a-z]\)/\U\1/')
APP_NAME_CAMEL=$(echo "$APP_NAME" | sed 's/-\([a-z]\)/\U\1/g')
APP_NAME_PKG=$(echo "$APP_NAME" | tr '[:upper:]' '[:lower:]' | sed 's/-//g' | sed 's/_//g')

echo "Creating app module: $APP_NAME"
echo "Service: $SERVICE_NAME"
echo "App package name: $APP_NAME_PKG"
echo ""

# Check if service exists
if [ ! -d "internal/$SERVICE_NAME" ]; then
    echo "Error: Service '$SERVICE_NAME' does not exist in internal/"
    exit 1
fi

# Check if app already exists
if [ -d "internal/$SERVICE_NAME/app/$APP_NAME" ]; then
    echo "Error: App '$APP_NAME' already exists in internal/$SERVICE_NAME/app/"
    exit 1
fi

# Get current module path from go.mod
MODULE_PATH=$(grep '^module ' go.mod | cut -d' ' -f2)

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
    
    # Process template with sed
    sed -e "s/{{\.ServiceName}}/$SERVICE_NAME/g" \
        -e "s/{{\.AppName}}/$APP_NAME/g" \
        -e "s/{{\.AppNameTitle}}/$APP_NAME_TITLE/g" \
        -e "s/{{\.AppNameCamel}}/$APP_NAME_CAMEL/g" \
        -e "s/{{\.AppNamePkg}}/$APP_NAME_PKG/g" \
        -e "s|{{\.ModulePath}}|$MODULE_PATH|g" \
        "$template_file" > "$output_file"
    
    echo "Generated: $output_file"
}

# Define template mappings for app files
declare -A templates=(
    ["templates/internal/app/module/app.go.tmpl"]="internal/$SERVICE_NAME/app/$APP_NAME/app.go"
    ["templates/internal/app/module/cmds.go.tmpl"]="internal/$SERVICE_NAME/app/$APP_NAME/cmds.go"
    ["templates/internal/app/module/qrys.go.tmpl"]="internal/$SERVICE_NAME/app/$APP_NAME/qrys.go"
)

# Process each template
for template_file in "${!templates[@]}"; do
    output_file="${templates[$template_file]}"
    process_template "$template_file" "$output_file"
done

echo ""
echo "✅ Successfully created app module '$APP_NAME' in service '$SERVICE_NAME'"
echo "📁 Created files:"
echo "  - internal/$SERVICE_NAME/app/$APP_NAME/app.go"
echo "  - internal/$SERVICE_NAME/app/$APP_NAME/cmds.go"
echo "  - internal/$SERVICE_NAME/app/$APP_NAME/qrys.go"
echo ""
echo "Next steps:"
echo "1. Update the service's app/module.go to include the new app module"
echo "2. Add your business logic to the generated files"
echo "3. Update any dependencies as needed" 