package logging

import (
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func logCompletion(log *zap.Logger, started time.Time, err error) {

	result, ok := status.FromError(err)
	if !ok {
		result = status.FromContextError(err)
	}
	code := result.Code()
	level := zapcore.DebugLevel
	switch code {
	case codes.DeadlineExceeded, codes.PermissionDenied, codes.ResourceExhausted, codes.Aborted:
		level = zapcore.WarnLevel
	case codes.Unknown, codes.Unimplemented, codes.Internal, codes.Unavailable, codes.DataLoss:
		level = zapcore.ErrorLevel
	}
	fields := []zap.Field{
		zap.String("grpc_code", code.String()),
		zap.Duration("duration", time.Since(started)),
	}
	if err != nil {
		fields = append(fields, zap.Error(err))
	}
	log.Log(level, "grpc request completed", fields...)
}
