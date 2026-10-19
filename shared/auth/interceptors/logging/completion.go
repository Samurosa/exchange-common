package logging

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func logCompletion(log *zap.Logger, started time.Time, err error) {
	code := status.Code(err)

	switch {
	case errors.Is(err, context.Canceled):
		code = codes.Canceled
	case errors.Is(err, context.DeadlineExceeded):
		code = codes.DeadlineExceeded
	}

	level := grpcLogLevel(code)

	fields := []zap.Field{
		zap.String("grpc_code", code.String()),
		zap.Duration("duration", time.Since(started)),
	}

	if err != nil {
		fields = append(fields, zap.Error(err))
	}

	log.Log(level, "grpc request completed", fields...)
}

func grpcLogLevel(code codes.Code) zapcore.Level {
	switch code {
	case codes.OK, codes.Canceled, codes.NotFound,
		codes.AlreadyExists:
		return zapcore.DebugLevel

	case codes.InvalidArgument, codes.FailedPrecondition,
		codes.OutOfRange:
		return zapcore.InfoLevel

	case codes.DeadlineExceeded, codes.PermissionDenied,
		codes.ResourceExhausted, codes.Aborted,
		codes.Unauthenticated:
		return zapcore.WarnLevel

	case codes.Unknown, codes.Unimplemented, codes.Internal,
		codes.Unavailable, codes.DataLoss:
		return zapcore.ErrorLevel

	default:
		return zapcore.WarnLevel
	}
}
