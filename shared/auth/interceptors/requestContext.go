package interceptors

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
)

func RequestContextInterceptor(
	log *zap.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		log := log.Named("RequestContextInterceptor")
		requestCtx := sharedContext.RequestContext{}

		ctx = sharedContext.WithRequestContext(
			ctx,
			requestCtx,
		)
		log.Debug("made request context")

		return handler(ctx, req)
	}
}
