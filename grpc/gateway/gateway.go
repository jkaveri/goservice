package grpcgateway

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/grpc-ecosystem/grpc-gateway/v2/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice"
)

// HTTPServiceCustomizer is an interface to customize http server
type HTTPServiceCustomizer interface {
	// GetMuxOptions returns the options that will be used for http service
	//
	// defaultOptions is the default options that will be used for http service
	// you can add or remove options from defaultOptions and the server will use
	// the returned options as the final options
	//
	// return value is the options that will be used forthe http service
	GetMuxOptions(
		defaultOptions []runtime.ServeMuxOption,
	) []runtime.ServeMuxOption
}

func CreateHandler(
	ctx context.Context,
	srv goservice.Service,
	grpcPort int,
	httpServer *http.Server,
	registerFunc func(ctx context.Context, mux *runtime.ServeMux, conn *grpc.ClientConn) error,
) http.Handler {
	log := golog.WithContext(ctx)

	conn, err := grpc.NewClient(
		fmt.Sprintf(":%d", grpcPort),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		log.WithError(err).
			Error("unable to dial grpc", golog.Int("port", grpcPort))
		panic(err)
	}

	go func() {
		<-ctx.Done()

		if err1 := conn.Close(); err1 != nil {
			log.WithError(err1).Error("unable to close grpc connection")
		}
	}()

	defaultMuxOptions := []runtime.ServeMuxOption{
		runtime.WithErrorHandler(ErrorHandler),
	}

	// customizer
	if customizer, ok := srv.(HTTPServiceCustomizer); ok {
		defaultMuxOptions = customizer.GetMuxOptions(defaultMuxOptions)
	}

	gwmux := runtime.NewServeMux(defaultMuxOptions...)

	err = registerFunc(
		ctx,
		gwmux,
		conn,
	)
	if err != nil {
		log.WithError(err).Error("register http error")
		panic(err)
	}

	return gwmux
}

func GetMux(h http.Handler) *runtime.ServeMux {
	mux, ok := h.(*runtime.ServeMux)
	if !ok {
		panic("handler is not a runtime.ServeMux")
	}

	return mux
}

type APIDocConfig struct {
	Enabled      bool
	Path         string
	SwaggerUIDir string
}

func ConfigureAPIDoc(mux *http.ServeMux, cfg APIDocConfig) {
	if !cfg.Enabled {
		return
	}

	path := cfg.Path
	if path[len(path)-1] != '/' {
		path += "/"
	}

	mux.HandleFunc(
		strings.TrimRight(path, "/"),
		func(w http.ResponseWriter, r *http.Request) {
			http.Redirect(w, r, path, http.StatusMovedPermanently)
		},
	)

	mux.Handle(
		path,
		http.StripPrefix(
			path,
			http.FileServer(http.Dir(cfg.SwaggerUIDir)),
		),
	)
}
