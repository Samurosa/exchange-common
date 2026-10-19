package sharedContext

import "context"

type requestContextKey struct{}

func UpdateRequestContext(
	ctx context.Context,
	update func(*RequestContext),
) context.Context {
	requestCtx, ok := GetRequestContext(ctx)
	if !ok {
		requestCtx = RequestContext{}
	}

	update(&requestCtx)

	return context.WithValue(
		ctx,
		requestContextKey{},
		requestCtx,
	)
}

func GetRequestContext(
	ctx context.Context,
) (RequestContext, bool) {
	requestCtx, ok := ctx.Value(requestContextKey{}).(RequestContext)

	return requestCtx, ok
}

func UserID(ctx context.Context) (string, bool) {
	requestCtx, ok := GetRequestContext(ctx)
	if !ok {
		return "", false
	}

	return requestCtx.Principal.UserID, true
}

func Role(ctx context.Context) (string, bool) {
	requestCtx, ok := GetRequestContext(ctx)
	if !ok {
		return "", false
	}

	return requestCtx.Principal.Role, true
}

func DeviceID(ctx context.Context) (string, bool) {
	requestCtx, ok := GetRequestContext(ctx)
	if !ok {
		return "", false
	}

	return requestCtx.Metadata.DeviceID, true
}

func JTI(ctx context.Context) (string, bool) {
	requestCtx, ok := GetRequestContext(ctx)
	if !ok {
		return "", false
	}

	return requestCtx.JTI, true
}
