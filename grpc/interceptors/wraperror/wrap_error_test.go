package wraperror_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	errors "github.com/jkaveri/goservice/errors"
	"github.com/jkaveri/goservice/grpc/interceptors/wraperror"
)

func Test_WrapError(t *testing.T) {
	t.Run("should not return error", func(t *testing.T) {
		mw := wraperror.UnaryInterceptor()

		h := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, "test", req)
			return "resp 1", nil
		})

		resp, err := mw(context.Background(), "test", &grpc.UnaryServerInfo{}, h)

		assert.NoError(t, err)

		assert.Equal(t, "resp 1", resp)
	})

	t.Run("should return structured error", func(t *testing.T) {
		mw := wraperror.UnaryInterceptor()

		h := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, "test req", req)
			return nil, errors.New("test err")
		})

		resp, err := mw(context.Background(), "test req", &grpc.UnaryServerInfo{}, h)

		assert.Error(t, err)
		assert.Nil(t, resp)

		assert.Equal(t, "test err", err.Error())
	})
}
