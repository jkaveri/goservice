package recovery_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/jkaveri/goservice/grpc/interceptors/recovery"
)

func Test_Recovery(t *testing.T) {
	t.Run("should catch panic and return internal error", func(t *testing.T) {
		mw := recovery.UnaryInterceptor()

		h := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
			panic("test panic")
		})

		resp, err := mw(
			context.Background(),
			"test req",
			&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
			h,
		)

		assert.Error(t, err)
		assert.Nil(t, resp)
		assert.Equal(t, codes.Internal, status.Code(err))
	})

	t.Run("should pass through normal requests", func(t *testing.T) {
		mw := recovery.UnaryInterceptor()

		h := grpc.UnaryHandler(func(ctx context.Context, req interface{}) (interface{}, error) {
			assert.Equal(t, "test req", req)
			return "test response", nil
		})

		resp, err := mw(
			context.Background(),
			"test req",
			&grpc.UnaryServerInfo{FullMethod: "/test.Service/Method"},
			h,
		)

		assert.NoError(t, err)
		assert.Equal(t, "test response", resp)
	})
}
