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

	assert.Equal(t, err, errors.Unwrap(errWithMetadata))
	assert.Equal(t, "test", errWithMetadata.Error())
	assert.Contains(t, fmt.Sprintf("%+v", errWithMetadata), "test")
	assert.Contains(t, fmt.Sprintf("%s", errWithMetadata), "test")
	assert.Contains(t, fmt.Sprintf("%q", errWithMetadata), "test")
}

func Test_GetMetadata(t *testing.T) {
	assert.Nil(t,
		errors.Metadata(
			errors.New("test"),
		),
	)

	assert.Equal(t,
		map[string]any{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.WithMetadata(
				errors.New("test"),
				map[string]any{
					"Email": "john@email.com",
				},
			),
		),
	)

	assert.Equal(t,
		map[string]any{
			"Email": "john@email.com",
		},
		errors.Metadata(
			errors.Wrap(
				errors.WithMetadata(
					errors.New("test"),
					map[string]any{
						"Email": "john@email.com",
					},
				),
				"wrapped",
			),
		),
	)
}

func TestWithMetadata_Error(t *testing.T) {
	type Args struct {
		inner error
	}
	type Expects struct {
		want string
	}

	inner := errors.New("db down")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "delegates-to-inner",
			args:    Args{inner: inner},
			expects: Expects{want: "db down"},
		},
		{
			name: "delegates-through-message-wrapper",
			args: Args{
				inner: errors.WithMessage(errors.New("leaf"), "annot"),
			},
			expects: Expects{want: "annot"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMetadata(tc.args.inner, map[string]any{"key": "val"})
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestWithMetadata_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	inner := errors.New("inner")
	err := errors.WithMetadata(inner, map[string]any{"key": "val"})

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-metadata",
			args: Args{err: err},
			expects: Expects{
				wantS:         "inner",
				wantQ:         "inner",
				wantV:         "inner",
				vPlusContains: []string{"inner"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.wantS, fmt.Sprintf("%s", tc.args.err))
			assert.Equal(t, tc.expects.wantQ, fmt.Sprintf("%q", tc.args.err))
			assert.Equal(t, tc.expects.wantV, fmt.Sprintf("%v", tc.args.err))

			gotPlus := fmt.Sprintf("%+v", tc.args.err)
			for _, sub := range tc.expects.vPlusContains {
				assert.Contains(t, gotPlus, sub)
			}
		})
	}
}
