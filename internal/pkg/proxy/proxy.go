package proxy

import (
	"errors"
	"io"
	"log/slog"

	"github.com/gofiber/fiber/v3"

	"github.com/maneko-group/image-proxy/internal/pkg/storage"
	"github.com/maneko-group/image-proxy/internal/pkg/transform"
)

type Handler struct {
	storage   storage.Getter
	logger    *slog.Logger
	processor *transform.Processor
}

func New(storage storage.Getter, logger *slog.Logger, processor *transform.Processor) *Handler {
	return &Handler{storage: storage, logger: logger, processor: processor}
}

func (h *Handler) Register(app *fiber.App) {
	app.Get("/*", h.Get)
}

func (h *Handler) Get(c fiber.Ctx) error {
	key := c.Params("*")
	if key == "" {
		return c.SendStatus(fiber.StatusNotFound)
	}
	params, err := transform.ParseParams(c.Queries())
	if err != nil {
		return c.SendStatus(fiber.StatusBadRequest)
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

	defer obj.Body.Close()
	source, err := io.ReadAll(obj.Body)
	if err != nil {
		return c.SendStatus(fiber.StatusBadGateway)
	}
	result, err := h.processor.Process(source, params, c.Get(fiber.HeaderAccept))
	if err != nil {
		h.logger.Error("image processing failed", slog.String("key", key), slog.String("error", err.Error()))
		return c.SendStatus(fiber.StatusInternalServerError)
	}
	setHeaders(c, result.ContentType, result.VaryAccept)
	return c.Send(result.Body)
}

func setHeaders(c fiber.Ctx, contentType string, varyAccept bool) {
	if contentType != "" {
		c.Set(fiber.HeaderContentType, contentType)
	}
	c.Set(fiber.HeaderCacheControl, "public, max-age=31536000, immutable")
	c.Set(fiber.HeaderAccessControlAllowOrigin, "*")
	if varyAccept {
		c.Set(fiber.HeaderVary, fiber.HeaderAccept)
	}
}
