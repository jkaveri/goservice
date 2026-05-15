package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestWithMessage_nilErr(t *testing.T) {
	assert.Nil(t, errors.WithMessage(nil, "ignored"))
}

func TestWithMessage(t *testing.T) {
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
			name:    "annotates-message",
			args:    Args{err: root, message: "while saving"},
			expects: Expects{want: "while saving"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessage(tc.args.err, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, root, errors.Unwrap(got))
		})
	}
}

func TestWithMessagef_nilErr(t *testing.T) {
	assert.Nil(t, errors.WithMessagef(nil, "ignored %s", "x"))
}

func TestWithMessagef(t *testing.T) {
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
				format: "user %s not found",
				args:   []any{"alice"},
			},
			expects: Expects{want: "user alice not found"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessagef(tc.args.err, tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, root, errors.Unwrap(got))
		})
	}
}

func TestWithMessage_Error(t *testing.T) {
	type Args struct {
		inner   error
		message string
	}
	type Expects struct {
		want      string
		wantInner string
	}

	inner := errors.New("inner")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "returns-annotation-only",
			args: Args{
				inner:   inner,
				message: "while saving",
			},
			expects: Expects{
				want:      "while saving",
				wantInner: "inner",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessage(tc.args.inner, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, tc.expects.wantInner, tc.args.inner.Error())
		})
	}
}

func TestWithMessagef_Error(t *testing.T) {
	type Args struct {
		inner  error
		format string
		args   []any
	}
	type Expects struct {
		want      string
		wantInner string
	}

	inner := errors.New("inner")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "formats-annotation-only",
			args: Args{
				inner:  inner,
				format: "user %s not found",
				args:   []any{"alice"},
			},
			expects: Expects{
				want:      "user alice not found",
				wantInner: "inner",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessagef(tc.args.inner, tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, tc.expects.wantInner, tc.args.inner.Error())
		})
	}
}

func TestWithMessage_Format(t *testing.T) {
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
	err := errors.WithMessage(inner, "annot")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-message",
			args: Args{err: err},
			expects: Expects{
				wantS:         "annot",
				wantQ:         "annot",
				wantV:         "annot",
				vPlusContains: []string{"inner", "annot"},
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
