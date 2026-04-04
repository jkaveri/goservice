# gRPC Proto Generation (v2)

This directory contains the protobuf definitions and configuration for generating Go gRPC code using buf v2.

## Files

- `service/v1/echo.proto` - The protobuf service definition (v2 structure)
- `buf.yaml` - Buf workspace configuration (v2)
- `buf.gen.yaml` - Buf generation configuration (v2)
- `generate.sh` - Script to generate Go code

## Directory Structure

Following buf v2 conventions:
```
proto/
├── service/
│   └── v1/
│       └── echo.proto
├── buf.yaml
├── buf.gen.yaml
└── gen/
    └── service/
        └── v1/
            ├── echo.pb.go
            └── echo_grpc.pb.go
```

## Generating Code

To generate Go code from the proto files:

```bash
# Option 1: Use the script
./generate.sh

# Option 2: Use buf directly
buf generate
```

## Generated Files

The generation creates the following files in the `gen/service/v1/` directory:

- `echo.pb.go` - Generated protobuf message types
- `echo_grpc.pb.go` - Generated gRPC service interfaces

## Configuration

The `buf.gen.yaml` file is configured to:

- Generate Go code using `protoc-gen-go` and `protoc-gen-go-grpc`
- Output files to the `gen/` directory
- Use source-relative import paths
- Set the Go package prefix to `github.com/jkaveri/goservice/examples/echogrpc/proto/gen`

## Usage in Go

To use the generated code in your Go service:

```go
import (
    "github.com/jkaveri/goservice/examples/echogrpc/proto/gen/service/v1"
)

// The generated service interface will be available as:
// v1.EchoServiceServer
// v1.EchoServiceClient
```
