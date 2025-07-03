package handlers

import (
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"strings"

	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BaseHandler contains common handler dependencies and methods
type BaseHandler struct {
	converter iface.Converter
	storage   *storage.FileStorage
	logger    utils.Logger
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(conv iface.Converter, storage *storage.FileStorage, logger utils.Logger) *BaseHandler {
	return &BaseHandler{
		converter: conv,
		storage:   storage,
		logger:    logger,
	}
}

// handleError is a helper for consistent error responses
func (h *BaseHandler) handleError(c *gin.Context, err error) {
	if h.logger != nil {
		h.logger.Error("Handler error", "error", err)
	}
	c.JSON(http.StatusInternalServerError, gin.H{
		"error": err.Error(),
	})
}

// getFileFromRequest handles file upload from form data
func (h *BaseHandler) getFileFromRequest(c *gin.Context, fieldName string) (*multipart.FileHeader, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, fmt.Errorf("no file uploaded")
		}
		return nil, fmt.Errorf("failed to get file from request: %w", err)
	}

	// TODO: Add file size validation if needed
	// if file.Size > maxSize {
	// 	return nil, fmt.Errorf("file too large")
	// }

	return file, nil
}

// getOutputFormat gets and validates the output format from the request
func (h *BaseHandler) getOutputFormat(c *gin.Context) (string, error) {
	format := c.PostForm("format")
	if format == "" {
		return "", fmt.Errorf("output format is required")
	}

	// Remove leading dot if present
	format = strings.TrimPrefix(format, ".")

	// TODO: Add format validation if needed
	// if !isValidFormat(format) {
	// 	return "", fmt.Errorf("unsupported output format: %s", format)
	// }

	return format, nil
}
