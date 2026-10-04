package authentication

import (
	"context"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"google.golang.org/grpc"
)

func AuthInterceptor(
	tokenParser *jwt.Parser,
	publicMethods map[string]struct{},
	sessionValidator session.Validator,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (interface{}, error) {
		if _, ok := publicMethods[info.FullMethod]; ok {
			return handler(ctx, req)
		}

		ctx, err := Authenticate(ctx, tokenParser, sessionValidator)
		if err != nil {
			return nil, err
		}

		return handler(ctx, req)
	}
}
