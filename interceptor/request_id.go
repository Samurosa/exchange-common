package interceptor

import (
	"context"

	"github.com/google/uuid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
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

func XRequestIDClient() grpc.UnaryClientInterceptor {

    return func(
        ctx context.Context,
        method string,
        req interface{},
        reply interface{},
        cc *grpc.ClientConn,
        invoker grpc.UnaryInvoker,
        opts ...grpc.CallOption,
    ) error {

        md, ok := metadata.FromIncomingContext(ctx)

        if ok {

            ctx = metadata.NewOutgoingContext(
                ctx,
                md,
            )
        }

        return invoker(
            ctx,
            method,
            req,
            reply,
            cc,
            opts...,
        )
    }
}
