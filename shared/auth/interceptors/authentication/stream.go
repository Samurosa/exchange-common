package authentication

import (
	"context"
	"time"

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
	log *zap.Logger,
	tokenParser *jwt.Parser,
	sessionValidator session.Validator,
) grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) error {
		log := log.Named("auth stream interceptor")

		ctx, err := Authenticate(
			stream.Context(),
			log,
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

		// Expiry cancels the stream at exp; session revocation is checked every second.
		// Handlers must observe stream.Context(), as with ordinary gRPC cancellation.
		go watchSession(ctx, cancel, requestCtx, sessionValidator)

		// Replace stream context with authenticated context.
		stream = &authenticatedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		err = handler(srv, stream)
		if cause := context.Cause(ctx); status.Code(cause) == codes.Unauthenticated {
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
			if err != nil {
				cancel(status.Error(codes.Unauthenticated, "invalid session"))
				return
			}
		}
	}
}
