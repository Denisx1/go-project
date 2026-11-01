package redisclient

import (
	"fmt"
	"time"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Host         string        `envconfig:"HOST" required:"true"` // localhost:6379 с хоста, user-redis:6379 в docker-сети
	Port         int           `envconfig:"PORT"`
	Password     string        `envconfig:"PASSWORD"`
	DB           int           `envconfig:"DB" default:"0"`
	DialTimeout  time.Duration `envconfig:"DIAL_TIMEOUT" default:"3s"`
	ReadTimeout  time.Duration `envconfig:"READ_TIMEOUT" default:"2s"`
	WriteTimeout time.Duration `envconfig:"WRITE_TIMEOUT" default:"2s"`
	PoolSize     int           `envconfig:"POOL_SIZE" default:"10"`
}

func NewRedisConfig() (Config, error) {
	var config Config
	if err := envconfig.Process("REDIS", &config); err != nil {
		return Config{}, fmt.Errorf("process envconfig:%w", err)
	}
	return config, nil
}
func NewRedisConfigMust() Config {
	config, err := NewRedisConfig()
	if err != nil {
		err = fmt.Errorf("get Postgres connection pool: %w", err)
		panic(err)
	}
	return config
}
func (c Config) SetAddr() string {
	return fmt.Sprintf("%s:%d", c.Host, c.Port)
}
