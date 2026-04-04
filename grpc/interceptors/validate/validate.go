package validate

import (
	"context"

	"github.com/jkaveri/goservice/errorcode"
	"google.golang.org/grpc"
)

type Validator interface {
	ValidateAll() error
}

func UnaryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		if v, ok := req.(Validator); ok {
			if err := v.ValidateAll(); err != nil {
				return nil, errorcode.InvalidRequest(
					friendlyValidationMessage(err),
				)
			}
		}

		return handler(ctx, req)
	}
}
