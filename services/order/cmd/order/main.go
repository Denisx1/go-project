package main

import (
	"context"
	"errors"
	"fmt"
	orderv1 "go_project/contracts/gen/go/order"
	"go_project/infrastructure/grpc/server"
	zap_logger "go_project/infrastructure/logging/zap"
	"go_project/infrastructure/messaging/component/consumer"
	kafka_reader "go_project/infrastructure/messaging/kafka/reader"
	kafka_writter "go_project/infrastructure/messaging/kafka/writter"
	pgx_pool "go_project/infrastructure/repository/postgres/pool/pgx"
	"go_project/order/internal/core/config"
	"go_project/order/internal/features/order/adapters/in/events"
	OrderGrpcHandler "go_project/order/internal/features/order/adapters/in/grpc"
	"go_project/order/internal/features/order/adapters/out/repository"
	"go_project/order/internal/features/order/application/service"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "fatal:", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()
	cfg := config.Load()
	logger, err := zap_logger.NewLogger(cfg.Log)
	if err != nil {
		return fmt.Errorf("init logger: %w", err)
	}
	defer logger.Close()
	logger.Debug("initializing postgres connection pool")
	pool, err := pgx_pool.NewPool(ctx, cfg.Postgres)
	if err != nil {
		return fmt.Errorf("init postgres pool: %w", err)
	}
	logger.Debug("Database successfully connected")
	defer pool.Close()

	logger.Debug("initializing kafka writer")

	writer, err := kafka_writter.NewWritter(cfg.KafkaWriter)
	if err != nil {
		return fmt.Errorf("init kafka writer: %w", err)
	}
	logger.Debug("Kafka writer successfully connected")
	defer writer.Close()

	logger.Debug("initializing kafka reader")
	reader, err := kafka_reader.NewReader(cfg.KafkaReader)
	if err != nil {
		return fmt.Errorf("init kafka reader: %w", err)
	}
	logger.Debug("Kafka reader successfully connected")
	defer reader.Close()
	repository := repository.NewRepository(pool)
	service := service.NewService(repository)
	handler := events.NewHandler(service)
	consumer := consumer.NewConsumer(reader, writer, handler.Handle, logger)

	grpcHandler := OrderGrpcHandler.NewOrderGRPCHandler()
	grpcServer, err := server.NewGrpcServer(cfg.GrpcServer, logger)
	if err != nil {
		return fmt.Errorf("create grpc server: %w", err)
	}
	orderv1.RegisterOrderServer(grpcServer, grpcHandler)

	errChan := make(chan error, 2)
	var wg sync.WaitGroup

	wg.Add(2)
	go func() {
		defer wg.Done()
		if err := consumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			errChan <- errors.New("kafka worker error: " + err.Error())
		}
		logger.Info("Kafka listener was stopped")
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
		logger.Info("shutdown signal received")
	case runErr = <-errChan:
		logger.Error("component failed", "error", runErr)
		stop() // отменяем ctx, чтобы остановить consumer
	}

	if err := grpcServer.Close(); err != nil {
		logger.Warn("grpc shutdown", "error", err)
	}
	logger.Info("Waiting for worker to gracefully shutdown...")
	wg.Wait()
	logger.Info("Application stopped successfully")
	return runErr
}
