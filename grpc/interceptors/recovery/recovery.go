package recovery

import (
	"context"
	"runtime/debug"

	"google.golang.org/grpc"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/errorcode"
	"github.com/jkaveri/goservice/grpc/interceptors/wraperror"
)

// UnaryInterceptor returns a unary server interceptor that recovers from
// panics.
// Panics are logged with stack trace and converted to internal server errors.
func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if p := recover(); p != nil {
				golog.WithContext(ctx).
					Error("panic recovered in gRPC handler",
						golog.String("method", info.FullMethod),
						golog.String("stack", string(debug.Stack())),
						golog.Any("panic", p),
					)

				err = wraperror.ToStructured(
					errorcode.InternalServer("internal server error"),
				)
			}
		}()

		return next(ctx, req)
	}
}
