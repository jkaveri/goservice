package interceptors_test

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"

	exclude "github.com/jkaveri/goservice/grpc/interceptors/exclude"
)

func Test_WithExcludes(t *testing.T) {
	t.Run("should-call-mw", func(t *testing.T) {
		ctx := context.Background()

		called := false
		mw := exclude.WithExcludes(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler,
			) (resp interface{}, err error) {
				called = true
				return handler(ctx, req)
			},
			"NotTestLogging",
		)

		resp, err := mw(
			ctx,
			"test",
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "test.proto/TestLogging",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				return "test response 1", nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "test response 1", resp)
		assert.True(t, called)
	})

	t.Run("should-not-call", func(t *testing.T) {
		ctx := context.Background()

		called := false
		mw := exclude.WithExcludes(
			func(
				ctx context.Context,
				req interface{},
				info *grpc.UnaryServerInfo,
				handler grpc.UnaryHandler,
			) (resp interface{}, err error) {
				called = true
				return handler(ctx, req)
			},
			"TestLogging",
		)

		resp, err := mw(
			ctx,
			"test",
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestLogging",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				return "test response 1", nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "test response 1", resp)
		assert.False(t, called)
	})
}
