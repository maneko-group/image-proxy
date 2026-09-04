package proxy

import (
	"errors"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/maneko-group/image-proxy/internal/pkg/storage"
)

type Handler struct {
	storage storage.Getter
	logger  *slog.Logger
}

func New(storage storage.Getter, logger *slog.Logger) *Handler {
	return &Handler{storage: storage, logger: logger}
}

func (h *Handler) Register(app *fiber.App) {
	app.Get("/*", h.Get)
}

func (h *Handler) Get(c fiber.Ctx) error {
	key := c.Params("*")
	if key == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}

	obj, err := h.storage.Get(c.Context(), key)
	if err != nil {
		if errors.Is(err, storage.ErrNotFound) {
			return c.SendStatus(fiber.StatusNotFound)
		}

		h.logger.Error("upstream fetch failed",
			slog.String("key", key),
			slog.String("error", err.Error()),
		)

		return c.SendStatus(fiber.StatusBadGateway)
	}

	if obj.ContentType != "" {
		c.Set(fiber.HeaderContentType, obj.ContentType)
	}

	return c.SendStream(obj.Body, int(obj.ContentLength))
}
