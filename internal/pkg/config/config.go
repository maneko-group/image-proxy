package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env      string     `env:"ENV" env-default:"production"`
	Port     int        `env:"PORT" env-default:"3000"`
	LogLevel slog.Level `env:"LOG_LEVEL" env-default:"info"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
