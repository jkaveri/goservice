package service

import (
	"context"

	"google.golang.org/grpc"

	"github.com/jkaveri/goservice"
	servicev1 "github.com/jkaveri/goservice/examples/echogrpc/proto/gen/service/v1"
)

var (
	_ goservice.GRPCService       = (*Service)(nil)
	_ servicev1.EchoServiceServer = (*Service)(nil)
)

func New(ctx context.Context) (goservice.Service, error) {
	return &Service{}, nil
}

type Service struct {
	servicev1.UnimplementedEchoServiceServer
}

// Echo implements servicev1.EchoServiceServer.
func (s *Service) Echo(context.Context, *servicev1.EchoRequest) (*servicev1.EchoResponse, error) {
	return &servicev1.EchoResponse{
		Message: "Hello, world!",
	}, nil
}

// GRPCServerOptions implements goservice.GRPCService.
func (s *Service) GRPCServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

// GRPCUnaryInterceptors implements goservice.GRPCService.
func (s *Service) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return []grpc.UnaryServerInterceptor{}
}

// RegisterGRPC implements goservice.GRPCService.
func (s *Service) RegisterGRPC(ctx context.Context, srv *grpc.Server) {
	servicev1.RegisterEchoServiceServer(srv, s)
}
