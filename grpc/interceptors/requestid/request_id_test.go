package requestid

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type stubGenerator struct {
	v     string
	calls int
}

func (s *stubGenerator) Generate() string {
	s.calls++
	return s.v
}

func Test_RequestID(t *testing.T) {
	t.Run("generate-new-id", func(t *testing.T) {
		ctx := context.Background()
		called := false
		uid := "744cfc23-da5f-4425-84e9-a5ff6d89408c"
		gen := &stubGenerator{v: uid}

		mw := UnaryInterceptor(
			gen,
		)

		resp, err := mw(
			ctx,
			"test",
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestLogging",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true

				md, exist := metadata.FromIncomingContext(ctx)
				assert.True(t, exist)

				val := md.Get(RequestIDKey)
				assert.NotEmpty(t, val)

				assert.Equal(t, uid, val[0])

				return "test response", nil
			},
		)

		assert.True(t, called)
		assert.NoError(t, err)
		assert.Equal(t, "test response", resp)
		assert.Equal(t, 1, gen.calls)
	})

	t.Run("reuse-existing-id", func(t *testing.T) {
		uid := "744cfc23-da5f-4425-84e9-a5ff6d89408c"
		ctx := metadata.NewIncomingContext(
			context.Background(),
			metadata.New(map[string]string{
				RequestIDKey: uid,
			}),
		)
		called := false

		gen := &stubGenerator{v: uid}

		mw := UnaryInterceptor(
			gen,
		)

		resp, err := mw(
			ctx,
			"test",
			&grpc.UnaryServerInfo{
				Server:     "test",
				FullMethod: "TestLogging",
			},
			func(ctx context.Context, req interface{}) (interface{}, error) {
				called = true

				md, exist := metadata.FromIncomingContext(ctx)
				assert.True(t, exist)

				val := md.Get(RequestIDKey)
				assert.NotEmpty(t, val)

				assert.Equal(t, uid, val[0])

				return "test response", nil
			},
		)

		assert.True(t, called)
		assert.NoError(t, err)
		assert.Equal(t, "test response", resp)
		assert.Zero(t, gen.calls)
	})
}
