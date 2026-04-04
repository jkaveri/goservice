package validate_test

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	validate "github.com/jkaveri/goservice/grpc/interceptors/validate"
)

type mockValidator struct {
	shouldError bool
}

func (m mockValidator) ValidateAll() error {
	if m.shouldError {
		return errors.New("validation error")
	}
	return nil
}

func Test_Validate(t *testing.T) {
	t.Run("valid-request", func(t *testing.T) {
		ctx := context.Background()
		called := false
		req := mockValidator{shouldError: false}

		mw := validate.UnaryInterceptor()

		resp, err := mw(
			ctx,
			req,
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestValidate",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true
				return "test response", nil
			},
		)

		assert.True(t, called)
		assert.NoError(t, err)
		assert.Equal(t, "test response", resp)
	})

	t.Run("invalid-request", func(t *testing.T) {
		ctx := context.Background()
		called := false
		req := mockValidator{shouldError: true}

		mw := validate.UnaryInterceptor()

		resp, err := mw(
			ctx,
			req,
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestValidate",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true
				return nil, nil
			},
		)

		assert.False(t, called)
		assert.Error(t, err)
		assert.Nil(t, resp)
	})

	t.Run("non-validator-request", func(t *testing.T) {
		ctx := context.Background()
		called := false
		req := "non-validator"

		mw := validate.UnaryInterceptor()

		resp, err := mw(
			ctx,
			req,
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestValidate",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true
				return "test response", nil
			},
		)

		assert.True(t, called)
		assert.NoError(t, err)
		assert.Equal(t, "test response", resp)
	})
}
