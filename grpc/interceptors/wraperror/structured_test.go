package wraperror_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/jkaveri/goservice/errorcode"
	errors "github.com/jkaveri/goservice/errors"
	"github.com/jkaveri/goservice/grpc/interceptors/wraperror"
)

func TestToStructured_Message(t *testing.T) {
	type Args struct {
		err error
	}
	type Expects struct {
		code    string
		message string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "uses-generic-message-when-no-with-message",
			args: Args{err: errorcode.NotFound("internal detail")},
			expects: Expects{
				code:    errorcode.CodeNotFound,
				message: "not found",
			},
		},
		{
			name: "prefers-user-message-over-generic",
			args: Args{
				err: errors.WithMessage(
					errorcode.NotFound("internal detail"),
					"that product is not available",
				),
			},
			expects: Expects{
				code:    errorcode.CodeNotFound,
				message: "that product is not available",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			se := wraperror.ToStructured(tc.args.err)

			assert.Equal(t, tc.expects.code, se.Code)
			assert.Equal(t, tc.expects.message, se.ErrorMessage)
		})
	}
}
