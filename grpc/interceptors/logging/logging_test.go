package logging_test

import (
	"context"
	"os"
	"strings"
	"testing"
	"time"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/grpc/interceptors/logging"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
)

type stubClock struct {
	times []time.Time
	idx   int
}

func (s *stubClock) Now() time.Time {
	t := s.times[s.idx]
	s.idx++
	return t
}

func (s *stubClock) Sleep(d time.Duration) {}

func Test_Logging(t *testing.T) {
	f, err := os.CreateTemp(t.TempDir(), "goservice-log")
	require.NoError(t, err)
	path := f.Name()
	require.NoError(t, f.Close())

	require.NoError(t, golog.InitDefault(golog.Config{
		Format: golog.FormatText,
		Output: path,
		Level:  golog.LevelInfo,
	}))

	now := time.Now()
	now2 := now.Add(400 * time.Millisecond)
	md := metadata.New(map[string]string{
		"x-request-id": "123",
	})
	ctx := metadata.NewIncomingContext(
		context.Background(),
		md,
	)

	c := &stubClock{times: []time.Time{now, now2}}

	mw := logging.UnaryInterceptor(c, true)
	resp, err := mw(
		ctx,
		"test",
		&grpc.UnaryServerInfo{
			Server:     "test",
			FullMethod: "TestLogging",
		},
		func(ctx context.Context, req interface{}) (interface{}, error) {
			return "test response", nil
		},
	)

	assert.NoError(t, err)
	assert.Equal(t, "test response", resp)

	data, err := os.ReadFile(path)
	require.NoError(t, err)
	out := string(data)
	assert.True(t, strings.Contains(out, "receive request"), "expected 'receive request' in output: %s", out)
	assert.True(t, strings.Contains(out, "response success"), "expected 'response success' in output: %s", out)
	assert.True(t, strings.Contains(out, "TestLogging"), "expected 'TestLogging' in output: %s", out)
}
