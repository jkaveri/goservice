package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestWrap_nilErr(t *testing.T) {
	assert.Nil(t, errors.Wrap(nil, "ignored"))
}

func TestWrap(t *testing.T) {
	type Args struct {
		err     error
		message string
	}
	type Expects struct {
		want string
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "wraps-message",
			args:    Args{err: root, message: "outer"},
			expects: Expects{want: "outer"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrap(tc.args.err, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, root, errors.Unwrap(got))
			assert.Contains(t, fmt.Sprintf("%+v", got), tc.expects.want)
			assert.Contains(t, fmt.Sprintf("%+v", got), root.Error())
		})
	}
}

func TestWrapf_nilErr(t *testing.T) {
	assert.Nil(t, errors.Wrapf(nil, "ignored %s", "x"))
}

func TestWrapf(t *testing.T) {
	type Args struct {
		err    error
		format string
		args   []any
	}
	type Expects struct {
		want string
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "formats-message",
			args: Args{
				err:    root,
				format: "step %d failed",
				args:   []any{3},
			},
			expects: Expects{want: "step 3 failed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrapf(tc.args.err, tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, root, errors.Unwrap(got))
		})
	}
}

func TestUnwrap_nilErr(t *testing.T) {
	assert.Nil(t, errors.Unwrap(nil))
}

func TestUnwrap(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		want error
	}

	root := errors.New("root")
	wrapped := errors.Wrap(root, "outer")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "leaf-err",
			args:    Args{err: root},
			expects: Expects{want: nil},
		},
		{
			name:    "one-level",
			args:    Args{err: wrapped},
			expects: Expects{want: root},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.want, errors.Unwrap(tc.args.err))
		})
	}
}

func TestWrap_Error(t *testing.T) {
	type Args struct {
		inner   error
		message string
	}
	type Expects struct {
		want      string
		wantInner string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "returns-wrap-message",
			args: Args{
				inner:   errors.New("root"),
				message: "outer",
			},
			expects: Expects{
				want:      "outer",
				wantInner: "root",
			},
		},
		{
			name: "plain-inner-unchanged",
			args: Args{
				inner:   &errNoUnwrap{msg: "root cause"},
				message: "outer",
			},
			expects: Expects{
				want:      "outer",
				wantInner: "root cause",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrap(tc.args.inner, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, tc.expects.wantInner, tc.args.inner.Error())
		})
	}
}

func TestWrapf_Error(t *testing.T) {
	type Args struct {
		inner  error
		format string
		args   []any
	}
	type Expects struct {
		want string
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "formats-wrap-message",
			args: Args{
				inner:  root,
				format: "step %d failed",
				args:   []any{3},
			},
			expects: Expects{want: "step 3 failed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrapf(tc.args.inner, tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestWrap_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	inner := &errNoUnwrap{msg: "root cause"}
	err := errors.Wrap(inner, "outer")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "wrap-plain-inner",
			args: Args{err: err},
			expects: Expects{
				wantS:         "outer",
				wantQ:         "outer",
				wantV:         "outer",
				vPlusContains: []string{"root cause", "outer"},
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
