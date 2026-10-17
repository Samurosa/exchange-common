package recovery

import (
	"context"
	"runtime/debug"

	"github.com/Samurosa/exchange-common/shared/auth/interceptors/logging"
	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func panicError(ctx context.Context, recovered interface{}) error {
	logging.FromContext(ctx).Error("grpc panic recovered",
		zap.Any("panic", recovered),
		zap.ByteString("stack", debug.Stack()),
	)

	return status.Error(codes.Internal, "internal server error")
}
