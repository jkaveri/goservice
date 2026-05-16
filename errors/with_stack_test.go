package errors_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestWrap_attachesStack(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantUnwrap   error
		wantStackLen int
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "wrap-adds-stack",
			args:    Args{err: root},
			expects: Expects{wantStackLen: 1},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrap(tc.args.err, "ctx")
			require.NotNil(t, got)
			assert.Equal(t, "ctx: root", got.Error())
			assert.True(t, errors.HasStack(got))

			var stackErr errors.StackError
			require.True(t, errors.As(got, &stackErr))
			assert.GreaterOrEqual(t, len(stackErr.StackTrace()), tc.expects.wantStackLen)
			assert.Contains(t, fmt.Sprintf("%+v", got), "root")
			assert.Contains(t, fmt.Sprintf("%+v", got), "with_stack_test.go")
		})
	}
}

func TestWrap_stackPreservesInnerError(t *testing.T) {
	type Args struct {
		inner error
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
			name:    "delegates-to-fundamental",
			args:    Args{inner: root},
			expects: Expects{want: "ctx: root"},
		},
		{
			name:    "delegates-to-plain-err",
			args:    Args{inner: &errNoUnwrap{msg: "db timeout"}},
			expects: Expects{want: "ctx: db timeout"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.Wrap(tc.args.inner, "ctx")
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
		})
	}
}

func TestWrap_stackFormat(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS         string
		wantQ         string
		wantV         string
		vPlusContains []string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "wrap-plain-err",
			args: Args{
				err: errors.Wrap(&errNoUnwrap{msg: "db timeout"}, "query"),
			},
			expects: Expects{
				wantS:         "query: db timeout",
				wantQ:         `"query"`,
				wantV:         "query: db timeout",
				vPlusContains: []string{"db timeout", "with_stack_test.go", "runtime."},
			},
		},
		{
			name: "wrap-fundamental-leaf",
			args: Args{
				err: errors.Wrap(errors.New("leaf"), "ctx"),
			},
			expects: Expects{
				wantS:         "ctx: leaf",
				wantQ:         `"ctx"`,
				wantV:         "ctx: leaf",
				vPlusContains: []string{"leaf", "with_stack_test.go", "runtime."},
			},
		},
		{
			name: "nested-wrap-keeps-first-stack-only",
			args: Args{
				err: errors.Wrap(
					errors.Wrap(errors.New("persist failed"), "inner"),
					"outer",
				),
			},
			expects: Expects{
				wantS:         "outer: inner: persist failed",
				wantQ:         `"outer"`,
				wantV:         "outer: inner: persist failed",
				vPlusContains: []string{"persist failed", "with_stack_test.go", "runtime."},
			},
		},
		{
			name: "with-message-chain",
			args: Args{
				err: errors.WithMessage(
					errors.Wrap(errors.New("inner"), "ctx"),
					"outer",
				),
			},
			expects: Expects{
				wantS:         "outer: ctx: inner",
				wantQ:         `"outer: ctx: inner"`,
				wantV:         "outer: ctx: inner",
				vPlusContains: []string{"outer", "inner", "with_stack_test.go", "runtime."},
			},
		},
		{
			name: "fmt-errorf-inner",
			args: Args{
				err: errors.Wrap(
					fmt.Errorf("call partner api: %s", "503 Service Unavailable"),
					"query",
				),
			},
			expects: Expects{
				wantS:         "query: call partner api: 503 Service Unavailable",
				wantQ:         `"query"`,
				wantV:         "query: call partner api: 503 Service Unavailable",
				vPlusContains: []string{"call partner api", "503 Service Unavailable", "with_stack_test.go", "runtime."},
			},
		},
		{
			name: "with-code-on-foreign-err",
			args: Args{
				err: errors.WithCode(
					fmt.Errorf("run backup: %w", context.DeadlineExceeded),
					"JOB_TIMEOUT",
				),
			},
			expects: Expects{
				wantS:         "[JOB_TIMEOUT] run backup: context deadline exceeded",
				wantQ:         `"[JOB_TIMEOUT] run backup: context deadline exceeded"`,
				wantV:         "[JOB_TIMEOUT] run backup: context deadline exceeded",
				vPlusContains: []string{"run backup", "context deadline exceeded", "with_stack_test.go", "runtime."},
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
