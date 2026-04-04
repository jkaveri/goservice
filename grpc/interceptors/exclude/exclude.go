package interceptors

import (
	"context"
	"strings"

	"golang.org/x/exp/slices"
	"google.golang.org/grpc"
)

func WithExcludes(
	mw grpc.UnaryServerInterceptor,
	methods ...string,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		method := getMethodName(info.FullMethod)

		if slices.Contains(methods, method) {
			return handler(ctx, req)
		}

		return mw(ctx, req, info, handler)
	}
}

func getMethodName(fullMethod string) string {
	i := strings.LastIndex(fullMethod, "/")
	if i > -1 && i < len(fullMethod)-1 {
		return fullMethod[i+1:]
	}

	return fullMethod
}
