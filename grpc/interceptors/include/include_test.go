package interceptors_test

import (
	"context"
	"testing"

	include "github.com/jkaveri/goservice/grpc/interceptors/include"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
)

func Test_WithInclude(t *testing.T) {
	t.Run("should-not-call-mw", func(t *testing.T) {
		ctx := context.Background()

		called := false
		mw := include.WithIncludes(
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
		assert.False(t, called)
	})

	t.Run("should-call", func(t *testing.T) {
		ctx := context.Background()

		called := false
		mw := include.WithIncludes(
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
			func(_ context.Context, req interface{}) (interface{}, error) {
				return "test response 1", nil
			},
		)

		assert.NoError(t, err)
		assert.Equal(t, "test response 1", resp)
		assert.True(t, called)
	})
}
