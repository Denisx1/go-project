package config

import (
	"fmt"
	"go_project/infrastructure/grpc/config.go"
	logger "go_project/infrastructure/logging"
	"go_project/infrastructure/messaging/kafka"
	kafka_reader "go_project/infrastructure/messaging/kafka/reader"
	kafka_writter "go_project/infrastructure/messaging/kafka/writter"
	pgx_pool "go_project/infrastructure/repository/postgres/pool/pgx"

	"github.com/kelseyhightower/envconfig"
)

type AppConfig struct {
	TimeZone     string `envconfig:"TIMEZONE" required:"true"`
	ServiceName  string `envconfig:"NAME" required:"true"`
	ServiceENV   string `envconfig:"ENV" requires:"true"`
	RequestIdKey string `envconfig:"ID_KEY" requires:"true"`
}

type Config struct {
	App         AppConfig
	Log         logger.Config
	GrpcServer  config.GrpcConfig
	kafka       kafka.Config
	KafkaReader kafka_reader.ReaderConfig
	KafkaWriter kafka_writter.WriterConfig
	Postgres    pgx_pool.Config
}

func NewAppConfig() (AppConfig, error) {
	var appConfig AppConfig
	if err := envconfig.Process("APP", &appConfig); err != nil {
		return AppConfig{}, fmt.Errorf("process app config: %w", err)
	}
	return appConfig, nil
}
func NewAppConfigMust() AppConfig {
	appConfig, err := NewAppConfig()
	if err != nil {
		err = fmt.Errorf("get App config: %w", err)
		panic(err)
	}
	return appConfig
}

func Load() Config {
	var cfg Config
	app := NewAppConfigMust()
	cfg.App = app
	postgres := pgx_pool.NewConfigMust()
	cfg.Postgres = postgres
	log := logger.NewLoggerConfigMust(app.ServiceName)
	cfg.Log = log
	server := config.NewServerConfigMust(cfg.App.ServiceENV, cfg.App.RequestIdKey)
	cfg.GrpcServer = server
	kf := kafka.NewKafkaConfigMust()
	cfg.kafka = kf
	kr := kafka_reader.NewReaderConfigMust(kf, cfg.App.RequestIdKey)
	cfg.KafkaReader = kr
	kw := kafka_writter.NewWriterConfigMust(kf)
	cfg.KafkaWriter = kw
	return cfg
}
