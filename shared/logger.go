package interceptor

import (
	"context"
	"time"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type Logger struct {
	log *zap.Logger
}

func NewLogger(log *zap.Logger) *Logger {
	return &Logger{log: log}
}

func (l *Logger) Unary(RequestIDKey any) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		start := time.Now()

		resp, err := handler(ctx, req)

		requestID, _ := ctx.Value(RequestIDKey).(string)

		l.log.Info("grpc request",
			zap.String("request_id", requestID),
			zap.String("method", info.FullMethod),
			zap.Duration("duration", time.Since(start)),
			zap.Error(err),
		)

		return resp, err
	}
}
