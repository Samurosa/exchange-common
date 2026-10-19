package recovery

import "google.golang.org/grpc"

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
