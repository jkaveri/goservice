package goservice

import (
	"context"
	"net/http"
	"os/signal"
	"sync"
	"syscall"

	"google.golang.org/grpc"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/check"
)

// Factory is a function that creates a service instance with a given context
// and returns the service along with any error that occurred during creation
type Factory = func(ctx context.Context) (Service, error)

// OnStartListener is an interface for services that need to perform actions
// when starting
type OnStartListener interface {
	OnStart(context.Context) error
}

// OnStopListener is an interface for services that need to perform cleanup
// actions when stopping
type OnStopListener interface {
	OnStop() error
}

// GRPCService is an interface for services that provide gRPC functionality
type GRPCService interface {
	// RegisterGRPC registers the gRPC service with the provided server
	RegisterGRPC(ctx context.Context, s *grpc.Server)

	// GRPCServerOptions returns the options that will be used for all gRPC
	// services
	GRPCServerOptions() []grpc.ServerOption

	// GRPCUnaryInterceptors returns the interceptors that will be used for all
	// gRPC services
	GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor
}

// HTTPService is an interface for services that provide HTTP functionality
type HTTPService interface {
	// RegisterHTTP registers the HTTP routes with the provided Fiber app
	RegisterHTTP(ctx context.Context, s *http.Server) http.Handler
}

// CronJobService is an interface for services that need to schedule periodic
// tasks
type CronJobService interface {
	// RegisterJobs registers cron jobs with the provided scheduler
	RegisterJobs(ctx context.Context, service Service)
}

// Service is a placeholder interface for any service implementation
// To provide HTTP functionality, implement HTTPService
// To provide gRPC functionality, implement GRPCService
type Service any

// shutdownWg is a wait group to wait for all servers to shutdown
var shutdownWg sync.WaitGroup

// Run initializes and starts all service components including gRPC, HTTP, and
// cron jobs
// It handles graceful shutdown on receiving interrupt signals
func Run(srvFactory Factory) {
	ctx, stop := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)

	defer stop()

	// load cfg
	cfg := loadConfig()

	// set current config
	SetCurrentConfig(&cfg)

	// init logger
	initLogger(&cfg)

	log := golog.WithContext(ctx)
	log.Info("log level", golog.String("level", cfg.Log.Level))

	log.Debug("creating service instance")

	// create service instance
	srv, err := srvFactory(ctx)
	check.PanicIfError(err)

	// start grpc server
	startGRPCServer(ctx, &cfg, srv)

	// start http server
	startHTTPServer(ctx, &cfg, srv)

	if listener, ok := srv.(OnStartListener); ok {
		err = listener.OnStart(ctx)
		check.PanicIfError(err)
	}

	// start health check server
	startHealthCheckServer(ctx, cfg.HealthServer.Port)

	<-ctx.Done()

	shutdownWg.Wait()

	if listener, ok := srv.(OnStopListener); ok {
		err = listener.OnStop()
		check.PanicIfError(err)
	}
}

var currentConfig *Config

func GetCurrentConfig() Config {
	return *currentConfig
}

func SetCurrentConfig(cfg *Config) {
	currentConfig = cfg
}
