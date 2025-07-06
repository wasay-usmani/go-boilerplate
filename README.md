# Go Boilerplate Template

A comprehensive Go boilerplate template with automatic initialization for building scalable microservices.

## 🚀 Quick Start

### Option 1: Use This Template (Recommended)

1. **Click "Use this template"** on GitHub
2. **Create your repository** with your desired name
3. **Automatic initialization** will run and configure everything for you
4. **Check the created issue** for next steps

The template automatically:
- ✅ Updates all import paths to match your repository
- ✅ Renames directories to match your project name
- ✅ Updates configuration files (Dockerfile, database configs, etc.)
- ✅ Runs tests to ensure everything works
- ✅ Creates a helpful issue with next steps

### Option 2: Manual Initialization

If you need to manually initialize or reinitialize:

```bash
# Clone the repository
git clone <your-repo-url>
cd <your-repo-name>

# Initialize with your module path
make init MODULE_PATH=github.com/your-org/your-project

# Build and test
make build
make test
```

## 📁 Project Structure

```
├── cmd/
│   ├── your-project/          # Main application entry point
│   └── microservice-name/     # Additional microservices
├── internal/
│   ├── your-project/          # Private application code
│   │   ├── app/               # Application logic
│   │   │   ├── module.go      # Main app module
│   │   │   ├── your-project/  # Default service app
│   │   │   └── app-name/      # Additional app modules
│   │   ├── config/            # Configuration management
│   │   ├── repository/        # Data access layer
│   │   └── server/            # HTTP/RPC server
│   └── microservice-name/     # Additional microservices
├── pkg/                       # Public packages
│   ├── cachekit/              # Caching utilities
│   ├── configkit/             # Configuration utilities
│   ├── dbkit/                 # Database utilities
│   ├── httpx/                 # HTTP utilities
│   ├── logkit/                # Logging utilities
│   └── utils/                 # General utilities
├── resources/
│   ├── scripts/               # Build and generation scripts
│   ├── your-project/          # Configuration files
│   │   ├── Dockerfile         # Container configuration
│   │   └── migrations/        # Database migrations
│   └── microservice-name/     # Additional microservice resources
└── Makefile                   # Build and development commands
```

## 🛠️ Available Commands

```bash
# Initialize boilerplate (updates all paths and configs)
make init MODULE_PATH=github.com/your-org/your-project

# Create a new microservice from boilerplate
make create-service SERVICE_NAME=my-service

# Add a new app module to existing service
make create-app SERVICE_NAME=my-service APP_NAME=my-app

# Build all services
make build

# Build a specific service
make build-service SERVICE_NAME=my-service

# Run a specific service
make run-service SERVICE_NAME=my-service

# Run tests
make test

# Run tests with coverage
make test-coverage

# Clean build artifacts
make clean

# Install dependencies
make deps

# Format code
make fmt

# Run linter
make lint

# Run all checks (format, lint, test)
make check
```

## 🚀 Running Services

### Running Individual Services

```bash
# Run the main service
make run-service SERVICE_NAME=go-boilerplate

# Run a specific microservice
make run-service SERVICE_NAME=user-service

# Run with custom environment
LISTEN_PORT=3000 make run-service SERVICE_NAME=go-boilerplate
```

### Running Multiple Services

For development with multiple services:

```bash
# Terminal 1: Run main service
make run-service SERVICE_NAME=go-boilerplate

# Terminal 2: Run user service
make run-service SERVICE_NAME=user-service

# Terminal 3: Run auth service
make run-service SERVICE_NAME=auth-service
```

## 🏗️ Microservice Development

### Creating New Microservices

The boilerplate includes powerful tools for generating new microservices:

```bash
# Create a new microservice
make create-service SERVICE_NAME=user-service

# This creates:
# - cmd/user-service/main.go
# - internal/user-service/ (complete service structure)
# - resources/user-service/Dockerfile
```

### Adding App Modules to Services

You can add multiple app modules within a single microservice:

```bash
# Add a new app module to an existing service
make create-app SERVICE_NAME=user-service APP_NAME=auth

# This creates:
# - internal/user-service/app/auth/app.go
# - internal/user-service/app/auth/cmds.go
# - internal/user-service/app/auth/qrys.go
```

### Example: Building a User Management System

```bash
# Create the user service
make create-service SERVICE_NAME=user-service

# Add authentication module
make create-app SERVICE_NAME=user-service APP_NAME=auth

# Add profile management module
make create-app SERVICE_NAME=user-service APP_NAME=profile

# Add notification module
make create-app SERVICE_NAME=user-service APP_NAME=notifications

# Build the services
make build-service SERVICE_NAME=user-service

# Run the services (in separate terminals)
make run-service SERVICE_NAME=user-service
make run-service SERVICE_NAME=auth-service
```

### Integration Steps

After creating app modules, you need to integrate them:

1. **Update the service's main app module** (`internal/user-service/app/module.go`):
```go
type Module struct {
    UserService userservice.App
    Auth        auth.App
    Profile     profile.App
    // ... other modules
}
```

2. **Add your business logic** to the generated files
3. **Update dependencies** in `go.mod` if needed

### Template System

The boilerplate uses a powerful template system located in `templates/`:

- **Service Templates**: Generate complete microservices
- **App Module Templates**: Generate app modules within services
- **Shell Scripts**: `resources/scripts/` contains generation scripts
- **Customizable**: All templates can be modified for your specific needs

Generated files automatically have:
- ✅ Valid Go package names (no hyphens/underscores)
- ✅ Correct import paths
- ✅ Proper module structure
- ✅ Docker configuration

## 🔧 Configuration

Update the database configuration in `resources/your-project/migrations/`:

```yaml
# resources/your-project/migrations/app/dbconfig.yml
development:
  dialect: mysql
  datasource: root:root@tcp(localhost:1444)/your_project?parseTime=true
  dir: resources/migrations/your-project/app/migrations
  table: migrations_app_your_project
```

### Environment Variables

Create a `.env` file with your configuration:

```env
# Server Configuration
LISTEN_HOST=0.0.0.0
LISTEN_PORT=8080
RPC_LISTEN_PORT=9090

# Database Configuration
WRITE_DB_URL=mysql://user:pass@localhost:3306/dbname
READ_DB_URL=mysql://user:pass@localhost:3306/dbname

# Cache Configuration
CACHE_URL=redis://localhost:6379
```

## 🐳 Docker

Build and run with Docker:

```bash
# Build the image
docker build -f resources/your-project/Dockerfile -t your-project .

# Run the container
docker run -p 8080:8080 your-project
```

## 📝 Features

- **🔄 Automatic Initialization**: Template automatically configures everything
- **🏗️ Clean Architecture**: Well-structured, maintainable codebase
- **🚀 Microservice Generation**: Create new services and app modules with simple commands
- **📊 Database Support**: MySQL/PostgreSQL with migrations
- **🔄 Caching**: Redis integration with Valkey
- **📝 Logging**: Structured logging with Zerolog
- **🔧 Configuration**: Environment-based configuration with Viper
- **🧪 Testing**: Comprehensive test setup
- **🐳 Docker**: Production-ready containerization
- **📈 Monitoring**: Built-in health checks and metrics
- **🔒 Security**: Input validation and secure defaults

## 🤝 Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests
5. Submit a pull request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## 🆘 Support

If you encounter any issues:

1. Check the [Issues](../../issues) page
2. Create a new issue with detailed information
3. Join our community discussions

---

**Happy coding! 🎯** 