// Package errorcode provides standardized error codes and error creation
// utilities
// for consistent error handling across the application.
package errorcode

import errors "github.com/jkaveri/goservice/errors"

// Error code constants use snake_case identifiers. HTTP and gRPC mapping lives
// in grpc/interceptors/wraperror.
const (
	// CodeNone represents no error (success state)
	CodeNone = "ok"

	// CodeInvalidRequest represents invalid request errors (400 Bad Request)
	CodeInvalidRequest = "invalid_request"

	// CodeNotFound represents resource not found errors (404 Not Found)
	CodeNotFound = "not_found"

	// CodeUnauthorized represents forbidden access (403 Forbidden):
	// authenticated
	// but not permitted.
	CodeUnauthorized = "forbidden"

	// CodeNotAuthenticated represents authentication required (401
	// Unauthorized)
	CodeNotAuthenticated = "not_authenticated"

	// CodeDuplicated represents duplicate resource errors (409 Conflict)
	CodeDuplicated = "conflict"

	// CodeInternalServer represents internal server errors (500 Internal Server
	// Error)
	CodeInternalServer = "internal_server_error"
)

// NewError creates a new error with the specified code and message.
// It wraps the error with the provided code using the underlying errors
// package.
func NewError(code string, message string) error {
	return errors.WithCode(errors.New(message), code)
}

// NotFound creates a new error with the CodeNotFound code and the specified
// message.
// Use this for 404 Not Found scenarios.
func NotFound(message string) error {
	return NewError(CodeNotFound, message)
}

// Unauthorized creates a new error with the CodeUnauthorized code and the
// specified message. Use this for 403 Forbidden scenarios where the caller is
// authenticated but not permitted.
func Unauthorized(message string) error {
	return NewError(CodeUnauthorized, message)
}

// Duplicated creates a new error with the CodeDuplicated code and the specified
// message.
// Use this for 409 Conflict scenarios where a resource already exists.
func Duplicated(message string) error {
	return NewError(CodeDuplicated, message)
}

// InternalServer creates a new error with the CodeInternalServer code and the
// specified message.
// Use this for 500 Internal Server Error scenarios.
func InternalServer(message string) error {
	return NewError(CodeInternalServer, message)
}

// InvalidRequest creates a new error with the CodeInvalidRequest code and the
// specified message. Use this for 400 Bad Request scenarios where the request
// is malformed or invalid.
func InvalidRequest(message string) error {
	return NewError(CodeInvalidRequest, message)
}

// NotAuthenticated creates a new error with the CodeNotAuthenticated code and
// the specified message. Use this for 401 Unauthorized scenarios where
// credentials are missing or invalid.
func NotAuthenticated(message string) error {
	return NewError(CodeNotAuthenticated, message)
}

// WithInvalidRequest wraps the error with the CodeInvalidRequest code.
func WithInvalidRequest(err error) error {
	return errors.WithCode(err, CodeInvalidRequest)
}

// WithNotFound wraps the error with the CodeNotFound code.
func WithNotFound(err error) error {
	return errors.WithCode(err, CodeNotFound)
}

// WithUnauthorized wraps the error with the CodeUnauthorized code.
func WithUnauthorized(err error) error {
	return errors.WithCode(err, CodeUnauthorized)
}

// WithDuplicated wraps the error with the CodeDuplicated code.
func WithDuplicated(err error) error {
	return errors.WithCode(err, CodeDuplicated)
}

// WithInternalServer wraps the error with the CodeInternalServer code.
func WithInternalServer(err error) error {
	return errors.WithCode(err, CodeInternalServer)
}
