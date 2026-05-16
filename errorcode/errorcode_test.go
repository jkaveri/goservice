package errorcode_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/errorcode"
)

func codedError(code, msg string) string {
	return fmt.Sprintf("[%s] %s", code, msg)
}

func TestErrorCodeConstants(t *testing.T) {
	tests := []struct {
		name     string
		expected string
		actual   string
	}{
		{
			name:     "CodeNone should be ok",
			expected: "ok",
			actual:   errorcode.CodeNone,
		},
		{
			name:     "CodeInvalidRequest should be invalid_request",
			expected: "invalid_request",
			actual:   errorcode.CodeInvalidRequest,
		},
		{
			name:     "CodeNotFound should be not_found",
			expected: "not_found",
			actual:   errorcode.CodeNotFound,
		},
		{
			name:     "CodeUnauthorized should be forbidden",
			expected: "forbidden",
			actual:   errorcode.CodeUnauthorized,
		},
		{
			name:     "CodeNotAuthenticated should be not_authenticated",
			expected: "not_authenticated",
			actual:   errorcode.CodeNotAuthenticated,
		},
		{
			name:     "CodeDuplicated should be conflict",
			expected: "conflict",
			actual:   errorcode.CodeDuplicated,
		},
		{
			name:     "CodeInternalServer should be internal_server_error",
			expected: "internal_server_error",
			actual:   errorcode.CodeInternalServer,
		},
		{
			name:     "CodeTooManyRequests should be too_many_requests",
			expected: "too_many_requests",
			actual:   errorcode.CodeTooManyRequests,
		},
		{
			name:     "CodeTimeout should be timeout",
			expected: "timeout",
			actual:   errorcode.CodeTimeout,
		},
		{
			name:     "CodeUnavailable should be unavailable",
			expected: "unavailable",
			actual:   errorcode.CodeUnavailable,
		},
		{
			name:     "CodeUnimplemented should be unimplemented",
			expected: "unimplemented",
			actual:   errorcode.CodeUnimplemented,
		},
		{
			name:     "CodeFailedPrecondition should be failed_precondition",
			expected: "failed_precondition",
			actual:   errorcode.CodeFailedPrecondition,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.actual)
		})
	}
}

func TestNewError(t *testing.T) {
	tests := []struct {
		name     string
		code     string
		message  string
		expected string
	}{
		{
			name:     "should create error with custom code and message",
			code:     "CUSTOM001",
			message:  "custom error message",
			expected: "custom error message",
		},
		{
			name:     "should create error with empty message",
			code:     "EMPTY001",
			message:  "",
			expected: "",
		},
		{
			name:     "should create error with empty code",
			code:     "",
			message:  "error with empty code",
			expected: "error with empty code",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.NewError(tt.code, tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(tt.code, tt.message), err.Error())
		})
	}
}

func TestNewErrorf(t *testing.T) {
	err := errorcode.NewErrorf("CUSTOM", "user %s not found", "alice")
	assert.Equal(t, codedError("CUSTOM", "user alice not found"), err.Error())

	// No args: format string is not passed to fmt.Sprintf; safe for arbitrary text.
	errLiteral := errorcode.NewErrorf("C", "plain literal")
	assert.Equal(t, codedError("C", "plain literal"), errLiteral.Error())

	// Dynamic text with "%" is fine when supplied as a formatted value.
	errPct := errorcode.NewErrorf("C", "%s", "100% done")
	assert.Equal(t, codedError("C", "100% done"), errPct.Error())
}

func TestNotFound(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create not found error",
			message:  "resource not found",
			expected: "resource not found",
		},
		{
			name:     "should create not found error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.NotFound(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeNotFound, tt.message), err.Error())
			assert.True(t, errorcode.IsNotFound(err))
		})
	}
}

func TestNotFoundf(t *testing.T) {
	err := errorcode.NotFoundf("widget %d missing", 42)
	assert.Equal(t, codedError(errorcode.CodeNotFound, "widget 42 missing"), err.Error())
	assert.True(t, errorcode.IsNotFound(err))
}

func TestUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create unauthorized error",
			message:  "unauthorized access",
			expected: "unauthorized access",
		},
		{
			name:     "should create unauthorized error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.Unauthorized(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeUnauthorized, tt.message), err.Error())
			assert.True(t, errorcode.IsUnauthorized(err))
		})
	}
}

func TestDuplicated(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create duplicated error",
			message:  "duplicate resource",
			expected: "duplicate resource",
		},
		{
			name:     "should create duplicated error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.Duplicated(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeDuplicated, tt.message), err.Error())
			assert.True(t, errorcode.IsDuplicated(err))
		})
	}
}

func TestInternalServer(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create internal server error",
			message:  "internal server error",
			expected: "internal server error",
		},
		{
			name:     "should create internal server error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.InternalServer(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeInternalServer, tt.message), err.Error())
			assert.True(t, errorcode.IsInternalServer(err))
		})
	}
}

func TestInvalidRequest(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create invalid request error",
			message:  "invalid request parameters",
			expected: "invalid request parameters",
		},
		{
			name:     "should create invalid request error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.InvalidRequest(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeInvalidRequest, tt.message), err.Error())
			assert.True(t, errorcode.IsInvalidRequest(err))
		})
	}
}

func TestNotAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create not authenticated error",
			message:  "user not authenticated",
			expected: "user not authenticated",
		},
		{
			name:     "should create not authenticated error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.NotAuthenticated(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeNotAuthenticated, tt.message), err.Error())
			assert.True(t, errorcode.IsNotAuthenticated(err))
		})
	}
}

func TestTooManyRequests(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create too many requests error",
			message:  "rate limit exceeded",
			expected: "rate limit exceeded",
		},
		{
			name:     "should create too many requests error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.TooManyRequests(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeTooManyRequests, tt.message), err.Error())
			assert.True(t, errorcode.IsTooManyRequests(err))
		})
	}
}

func TestTooManyRequestsf(t *testing.T) {
	err := errorcode.TooManyRequestsf("rate limit exceeded for user %s", "alice")
	assert.Equal(t, codedError(errorcode.CodeTooManyRequests, "rate limit exceeded for user alice"), err.Error())
	assert.True(t, errorcode.IsTooManyRequests(err))
}

func TestTimeout(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create timeout error",
			message:  "request timed out",
			expected: "request timed out",
		},
		{
			name:     "should create timeout error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.Timeout(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeTimeout, tt.message), err.Error())
			assert.True(t, errorcode.IsTimeout(err))
		})
	}
}

func TestTimeoutf(t *testing.T) {
	err := errorcode.Timeoutf("request timed out after %ds", 30)
	assert.Equal(t, codedError(errorcode.CodeTimeout, "request timed out after 30s"), err.Error())
	assert.True(t, errorcode.IsTimeout(err))
}

func TestUnavailable(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create unavailable error",
			message:  "service unavailable",
			expected: "service unavailable",
		},
		{
			name:     "should create unavailable error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.Unavailable(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeUnavailable, tt.message), err.Error())
			assert.True(t, errorcode.IsUnavailable(err))
		})
	}
}

func TestUnavailablef(t *testing.T) {
	err := errorcode.Unavailablef("service %s unavailable", "billing")
	assert.Equal(t, codedError(errorcode.CodeUnavailable, "service billing unavailable"), err.Error())
	assert.True(t, errorcode.IsUnavailable(err))
}

func TestUnimplemented(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create unimplemented error",
			message:  "method not implemented",
			expected: "method not implemented",
		},
		{
			name:     "should create unimplemented error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.Unimplemented(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeUnimplemented, tt.message), err.Error())
			assert.True(t, errorcode.IsUnimplemented(err))
		})
	}
}

func TestUnimplementedf(t *testing.T) {
	err := errorcode.Unimplementedf("method %s not implemented", "GetUser")
	assert.Equal(t, codedError(errorcode.CodeUnimplemented, "method GetUser not implemented"), err.Error())
	assert.True(t, errorcode.IsUnimplemented(err))
}

func TestFailedPrecondition(t *testing.T) {
	tests := []struct {
		name     string
		message  string
		expected string
	}{
		{
			name:     "should create failed precondition error",
			message:  "resource not in required state",
			expected: "resource not in required state",
		},
		{
			name:     "should create failed precondition error with empty message",
			message:  "",
			expected: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := errorcode.FailedPrecondition(tt.message)
			assert.NotNil(t, err)
			assert.Equal(t, codedError(errorcode.CodeFailedPrecondition, tt.message), err.Error())
			assert.True(t, errorcode.IsFailedPrecondition(err))
		})
	}
}

func TestFailedPreconditionf(t *testing.T) {
	err := errorcode.FailedPreconditionf("order %d already shipped", 7)
	assert.Equal(t, codedError(errorcode.CodeFailedPrecondition, "order 7 already shipped"), err.Error())
	assert.True(t, errorcode.IsFailedPrecondition(err))
}
