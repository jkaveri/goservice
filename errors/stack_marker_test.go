package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestHasStack(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		want bool
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "nil",
			args:    Args{err: nil},
			expects: Expects{want: false},
		},
		{
			name:    "leaf-new",
			args:    Args{err: errors.New("leaf")},
			expects: Expects{want: false},
		},
		{
			name:    "wrap-attaches-stack",
			args:    Args{err: errors.Wrap(errors.New("leaf"), "ctx")},
			expects: Expects{want: true},
		},
		{
			name: "outer-wrap-propagates-inner-stack",
			args: Args{
				err: errors.Wrap(
					errors.Wrap(errors.New("leaf"), "inner"),
					"outer",
				),
			},
			expects: Expects{want: true},
		},
		{
			name: "with-message-propagates-stack",
			args: Args{
				err: errors.WithMessage(
					errors.Wrap(errors.New("leaf"), "inner"),
					"annot",
				),
			},
			expects: Expects{want: true},
		},
		{
			name: "with-code-attaches-stack",
			args: Args{
				err: errors.WithCode(errors.New("leaf"), "E1"),
			},
			expects: Expects{want: true},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.want, errors.HasStack(tc.args.err))
		})
	}
}

func TestWrap_skipsSecondStackCapture(t *testing.T) {
	inner := errors.Wrap(errors.New("persist failed"), "inner")
	outer := errors.Wrap(inner, "outer")

	assert.NotSame(t, inner, outer)
	assert.True(t, errors.HasStack(outer))

	var stackErr errors.StackError
	require.True(t, errors.As(outer, &stackErr))
}
