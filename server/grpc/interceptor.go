package grpc

import (
	"context"
	"strings"
	"time"

	"github.com/pkg/errors"
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

		logger.Debug("grpc unary call incoming", append(fields, "request", req)...)

		resp, err = handler(ctx, req)

		code := status.Code(err)
		fields = append(fields, "duration", time.Since(startTime), "code", code.String(), "request", req)
		if err == nil {
			fields = append(fields, "response", resp)
		} else if errors.Is(err, context.DeadlineExceeded) || errors.Is(err, context.Canceled) {
			fields = append(fields, "error", err.Error())
		} else if s, ok := status.FromError(err); ok {
			fields = append(fields, "error", s.Message())
		} else {
			fields = append(fields, "error", err.Error())
		}
		switch code {
		case codes.Canceled, codes.Unknown, codes.DeadlineExceeded, codes.ResourceExhausted,
			codes.Unimplemented, codes.Unavailable, codes.Unauthenticated:
			logger.Warn("grpc unary error from framework", fields...)
		case codes.InvalidArgument, codes.NotFound, codes.AlreadyExists, codes.PermissionDenied,
			codes.FailedPrecondition, codes.Aborted, codes.OutOfRange, codes.Internal, codes.DataLoss:
			logger.Warn("grpc unary error from service", fields...)
		default:
			logger.Info("grpc unary call finished", fields...)
		}
		return resp, err
	}
}

func ErrorInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
		resp, err = handler(ctx, req)
		if err != nil {
			return
		}
		if _, ok := status.FromError(err); !ok && !errors.Is(err, context.DeadlineExceeded) && !errors.Is(err, context.Canceled) {
			return resp, status.Error(codes.Internal, err.Error())
		}
		return
	}
}
