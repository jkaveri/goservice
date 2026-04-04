package requestid

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	"github.com/jkaveri/goservice/idgen"
)

const RequestIDKey = "request_id"

func UnaryInterceptor(
	idGen idgen.Generator,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		md, exist := metadata.FromIncomingContext(ctx)
		if !exist {
			md = metadata.New(map[string]string{})
		}

		currentID := md.Get(RequestIDKey)

		if len(currentID) == 0 {
			// fallback to trace_id
			traceID := md.Get("x-trace-id")
			if len(traceID) > 0 {
				currentID = traceID
				md.Set(RequestIDKey, currentID...)
			}
		}

		if len(currentID) == 0 {
			currentID = []string{idGen.Generate()}
			md.Set(RequestIDKey, currentID...)
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		resp, err = handler(ctx, req)

		_ = grpc.SetHeader(ctx, metadata.Pairs(RequestIDKey, currentID[0]))

		return resp, err
	}
}

func GetRequestID(ctx context.Context) string {
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return ""
	}

	a := md.Get(RequestIDKey)
	if len(a) == 0 {
		return ""
	}

	return a[0]
}
