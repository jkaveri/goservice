# Error Code Package

The `errorcode` package provides standardized error codes and utilities for consistent error handling across the application. It integrates with the `errors` package to provide structured error handling with error codes.

## Features

- **Standardized Error Codes**: Predefined error codes for common scenarios
- **Error Creation Utilities**: Helper functions to create typed errors
- **Error Checking Functions**: Utilities to check error types
- **HTTP Status Code Mapping**: Error codes map to appropriate HTTP status codes

## Error Codes

| Code     | Constant                 | HTTP Status | Description                |
| -------- | ------------------------ | ----------- | -------------------------- |
| `000000` | `CodeNone`               | N/A         | No error (success state)   |
| `300000` | `CodeInvalidRequest`     | 400         | Invalid request parameters |
| `300001` | `CodeInternalServer`     | 500         | Internal server error      |
| `300003` | `CodeDuplicated`         | 409         | Duplicate resource         |
| `330001` | `CodeNotAuthenticated`   | 401         | Authentication required    |
| `340001` | `CodeUnauthorized`       | 401         | Unauthorized access        |
| `350002` | `CodeNotFound`           | 404         | Resource not found         |
| N/A      | `CodeTooManyRequests`    | 429         | Rate limit exceeded        |
| N/A      | `CodeTimeout`            | 504         | Deadline exceeded          |
| N/A      | `CodeUnavailable`        | 503         | Service unavailable        |
| N/A      | `CodeUnimplemented`      | 501         | Operation not implemented  |
| N/A      | `CodeFailedPrecondition` | 412         | Failed precondition        |

## Usage

### Creating Errors

```go
import "github.com/jkaveri/goservice/errorcode"

// Create a custom error with specific code
err := errorcode.NewError("CUSTOM001", "custom error message")

// Create typed errors
notFoundErr := errorcode.NotFound("user not found")
unauthorizedErr := errorcode.Unauthorized("access denied")
duplicatedErr := errorcode.Duplicated("email already exists")
internalErr := errorcode.InternalServer("database connection failed")
invalidReqErr := errorcode.InvalidRequest("missing required field")
notAuthErr := errorcode.NotAuthenticated("login required")
rateLimitErr := errorcode.TooManyRequests("rate limit exceeded")
timeoutErr := errorcode.Timeout("request timed out")
unavailableErr := errorcode.Unavailable("service unavailable")
unimplementedErr := errorcode.Unimplemented("method not implemented")
preconditionErr := errorcode.FailedPrecondition("resource not in required state")
```

### Checking Error Types

```go
import "github.com/jkaveri/goservice/errorcode"

// Check specific error code
if errorcode.IsErrorCode(err, "300000") {
    // Handle invalid request
}

// Check error types
if errorcode.IsNotFound(err) {
    // Handle not found error
}

if errorcode.IsUnauthorized(err) {
    // Handle unauthorized error
}

if errorcode.IsDuplicated(err) {
    // Handle duplicate resource error
}

if errorcode.IsInternalServer(err) {
    // Handle internal server error
}

if errorcode.IsInvalidRequest(err) {
    // Handle invalid request error
}

if errorcode.IsNotAuthenticated(err) {
    // Handle authentication error
}

if errorcode.IsTooManyRequests(err) {
    // Handle rate limit error
}

if errorcode.IsTimeout(err) {
    // Handle deadline-exceeded error
}

if errorcode.IsUnavailable(err) {
    // Handle service unavailable error
}

if errorcode.IsUnimplemented(err) {
    // Handle unimplemented method error
}

if errorcode.IsFailedPrecondition(err) {
    // Handle failed precondition error
}
```

### Example: HTTP Handler

```go
func handleGetUser(w http.ResponseWriter, r *http.Request) {
    user, err := getUserByID(r.URL.Query().Get("id"))
    if err != nil {
        if errorcode.IsNotFound(err) {
            http.Error(w, err.Error(), http.StatusNotFound)
            return
        }
        if errorcode.IsUnauthorized(err) {
            http.Error(w, err.Error(), http.StatusUnauthorized)
            return
        }
        // Default to internal server error
        http.Error(w, "Internal server error", http.StatusInternalServerError)
        return
    }

    json.NewEncoder(w).Encode(user)
}

func getUserByID(id string) (*User, error) {
	if id == "" {
		return nil, errorcode.InvalidRequest("user ID is required")
	}

	user, err := db.GetUser(id)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errorcode.NotFound("user not found")
		}
		return nil, errorcode.InternalServer("database error")
	}

	return user, nil
}
```

### Example: gRPC Service

```go
func (s *UserService) GetUser(ctx context.Context, req *pb.GetUserRequest) (*pb.GetUserResponse, error) {
	if req.UserId == "" {
		return nil, errorcode.InvalidRequest("user ID is required")
	}

	user, err := s.repo.GetUser(req.UserId)
	if err != nil {
		if errorcode.IsNotFound(err) {
			return nil, status.Error(codes.NotFound, err.Error())
		}
		return nil, status.Error(codes.Internal, "internal server error")
	}

	return &pb.GetUserResponse{User: user}, nil
}
```

## Testing

The package includes comprehensive tests for all functions and constants:

```bash
# Run all tests
go test ./errorcode

# Run tests with coverage
go test -cover ./errorcode

# Run tests with verbose output
go test -v ./errorcode
```

## Integration with errors

This package is built on top of the `errors` package and provides a higher-level interface for common error scenarios. The underlying `errors` package handles the actual error code extraction and management.

## Best Practices

1. **Use Typed Errors**: Prefer the typed error creation functions over `NewError` for common scenarios
2. **Check Error Types**: Use the `Is*` functions to check error types rather than comparing error codes directly
3. **Consistent Error Messages**: Provide clear, user-friendly error messages
4. **HTTP Status Mapping**: Map error codes to appropriate HTTP status codes in your handlers
5. **Error Logging**: Log errors with their codes for debugging and monitoring

## Contributing

When adding new error codes:

1. Add the constant to the `const` block in `errorcode.go`
2. Add a creation function if it's a common scenario
3. Add a checking function in `check.go`
4. Add comprehensive tests
5. Update this README with the new error code
