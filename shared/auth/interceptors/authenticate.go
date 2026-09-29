package interceptors

import (
	"context"
	"strings"

	"github.com/Samurosa/exchange-common/shared/auth/jwt"
	"github.com/Samurosa/exchange-common/shared/auth/session"
	"github.com/Samurosa/exchange-common/shared/auth/sharedContext"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func Authenticate(
	ctx context.Context,
	log *zap.Logger,
	tokenParser *jwt.Parser,
	sessionValidator session.Validator,
) (context.Context, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		log.Debug("missing metadata")
		return nil, status.Error(
			codes.Unauthenticated,
			"missing authentication metadata",
		)
	}

	values := md.Get("authorization")
	if len(values) == 0 {
		log.Debug("missing authorization header")
		return nil, status.Error(
			codes.Unauthenticated,
			"missing authorization header",
		)
	}

	token, ok := extractBearerToken(values[0])
	if !ok {
		log.Debug("invalid authorization header format")
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid authorization header",
		)
	}

	claims, err := tokenParser.ParseAccessToken(token)
	if err != nil {
		log.Debug(
			"invalid access token",
			zap.Error(err),
		)
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid access token",
		)
	}

	if err := sessionValidator.Validate(ctx, claims.ID); err != nil {
		log.Debug(
			"session validation failed",
			zap.Error(err),
		)
		return nil, status.Error(
			codes.Unauthenticated,
			"invalid session",
		)
	}

	ctx = sharedContext.UpdateRequestContext(
		ctx,
		func(requestCtx *sharedContext.RequestContext) {
			requestCtx.Principal.UserID = claims.UserID
			requestCtx.Principal.Role = claims.Role
			requestCtx.Metadata.DeviceID = claims.DeviceID
			requestCtx.JTI = claims.ID
		},
	)

	log.Debug(
		"authentication successful",
		zap.String("user_id", claims.UserID),
		zap.String("jti", claims.ID),
	)

	return ctx, nil
}

func extractBearerToken(value string) (string, bool) {
	const prefix = "Bearer "

	if !strings.HasPrefix(value, prefix) {
		return "", false
	}

	token := strings.TrimSpace(strings.TrimPrefix(value, prefix))
	if token == "" {
		return "", false
	}

	return token, true
}
