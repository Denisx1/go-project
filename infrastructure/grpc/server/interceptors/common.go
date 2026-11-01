package interceptors

import (
	"context"
	"errors"
	"fmt"
	logger "go_project/infrastructure/logging"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

type RequestIdKey struct{}

func RequestId(requestIdKey string) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		var requestID string

		if md, ok := metadata.FromIncomingContext(ctx); ok {
			if v := md.Get(requestIdKey); len(v) > 0 {
				requestID = v[0]
			}
		}
		if requestID == "" {
			requestID = uuid.New().String()
		}
		ctx = context.WithValue(ctx, RequestIdKey{}, requestID)
		return handler(ctx, req)
	}
}

func Logger(
	log logger.Logger,
) grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req interface{},
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		requestId := ctx.Value(RequestIdKey{}).(string)
		l := log.With(
			"request_id", requestId,
			"method", info.FullMethod,
		)
		ctx = logger.ToContext(ctx, l)
		return handler(ctx, req)

	}
}

func Trace() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (any, error) {
		log := logger.FromContext(ctx)

		start := time.Now()
		log.Debug(
			">>> Incoming Grpc request",
			"StartTime", start.UTC(),
		)
		res, err := handler(ctx, req)
		code := status.Code(err)
		attrs := []any{
			"code", code.String(),
			"Latency", time.Since(start).String(),
		}
		switch {
		case err == nil:
			log.Info("<<< grpc request completed", attrs...)
		case isServerFault(code):
			log.Error("grpc request failed", append(attrs, "err", err)...)
		default:
			log.Warn("grpc request rejected", append(attrs, "err", err)...)
		}
		return res, err
	}
}
func Errors() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		res, err := handler(ctx, req)
		log := logger.FromContext(ctx)
		if err == nil {
			return res, nil
		}

		if _, ok := status.FromError(err); ok {
			return nil, err
		}
		code, clientMsg := mapDomainError(err)

		if isServerFault(code) {
			log.Error("request failed", "code", code, "error", err)
		} else {
			log.Warn("request rejected", "code", code, "error", err)
		}
		return nil, status.Error(code, clientMsg)

	}
}

func Panic() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		var ErrPanic = errors.New("panic recovered")
		defer func() {
			if r := recover(); r != nil {
				logger.FromContext(ctx).Error("handler panic",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				resp = nil
				err = fmt.Errorf("%w: %v", ErrPanic, r)
			}
		}()
		return handler(ctx, req)
	}
}
