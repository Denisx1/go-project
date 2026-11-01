package server

import (
	"context"
	"errors"
	"fmt"
	grpcconfig "go_project/infrastructure/grpc/config.go"
	"go_project/infrastructure/grpc/server/interceptors"
	logger "go_project/infrastructure/logging"
	"net"
	"sync"
	"time"

	"buf.build/go/protovalidate"

	pvmw "github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/protovalidate"
	"google.golang.org/grpc"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	"google.golang.org/grpc/reflection"
)

type Server struct {
	cfg      grpcconfig.Config
	log      logger.Logger
	raw      *grpc.Server
	health   *health.Server
	listener net.Listener

	closeOnce sync.Once
	closeErr  error
}

type Check struct {
	Name string
	Fn   func(ctx context.Context) error
}

func NewGrpcServer(
	cfg grpcconfig.Config,
	log logger.Logger,
) (*Server, error) {
	lis, err := net.Listen("tcp", ":"+cfg.Port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen port %s: %w", cfg.Port, err)
	}
	validator, err := protovalidate.New() // один раз при старте
	if err != nil {
		return nil, fmt.Errorf("create validator: %w", err)
	}
	chain := []grpc.UnaryServerInterceptor{
		interceptors.RequestId(cfg.RequestIdKey),
		interceptors.Logger(log),
		interceptors.Trace(),
		interceptors.Errors(),
		pvmw.UnaryServerInterceptor(validator),
		interceptors.Panic(),
	}

	base := []grpc.ServerOption{
		grpc.ChainUnaryInterceptor(chain...),
		grpc.KeepaliveParams(cfg.Keepalive.ServerParameters()),
		grpc.KeepaliveEnforcementPolicy(cfg.Keepalive.EnforcementPolicy()),
	}
	if cfg.MaxRecvMsgSize > 0 {
		base = append(base, grpc.MaxRecvMsgSize(cfg.MaxRecvMsgSize))
	}
	if cfg.MaxSendMsgSize > 0 {
		base = append(base, grpc.MaxSendMsgSize(cfg.MaxSendMsgSize))
	}

	raw := grpc.NewServer(base...)
	hs := health.NewServer()
	hs.SetServingStatus("", grpc_health_v1.HealthCheckResponse_NOT_SERVING)
	grpc_health_v1.RegisterHealthServer(raw, hs)

	if cfg.EnableReflection {
		reflection.Register(raw)
	}
	return &Server{cfg: cfg, log: log, raw: raw, health: hs, listener: lis}, nil
}

func (s *Server) RegisterService(
	desc *grpc.ServiceDesc,
	impl interface{},
) {
	s.raw.RegisterService(desc, impl)
}

func (s *Server) Addr() net.Addr { return s.listener.Addr() }

func (s *Server) Start() error {
	s.health.SetServingStatus("", grpc_health_v1.HealthCheckResponse_SERVING)
	s.log.Info("grpc server started", "addr", s.listener.Addr().String())
	err := s.raw.Serve(s.listener)
	if err == nil || errors.Is(err, grpc.ErrServerStopped) {
		return nil
	}
	return fmt.Errorf("grpc serve: %w", err)
}

// Close: health -> NOT_SERVING, GracefulStop, по таймауту принудительный Stop.
// Повторные вызовы безопасны.
func (s *Server) Close() error {
	s.closeOnce.Do(func() {
		s.health.Shutdown()

		done := make(chan struct{})
		go func() {
			s.raw.GracefulStop()
			close(done)
		}()

		timer := time.NewTimer(s.cfg.Timeout)
		defer timer.Stop()

		select {
		case <-done:
		case <-timer.C:
			s.log.Warn("graceful stop timed out, forcing", "timeout", s.cfg.Timeout.String())
			s.raw.Stop()
			<-done
			s.closeErr = fmt.Errorf("graceful stop exceeded %s", s.cfg.Timeout)
		}

		_ = s.listener.Close()
	})
	return s.closeErr
}

func (s *Server) WatchHealth(ctx context.Context, checks ...Check) {
	go func() {
		t := time.NewTicker(s.cfg.HeathInterval)
		defer t.Stop()
		for {
			status := grpc_health_v1.HealthCheckResponse_SERVING
			for _, c := range checks {
				cctx, cancel := context.WithTimeout(ctx, 2*time.Second)
				err := c.Fn(cctx)
				cancel()
				if err != nil {
					s.log.Warn("health check failed", "check", c.Name, "error", err)
					status = grpc_health_v1.HealthCheckResponse_NOT_SERVING
				}
			}
			s.health.SetServingStatus("", status)

			select {
			case <-ctx.Done():
				return
			case <-t.C:
			}
		}
	}()
}
