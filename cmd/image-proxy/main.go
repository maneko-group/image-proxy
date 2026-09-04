package main

import (
	"context"
	"fmt"
	"log"
	"log/slog"

	"github.com/aws/aws-sdk-go-v2/aws"
	awsconfig "github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/gofiber/fiber/v3"
	"github.com/maneko-group/image-proxy/internal/pkg/config"
	"github.com/maneko-group/image-proxy/internal/pkg/logger"
	"github.com/maneko-group/image-proxy/internal/pkg/proxy"
	"github.com/maneko-group/image-proxy/internal/pkg/storage"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		log.Fatal(err)
	}

	logger.SetDefault(logger.Options{
		AppName:     "image-proxy",
		Environment: logger.Environment(cfg.Env),
		Level:       slog.Level(cfg.LogLevel),
	})

	awsCfg, err := awsconfig.LoadDefaultConfig(context.Background(),
		awsconfig.WithRegion(cfg.S3Region),
		awsconfig.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(
			cfg.S3AccessKeyID,
			cfg.S3SecretAccessKey,
			"",
		)),
	)
	if err != nil {
		log.Fatal(fmt.Errorf("load aws config: %w", err))
	}

	s3Client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		if cfg.S3Endpoint != "" {
			o.BaseEndpoint = new(cfg.S3Endpoint)
			o.UsePathStyle = true
			o.ResponseChecksumValidation = aws.ResponseChecksumValidationWhenRequired
		}
	})

	handler := proxy.New(
		storage.New(s3Client, cfg.S3Bucket),
		logger.WithComponent(slog.Default(), "proxy"),
	)

	app := fiber.New()
	handler.Register(app)

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		slog.Info("server started", slog.String("addr", fmt.Sprintf(":%s", data.Port)))
		return nil
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
