package handlers

import (
	"log/slog"
	"sync"
	"time"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter/factory"
	"github.com/gofiber/fiber/v2"
)

// Handler provides application handlers
type Handler struct {
	config           *config.Config
	storage          storage.Storage
	converterFactory *factory.ConverterFactory
	logger           *slog.Logger
	mu              sync.RWMutex
	conversions     map[string]*Conversion
}

// NewHandler creates a new handler instance
func NewHandler(cfg *config.Config, store storage.Storage, factory *factory.ConverterFactory, log *slog.Logger) *Handler {
	if log == nil {
		log = slog.Default()
	}

	return &Handler{
		config:           cfg,
		storage:          store,
		converterFactory: factory,
		logger:           log,
		conversions:      make(map[string]*Conversion),
	}
}

// HealthCheck handles health check requests
func (h *Handler) HealthCheck(c *fiber.Ctx) error {
	return c.JSON(fiber.Map{
		"status":  "ok",
		"version": h.config.App.Version,
	})
}

// errorResponse is a helper for error responses
func (h *Handler) errorResponse(c *fiber.Ctx, status int, code, message string, err error) error {
	response := fiber.Map{
		"success": false,
		"error":   code,
		"message": message,
	}

	// Include error details in development
	if h.config.App.Env == "development" && err != nil {
		response["details"] = err.Error()
	}

	return c.Status(status).JSON(response)
}

// updateConversionStatus updates the status of a conversion
func (h *Handler) updateConversionStatus(conversionID, status string) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conv, ok := h.conversions[conversionID]; ok {
		conv.Status = status
	}
}

// updateConversionError updates a conversion with an error status
func (h *Handler) updateConversionError(conversionID string, err error) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conv, ok := h.conversions[conversionID]; ok {
		now := time.Now()
		conv.Status = "error"
		conv.Error = err.Error()
		conv.CompletedAt = now
	}
}

// updateConversionSuccess updates a conversion with success status and file info
func (h *Handler) updateConversionSuccess(conversionID, convertedName string, fileSize int64) {
	h.mu.Lock()
	defer h.mu.Unlock()
	if conv, ok := h.conversions[conversionID]; ok {
		now := time.Now()
		conv.Status = "completed"
		conv.ConvertedName = convertedName
		conv.FileSize = fileSize
		conv.CompletedAt = now
		conv.DownloadURL = "/download/" + convertedName
	}
}

// successResponse is a helper for success responses
func (h *Handler) successResponse(c *fiber.Ctx, status int, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"success": true,
		"data":    data,
	})
}
