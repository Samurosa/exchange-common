package logging

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

// LoggerInterceptor adds a logger with the gRPC method to each request context.
// The base logger is initialized once by the application and reused for requests.
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
		requestLog := log.With(zap.String("grpc_method", info.FullMethod))
		ctx = context.WithValue(ctx, loggerContextKey{}, requestLog)

		return handler(ctx, req)
	}
}
