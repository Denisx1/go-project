package kafka_writter

import (
	"fmt"
	"go_project/infrastructure/messaging/kafka"

	"time"

	"github.com/kelseyhightower/envconfig"
)

type WriterConfig struct {
	CommonConfig kafka.Config

	DLTSuffix string   `envconfig:"DLT_SUFFIX"`
	Topics    []string `envconfig:"TOPICS"`

	// WriteTimeout - таймаут одной записи. По умолчанию 10s.
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT"`

	// BatchTimeout - сколько Writer ждёт добора батча перед отправкой.
	// У kafka-go по умолчанию 1s, что даёт до секунды задержки на каждый синхронный Publish.
	// Здесь по умолчанию 10ms.
	BatchTimeout time.Duration `envconfig:"BATCH_TIMEOUT"`
	BatchSize    int           `envconfig:"BATCH_SIZE"`
	// AllowAutoTopicCreation - только для локальной разработки.
	// В проде топики создаются заранее (см. EnsureTopic или миграции инфраструктуры).
	AllowAutoTopicCreation bool `envconfig:"ALLOW_TOPIC_AUTO_CREATION"`
}

func NewWiriterConfig(kafkaConfig kafka.Config) (WriterConfig, error) {
	var writerConfig WriterConfig

	if err := envconfig.Process("KAFKA_WRITER", &writerConfig); err != nil {
		return WriterConfig{}, fmt.Errorf("process envconfig:%w", err)
	}
	if err := writerConfig.Validate(); err != nil {
		return writerConfig, fmt.Errorf("%s: %w", "KAFKA_WRITER", err)
	}

	return WriterConfig{
		CommonConfig:           kafkaConfig,
		WriteTimeout:           writerConfig.WriteTimeout,
		BatchTimeout:           writerConfig.BatchTimeout,
		AllowAutoTopicCreation: writerConfig.AllowAutoTopicCreation,
	}, nil
}
func NewWriterConfigMust(kafkaConfig kafka.Config) WriterConfig {
	writerConfig, err := NewWiriterConfig(kafkaConfig)
	if err != nil {
		err = fmt.Errorf("initiate Kafka writer config: %w", err)
		panic(err)

	}
	return writerConfig
}

func (c WriterConfig) Validate() error {

	if c.DLTSuffix == "" {
		return fmt.Errorf("kafka write DLT must be")
	}
	if len(c.Topics) == 0 {
		return fmt.Errorf("kafka write Topics must be")
	}
	if c.WriteTimeout <= 0 {
		return fmt.Errorf("kafka write timeout must be positive")
	}

	if c.BatchTimeout < 0 {
		return fmt.Errorf("kafka batch timeout must be positive")
	}
	return nil
}
