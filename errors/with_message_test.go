package errors_test

import (
	stderrors "errors"
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestNewMessage(t *testing.T) {
	got := errors.NewMessage("slug already in use")
	require.NotNil(t, got)
	assert.Equal(t, "slug already in use", got.Error())
	assert.Equal(t, "slug already in use", errors.Message(got))
	assert.Nil(t, errors.Unwrap(got))
	assert.False(t, errors.HasStack(got))
}

func TestNewMessagef(t *testing.T) {
	got := errors.NewMessagef("step %d failed", 3)
	require.NotNil(t, got)
	assert.Equal(t, "step 3 failed", got.Error())
	assert.Equal(t, "step 3 failed", errors.Message(got))
}

func TestWithMessage_nilErr(t *testing.T) {
	got := errors.WithMessage(nil, "slug already in use")
	require.NotNil(t, got)
	assert.Equal(t, "slug already in use", got.Error())
	assert.Equal(t, "slug already in use", errors.Message(got))
	assert.Nil(t, errors.Unwrap(got))
	assert.False(t, errors.HasStack(got))
}

func TestWithMessage(t *testing.T) {
	type Args struct {
		err     error
		message string
	}
	type Expects struct {
		wantError string
	}

	root := errors.New("root")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name:    "prefixes-message",
			args:    Args{err: root, message: "outer"},
			expects: Expects{wantError: "outer: root"},
		},
		{
			name:    "nil-err-message-only",
			args:    Args{err: nil, message: "slug already in use"},
			expects: Expects{wantError: "slug already in use"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessage(tc.args.err, tc.args.message)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.wantError, got.Error())
			assert.Equal(t, tc.args.message, errors.Message(got))

			if tc.args.err == nil {
				assert.Nil(t, errors.Unwrap(got))
				assert.False(t, errors.HasStack(got))
				return
			}

			assert.True(t, errors.Is(got, tc.args.err))
			assert.True(t, errors.HasStack(got))
			assert.Contains(t, fmt.Sprintf("%+v", got), "root")
			assert.Contains(t, fmt.Sprintf("%+v", got), tc.args.message)
			assert.Equal(t, fmt.Sprintf("%q", got.Error()), fmt.Sprintf("%q", got))
		})
	}
}

func TestWithMessagef_nilErr(t *testing.T) {
	got := errors.WithMessagef(nil, "step %d failed", 3)
	require.NotNil(t, got)
	assert.Equal(t, "step 3 failed", got.Error())
	assert.Equal(t, "step 3 failed", errors.Message(got))
}

func TestWithMessagef(t *testing.T) {
	type Args struct {
		err    error
		format string
		args   []any
	}
	type Expects struct {
		wantError string
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
			expects: Expects{wantError: "step 3 failed: root"},
		},
		{
			name: "nil-err-formatted-message-only",
			args: Args{
				err:    nil,
				format: "step %d failed",
				args:   []any{3},
			},
			expects: Expects{wantError: "step 3 failed"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithMessagef(tc.args.err, tc.args.format, tc.args.args...)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.wantError, got.Error())

			if tc.args.err == nil {
				assert.Nil(t, errors.Unwrap(got))
				return
			}

			assert.True(t, errors.Is(got, tc.args.err))
			assert.True(t, errors.HasStack(got))
		})
	}
}

func TestMessage(t *testing.T) {
	type Args struct {
		err error
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
			name:    "nil",
			args:    Args{err: nil},
			expects: Expects{want: ""},
		},
		{
			name:    "plain-fundamental",
			args:    Args{err: errors.New("x")},
			expects: Expects{want: ""},
		},
		{
			name:    "message-only-new-message",
			args:    Args{err: errors.NewMessage("slug already in use")},
			expects: Expects{want: "slug already in use"},
		},
		{
			name: "stdlib-without-message-interface",
			args: Args{err: fmt.Errorf("stdlib: %w", stderrors.New("inner"))},
			expects: Expects{
				want: "",
			},
		},
		{
			name: "single-with-message",
			args: Args{
				err: errors.WithMessage(errors.New("leaf"), "annot"),
			},
			expects: Expects{want: "annot"},
		},
		{
			name: "nested-create-tenant-duplicate-slug",
			args: Args{
				err: errors.WithMessage(
					errors.WithMessage(
						errors.New(`duplicate key value violates unique constraint "tenants_slug_key"`),
						"slug already in use",
					),
					"create tenant",
				),
			},
			expects: Expects{want: "create tenant: slug already in use"},
		},
		{
			name: "skips-wrap-without-message-interface",
			args: Args{
				err: errors.Wrap(
					errors.WithMessage(errors.New("leaf"), "annot"),
					"wrapped",
				),
			},
			expects: Expects{want: "annot"},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.want, errors.Message(tc.args.err))
		})
	}
}

func TestWithMessage_Format(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		wantS           string
		wantV           string
		wantQ           string
		wantVPlus       string
		wantVPlusPrefix string
		vPlusContains   []string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-message",
			args: Args{
				err: errors.WithMessage(errors.New("inner"), "outer"),
			},
			expects: Expects{
				wantS:           "outer: inner",
				wantV:           "outer: inner",
				wantQ:           "\"outer: inner\"",
				wantVPlusPrefix: "outer\ninner",
				vPlusContains:   []string{"runtime."},
			},
		},
		{
			name: "with-messagef",
			args: Args{
				err: errors.WithMessagef(errors.New("inner"), "step %d", 9),
			},
			expects: Expects{
				wantS:           "step 9: inner",
				wantV:           "step 9: inner",
				wantQ:           "\"step 9: inner\"",
				wantVPlusPrefix: "step 9\ninner",
				vPlusContains:   []string{"runtime."},
			},
		},
		{
			name: "nested-with-message",
			args: Args{
				err: errors.WithMessage(
					errors.WithMessage(errors.New("leaf"), "mid"),
					"outer",
				),
			},
			expects: Expects{
				wantS:           "outer: mid: leaf",
				wantV:           "outer: mid: leaf",
				wantQ:           "\"outer: mid: leaf\"",
				wantVPlusPrefix: "outer\nmid\nleaf",
				vPlusContains:   []string{"runtime."},
			},
		},
		{
			name: "message-contains-colon",
			args: Args{
				err: errors.WithMessage(errors.New("inner"), "rpc: failed"),
			},
			expects: Expects{
				wantS:           "rpc: failed: inner",
				wantV:           "rpc: failed: inner",
				wantQ:           "\"rpc: failed: inner\"",
				wantVPlusPrefix: "rpc: failed\ninner",
				vPlusContains:   []string{"runtime."},
			},
		},
		{
			name: "stdlib-cause",
			args: Args{
				err: errors.WithMessage(fmt.Errorf("stdlib leaf"), "annot"),
			},
			expects: Expects{
				wantS:           "annot: stdlib leaf",
				wantV:           "annot: stdlib leaf",
				wantQ:           "\"annot: stdlib leaf\"",
				wantVPlusPrefix: "annot\nstdlib leaf",
				vPlusContains:   []string{"runtime."},
			},
		},
		{
			name: "nil-cause-message-only",
			args: Args{
				err: errors.WithMessage(nil, "slug already in use"),
			},
			expects: Expects{
				wantS:     "slug already in use",
				wantV:     "slug already in use",
				wantQ:     "\"slug already in use\"",
				wantVPlus: "slug already in use",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.Equal(t, tc.expects.wantS, fmt.Sprintf("%s", tc.args.err))
			assert.Equal(t, tc.expects.wantQ, fmt.Sprintf("%q", tc.args.err))
			assert.Equal(t, tc.expects.wantV, fmt.Sprintf("%v", tc.args.err))

			gotPlus := fmt.Sprintf("%+v", tc.args.err)
			switch {
			case tc.expects.wantVPlusPrefix != "":
				assert.True(t, len(gotPlus) >= len(tc.expects.wantVPlusPrefix))
				assert.Equal(t, tc.expects.wantVPlusPrefix, gotPlus[:len(tc.expects.wantVPlusPrefix)])
				for _, sub := range tc.expects.vPlusContains {
					assert.Contains(t, gotPlus, sub)
				}
			default:
				assert.Equal(t, tc.expects.wantVPlus, gotPlus)
			}
		})
	}
}
