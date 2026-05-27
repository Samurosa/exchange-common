package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type ctxKey string

const RequestIDKey ctxKey = "x-request-id"

func RequestID() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {

		requestID := uuid.New().String()
		ctx = context.WithValue(ctx, RequestIDKey, requestID)

		return handler(ctx, req)
	}
}
