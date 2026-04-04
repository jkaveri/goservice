package errorcode_test

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/errorcode"
)

func TestIsErrorCode(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		code     string
		expected bool
	}{
		{
			name:     "should return true when error code matches",
			err:      errorcode.NewError("invalid_request", "invalid request"),
			code:     "invalid_request",
			expected: true,
		},
		{
			name:     "should return false when error code doesn't match",
			err:      errorcode.NewError("invalid_request", "invalid request"),
			code:     "not_found",
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			code:     "invalid_request",
			expected: false,
		},
		{
			name:     "should return false for regular error without code",
			err:      errors.New("regular error"),
			code:     "invalid_request",
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsErrorCode(tt.err, tt.code)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNotFound(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for not found error",
			err:      errorcode.NotFound("resource not found"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsNotFound(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsUnauthorized(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for unauthorized error",
			err:      errorcode.Unauthorized("unauthorized access"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsUnauthorized(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsDuplicated(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for duplicated error",
			err:      errorcode.Duplicated("duplicate resource"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsDuplicated(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInternalServer(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for internal server error",
			err:      errorcode.InternalServer("internal server error"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsInternalServer(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsInvalidRequest(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for invalid request error",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.NotFound("not found"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsInvalidRequest(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestIsNotAuthenticated(t *testing.T) {
	tests := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "should return true for not authenticated error",
			err:      errorcode.NotAuthenticated("not authenticated"),
			expected: true,
		},
		{
			name:     "should return false for other error types",
			err:      errorcode.InvalidRequest("invalid request"),
			expected: false,
		},
		{
			name:     "should return false for nil error",
			err:      nil,
			expected: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := errorcode.IsNotAuthenticated(tt.err)
			assert.Equal(t, tt.expected, result)
		})
	}
}
