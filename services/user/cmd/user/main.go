package main

import (
	"context"
	"errors"
	"fmt"
	userv1 "go_project/contracts/gen/go/user"
	bcrypt_hash "go_project/infrastructure/bcrypt"
	redisclient "go_project/infrastructure/cache/redis"
	"go_project/infrastructure/grpc/server"
	zap_logger "go_project/infrastructure/logging/zap"
	"go_project/infrastructure/messaging/component/consumer"
	kafka_reader "go_project/infrastructure/messaging/kafka/reader"
	kafka_writter "go_project/infrastructure/messaging/kafka/writter"
	"go_project/infrastructure/messaging/router"
	pgx_pool "go_project/infrastructure/repository/postgres/pool/pgx"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"user/internal/core/config"
	useroutbox "user/internal/core/outbox"
	event "user/internal/features/candidate/adapters/in/events"
	candidategrpchandler "user/internal/features/candidate/adapters/in/grpc"
	candidatecache "user/internal/features/candidate/adapters/out/cache"
	candidaterepository "user/internal/features/candidate/adapters/out/repository"
	candidateservice "user/internal/features/candidate/application/service"
	profilerepository "user/internal/features/profile/adapters/out/repository"
	rolegrpchandler "user/internal/features/role/adapters/in/grpc"
	rolecache "user/internal/features/role/adapters/out/cache"
	rolerepository "user/internal/features/role/adapters/out/repository"
	roleservice "user/internal/features/role/application/service"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	cfg := config.Load()

	log, err := zap_logger.NewLogger(cfg.Log)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer log.Close()

	pool, err := pgx_pool.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres pool: %w", err)
	}
	log.Debug("database successfully connected")
	defer pool.Close()
	cache, err := redisclient.NewClient(ctx, cfg.Cache)
	if err != nil {
		return fmt.Errorf("init redis: %w", err)
	}
	defer cache.Close()
	writer, err := kafka_writter.NewWritter(cfg.KafkaWriter)
	if err != nil {
		return fmt.Errorf("init kafka writer: %w", err)
	}
	log.Debug("kafka writer successfully connected")
	defer writer.Close()

	reader, err := kafka_reader.NewReader(cfg.KafkaReader)
	if err != nil {
		return fmt.Errorf("init kafka reader: %w", err)
	}
	log.Debug("kafka reader successfully connected")
	defer reader.Close()
	bcrypt := bcrypt_hash.NewBcryptHasher(10)
	outbox := useroutbox.NewUserOutbox(pool)
	log.Debug("outbox connected successfuly")

	r := router.NewRouter()

	roleRepository := rolerepository.NewRoleRepository(pool)
	roleCache := rolecache.NewRoleCache(cache)
	roleService := roleservice.NewRoleService(roleRepository, roleCache, pool)
	roleGrpcHandler := rolegrpchandler.NewCandidateGrpcHandler(roleService)

	profileRepository := profilerepository.NewProfileRepository(pool)

	candidateRepository := candidaterepository.NewCandidateRepository(pool)
	candidateCacheRepository := candidatecache.NewCandidateCache(cache, cfg.App.SessionKey)
	candidateService := candidateservice.NewCandidateService(
		candidateRepository,
		profileRepository,
		outbox,
		roleRepository,
		candidateCacheRepository,
		pool,
		bcrypt)
	candidateGrpcHandler := candidategrpchandler.NewCandidateGrpcHandler(candidateService)
	event.NewCandidateEventHandler(candidateService).Register(r)
	log.Debug("feature candidate initialized")

	consumer := consumer.NewConsumer(reader, writer, r.Dispatch, log)

	grpcServer, err := server.NewGrpcServer(cfg.GrpcServer, log)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	userv1.RegisterCandidateServiceServer(grpcServer, candidateGrpcHandler)
	userv1.RegisterRoleServiceServer(grpcServer, roleGrpcHandler)
	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	grpcServer.WatchHealth(ctx,
		server.Check{Name: "postgres", Fn: pool.Ping},
		server.Check{Name: "redis", Fn: cache.Ping},
	)
	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := consumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errChan <- errors.New("kafka consumer error: %w" + err.Error())
		}
		log.Info("Kafka listener was stopped")
	}()
	go func() {
		defer wg.Done()
		if err := grpcServer.Start(); err != nil {
			errChan <- fmt.Errorf("grpc server: %w", err)
		}
	}()
	var runErr error
	select {
	case <-ctx.Done():
		log.Info("shutdown signal received")
	case runErr = <-errChan:
		log.Error("component failed", "error", runErr)
		stop() // отменяем ctx, чтобы остановить consumer
	}
	if err := grpcServer.Close(); err != nil {
		log.Warn("grpc shutdown", "error", err)
	}
	log.Info("Waiting for worker to gracefully shutdown...")
	wg.Wait()
	log.Info("Application stopped successfully")
	return runErr
}
