package recovery

import (
	"context"

	"google.golang.org/grpc"
)

// RecoveryInterceptor перехватывает паники в следующих интерсепторах и обработчике
// в горутине unary RPC. Для записи паники в лог подключается после LoggerInterceptor.
func RecoveryInterceptor() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp interface{}, err error) {
		// Именованные результаты позволяют вернуть gRPC-ошибку после recover.
		defer func() {
			if recovered := recover(); recovered != nil {
				resp = nil
				err = panicError(ctx, recovered)
			}
		}()

		return handler(ctx, req)
	}
}
