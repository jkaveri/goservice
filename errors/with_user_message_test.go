package errors_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

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
