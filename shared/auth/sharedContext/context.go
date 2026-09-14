package sharedContext

import (
	"context"
)

type requestContextKey struct {
}

func WithRequestContext(
	ctx context.Context,
	requestCtx RequestContext,
) context.Context {
	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	)
}

func GetRequestContext(
	ctx context.Context,
) (RequestContext, error) {

	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(RequestContext)

	if !ok {
		return RequestContext{}, ErrInvalidContext
	}

	return requestCtx, nil
}

func UserID(ctx context.Context) (string, error) {
	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(RequestContext)

	if !ok {
		return "", ErrInvalidContext
	}

	return requestCtx.Principal.UserID, nil
}

func DeviceID(ctx context.Context) (string, error) {
	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(RequestContext)

	if !ok {
		return "", ErrInvalidContext
	}

	return requestCtx.Metadata.DeviceID, nil
}

func JTI(ctx context.Context) (string, error) {
	requestCtx, ok := ctx.Value(
		requestContextKey{},
	).(RequestContext)

	if !ok {
		return "", ErrInvalidContext
	}

	return requestCtx.JTI, nil
}

func UpdateRequestContext(
	ctx context.Context,
	update func(*RequestContext),
) (
	context.Context,
	error,
) {

	requestCtx, ok := ctx.Value(requestContextKey{}).(RequestContext)

	if !ok {
		return ctx, ErrCtxForUpdateNotFound
	}

	update(&requestCtx)

	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	), nil
}
