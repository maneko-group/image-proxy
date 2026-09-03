package logger

import (
	"io"
	"log/slog"
	"os"

	"github.com/rs/zerolog"
	slogzerolog "github.com/samber/slog-zerolog/v2"
)

type Environment string

const (
	Development Environment = "development"
	Production  Environment = "production"
)

type Options struct {
	AppName     string
	Environment Environment
	Level       slog.Level
}

func New(options Options, additionalHandlers ...slog.Handler) *slog.Logger {
	handlers := []slog.Handler{
		newZeroLogHandler(options.Level, options.Environment),
	}

	handlers = append(handlers, additionalHandlers...)

	logger := slog.New(
		slog.NewMultiHandler(
			handlers...,
		),
	)

	if options.AppName != "" {
		logger = logger.With(
			slog.String("app", options.AppName),
		)
	}

	return logger
}

func SetDefault(options Options, additionalHandlers ...slog.Handler) {
	logger := New(options, additionalHandlers...)
	slog.SetDefault(logger)
}

func WithComponent(logger *slog.Logger, name string) *slog.Logger {
	return logger.With(
		slog.String("component", name),
	)
}

func Error(err error) slog.Attr {
	if err == nil {
		return slog.Any("error", err)
	}

	return slog.String("error", err.Error())
}

func newZeroLogHandler(level slog.Level, environment Environment) slog.Handler {
	var writer io.Writer

	switch environment {
	case Production:
		writer = os.Stderr
	default:
		writer = zerolog.ConsoleWriter{
			Out: os.Stderr,
		}
	}

	logger := zerolog.New(writer)
	option := slogzerolog.Option{
		Level:     level,
		Logger:    &logger,
		AddSource: true,
	}

	return option.NewZerologHandler()
}
