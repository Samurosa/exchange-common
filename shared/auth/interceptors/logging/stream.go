package logging

import (
	"context"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type loggerServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *loggerServerStream) Context() context.Context {
	return s.ctx
}

// LoggerStreamInterceptor adds a logger with the gRPC method to a stream context.
// The base logger is initialized once by the application and reused for streams.
func LoggerStreamInterceptor(log *zap.Logger) grpc.StreamServerInterceptor {
	if log == nil {
		log = noopLogger
	}

	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		requestLog := log.With(zap.String("grpc_method", info.FullMethod))
		ctx := context.WithValue(stream.Context(), loggerContextKey{}, requestLog)

		return handler(srv, &loggerServerStream{
			ServerStream: stream,
			ctx:          ctx,
		})
	}
}
