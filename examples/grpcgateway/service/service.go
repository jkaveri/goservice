package service

import (
	"context"
	"fmt"
	"net/http"

	"google.golang.org/grpc"

	"github.com/jkaveri/goservice"
	servicev1 "github.com/jkaveri/goservice/examples/grpcgateway/proto/gen/go/service/v1"
	"github.com/jkaveri/goservice/errorcode"
	apperrors "github.com/jkaveri/goservice/errors"
	grcpgateway "github.com/jkaveri/goservice/grpc/gateway"
)

var (
	_ goservice.GRPCService       = (*Service)(nil)
	_ servicev1.EchoServiceServer = (*Service)(nil)
	_ goservice.HTTPService       = (*Service)(nil)
)

func New(ctx context.Context) (goservice.Service, error) {
	return &Service{}, nil
}

type Service struct {
	servicev1.UnimplementedEchoServiceServer
}

// RegisterHTTP implements goservice.HTTPService.
func (s *Service) RegisterHTTP(
	ctx context.Context,
	httpServer *http.Server,
) http.Handler {
	cfg := goservice.GetCurrentConfig()

	return grcpgateway.CreateHandler(
		ctx,
		s,
		cfg.GRPCServer.Port,
		httpServer,
		servicev1.RegisterEchoServiceHandler,
	)
}

// Echo implements servicev1.EchoServiceServer.
func (s *Service) Echo(ctx context.Context, req *servicev1.EchoRequest) (*servicev1.EchoResponse, error) {
	return &servicev1.EchoResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Message),
	}, nil
}

// EchoError implements servicev1.EchoServiceServer; returns a coded error for gateway error-handler demos/tests.
func (s *Service) EchoError(ctx context.Context, _ *servicev1.EchoErrorRequest) (*servicev1.EchoErrorResponse, error) {
	return nil, errorcode.NotFound("example resource is missing")
}

// EchoWrappedError returns an error built from several errors.Wrap layers (full chain + stacks for logs),
// then WithCode and WithUserMessage so the grpc-gateway JSON body shows a stable app code and a safe message
// instead of the long "msg1: msg2: msg3" string from err.Error().
func (s *Service) EchoWrappedError(
	ctx context.Context,
	_ *servicev1.EchoWrappedErrorRequest,
) (*servicev1.EchoWrappedErrorResponse, error) {
	// Low-level failure (e.g. driver / IO).
	root := apperrors.New("connection reset by peer")
	// Each layer adds context the way service code typically does.
	repoErr := apperrors.Wrap(root, "load order by id from repository")
	svcErr := apperrors.Wrap(repoErr, "getOrderForCheckout")
	// Attach business code and the string clients should see; StructuredError.Error() still returns the full chain for logging.
	return nil, apperrors.WithUserMessage(
		apperrors.WithCode(svcErr, errorcode.CodeInternalServer),
		"Checkout could not be completed. Please try again later.",
	)
}

// GRPCServerOptions implements goservice.GRPCService.
func (s *Service) GRPCServerOptions() []grpc.ServerOption {
	return []grpc.ServerOption{}
}

// GRPCUnaryInterceptors implements goservice.GRPCService.
func (s *Service) GRPCUnaryInterceptors() []grpc.UnaryServerInterceptor {
	return goservice.DefaultInterceptors()
}

// RegisterGRPC implements goservice.GRPCService.
func (s *Service) RegisterGRPC(ctx context.Context, srv *grpc.Server) {
	servicev1.RegisterEchoServiceServer(srv, s)
}
