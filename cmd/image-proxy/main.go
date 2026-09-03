package main

import (
	"fmt"
	"log"
	"log/slog"

	"github.com/gofiber/fiber/v3"
	"github.com/maneko-group/image-proxy/internal/pkg/config"
	"github.com/maneko-group/image-proxy/internal/pkg/logger"
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

	app := fiber.New()

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		slog.Info("server started", slog.String("addr", fmt.Sprintf(":%s", data.Port)))
		return nil
	})

	log.Fatal(app.Listen(fmt.Sprintf(":%d", cfg.Port)))
}
