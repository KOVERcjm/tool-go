package grpc

import (
	"context"
	"strings"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/kovercjm/tool-go/logger"
)

func LoggerInterceptor(logger logger.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		startTime := time.Now()
		fullMethodParts := strings.SplitN(info.FullMethod, "/", 3)
		service, method := info.FullMethod, info.FullMethod
		if len(fullMethodParts) == 3 {
			service, method = fullMethodParts[1], fullMethodParts[2]
		}

		// TODO trace id from context?
		fields := []interface{}{"service", service, "method", method}
		if deadline, ok := ctx.Deadline(); ok {
			fields = append(fields, "request.deadline", deadline.Format(time.RFC3339))
		}

		logger.Info("grpc call received", append(fields, "request", req)...)

		resp, err = handler(ctx, req)

		code := status.Code(err)
		fields = append(fields, "duration", time.Since(startTime), "code", code.String(), "request", req)
		if err != nil {
			fields = append(fields, "error", err)
		} else {
			fields = append(fields, "response", resp)
		}
		switch code {
		case codes.Unknown, codes.DeadlineExceeded, codes.AlreadyExists, codes.Internal, codes.Unavailable, codes.DataLoss:
			logger.Error("grpc call finished", fields...)
		case codes.Canceled, codes.InvalidArgument, codes.NotFound, codes.Aborted, codes.PermissionDenied,
			codes.Unauthenticated, codes.ResourceExhausted, codes.FailedPrecondition, codes.OutOfRange, codes.Unimplemented:
			logger.Warn("grpc call finished", fields...)
		default:
			logger.Info("grpc call finished", fields...)
		}
		return resp, err
	}
}
