package main

import (
	"fmt"
	"log"
	"log/slog"
	"os"

	"github.com/gofiber/fiber/v3"
	"github.com/maneko-group/image-proxy/internal/pkg/logger"
)

func main() {
	logger.SetDefault(logger.Options{
		AppName:     "image-proxy",
		Environment: logger.Environment(os.Getenv("ENV")),
		Level:       slog.LevelInfo,
	})

	app := fiber.New()

	app.Get("/", func(c fiber.Ctx) error {
		return c.SendString("Hello, world!")
	})

	app.Hooks().OnListen(func(data fiber.ListenData) error {
		slog.Info("server started", slog.String("addr", fmt.Sprintf(":%s", data.Port)))
		return nil
	})

	log.Fatal(app.Listen(":3000"))
}
