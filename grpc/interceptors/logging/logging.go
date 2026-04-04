package logging

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"

	golog "github.com/jkaveri/golog/v2"
	"github.com/jkaveri/goservice/clock"
)

func UnaryInterceptor(
	c clock.Clock,
	pruneLogMetadata bool,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		next grpc.UnaryHandler,
	) (resp interface{}, err error) {
		var (
			start = c.Now()
			log   = golog.WithContext(ctx).
				With(golog.String("endpoint", info.FullMethod))
		)

		md, _ := metadata.FromIncomingContext(ctx)

		receiveAttrs := []golog.Attr{}
		if !pruneLogMetadata {
			receiveAttrs = append(receiveAttrs, golog.Int("metadata_keys", len(md)))
		}
		log.Info("receive request", receiveAttrs...)

		// Handle request
		resp, err = next(ctx, req)

		latency := c.Now().Sub(start).Milliseconds()

		log = log.With(golog.Int64("latency_ms", latency))

		if err != nil {
			// Outer wraperror interceptor logs Error for the same failure.
			log.WithError(err).Info("handler returned error")
		} else {
			log.Info(
				"response success",
				golog.String("response_type", fmt.Sprintf("%T", resp)),
			)
		}

		return resp, err
	}
}
