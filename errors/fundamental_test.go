package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestNew(t *testing.T) {
	type Args struct {
		message string
	}
	type Expects struct{}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "empty-message",
			args: Args{message: ""},
		},
		{
			name: "returns-message",
			args: Args{message: "something failed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.New(tc.args.message)
			require.NotNil(t, got)
		})
	}
}

func TestErrorf(t *testing.T) {
	type Args struct {
		format string
		args   []any
	}
	testCases := []struct {
		name string
		args Args
	}{
		{
			name: "empty-format",
			args: Args{
				format: "",
				args:   []any{},
			},
		},
		{
			name: "empty-args",
			args: Args{
				format: "code %d: %s",
				args:   []any{},
			},
		},
		{
			name: "formats-message",
			args: Args{
				format: "code %d: %s",
				args:   []any{42, "timeout"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Errorf(tc.args.format, tc.args.args...)
			require.NotNil(t, got)
		})
	}
}

func TestNew_Error(t *testing.T) {
	type Args struct {
		message string
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
			name:    "returns-message",
			args:    Args{message: "something failed"},
			expects: Expects{want: "something failed"},
		},
		{
			name:    "empty-message",
			args:    Args{message: ""},
			expects: Expects{want: ""},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.New(tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestErrorf_Error(t *testing.T) {
	type Args struct {
		format string
		args   []any
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
			name: "formats-message",
			args: Args{
				format: "code %d: %s",
				args:   []any{42, "timeout"},
			},
			expects: Expects{want: "code 42: timeout"},
		},
		{
			name: "empty-format",
			args: Args{
				format: "",
				args:   []any{},
			},
			expects: Expects{want: ""},
		},
		{
			name: "empty-args",
			args: Args{
				format: "code %d: %s",
				args:   []any{},
			},
			expects: Expects{want: "code %!d(MISSING): %!s(MISSING)"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Errorf(tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestNew_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS     string
		wantQ     string
		wantV     string
		wantVPlus string
	}

	leaf := errors.New("boom")
	root := errors.New("root")
	mid := errors.Wrap(root, "mid")
	top := errors.Wrap(mid, "top")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "leaf",
			args: Args{err: leaf},
			expects: Expects{
				wantS:     "boom",
				wantQ:     `"boom"`,
				wantV:     "boom",
				wantVPlus: "boom",
			},
		},
		{
			name: "three-level-chain",
			args: Args{err: top},
			expects: Expects{
				wantS:     "top: mid: root",
				wantQ:     `"top"`,
				wantV:     "top: mid: root",
				wantVPlus: "top\n\tmid\n\troot",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.wantS, fmt.Sprintf("%s", tc.args.err))
			assert.Equal(t, tc.expects.wantQ, fmt.Sprintf("%q", tc.args.err))
			assert.Equal(t, tc.expects.wantV, fmt.Sprintf("%v", tc.args.err))

			gotPlus := fmt.Sprintf("%+v", tc.args.err)
			if tc.name == "three-level-chain" {
				assert.True(t, len(gotPlus) >= len(tc.expects.wantVPlus))
				assert.Equal(t, tc.expects.wantVPlus, gotPlus[:len(tc.expects.wantVPlus)])
				assert.Contains(t, gotPlus, "runtime.")
				assert.True(t, errors.HasStack(tc.args.err))
				return
			}

			assert.Equal(t, tc.expects.wantVPlus, gotPlus)
		})
	}
}

func TestWrap_Error(t *testing.T) {
	type Args struct {
		cause   error
		message string
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
			name:    "joins-message-and-cause",
			args:    Args{cause: errors.New("inner"), message: "outer"},
			expects: Expects{want: "outer: inner"},
		},
		{
			name:    "nil-cause-uses-message-only",
			args:    Args{cause: nil, message: "orphan"},
			expects: Expects{want: "orphan"},
		},
		{
			name: "three-level-joins-each-layer",
			args: Args{
				cause: errors.Wrap(
					errors.Wrap(errors.New("root"), "mid"),
					"inner",
				),
				message: "outer",
			},
			expects: Expects{want: "outer: inner: mid: root"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrap(tc.args.cause, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestWrapf_Error(t *testing.T) {
	type Args struct {
		cause   error
		format  string
		fmtArgs []any
	}
	type Expects struct {
		want string
	}

	root := errors.New("root")
	mid := errors.Wrapf(root, "errno=%d", 5)

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "single-layer",
			args: Args{
				cause:   errors.New("inner"),
				format:  "code %d",
				fmtArgs: []any{9},
			},
			expects: Expects{want: "code 9: inner"},
		},
		{
			name: "nested-formatted-layers",
			args: Args{
				cause:   mid,
				format:  "svc %s",
				fmtArgs: []any{"api"},
			},
			expects: Expects{want: "svc api: errno=5: root"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrapf(tc.args.cause, tc.args.format, tc.args.fmtArgs...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestUnwrap(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		want error
	}

	inner := errors.New("inner")
	root := errors.New("root")
	mid := errors.Wrap(root, "mid")
	top := errors.Wrap(mid, "top")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "returns-cause-from-wrap",
			args:    Args{err: errors.Wrap(inner, "outer")},
			expects: Expects{want: inner},
		},
		{
			name:    "nested-returns-immediate-cause-only",
			args:    Args{err: top},
			expects: Expects{want: mid},
		},
		{
			name:    "leaf-new",
			args:    Args{err: errors.New("leaf")},
			expects: Expects{want: nil},
		},
		{
			name:    "plain-error-without-unwrap",
			args:    Args{err: &errNoUnwrap{msg: "plain"}},
			expects: Expects{want: nil},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.name == "returns-cause-from-wrap" {
				assert.True(t, errors.Is(tc.args.err, inner))
				assert.True(t, errors.HasStack(tc.args.err))
				return
			}

			assert.Equal(t, tc.expects.want, errors.Unwrap(tc.args.err))
		})
	}
}

func TestWrap_errorsIs(t *testing.T) {
	type Args struct {
		chain  error
		target error
	}

	inner := errors.New("inner")
	singleWrapped := errors.Wrap(inner, "outer")

	root := errors.New("root")
	mid := errors.Wrap(root, "mid")
	top := errors.Wrap(mid, "top")

	testCases := []struct {
		name string
		args Args
	}{
		{
			name: "single-wrap",
			args: Args{chain: singleWrapped, target: inner},
		},
		{
			name: "three-level-finds-mid",
			args: Args{chain: top, target: mid},
		},
		{
			name: "three-level-finds-root",
			args: Args{chain: top, target: root},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.True(t, stderrors.Is(tc.args.chain, tc.args.target))
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

	inner := errors.New("inner")
	root := errors.New("root")
	mid := errors.Wrap(root, "mid")
	top := errors.Wrap(mid, "top")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "wrapped-non-nil-cause",
			args: Args{err: errors.Wrap(inner, "outer")},
			expects: Expects{
				wantS:         "outer: inner",
				wantQ:         `"outer"`,
				wantV:         "outer: inner",
				vPlusContains: []string{"outer", "inner"},
			},
		},
		{
			name: "three-level-chain",
			args: Args{err: top},
			expects: Expects{
				wantS:         "top: mid: root",
				wantQ:         `"top"`,
				wantV:         "top: mid: root",
				vPlusContains: []string{"top", "mid", "root"},
			},
		},
		{
			name: "wrapped-nil-cause",
			args: Args{err: errors.Wrap(nil, "only")},
			expects: Expects{
				wantS:         "only",
				wantQ:         `"only"`,
				wantV:         "only",
				vPlusContains: []string{"only"},
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
