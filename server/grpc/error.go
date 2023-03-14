package grpc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type Error struct {
	Code    codes.Code
	Message string
}

func (e Error) Error() string {
	return e.Message
}

func (e Error) GRPCStatus() *status.Status {
	return status.New(e.Code, e.Message)
}
