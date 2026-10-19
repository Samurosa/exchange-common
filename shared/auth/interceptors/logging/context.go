package logging

import (
	"context"

	"go.uber.org/zap"
)

type loggerContextKey struct{}

var noopLogger = zap.NewNop()

func FromContext(ctx context.Context) *zap.Logger {
	log, ok := ctx.Value(loggerContextKey{}).(*zap.Logger)
	if !ok || log == nil {
		return noopLogger
	}

	return WithTrace(ctx, log)
}
