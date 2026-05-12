package errorcode_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/errorcode"
)

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
			assert.Equal(t, tt.expected, err.Error())
		})
	}
}

func TestNewErrorf(t *testing.T) {
	err := errorcode.NewErrorf("CUSTOM", "user %s not found", "alice")
	assert.Equal(t, "user alice not found", err.Error())

	// No args: format string is not passed to fmt.Sprintf; safe for arbitrary text.
	errLiteral := errorcode.NewErrorf("C", "plain literal")
	assert.Equal(t, "plain literal", errLiteral.Error())

	// Dynamic text with "%" is fine when supplied as a formatted value.
	errPct := errorcode.NewErrorf("C", "%s", "100% done")
	assert.Equal(t, "100% done", errPct.Error())
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
			assert.Equal(t, tt.expected, err.Error())
			assert.True(t, errorcode.IsNotFound(err))
		})
	}
}

func TestNotFoundf(t *testing.T) {
	err := errorcode.NotFoundf("widget %d missing", 42)
	assert.Equal(t, "widget 42 missing", err.Error())
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
			assert.Equal(t, tt.expected, err.Error())
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
			assert.Equal(t, tt.expected, err.Error())
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
			assert.Equal(t, tt.expected, err.Error())
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
			assert.Equal(t, tt.expected, err.Error())
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
			assert.Equal(t, tt.expected, err.Error())
			assert.True(t, errorcode.IsNotAuthenticated(err))
		})
	}
}
