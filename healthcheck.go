package goservice

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"time"

	golog "github.com/jkaveri/golog/v2"
	errors "github.com/jkaveri/goservice/errors"
)

func startHealthCheckServer(ctx context.Context, port int) {
	log := golog.WithContext(ctx)

	log.Debug("starting health check server")

	s := &http.Server{
		Addr:              fmt.Sprintf(":%d", port),
		ReadHeaderTimeout: 10 * time.Second,
		BaseContext: func(_ net.Listener) context.Context {
			return ctx
		},
		Handler: http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
			w.WriteHeader(http.StatusOK)
		}),
	}

	shutdownWg.Add(1)

	go func() {
		<-ctx.Done()

		if err := s.Shutdown(ctx); err != nil {
			log.WithError(err).Error("cannot shutdown health check server")
			return
		}

		log.Info("health check server is shutdown")
		shutdownWg.Done()
	}()

	go func() {
		log.Info("health check server is listening", golog.Int("port", port))

		if err := s.ListenAndServe(); err != nil &&
			!errors.Is(err, http.ErrServerClosed) {
			log.WithError(err).Error(
				"cannot start health check server",
				golog.Int("port", port),
			)

			return
		}
	}()
}
