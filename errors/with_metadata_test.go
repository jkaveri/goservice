package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func Test_WithExtra(t *testing.T) {
	err := errors.New("test")

	errWithMetadata := errors.WithMetadata(err, map[string]interface{}{
		"Email": "john@email.com",
	})

	metadataGetter, ok := errWithMetadata.(errors.MetadataError)
	require.True(t, ok)

	assert.True(t, ok)
	assert.Equal(t, map[string]interface{}{
		"Email": "john@email.com",
	}, metadataGetter.Metadata())

	assert.Equal(t, err, errors.Cause(errWithMetadata))
	assert.Equal(t, err, errors.Unwrap(errWithMetadata))
	assert.Equal(t, "test", errWithMetadata.Error())
	assert.Contains(t, fmt.Sprintf("%+v", errWithMetadata), "test")
	assert.Contains(t, fmt.Sprintf("%s", errWithMetadata), "test")
	assert.Contains(t, fmt.Sprintf("%q", errWithMetadata), "test")
}

func Test_GetMetadata(t *testing.T) {
	assert.Equal(t,
		map[string]interface{}{},
		errors.Metadata(
			errors.New("test"),
		),
	)

	assert.Equal(t,
		map[string]interface{}{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.WithMetadata(
				errors.New("test"),
				map[string]interface{}{
					"Email": "john@email.com",
				},
			),
		),
	)

	assert.Equal(t,
		map[string]interface{}{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.Wrap(
				errors.WithMetadata(
					errors.New("test"),
					map[string]interface{}{
						"Email": "john@email.com",
					},
				),
				"wrapped",
			),
		),
	)
}
