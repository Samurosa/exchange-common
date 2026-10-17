package recovery

import (
	"context"

	"google.golang.org/grpc"
)

// RecoveryInterceptor
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {

		defer func() {
			if recovered := recover(); recovered != nil {
				resp = nil
				err = panicError(ctx, recovered)
			}
		}()

		return handler(ctx, req)
	}
}
