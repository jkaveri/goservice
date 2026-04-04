package goservice

import (
	"context"
	"net"
	"net/http"
	"time"

	golog "github.com/jkaveri/golog/v2"
	errors "github.com/jkaveri/goservice/errors"
)

type ServerCustomizer interface {
	SetupHTTPServer(server *http.Server)
}

func startHTTPServer(ctx context.Context, cfg *Config, svc any) {
	log := golog.WithContext(ctx)

	httpService, ok := svc.(HTTPService)
	if !ok {
		return
	}

	log.Info("starting http server")

	// register http routes
	server := &http.Server{
		Addr:              tcpAddrFromPort(cfg.HTTPServer.Port),
		ReadHeaderTimeout: 5 * time.Second,
		ReadTimeout:       15 * time.Second,
		WriteTimeout:      15 * time.Second,
		IdleTimeout:       120 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
	}

	if customizer, ok := httpService.(ServerCustomizer); ok {
		customizer.SetupHTTPServer(server)
	}

	// register http handler
	server.Handler = httpService.RegisterHTTP(ctx, server)

	shutdownWg.Add(1)

	go func() {
		<-ctx.Done()

		if err := server.Shutdown(ctx); err != nil {
			log.WithError(err).Error("cannot shutdown http server")
			return
		}

		log.Info("http server is shutdown")
		shutdownWg.Done()
	}()

	go func() {
		log.Info(
			"http server is listening",
			golog.Int("port", cfg.HTTPServer.Port),
		)

		if err := server.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error(
				"cannot start http server",
				golog.Int("port", cfg.HTTPServer.Port),
			)
		}
	}()
}
