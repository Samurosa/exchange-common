package logging

import (
	"context"
	"time"

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
		started := time.Now()
		requestLog := log.With(zap.String("grpc_method", info.FullMethod), zap.String("grpc_type", "stream"))
		ctx := context.WithValue(stream.Context(), loggerContextKey{}, requestLog)
		requestLog.Debug("grpc stream started")

		err := handler(srv, &loggerServerStream{
			ServerStream: stream,
			ctx:          ctx,
		})
		logCompletion(requestLog, started, err)
		return err
	}
}
