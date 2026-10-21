package logger

import (
	"fmt"
	"strings"

	"github.com/kelseyhightower/envconfig"
)

type Config struct {
	Level       string `envconfig:"LEVEL"`
	Format      string `envconfig:"FORMAT"`
	Folder      string `envconfig:"FOLDER"` // пусто = только stdout
	ServiceName string `ignored:"true"`     // заполняется кодом из APP_NAME, не из env
}

func (c Config) Validate() error {
	switch strings.ToUpper(c.Format) {
	case "JSON", "CONSOLE":
		return nil
	}
	return fmt.Errorf("LOGGER_FORMAT %q: want JSON|CONSOLE", c.Format)
}

func NewLoggerConfig(serviceName string) (Config, error) {
	var config Config
	if err := envconfig.Process("LOGGER", &config); err != nil {
		return Config{}, fmt.Errorf("process logger config: %w", err)
	}
	return Config{
		Level:       config.Level,
		Format:      config.Format,
		Folder:      config.Folder,
		ServiceName: serviceName,
	}, nil
}
func NewLoggerConfigMust(serviceName string) Config {
	config, err := NewLoggerConfig(serviceName)
	if err != nil {
		err = fmt.Errorf("get Logger config: %w", err)
		panic(err)
	}
	return config
}
