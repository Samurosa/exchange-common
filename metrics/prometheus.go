package interceptor

import (
	"context"
	"time"

	"google.golang.org/grpc"
)

var (
	grpcRequestsTotal = 0
	grpcLatencySumMs  = int64(0)
)

func Prometheus() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()

		resp, err := handler(ctx, req)

		duration := time.Since(start).Milliseconds()

		grpcRequestsTotal++
		grpcLatencySumMs += duration

		return resp, err
	}
}
