package grpcx

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

var (
	ErrNotFound            = errors.New("not found")
	ErrInternalServerError = errors.New("internal server error")
)

type GRPCError struct {
	msg  string
	kind error
}

func (e *GRPCError) Error() string {
	return e.msg
}

func (e *GRPCError) Is(target error) bool {
	return e.kind == target
}

func TranslateErr(err error) *GRPCError {
	sts, ok := status.FromError(err)

	if !ok {
		return nil
	}

	var kind error
	switch sts.Code() {
	case codes.NotFound:
		kind = ErrNotFound
	case codes.Internal:
		kind = ErrInternalServerError
	default:
		kind = ErrInternalServerError
	}

	return &GRPCError{
		msg:  sts.Message(),
		kind: kind,
	}
}
