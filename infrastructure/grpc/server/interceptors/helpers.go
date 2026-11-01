package interceptors

import (
	"errors"
	error_type "go_project/infrastructure/error"

	"google.golang.org/grpc/codes"
)

func isServerFault(c codes.Code) bool {
	switch c {
	case codes.Internal, codes.Unknown, codes.DataLoss, codes.Unavailable:
		return true
	}
	return false
}

func mapDomainError(err error) (codes.Code, string) {
	switch {
	case errors.Is(err, error_type.ErrAlreadyExist):
		return codes.AlreadyExists, err.Error()
	case errors.Is(err, error_type.ErrConflict):
		return codes.AlreadyExists, err.Error()
	case errors.Is(err, error_type.ErrNotFound):
		return codes.NotFound, err.Error()
	case errors.Is(err, error_type.ErrInvalidArgument):
		return codes.InvalidArgument, err.Error()
	default:
		return codes.Internal, "internal error"
	}
}
