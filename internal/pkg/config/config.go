package config

import (
	"fmt"
	"log/slog"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Env               string     `env:"ENV" env-default:"production"`
	Port              int        `env:"PORT" env-default:"3000"`
	LogLevel          slog.Level `env:"LOG_LEVEL" env-default:"info"`
	S3Region          string     `env:"S3_REGION" env-default:"us-east-1"`
	S3Endpoint        string     `env:"S3_ENDPOINT"`
	S3Bucket          string     `env:"S3_BUCKET" env-required:"true"`
	S3AccessKeyID     string     `env:"S3_ACCESS_KEY_ID" env-required:"true"`
	S3SecretAccessKey string     `env:"S3_SECRET_ACCESS_KEY" env-required:"true"`
}

func New() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &cfg, nil
}
