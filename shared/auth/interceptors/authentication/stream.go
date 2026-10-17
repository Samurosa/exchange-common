package authentication

import (
	"context"
	"time"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"

	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type authenticatedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
}

func (s *authenticatedServerStream) SendMsg(message interface{}) error {
	if err := context.Cause(s.ctx); err != nil {
		return err
	}
	return s.ServerStream.SendMsg(message)
}

func (s *authenticatedServerStream) RecvMsg(message interface{}) error {
	if err := context.Cause(s.ctx); err != nil {
		return err
	}
	return s.ServerStream.RecvMsg(message)
}

func AuthStreamInterceptor(
	tokenParser *jwt.Parser,
	sessionValidator session.Validator,
) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		ctx, err := Authenticate(
			stream.Context(),
			tokenParser,
			sessionValidator,
		)
		if err != nil {
			return err
		}

		requestCtx, _ := sharedContext.GetRequestContext(ctx)
		ctx, stopExpiry := context.WithDeadlineCause(ctx, requestCtx.AccessTokenExpiresAt,
			status.Error(codes.Unauthenticated, "access token expired"))
		defer stopExpiry()
		ctx, cancel := context.WithCancelCause(ctx)
		defer cancel(nil)

		go watchSession(ctx, cancel, requestCtx, sessionValidator)

		stream = &authenticatedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		err = handler(srv, stream)
		if ctx.Err() != nil {
			cause := contextError(ctx)
			logging.FromContext(ctx).Debug("authenticated stream context ended", zap.Error(cause))

			return cause
		}
		return err
	}
}

func watchSession(ctx context.Context, cancel context.CancelCauseFunc, requestCtx sharedContext.RequestContext, validator session.Validator) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			validationCtx, stop := context.WithTimeout(ctx, time.Second)
			err := validator.Validate(validationCtx, requestCtx.JTI)
			stop()
			if ctx.Err() != nil {
				return
			}
			if err != nil {
				cancel(sessionValidationError(ctx, err))
				return
			}
		}
	}
}
