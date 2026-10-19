package grpcerrors

import (
	"context"
	"errors"
	"io"
	"net"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func Map(err error) (error, bool) {
	switch {
	case err == nil:
		return nil, true
	case errors.Is(err, context.Canceled):
		return status.Error(codes.Canceled, "request canceled"), true
	case errors.Is(err, context.DeadlineExceeded):
		return status.Error(codes.DeadlineExceeded, "request deadline exceeded"), true
	}

	if st, ok := status.FromError(err); ok {
		code := st.Code()
		if code == codes.Unknown || code == codes.OK {
			code = codes.Internal
		}
		return status.Error(code, "dependency request failed"), true
	}

	var sqlErr interface{ SQLState() string }
	if errors.As(err, &sqlErr) {
		state := sqlErr.SQLState()
		switch {
		case state == "40001" || state == "40P01":
			return status.Error(codes.Aborted, "database transaction must be retried"), true
		case strings.HasPrefix(state, "08"), strings.HasPrefix(state, "53"),
			state == "57P01", state == "57P02", state == "57P03":
			return status.Error(codes.Unavailable, "database unavailable"), true
		}
	}

	var networkErr net.Error
	if errors.As(err, &networkErr) || errors.Is(err, io.EOF) ||
		errors.Is(err, io.ErrUnexpectedEOF) || errors.Is(err, io.ErrClosedPipe) ||
		errors.Is(err, net.ErrClosed) {
		return status.Error(codes.Unavailable, "dependency connection unavailable"), true
	}
	return nil, false
}
