package errors_test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	errors "github.com/jkaveri/goservice/errors"
)

func TestWithUserMessage_GetUserMessage(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		userMsg string
		errStr  string
	}

	base := errors.New("db connection refused")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "nil-error-returns-nil",
			args: Args{err: nil},
			expects: Expects{
				userMsg: "",
				errStr:  "",
			},
		},
		{
			name: "plain-error-has-no-user-message",
			args: Args{err: base},
			expects: Expects{
				userMsg: "",
				errStr:  "db connection refused",
			},
		},
		{
			name: "user-message-outermost",
			args: Args{
				err: errors.WithUserMessage(
					errors.WithCode(base, "not_found"),
					"we could not find that item",
				),
			},
			expects: Expects{
				userMsg: "we could not find that item",
				errStr:  "db connection refused",
			},
		},
		{
			name: "user-message-inside-code-wrapper",
			args: Args{
				err: errors.WithCode(
					errors.WithUserMessage(base, "try again later"),
					"internal_server_error",
				),
			},
			expects: Expects{
				userMsg: "try again later",
				errStr:  "db connection refused",
			},
		},
		{
			name: "empty-user-message-skips-wrapper",
			args: Args{err: errors.WithUserMessage(base, "")},
			expects: Expects{
				userMsg: "",
				errStr:  "db connection refused",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			gotMsg := errors.GetUserMessage(tc.args.err)
			assert.Equal(t, tc.expects.userMsg, gotMsg)

			if tc.args.err == nil {
				assert.Equal(t, tc.expects.errStr, "")
				return
			}

			assert.Equal(t, tc.expects.errStr, tc.args.err.Error())
		})
	}
}

func TestWithUserMessage_Error(t *testing.T) {
	type Args struct {
		inner       error
		userMessage string
	}
	type Expects struct {
		want          string
		wantUserMsg   string
	}

	inner := errors.New("technical failure")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "returns-inner-not-user-message",
			args: Args{
				inner:       inner,
				userMessage: "please try again",
			},
			expects: Expects{
				want:        "technical failure",
				wantUserMsg: "please try again",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			got := errors.WithUserMessage(tc.args.inner, tc.args.userMessage)
			require.NotNil(t, got)
			assert.Equal(t, tc.expects.want, got.Error())
			assert.Equal(t, tc.expects.wantUserMsg, errors.GetUserMessage(got))
			assert.Equal(t, tc.expects.want, tc.args.inner.Error())
		})
	}
}

func TestWithUserMessage_Format(t *testing.T) {
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
	err := errors.WithUserMessage(inner, "try again")

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "with-user-message",
			args: Args{err: err},
			expects: Expects{
				wantS:         "inner",
				wantQ:         "inner",
				wantV:         "inner",
				vPlusContains: []string{`user_message="try again"`, "inner"},
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
