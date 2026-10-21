package kafka_reader

import (
	"errors"
	"fmt"
	"go_project/infrastructure/messaging/kafka"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type ReaderConfig struct {
	CommonConfig kafka.Config

	Topics  []string `envconfig:"TOPICS" required:"true"`
	GroupID string   `envconfig:"GROUP_ID" required:"true"`

	StartFromLatest bool `envconfig:"START_FROM_LATEST"`

	MinBytes     int           `envconfig:"MIN_BYTES"`
	MaxBytes     int           `envconfig:"MAX_BYTES"`
	MaxWait      time.Duration `envconfig:"MAX_WAIT"`
	RequestIDKey string
}

func NewReaderConfig(
	kafkaConfig kafka.Config,
	requestIdKey string,
) (ReaderConfig, error) {
	var readerConfig ReaderConfig
	if err := envconfig.Process("KAFKA_READER", &readerConfig); err != nil {
		return ReaderConfig{}, fmt.Errorf("process envconfig:%w", err)
	}

	if err := readerConfig.Validate(); err != nil {
		return readerConfig, fmt.Errorf("%s: %w", "KAFKA_READER", err)
	}
	return ReaderConfig{
		CommonConfig:    kafkaConfig,
		Topics:          readerConfig.Topics,
		GroupID:         readerConfig.GroupID,
		StartFromLatest: readerConfig.StartFromLatest,
		MinBytes:        readerConfig.MinBytes,
		MaxBytes:        readerConfig.MaxBytes,
		MaxWait:         readerConfig.MaxWait,
		RequestIDKey:    requestIdKey,
	}, nil
}
func NewReaderConfigMust(
	KafkaConfig kafka.Config,
	requestIdKey string,
) ReaderConfig {

	readerConfig, err := NewReaderConfig(KafkaConfig, requestIdKey)

	if err != nil {
		err = fmt.Errorf("initiate Kafka reader config: %w", err)
		panic(err)

	}
	return readerConfig
}

func (c ReaderConfig) Validate() error {
	if c.MinBytes <= 0 {
		c.MinBytes = 1
	}
	if c.MaxBytes <= 0 {
		c.MaxBytes = 10 << 20
	}
	if c.MaxWait <= 0 {
		c.MaxWait = 500 * time.Millisecond
	}
	if len(c.Topics) == 0 {
		return errors.New("kafka: subscriber topics are required")
	}
	if c.GroupID == "" {
		return errors.New("kafka: subscriber group id is required")
	}
	return nil
}
