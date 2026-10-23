package GrpcResponse

import (
	"errors"
	"fmt"
	error_type "go_project/infrastructure/error"
	logger "go_project/infrastructure/logging"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GrpcResponseHandler struct {
	log logger.Logger
}

func NewGrpcResponseHandler(
	log logger.Logger,
) *GrpcResponseHandler {
	return &GrpcResponseHandler{log: log}
}

func (r *GrpcResponseHandler) ErrorResponse(err error, msg string) error {
	if err == nil {
		return nil
	}

	var (
		grpcCode codes.Code
		logFunc  func(string, ...any)
	)
	switch {
	case errors.Is(err, error_type.ErrConflict):
		logFunc = r.log.Warn
		grpcCode = codes.AlreadyExists
	case errors.Is(err, error_type.ErrNotFound):
		logFunc = r.log.Debug
		grpcCode = codes.NotFound
	case errors.Is(err, error_type.ErrInvalidArgument):
		logFunc = r.log.Warn
		grpcCode = codes.InvalidArgument
	default:
		logFunc = r.log.Warn
		grpcCode = codes.Internal
	}

	logFunc(msg, err)
	return status.Error(grpcCode, fmt.Sprintf("%s: %v", msg, err))
}

func (h *GrpcResponseHandler) PanicResponse(p any, msg string) error {
	statusCode := codes.Internal
	err := fmt.Errorf("unexpected panic: %v", p)

	h.log.Error(msg, zap.Error(err))

	return status.Error(statusCode, msg)
}
