package handlers

import (
	"context"
	"fmt"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ConvertHandler handles file conversion requests
type ConvertHandler struct {
	*BaseHandler
}

// NewConvertHandler creates a new ConvertHandler
func NewConvertHandler(conv *converter.Converter, storage *storage.FileStorage, logger *utils.Logger) *ConvertHandler {
	return &ConvertHandler{
		BaseHandler: NewBaseHandler(conv, storage, logger),
	}
}

// ConvertFile handles file conversion requests
func (h *ConvertHandler) ConvertFile(c *gin.Context) {
	// Get the uploaded file
	file, err := h.getFileFromRequest(c, "file")
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Get the output format
	outputFormat, err := h.getOutputFormat(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	// Save the uploaded file
	inputFilename, err := h.storage.SaveFile(file)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to save file: %w", err))
		return
	}
	defer h.storage.DeleteFile(inputFilename)

	// Generate output filename
	ext := filepath.Ext(inputFilename)
	outputFilename := inputFilename[:len(inputFilename)-len(ext)] + "." + outputFormat

	// Get full paths
	inputPath, err := h.storage.GetFilePath(inputFilename)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get input file path: %w", err))
		return
	}

	outputPath, err := h.storage.GetFilePath(outputFilename)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to get output file path: %w", err))
		return
	}

	// Create a context with timeout
	ctx, cancel := context.WithTimeout(c.Request.Context(), 5*time.Minute)
	defer cancel()

	// Perform the conversion
	h.logger.Info("Starting conversion", 
		"input", inputPath, 
		"output", outputPath,
		"format", outputFormat)

	err = h.converter.Convert(ctx, inputPath, outputPath, map[string]interface{}{
		"format": outputFormat,
	})

	if err != nil {
		h.handleError(c, fmt.Errorf("conversion failed: %w", err))
		return
	}

	// Check if output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		h.handleError(c, utils.NewAppError(
			utils.ErrFileConversion,
			"Conversion completed but output file was not created",
			http.StatusInternalServerError,
		))
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", filepath.Base(outputPath)))
	c.Header("Content-Type", "application/octet-stream")

	// Serve the file
	c.File(outputPath)

	// Clean up the output file
	go func() {
		time.Sleep(30 * time.Second) // Give time for the file to be downloaded
		if err := h.storage.DeleteFile(outputFilename); err != nil {
			h.logger.Error(fmt.Errorf("failed to delete output file: %w", err), 
				"filename", outputFilename)
		}
	}()
}

// GetSupportedFormats returns the list of supported conversion formats
func GetSupportedFormats(c *gin.Context) {
	// This would be populated from the converter package
	formats := map[string][]string{
		"documents": {"pdf", "docx", "odt", "txt", "rtf"},
		"images": {"jpg", "jpeg", "png", "gif", "webp", "bmp"},
		"archives": {"zip", "tar", "gz", "7z", "rar"},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formats,
	})
}

// HealthCheck returns the health status of the service
func HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}
