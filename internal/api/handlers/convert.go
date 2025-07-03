package handlers

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/amannvl/freefileconverterz/internal/storage"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/gin-gonic/gin"
)

// ConvertHandler handles file conversion requests
type ConvertHandler struct {
	*BaseHandler
}

// NewConvertHandler creates a new ConvertHandler
func NewConvertHandler(conv iface.Converter, storage *storage.FileStorage, logger utils.Logger) *ConvertHandler {
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

	h.logger.Info("Received file upload",
		"filename", file.Filename,
		"size", file.Size,
		"content_type", file.Header.Get("Content-Type"))

	// Get the output format
	outputFormat, err := h.getOutputFormat(c)
	if err != nil {
		h.handleError(c, err)
		return
	}

	h.logger.Info("Requested output format", "format", outputFormat)

	// Save the uploaded file
	inputFilename, err := h.storage.SaveFile(file)
	if err != nil {
		h.handleError(c, fmt.Errorf("failed to save file: %w", err))
		return
	}
	defer h.storage.DeleteFile(inputFilename)

	// Generate output filename using only the original filename without extension
	ext := filepath.Ext(file.Filename)
	originalName := file.Filename[:len(file.Filename)-len(ext)]
	outputFilename := originalName + "." + outputFormat

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

	h.logger.Info("File paths prepared",
		"input_path", inputPath,
		"output_path", outputPath,
		"output_dir", filepath.Dir(outputPath))

	// Ensure output directory exists
	outputDir := filepath.Dir(outputPath)
	h.logger.Info("Ensuring output directory exists", "directory", outputDir)
	
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		h.logger.Error("Failed to create output directory",
			"error", err,
			"directory", outputDir)
		h.handleError(c, fmt.Errorf("failed to create output directory: %w", err))
		return
	}

	// Verify directory permissions
	if err := os.Chmod(outputDir, 0755); err != nil {
		h.logger.Info("Failed to set directory permissions",
			"error", err,
			"directory", outputDir)
	}

	// Remove any existing output file to avoid conflicts
	h.logger.Info("Checking for existing output file", "path", outputPath)
	if err := os.Remove(outputPath); err != nil && !os.IsNotExist(err) {
		h.logger.Error("Failed to remove existing output file",
			"error", err,
			"path", outputPath)
		h.handleError(c, fmt.Errorf("failed to remove existing output file: %w", err))
		return
	}

	// Create a context with timeout for the conversion
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute) // Increased timeout to 2 minutes
	defer cancel()

	// Log detailed information before conversion
	h.logger.Info("Starting document conversion", 
		"input_path", inputPath, 
		"output_path", outputPath,
		"output_dir", filepath.Dir(outputPath))

	// Log directory permissions
	if dirInfo, err := os.Stat(filepath.Dir(outputPath)); err == nil {
		h.logger.Info("Output directory info",
			"path", filepath.Dir(outputPath),
			"mode", dirInfo.Mode().String())
	}

	// Log input file permissions
	if fileInfo, err := os.Stat(inputPath); err == nil {
		h.logger.Info("Input file info",
			"path", inputPath,
			"size", fileInfo.Size(),
			"mode", fileInfo.Mode().String())
	}

	// Perform the conversion
	h.logger.Info("Calling converter.Convert")
	conversionStart := time.Now()
	err = h.converter.Convert(ctx, inputPath, outputPath)
	conversionDuration := time.Since(conversionStart)

	h.logger.Info("Converter.Convert completed", 
		"duration", conversionDuration.String(),
		"error", err)

	if err != nil {
		h.logger.Error("Conversion failed", 
			"error", err, 
			"input_path", inputPath, 
			"output_path", outputPath,
			"conversion_duration", conversionDuration.String())
		
		// Check if output file exists despite the error
		if stat, statErr := os.Stat(outputPath); statErr == nil {
			h.logger.Info("Output file exists but conversion reported error",
				"path", outputPath,
				"size", stat.Size())
		}

		// List contents of output directory for debugging
		if files, readErr := os.ReadDir(filepath.Dir(outputPath)); readErr == nil {
			fileList := make([]string, 0, len(files))
			for _, file := range files {
				fileList = append(fileList, file.Name())
			}
			h.logger.Info("Contents of output directory",
				"directory", filepath.Dir(outputPath),
				"files", fileList)
		}
		
		// Clean up output file if conversion fails
		if removeErr := os.Remove(outputPath); removeErr != nil && !os.IsNotExist(removeErr) {
			h.logger.Error("Failed to delete output file", 
				"error", removeErr, 
				"path", outputPath)
		}
		
		h.handleError(c, fmt.Errorf("conversion failed: %w", err))
		return
	}

	h.logger.Info("Conversion completed successfully", 
		"input_path", inputPath, 
		"output_path", outputPath)

	// Verify the output file was created and has content
	info, err := os.Stat(outputPath)
	if err != nil || info.Size() == 0 {
		if err == nil {
			err = fmt.Errorf("conversion produced empty file")
		}
		// Clean up the input file
		if removeErr := os.Remove(inputPath); removeErr != nil && !os.IsNotExist(removeErr) {
			h.logger.Error("Failed to delete input file", "error", removeErr, "path", inputPath)
		}
		h.handleError(c, fmt.Errorf("conversion failed: %w", err))
		return
	}

	// Set headers for file download
	c.Header("Content-Description", "File Transfer")
	c.Header("Content-Transfer-Encoding", "binary")
	c.Header("Content-Disposition", fmt.Sprintf("attachment; filename=%s", outputFilename))
	c.Header("Content-Type", "application/octet-stream")

	// Stream the file
	c.File(outputPath)

	// Clean up the output file after sending
	go func() {
		// Wait a bit longer to ensure file transfer
		time.Sleep(2 * time.Second)
		if err := os.Remove(outputPath); err != nil {
			h.logger.Error("Failed to delete output file", "error", err, "path", outputPath)
		} else {
			h.logger.Info("Cleaned up temporary file", "path", outputPath)
		}
	}()
}

// GetSupportedFormats returns the list of supported conversion formats
func (h *ConvertHandler) GetSupportedFormats(c *gin.Context) {
	// This would be populated from the converter package
	formats := map[string]map[string][]string{
		"documents": {
			"input":  {"pdf", "docx", "doc", "odt", "txt", "rtf"},
			"output": {"pdf", "docx", "odt", "txt", "rtf"},
		},
		"images": {
			"input":  {"jpg", "jpeg", "png", "gif", "webp", "bmp"},
			"output": {"jpg", "jpeg", "png", "gif", "webp", "bmp"},
		},
		"archives": {
			"input":  {"zip", "tar", "gz", "7z", "rar"},
			"output": {"zip", "tar", "gz", "7z"},
		},
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    formats,
	})
}

// HealthCheck returns the health status of the service
func (h *ConvertHandler) HealthCheck(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"status":  "ok",
		"version": "1.0.0",
		"time":    time.Now().Format(time.RFC3339),
	})
}
