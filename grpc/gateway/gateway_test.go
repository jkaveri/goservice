package grpcgateway

import (
	"context"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"github.com/jkaveri/goservice"
)

func TestConfigureAPIDoc(t *testing.T) {
	type Args struct {
		config     APIDocConfig
		method     string
		path       string
		setupFiles func(t *testing.T, dir string)
	}
	type Expects struct {
		statusCode int
		location   string
		body       string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "disabled-no-handlers-registered",
			args: Args{
				config: APIDocConfig{
					Enabled:      false,
					Path:         "/api-doc/",
					SwaggerUIDir: "/tmp",
				},
				method:     http.MethodGet,
				path:       "/api-doc/",
				setupFiles: nil,
			},
			expects: Expects{
				statusCode: http.StatusNotFound,
				location:   "",
				body:       "404 page not found\n",
			},
		},
		{
			name: "redirects-path-without-trailing-slash-to-path-with-trailing-slash",
			args: Args{
				config: APIDocConfig{
					Enabled:      true,
					Path:         "/api-doc",
					SwaggerUIDir: "", // set in Arrange
				},
				method: http.MethodGet,
				path:   "/api-doc",
				setupFiles: func(t *testing.T, dir string) {
					require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("swagger ui"), 0o644))
				},
			},
			expects: Expects{
				statusCode: http.StatusMovedPermanently,
				location:   "/api-doc/",
				body:       "<a href=\"/api-doc/\">Moved Permanently</a>.\n\n",
			},
		},
		{
			name: "serves-files-from-swagger-ui-dir-at-path-with-trailing-slash",
			args: Args{
				config: APIDocConfig{
					Enabled:      true,
					Path:         "/api-doc/",
					SwaggerUIDir: "", // set in Arrange
				},
				method: http.MethodGet,
				path:   "/api-doc/",
				setupFiles: func(t *testing.T, dir string) {
					require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("swagger ui"), 0o644))
				},
			},
			expects: Expects{
				statusCode: http.StatusOK,
				location:   "",
				body:       "swagger ui",
			},
		},
		{
			name: "serves-nested-file-with-path-without-trailing-slash-in-config",
			args: Args{
				config: APIDocConfig{
					Enabled:      true,
					Path:         "/docs",
					SwaggerUIDir: "", // set in Arrange
				},
				method: http.MethodGet,
				path:   "/docs/static/swagger.json",
				setupFiles: func(t *testing.T, dir string) {
					require.NoError(t, os.MkdirAll(filepath.Join(dir, "static"), 0o755))
					require.NoError(t, os.WriteFile(filepath.Join(dir, "static", "swagger.json"), []byte(`{"openapi":"3.0"}`), 0o644))
				},
			},
			expects: Expects{
				statusCode: http.StatusOK,
				location:   "",
				body:       `{"openapi":"3.0"}`,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 1) Arrange
			args := tc.args
			if args.setupFiles != nil {
				tmpDir := t.TempDir()
				args.setupFiles(t, tmpDir)
				args.config.SwaggerUIDir = tmpDir
			}

			// 2) Construct
			mux := http.NewServeMux()
			ConfigureAPIDoc(mux, args.config)

			// 3) Act
			req := httptest.NewRequest(args.method, args.path, nil)
			rec := httptest.NewRecorder()
			mux.ServeHTTP(rec, req)

			// 4) Assert
			got := Expects{
				statusCode: rec.Code,
				location:   rec.Header().Get("Location"),
				body:       rec.Body.String(),
			}
			assert.Equal(t, tc.expects, got)
		})
	}
}

func TestGetMux(t *testing.T) {
	type Args struct {
		handler http.Handler
	}
	type Expects struct {
		wantPanic bool
		panicMsg  string
	}

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "returns-mux-when-handler-is-runtime-serve-mux",
			args: Args{
				handler: runtime.NewServeMux(),
			},
			expects: Expects{
				wantPanic: false,
				panicMsg:  "",
			},
		},
		{
			name: "panics-when-handler-is-not-runtime-serve-mux",
			args: Args{
				handler: http.NewServeMux(),
			},
			expects: Expects{
				wantPanic: true,
				panicMsg:  "handler is not a runtime.ServeMux",
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.expects.wantPanic {
				require.PanicsWithValue(t, tc.expects.panicMsg, func() {
					GetMux(tc.args.handler)
				})
				return
			}
			got := GetMux(tc.args.handler)
			assert.NotNil(t, got)
			assert.Same(t, tc.args.handler, got)
		})
	}
}

func TestCreateHandler(t *testing.T) {
	type Args struct {
		ctx        context.Context
		srv        goservice.Service
		grpcPort   int
		httpServer *http.Server
	}
	type Expects struct {
		wantErr bool
	}

	lis, err := net.Listen("tcp", ":0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = lis.Close() })

	grpcPort := lis.Addr().(*net.TCPAddr).Port
	gsrv := grpc.NewServer()
	go func() { _ = gsrv.Serve(lis) }()
	t.Cleanup(func() { gsrv.Stop() })

	testCases := []struct {
		name    string
		args    Args
		expects Expects
	}{
		{
			name: "returns-handler-when-register-succeeds",
			args: Args{
				ctx:        context.Background(),
				srv:        nil,
				grpcPort:   grpcPort,
				httpServer: nil,
			},
			expects: Expects{
				wantErr: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			ctx, cancel := context.WithCancel(tc.args.ctx)
			t.Cleanup(cancel)

			registerFunc := func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error {
				return nil
			}

			got := CreateHandler(
				ctx,
				tc.args.srv,
				tc.args.grpcPort,
				tc.args.httpServer,
				registerFunc,
			)

			require.NotNil(t, got)

			mux := GetMux(got)
			assert.NotNil(t, mux)
		})
	}
}
