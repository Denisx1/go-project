package interceptors

import (
	"context"
	logger "go_project/infrastructure/logging"
	"runtime/debug"
	"time"

	"github.com/google/uuid"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
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
			"server", info.Server,
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
			"method", info.FullMethod,
			"code", code.String(),
			"duration_ms", time.Since(start).Milliseconds(),
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

func Panic() grpc.UnaryServerInterceptor {
	return func(
		ctx context.Context,
		req any,
		info *grpc.UnaryServerInfo,
		handler grpc.UnaryHandler,
	) (resp any, err error) {
		defer func() {
			if r := recover(); r != nil {
				logger.FromContext(ctx).Error("handler panic",
					"method", info.FullMethod,
					"panic", r,
					"stack", string(debug.Stack()),
				)
				resp = nil
				err = status.Error(codes.Internal, "internal error")
			}
		}()
		return handler(ctx, req)
	}
}
