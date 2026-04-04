# GoService Framework

A comprehensive Go service framework that provides a robust foundation for building microservices with gRPC and HTTP support, including built-in interceptors, health checks, and development tools.

## Features

- **Multi-Protocol Support**: Built-in support for both gRPC and HTTP servers
- **gRPC Gateway Integration**: Automatic HTTP-to-gRPC proxy with customizable options
- **Built-in Interceptors**: Request ID tracking, logging, validation, and error wrapping
- **Health Check Server**: Dedicated health check endpoint for monitoring
- **Graceful Shutdown**: Proper signal handling and resource cleanup
- **Configuration Management**: Flexible configuration system with environment support
- **Development Tools**: CLI commands for linting, testing, and code formatting
- **Mock Support**: Comprehensive mock interfaces for testing

## Quick Start

### Basic gRPC Service

```go
package main

import (
    "context"
    "github.com/jkaveri/goservice"
    "google.golang.org/grpc"
)

type MyService struct {
    // Your service implementation
}

func (s *MyService) RegisterGRPC(ctx context.Context, srv *grpc.Server) {
    // Register your gRPC services
}

func (s *MyService) GRPCServerOptions() []grpc.ServerOption {
    return []grpc.ServerOption{}
}

func (s *MyService) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
    return goservice.DefaultInterceptors()
}

func New(ctx context.Context) (goservice.Service, error) {
    return &MyService{}, nil
}

func main() {
    goservice.Run(New)
}
```

### gRPC + HTTP Gateway Service

```go
package main

import (
    "context"
    "net/http"
    "github.com/jkaveri/goservice"
    "google.golang.org/grpc"
)

type MyService struct {
    // Your service implementation
}

// Implement gRPC interface
func (s *MyService) RegisterGRPC(ctx context.Context, srv *grpc.Server) {
    // Register your gRPC services
}

func (s *MyService) GRPCServerOptions() []grpc.ServerOption {
    return []grpc.ServerOption{}
}

func (s *MyService) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
    return goservice.DefaultInterceptors()
}

// Implement HTTP interface
func (s *MyService) RegisterHTTP(ctx context.Context, srv *http.Server) http.Handler {
    // Return your HTTP handler
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte("Hello, World!"))
    })
}

func New(ctx context.Context) (goservice.Service, error) {
    return &MyService{}, nil
}

func main() {
    goservice.Run(New)
}
```

## Configuration

The framework uses a flexible configuration system. Default configuration:

```go
type Config struct {
    Debug         bool   // Enable debug mode
    DeploymentEnv string // Deployment environment (dev, staging, prod)

    GRPCServer ServerConfig // gRPC server configuration
    HTTPServer ServerConfig // HTTP server configuration
    HealthServer ServerConfig // Health check server configuration
    Log LogConfig // Logging configuration (uses golog)
}

type LogConfig struct {
    Format    string // "text" or "json"
    Level     string // "debug", "info", "warn", "error"
    AddSource bool   // Add file:line to each log record
}

type ServerConfig struct {
    Port int // Server port
}
```

Default ports:
- gRPC Server: 9000
- HTTP Server: 8080
- Health Check: 8081

Configuration can be loaded from environment variables or configuration files using the `goconfig` package.

## Built-in Interceptors

The framework provides several built-in gRPC interceptors:

- **Request ID**: Automatically generates and tracks request IDs
- **Logging**: Comprehensive request/response logging
- **Validation**: Automatic request validation using protobuf validation rules
- **Error Wrapping**: Structured error handling and response formatting

### Using Default Interceptors

```go
func (s *MyService) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
    return goservice.DefaultInterceptors()
}
```

### Custom Interceptors

```go
func (s *MyService) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
    interceptors := goservice.DefaultInterceptors()

    // Add your custom interceptors
    interceptors = append(interceptors, myCustomInterceptor)

    return interceptors
}
```

## gRPC Gateway

For HTTP-to-gRPC proxy functionality:

```go
import (
    "github.com/jkaveri/goservice/grpc/gateway"
)

func (s *MyService) RegisterHTTP(ctx context.Context, srv *http.Server) http.Handler {
    cfg := goservice.GetCurrentConfig()

    return grcpgateway.CreateHandler(
        ctx,
        s,
        cfg.GRPCServer.Port,
        srv,
        yourpb.RegisterYourServiceHandler,
    )
}
```

## Lifecycle Hooks

Implement lifecycle interfaces for custom startup/shutdown logic:

```go
type MyService struct {
    // Your service
}

// OnStartListener - called when service starts
func (s *MyService) OnStart(ctx context.Context) error {
    // Initialize resources
    return nil
}

// OnStopListener - called when service stops
func (s *MyService) OnStop() error {
    // Cleanup resources
    return nil
}
```

## Development Tools

The framework includes a CLI tool for development tasks:

```bash
# Run linting
goservice lint

# Format code
goservice pretty

# Run tests with coverage
goservice test

# Generate coverage report
goservice coverage
```

## Examples

### Echo gRPC Service

A simple echo service demonstrating basic gRPC functionality:

```bash
cd examples/echogrpc
go run main.go
```

### gRPC Gateway Service

A service with both gRPC and HTTP gateway support:

```bash
cd examples/grpcgateway
go run main.go
```

## Project Structure

```
goservice/
├── cmd/goservice/          # CLI tool for development
├── grpc/                   # gRPC utilities and interceptors
│   ├── gateway/           # gRPC gateway functionality
│   └── interceptors/      # Built-in interceptors
├── examples/             # Example implementations
│   ├── echogrpc/         # Basic gRPC service
│   └── grpcgateway/      # gRPC + HTTP gateway
├── mock/                  # Mock interfaces for testing
├── check/                 # Error checking utilities
├── clock/                 # Time utilities
├── idgen/                 # ID generation utilities
└── validate/              # Validation utilities
```

## Dependencies

- **gRPC**: Google's gRPC framework
- **grpc-gateway**: HTTP-to-gRPC proxy
- **golog**: Structured logging (git.toolsfdg.net/fe/golog)
- **goconfig**: Configuration management
- **validator**: Request validation
- **cobra**: CLI framework

## Testing

The framework provides comprehensive mock interfaces for testing:

```go
import (
    "github.com/jkaveri/goservice/mock"
)

// Use mock services in tests
mockService := &mock.MockService{}
```

## Contributing

1. Fork the repository
2. Create a feature branch
3. Make your changes
4. Add tests for new functionality
5. Submit a pull request

## License

This project is licensed under the MIT License.
