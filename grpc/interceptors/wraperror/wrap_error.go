package wraperror

import (
	"context"

	"google.golang.org/grpc"

	golog "github.com/jkaveri/golog/v2"
)

//nolint:revive
func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Handle request
		resp, err = next(ctx, req)
		if err == nil {
			return resp, nil
		}

		golog.WithContext(ctx).
			With(golog.String("endpoint", info.FullMethod)).
			WithError(err).
			Error("rpc error")

		return nil, ToStructured(err)
	}
}
