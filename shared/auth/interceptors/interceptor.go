package interceptors

import (
	"context"
	"strings"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"go.uber.org/zap"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

var publicMethods = map[string]struct{}{
	"/user.UserService/Login":        {},
	"/user.UserService/Registration": {},
	"/user.UserService/RefreshToken": {},
}

func AuthInterceptor(
	log *zap.Logger,
	tokenManager jwt.Parser,
) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		log := log.Named("auth interceptor")

		if _, ok := publicMethods[info.FullMethod]; ok {
			log.Debug("A verification token is not required for this RPS.")
			return handler(ctx, req)
		}

		meta, ok := metadata.FromIncomingContext(ctx)
		if !ok {
			log.Error("missing metadata")
			return nil, status.Errorf(codes.Unauthenticated, "Unauthenticated")
		}

		values := meta.Get("authorization")
		if len(values) == 0 {
			log.Error("missing token")
			return nil, status.Errorf(codes.Unauthenticated, "missing token")
		}

		token := strings.TrimPrefix(values[0], "Bearer ")

		claims, err := tokenManager.ParseAccessToken(token)
		if err != nil {
			log.Error("invalid token", zap.Error(err))
			return nil, status.Error(codes.Unauthenticated, "invalid token")
		}

		ctx, err = sharedContext.UpdateRequestContext(ctx,
			func(baseContext *sharedContext.RequestContext) {
				baseContext.Principal.UserID = claims.UserID
				baseContext.Principal.Role = claims.Role
				baseContext.Metadata.DeviceID = claims.Device
				baseContext.JTI = claims.JTI
			})
		if err != nil {
			log.Error("failed to update request context", zap.Error(err))
			return nil, status.Error(codes.Internal, "failed to update request context")
		}
		log.Debug("success pulling data in request context from access token", zap.String("id", claims.JTI))

		return handler(ctx, req)
	}
}
