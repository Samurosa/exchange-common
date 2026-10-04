package authentication

import (
	"context"
	"errors"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func sessionValidationError(ctx context.Context, err error) error {
	if ctx.Err() != nil {
		return contextError(ctx)
	}

	log := logging.FromContext(ctx).Named("authenticate")
	if errors.Is(err, session.ErrNotFound) {
		log.Debug("session not found", zap.Error(err))
		return status.Error(codes.Unauthenticated, "invalid session")
	}

	log.Error("session validation unavailable", zap.Error(err))
	return status.Error(codes.Unavailable, "session validation unavailable")
}

func contextError(ctx context.Context) error {
	if cause := context.Cause(ctx); cause != nil {
		if _, ok := status.FromError(cause); ok {
			return cause
		}
	}
	return status.FromContextError(ctx.Err()).Err()
}
