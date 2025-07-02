package handlers

import (
	"net/http"

	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-gonic/gin"
)

// BaseHandler contains common handler dependencies and methods
type BaseHandler struct {
	converter  *converter.Converter
	storage    *storage.FileStorage
	logger     *utils.Logger
}

// NewBaseHandler creates a new BaseHandler
func NewBaseHandler(converter *converter.Converter, storage *storage.FileStorage, logger *utils.Logger) *BaseHandler {
	return &BaseHandler{
		converter:  converter,
		storage:    storage,
		logger:     logger,
	}
}

// handleError is a helper for consistent error responses
func (h *BaseHandler) handleError(c *gin.Context, err error) {
	switch e := err.(type) {
	case *utils.AppError:
		c.JSON(e.StatusCode, e.ToResponse())
	default:
		h.logger.Error(err, "Internal server error")
		c.JSON(http.StatusInternalServerError, utils.ErrorResponse{
			Error: utils.ErrorInfo{
				Code:    utils.ErrInternal,
				Message: "Internal server error",
			},
		})
	}
}

// getFileFromRequest handles file upload from form data
func (h *BaseHandler) getFileFromRequest(c *gin.Context, fieldName string) (*multipart.FileHeader, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		if errors.Is(err, http.ErrMissingFile) {
			return nil, utils.NewAppError(
				utils.ErrInvalidInput,
				"No file uploaded",
				http.StatusBadRequest,
			)
		}
		return nil, fmt.Errorf("failed to get file from request: %w", err)
	}

	// Validate file size
	if file.Size > h.storage.MaxFileSize() {
		return nil, utils.ErrFileTooLarge
	}

	return file, nil
}

// getOutputFormat gets and validates the output format from the request
func (h *BaseHandler) getOutputFormat(c *gin.Context) (string, error) {
	format := c.PostForm("format")
	if format == "" {
		return "", utils.NewAppError(
			utils.ErrInvalidInput,
			"Output format is required",
			http.StatusBadRequest,
		)
	}

	// Remove leading dot if present
	format = strings.TrimPrefix(format, ".")

	// Validate the format
	if !h.converter.IsSupportedOutputFormat(format) {
		return "", utils.NewAppError(
			utils.ErrInvalidInput,
			fmt.Sprintf("Unsupported output format: %s", format),
			http.StatusBadRequest,
		)
	}

	return format, nil
}
