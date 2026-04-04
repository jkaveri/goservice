package main

import (
	"context"
	"encoding/json"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice"
	servicev1 "github.com/jkaveri/goservice/examples/grpcgateway/proto/gen/go/service/v1"
	"github.com/jkaveri/goservice/examples/grpcgateway/service"
	"github.com/jkaveri/goservice/errorcode"
	grpcgateway "github.com/jkaveri/goservice/grpc/gateway"
)

func TestGateway_ErrorHandler_StructuredJSON(t *testing.T) {
	type Args struct {
		path string
	}
	type Expects struct {
		statusCode int
		code       string
		message    string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "get-echo-error-returns-structured-body",
			args: Args{path: "/echo/error"},
			expects: Expects{
				statusCode: http.StatusNotFound,
				code:       errorcode.CodeNotFound,
				message:    "example resource is missing",
			},
		},
		{
			name: "get-echo-wrapped-error-user-message-not-full-wrap-chain",
			args: Args{path: "/echo/error/wrapped"},
			expects: Expects{
				statusCode: http.StatusInternalServerError,
				code:       errorcode.CodeInternalServer,
				message:    "Checkout could not be completed. Please try again later.",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			require.NoError(t, golog.InitDefault(golog.Config{
				Level:  golog.LevelError,
				Format: golog.FormatText,
				Output: os.DevNull,
			}))

			ctx, cancel := context.WithCancel(context.Background())
			t.Cleanup(cancel)

			lis, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			t.Cleanup(func() { _ = lis.Close() })

			grpcPort := lis.Addr().(*net.TCPAddr).Port

			gsrv := grpc.NewServer(grpc.ChainUnaryInterceptor(goservice.DefaultInterceptors()...))

			svc, err := service.New(ctx)
			require.NoError(t, err)

			servicev1.RegisterEchoServiceServer(gsrv, svc.(servicev1.EchoServiceServer))

			go func() { _ = gsrv.Serve(lis) }()
			t.Cleanup(gsrv.Stop)

			h := grpcgateway.CreateHandler(ctx, svc, grpcPort, nil, servicev1.RegisterEchoServiceHandler)
			ts := httptest.NewServer(h)
			t.Cleanup(ts.Close)

			res, err := http.Get(ts.URL + tc.args.path)
			require.NoError(t, err)
			t.Cleanup(func() { _ = res.Body.Close() })

			assert.Equal(t, tc.expects.statusCode, res.StatusCode)

			var body map[string]interface{}
			require.NoError(t, json.NewDecoder(res.Body).Decode(&body))

			assert.Equal(t, tc.expects.code, body["code"])
			assert.Equal(t, tc.expects.message, body["message"])
			_, hasMeta := body["metadata"]
			assert.True(t, hasMeta)
		})
	}
}
