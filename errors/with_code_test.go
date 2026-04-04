package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"

	errors "github.com/jkaveri/goservice/errors"
)

func Test_WithCode(t *testing.T) {
	err := errors.New("test")

	errWithCode := errors.WithCode(err, "T123")

	assert.Nil(t, errors.WithCode(nil, "T123"))

	codeContainer, ok := errWithCode.(errors.CodeError)
	assert.True(t, ok)
	assert.Equal(t, "T123", codeContainer.Code())
	assert.Equal(t, err, errors.Cause(errWithCode))
	assert.Equal(t, err, errors.Unwrap(errWithCode))
	assert.Equal(t, "test", errWithCode.Error())
	assert.Contains(t, fmt.Sprintf("%+v", errWithCode), "[T123] test")
	assert.Contains(t, fmt.Sprintf("%s", errWithCode), "test")
	assert.Contains(t, fmt.Sprintf("%q", errWithCode), "test")
}

func Test_GetCode(t *testing.T) {
	assert.Equal(t,
		"",
		errors.Code(
			errors.New("test"),
		),
	)

	assert.Equal(t,
		"T123",
		errors.Code(
			errors.WithCode(errors.New("test"), "T123"),
		),
	)

	assert.Equal(t,
		"T123",
		errors.Code(
			errors.Wrap(
				errors.WithCode(errors.New("test"), "T123"),
				"wrapped",
			),
		),
	)
}

func Test_ContainsCode(t *testing.T) {
	assert.False(t,
		errors.ContainsCode(
			errors.New("test"),
			"T123",
		),
	)

	assert.True(t,
		errors.ContainsCode(
			errors.WithCode(errors.New("test"), "T123"),
			"T123",
		),
	)

	assert.True(t,
		errors.ContainsCode(
			errors.Wrap(
				errors.WithCode(errors.New("test"), "T123"),
				"wrapped",
			),
			"T123",
		),
	)
}
