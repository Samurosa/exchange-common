package logging

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func LoggerInterceptor(log *zap.Logger) grpc.UnaryServerInterceptor {
	if log == nil {
		log = noopLogger
	}

	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		started := time.Now()

		requestLog := log.With(
			zap.String("grpc_method", info.FullMethod),
			zap.String("grpc_type", "unary"),
		)

		ctx = context.WithValue(ctx, loggerContextKey{}, requestLog)

		resp, err := handler(ctx, req)

		logCompletion(WithTrace(ctx, requestLog), started, err)

		return resp, err
	}
}
