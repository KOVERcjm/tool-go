package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func InvalidArgument(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.InvalidArgument, "")
	}
	return status.New(codes.InvalidArgument, message[0])
}

func NotFound(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.NotFound, "")
	}
	return status.New(codes.NotFound, message[0])
}

func AlreadyExists(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.AlreadyExists, "")
	}
	return status.New(codes.AlreadyExists, message[0])
}

func PermissionDenied(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.PermissionDenied, "")
	}
	return status.New(codes.PermissionDenied, message[0])
}

func FailedPrecondition(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.FailedPrecondition, "")
	}
	return status.New(codes.FailedPrecondition, message[0])
}

func Aborted(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.Aborted, "")
	}
	return status.New(codes.Aborted, message[0])
}

func OutOfRange(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.OutOfRange, "")
	}
	return status.New(codes.OutOfRange, message[0])
}

func Internal(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.Internal, "")
	}
	return status.New(codes.Internal, message[0])
}

func DataLoss(message ...string) *status.Status {
	if len(message) == 0 {
		return status.New(codes.DataLoss, "")
	}
	return status.New(codes.DataLoss, message[0])
}
