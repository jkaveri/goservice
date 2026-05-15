package errors_test

import (
	stderrors "errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestIs_nilErr(t *testing.T) {
	assert.False(t, errors.Is(nil, errors.New("root")))
}

func TestIs(t *testing.T) {
	type Args struct {
		err    error
		target error
	}
	type Expects struct {
		want bool
	}

	root := errors.New("root")
	wrapped := errors.Wrap(root, "outer")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "same-pointer",
			args:    Args{err: root, target: root},
			expects: Expects{want: true},
		},
		{
			name:    "wrapped-root",
			args:    Args{err: wrapped, target: root},
			expects: Expects{want: true},
		},
		{
			name:    "different-error",
			args:    Args{err: wrapped, target: errors.New("other")},
			expects: Expects{want: false},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.want, errors.Is(tc.args.err, tc.args.target))
		})
	}
}

func TestAs_nilErr(t *testing.T) {
	var codeErr errors.CodeError
	assert.False(t, errors.As(nil, &codeErr))
}

func TestAs_findsCodeInChain(t *testing.T) {
	coded := errors.WithCode(errors.New("root"), "E1")
	err := errors.Wrap(coded, "outer")

	var codeErr errors.CodeError
	require.True(t, errors.As(err, &codeErr))
	assert.Equal(t, "E1", codeErr.Code())
}

func TestAs_plainError(t *testing.T) {
	var codeErr errors.CodeError
	assert.False(t, errors.As(errors.New("plain"), &codeErr))
}

func TestJoin_allNil(t *testing.T) {
	assert.Nil(t, errors.Join(nil, nil))
}

func TestJoin(t *testing.T) {
	type Args struct {
		errs []error
	}
	type Expects struct {
		want string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "skips-nil",
			args:    Args{errs: []error{errors.New("a"), nil, errors.New("b")}},
			expects: Expects{want: "a\nb"},
		},
		{
			name:    "single-error",
			args:    Args{errs: []error{errors.New("only")}},
			expects: Expects{want: "only"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Join(tc.args.errs...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.True(t, stderrors.Is(got, tc.args.errs[0]))
		})
	}
}
