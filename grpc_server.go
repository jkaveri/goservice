package goservice

import (
	"context"
	"net"

	"google.golang.org/grpc"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/clock"
	"github.com/jkaveri/goservice/grpc/interceptors/logging"
	"github.com/jkaveri/goservice/grpc/interceptors/recovery"
	"github.com/jkaveri/goservice/grpc/interceptors/requestid"
	"github.com/jkaveri/goservice/grpc/interceptors/validate"
	"github.com/jkaveri/goservice/grpc/interceptors/wraperror"
	"github.com/jkaveri/goservice/idgen"
)

func startGRPCServer(ctx context.Context, cfg *Config, srv any) {
	log := golog.WithContext(ctx)

	grpcService, ok := srv.(GRPCService)
	if !ok {
		return
	}

	log.Debug("starting grpc server")

	// Create a listener on TCP port
	lis, err := net.Listen("tcp", tcpAddrFromPort(cfg.GRPCServer.Port))
	if err != nil {
		log.WithError(err).
			Error("could not listen on port", golog.Int("port", cfg.GRPCServer.Port))
		panic(err)
	}

	// interceptors
	interceptors := grpcService.GRPCUnaryInterceptors()
	serverOptions := append(
		[]grpc.ServerOption{
			grpc.ChainUnaryInterceptor(interceptors...),
		},
		grpcService.GRPCServerOptions()...,
	)

	// Create a gRPC server object
	s := grpc.NewServer(serverOptions...)

	grpcService.RegisterGRPC(ctx, s)

	shutdownWg.Add(1)

	go func() {
		<-ctx.Done()

		s.GracefulStop()

		if err1 := lis.Close(); err1 != nil {
			log.WithError(err1).Error("could not close grpc listener")
		}

		log.Info("grpc server is shutdown")
		shutdownWg.Done()
	}()

	// Serve gRPC server
	log.Info("grpc listening", golog.Int("port", cfg.GRPCServer.Port))

	go func() {
		if err := s.Serve(lis); err != nil {
			log.WithError(err).
				Error("serve grpc error")
		}
	}()
}

func DefaultInterceptors() []grpc.UnaryServerInterceptor {
	items := []grpc.UnaryServerInterceptor{
		recovery.UnaryInterceptor(),
		wraperror.UnaryInterceptor(),

		requestid.UnaryInterceptor(
			idgen.UUIDV4(),
		),

		logging.UnaryInterceptor(
			clock.Default,
			true,
		),

		validate.UnaryInterceptor(),
	}

	return items
}
