package recovery

import "google.golang.org/grpc"

// RecoveryStreamInterceptor перехватывает паники в следующих интерсепторах
// и обработчике в горутине stream RPC. Для записи паники в лог подключается
// после LoggerStreamInterceptor.
func RecoveryStreamInterceptor() grpc.StreamServerInterceptor {
	return func(
		srv interface{},
		stream grpc.ServerStream,
		info *grpc.StreamServerInfo,
		handler grpc.StreamHandler,
	) (err error) {
		defer func() {
			if recovered := recover(); recovered != nil {
				err = panicError(stream.Context(), recovered)
			}
		}()

		return handler(srv, stream)
	}
}
