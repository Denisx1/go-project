package grpcconfig

import (
	"fmt"
	"strconv"
	"time"

	"github.com/kelseyhightower/envconfig"
	"google.golang.org/grpc/keepalive"
)

type KeepaliveConfig struct {
	// ServerParameters
	MaxConnectionIdle     time.Duration `envconfig:"MAX_CONNECTION_IDLE"`
	MaxConnectionAge      time.Duration `envconfig:"MAX_CONNECTION_AGE"`
	MaxConnectionAgeGrace time.Duration `envconfig:"MAX_CONNECTION_AGE_GRACE"`

	Time    time.Duration `envconfig:"TIME"`
	Timeout time.Duration `envconfig:"TIMEOUT"`

	// EnforcementPolicy
	MinTime             time.Duration `envconfig:"MIN_TIME" default:"10s"`
	PermitWithoutStream bool          `envconfig:"PERMIT_WITHOUT_STREAM" default:"true"`
}

type Config struct {
	Port             string        `envconfig:"PORT" required:"true"`
	Timeout          time.Duration `envconfig:"TIMEOUT" required:"true"`
	EnableReflection bool          `envconfig:"REFLECTION" required:"true"`

	MaxRecvMsgSize int           `envconfig:"MAX_RECV_MSG_SIZE" default:"4194304"`
	MaxSendMsgSize int           `envconfig:"MAX_SEND_MSG_SIZE" default:"4194304"`
	HeathInterval  time.Duration `envconfig:"HEALTH_INTERVAL" default:"10s"`

	Keepalive KeepaliveConfig `ignored:"true"`

	RequestIdKey string `ignored:"true"`
}

func NewServerConfig(
	env string,
	requestIdKey string,
) (Config, error) {
	var keepaliveConfig KeepaliveConfig
	var config Config
	if err := envconfig.Process("GRPC_KEEPALIVE", &keepaliveConfig); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}
	if err := envconfig.Process("GRPC", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}
	switch env {
	case "local":
		config.EnableReflection = true
	case "prod":
		if config.EnableReflection {
			return Config{}, fmt.Errorf("grpc reflection must be disabled in prod")
		}
	}
	config.Keepalive = keepaliveConfig
	config.RequestIdKey = requestIdKey
	if err := config.validate(); err != nil {
		return Config{}, fmt.Errorf("validate grpc config: %w", err)
	}
	return Config{
		Port:             config.Port,
		Timeout:          config.Timeout,
		EnableReflection: config.EnableReflection,
		MaxRecvMsgSize:   config.MaxRecvMsgSize,
		MaxSendMsgSize:   config.MaxSendMsgSize,
		HeathInterval:    config.HeathInterval,
		Keepalive:        keepaliveConfig,
		RequestIdKey:     requestIdKey,
	}, nil
}

func NewServerConfigMust(
	env string,
	requestIdKey string,
) Config {
	config, err := NewServerConfig(env, requestIdKey)
	if err != nil {
		err = fmt.Errorf("get Grpc config: %w", err)
		panic(err)
	}
	return config

}

func (c *Config) validate() error {
	port, err := strconv.Atoi(c.Port)
	if err != nil || port < 0 || port > 65535 {
		return fmt.Errorf("GRPC_PORT: invalid port %q", c.Port)
	}
	if c.Timeout <= 0 {
		return fmt.Errorf("GRPC_SHUTDOWN_TIMEOUT must be > 0, got %s", c.Timeout)
	}
	if c.MaxRecvMsgSize < 0 {
		return fmt.Errorf("GRPC_MAX_RECV_MSG_SIZE must be >= 0, got %d", c.MaxRecvMsgSize)
	}
	if c.MaxSendMsgSize < 0 {
		return fmt.Errorf("GRPC_MAX_SEND_MSG_SIZE must be >= 0, got %d", c.MaxSendMsgSize)
	}
	durations := []struct {
		name  string
		value time.Duration
	}{
		{"GRPC_KEEPALIVE_MAX_CONNECTION_IDLE", c.Keepalive.MaxConnectionIdle},
		{"GRPC_KEEPALIVE_MAX_CONNECTION_AGE", c.Keepalive.MaxConnectionAge},
		{"GRPC_KEEPALIVE_MAX_CONNECTION_AGE_GRACE", c.Keepalive.MaxConnectionAgeGrace},
		{"GRPC_KEEPALIVE_TIME", c.Keepalive.Time},
		{"GRPC_KEEPALIVE_TIMEOUT", c.Keepalive.Timeout},
		{"GRPC_KEEPALIVE_MIN_TIME", c.Keepalive.MinTime},
	}
	for _, d := range durations {
		if d.value < 0 {
			return fmt.Errorf("%s must be >= 0, got %s", d.name, d.value)
		}
	}
	return nil

}

// ServerParameters преобразует конфиг в структуру grpc-go.
func (k KeepaliveConfig) ServerParameters() keepalive.ServerParameters {
	return keepalive.ServerParameters{
		MaxConnectionIdle:     k.MaxConnectionIdle,
		MaxConnectionAge:      k.MaxConnectionAge,
		MaxConnectionAgeGrace: k.MaxConnectionAgeGrace,
		Time:                  k.Time,
		Timeout:               k.Timeout,
	}
}

// EnforcementPolicy преобразует конфиг в политику grpc-go.
func (k KeepaliveConfig) EnforcementPolicy() keepalive.EnforcementPolicy {
	return keepalive.EnforcementPolicy{
		MinTime:             k.MinTime,
		PermitWithoutStream: k.PermitWithoutStream,
	}
}
