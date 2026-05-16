package validate

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/jkaveri/goservice/errorcode"
	svcerrors "github.com/jkaveri/goservice/errors"
)

type stubValidator struct {
	err error
}

func (s stubValidator) ValidateAll() error {
	return s.err
}

func TestUnaryInterceptor_usesFriendlyValidationMessage(t *testing.T) {
	ctx := context.Background()
	mw := UnaryInterceptor()

	called := false
	_, err := mw(
		ctx,
		stubValidator{
			err: testFieldErr{
				field:  "Message",
				reason: "value length must be at least 1 runes",
			},
		},
		&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
		func(context.Context, interface{}) (interface{}, error) {
			called = true
			return nil, nil
		},
	)

	require.Error(t, err)
	assert.False(t, called)
	assert.Equal(t, errorcode.CodeInvalidRequest, svcerrors.Code(err))
	assert.Equal(t, "[invalid_request] Message: value length must be at least 1 runes", err.Error())
}
