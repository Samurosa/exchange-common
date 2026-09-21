package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"

	"go.uber.org/zap"
	"google.golang.org/grpc"
)

type authenticatedServerStream struct {
	grpc.ServerStream
	ctx context.Context
}

func (s *authenticatedServerStream) Context() context.Context {
	return s.ctx
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

		// Replace stream context with authenticated context.
		stream = &authenticatedServerStream{
			ServerStream: stream,
			ctx:          ctx,
		}

		return handler(srv, stream)
	}
}
