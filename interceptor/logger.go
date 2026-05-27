package interceptor

import (
	"context"
	"fmt"

	"google.golang.org/grpc"
)

func Logger() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {

		resp, err := handler(ctx, req)

		// позже сюда добавишь zap
		fmt.Println("method:", info.FullMethod, "err:", err)

		return resp, err
	}
}
