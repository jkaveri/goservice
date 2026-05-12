// Package errorcode provides standardized error codes and error creation
// utilities
// for consistent error handling across the application.
package errorcode

import (
	"fmt"

	errors "github.com/jkaveri/goservice/errors"
)

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

func formatMessage(format string, args ...any) string {
	if len(args) == 0 {
		return format
	}

	return fmt.Sprintf(format, args...)
}

func newCoded(code, msg string) error {
	return errors.WithCode(errors.New(msg), code)
}

func newCodedf(code, format string, args ...any) error {
	return errors.WithCode(errors.New(formatMessage(format, args...)), code)
}

// NewError creates a new error with the specified code and message.
// It wraps the error with the provided code using the underlying errors
// package.
func NewError(code string, message string) error {
	return newCoded(code, message)
}

// NewErrorf is like [NewError] but formats the message with fmt.Sprintf when
// args is non-empty; with no args, format is used as the literal message (so
// "%" is not treated as a verb).
func NewErrorf(code string, format string, args ...any) error {
	return newCodedf(code, format, args...)
}

// NotFound creates a new error with the CodeNotFound code and the specified
// message.
// Use this for 404 Not Found scenarios.
func NotFound(message string) error {
	return newCoded(CodeNotFound, message)
}

// NotFoundf is like [NotFound] but supports fmt-style formatting; see
// [NewErrorf].
func NotFoundf(format string, args ...any) error {
	return newCodedf(CodeNotFound, format, args...)
}

// Unauthorized creates a new error with the CodeUnauthorized code and the
// specified message. Use this for 403 Forbidden scenarios where the caller is
// authenticated but not permitted.
func Unauthorized(message string) error {
	return newCoded(CodeUnauthorized, message)
}

// Unauthorizedf is like [Unauthorized] but supports fmt-style formatting; see
// [NewErrorf].
func Unauthorizedf(format string, args ...any) error {
	return newCodedf(CodeUnauthorized, format, args...)
}

// Duplicated creates a new error with the CodeDuplicated code and the specified
// message.
// Use this for 409 Conflict scenarios where a resource already exists.
func Duplicated(message string) error {
	return newCoded(CodeDuplicated, message)
}

// Duplicatedf is like [Duplicated] but supports fmt-style formatting; see
// [NewErrorf].
func Duplicatedf(format string, args ...any) error {
	return newCodedf(CodeDuplicated, format, args...)
}

// InternalServer creates a new error with the CodeInternalServer code and the
// specified message.
// Use this for 500 Internal Server Error scenarios.
func InternalServer(message string) error {
	return newCoded(CodeInternalServer, message)
}

// InternalServerf is like [InternalServer] but supports fmt-style formatting;
// see [NewErrorf].
func InternalServerf(format string, args ...any) error {
	return newCodedf(CodeInternalServer, format, args...)
}

// InvalidRequest creates a new error with the CodeInvalidRequest code and the
// specified message. Use this for 400 Bad Request scenarios where the request
// is malformed or invalid.
func InvalidRequest(message string) error {
	return newCoded(CodeInvalidRequest, message)
}

// InvalidRequestf is like [InvalidRequest] but supports fmt-style formatting;
// see [NewErrorf].
func InvalidRequestf(format string, args ...any) error {
	return newCodedf(CodeInvalidRequest, format, args...)
}

// NotAuthenticated creates a new error with the CodeNotAuthenticated code and
// the specified message. Use this for 401 Unauthorized scenarios where
// credentials are missing or invalid.
func NotAuthenticated(message string) error {
	return newCoded(CodeNotAuthenticated, message)
}

// NotAuthenticatedf is like [NotAuthenticated] but supports fmt-style
// formatting; see [NewErrorf].
func NotAuthenticatedf(format string, args ...any) error {
	return newCodedf(CodeNotAuthenticated, format, args...)
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
