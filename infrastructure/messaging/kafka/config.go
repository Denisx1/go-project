// Package kafka - ЕДИНСТВЕННОЕ место в репозитории, которое знает про Kafka.
// Наружу торчат только конструкторы, возвращающие messaging.Publisher/Subscriber.
package kafka

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/kelseyhightower/envconfig"
	kafkago "github.com/segmentio/kafka-go"
)

// Config - параметры подключения, общие для продюсера, консьюмера и админ-хелперов.
// Теги env рассчитаны на caarlos0/env (или любой совместимый загрузчик).
type Config struct {
	Brokers  []string `envconfig:"BROKERS"`
	ClientID string   `envconfig:"CLIENT_ID"`

	DialTimeout time.Duration `envconfig:"KAFKA_DIAL_TIMEOUT" envDefault:"10s"`

	// TLS включает TLS с системными корневыми сертификатами.
	TLS bool `env:"TLS"`

	// // SASL/PLAIN. Пустой username - SASL выключен.
	// SASLUsername string `envconfig:"KAFKA_SASL_USERNAME"`
	// SASLPassword string `envconfig:"KAFKA_SASL_PASSWORD"`
}

func (c Config) Validate() error {
	if len(c.Brokers) == 0 {
		return errors.New("kafka: brokers are required")
	}
	for _, b := range c.Brokers {
		if strings.TrimSpace(b) == "" {
			return errors.New("kafka: empty broker address")
		}
	}
	// if (c.SASLUsername == "") != (c.SASLPassword == "") {
	// 	return errors.New("kafka: sasl username and password must be set together")
	// }
	return nil
}

func (c Config) dialTimeout() time.Duration {
	if c.DialTimeout > 0 {
		return c.DialTimeout
	}
	return 10 * time.Second
}

func (c Config) tlsConfig() *tls.Config {
	if !c.TLS {
		return nil
	}
	return &tls.Config{MinVersion: tls.VersionTLS12}
}

// dialer используется Reader'ом и админ-хелперами.
func (c Config) Dialer() *kafkago.Dialer {
	d := &kafkago.Dialer{
		ClientID:  c.ClientID,
		Timeout:   c.dialTimeout(),
		DualStack: true,
		TLS:       c.tlsConfig(),
	}
	// if c.SASLUsername != "" {
	// 	d.SASLMechanism = plain.Mechanism{Username: c.SASLUsername, Password: c.SASLPassword}
	// }
	return d
}

// transport используется Writer'ом.
func (c Config) Transport() *kafkago.Transport {
	t := &kafkago.Transport{
		ClientID:    c.ClientID,
		DialTimeout: c.dialTimeout(),
		TLS:         c.tlsConfig(),
	}
	// if c.SASLUsername != "" {
	// 	t.SASL = plain.Mechanism{Username: c.SASLUsername, Password: c.SASLPassword}
	// }
	return t
}

func NewKafkaConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("KAFKA", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig:%w", err)
	}
	if err := config.Validate(); err != nil {
		return config, fmt.Errorf("%s: %w", "KAFKA", err)
	}
	return config, nil
}

func NewKafkaConfigMust() Config {
	config, err := NewKafkaConfig()
	if err != nil {
		err = fmt.Errorf("connect to kafka: %w", err)
		panic(err)
	}
	return config
}
