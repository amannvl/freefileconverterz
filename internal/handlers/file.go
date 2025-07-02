package handlers

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/utils"
)

// ConvertFile handles file conversion requests
// @Summary Convert a file to another format
// @Description Converts an uploaded file to the specified target format
// @Tags files
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "File to convert"
// @Param format formData string true "Target format to convert to"
// @Success 200 {object} map[string]interface{}
// @Failure 400 {object} map[string]interface{}
// @Failure 500 {object} map[string]interface{}
// @Router /api/v1/convert [post]
func (h *Handler) ConvertFile(c *fiber.Ctx) error {
	// Parse the multipart form
	form, err := c.MultipartForm()
	if err != nil {
		h.logger.Error("Failed to parse form", "error", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Failed to parse form: " + err.Error(),
		})
	}

	// Get the file from the form
	files := form.File["file"]
	if len(files) == 0 {
		h.logger.Error("No file uploaded")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "No file uploaded",
		})
	}

	fileHeader := files[0]

	// Get target format from form
	targetFormat := c.FormValue("format")
	if targetFormat == "" {
		h.logger.Error("No target format provided")
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Target format is required (use 'format' parameter)",
		})
	}

	h.logger.Info("Starting file conversion", 
		"filename", fileHeader.Filename, 
		"size", fileHeader.Size,
		"targetFormat", targetFormat,
	)

	// Create a new conversion
	conversionID := utils.UUIDv4()
	now := time.Now()
	conversion := &Conversion{
		ID:           conversionID,
		Status:       "pending",
		CreatedAt:    now,
		TargetFormat: targetFormat,
		OriginalName: fileHeader.Filename,
		FileSize:     fileHeader.Size,
	}

	// Store the conversion
	h.mu.Lock()
	h.conversions[conversionID] = conversion
	h.mu.Unlock()

	// Process the conversion in the background
	go h.processConversion(conversionID, &FileHeader{fileHeader}, targetFormat)

	h.logger.Info("Conversion started", "conversionID", conversionID)

	// Return the conversion ID to track progress
	return c.JSON(fiber.Map{
		"id":     conversionID,
		"status": "pending",
	})
}



// FileHeader wraps multipart.FileHeader to implement our interface
type FileHeader struct {
	*multipart.FileHeader
}

// Filename returns the name of the file
func (f *FileHeader) Filename() string {
	return f.FileHeader.Filename
}

// Size returns the size of the file in bytes
func (f *FileHeader) Size() int64 {
	return f.FileHeader.Size
}

// Open returns a ReadCloser for the file
func (f *FileHeader) Open() (multipart.File, error) {
	return f.FileHeader.Open()
}

// processConversion handles the actual file conversion in a separate goroutine
func (h *Handler) processConversion(conversionID string, fileHeader *FileHeader, targetFormat string) {
	// Update status to processing
	h.updateConversionStatus(conversionID, "processing")

	// Log the start of conversion
	h.logger.Info("Starting conversion",
		"conversionID", conversionID,
		"filename", fileHeader.Filename(),
		"targetFormat", targetFormat,
	)

	// Open the uploaded file
	srcFile, err := fileHeader.Open()
	if err != nil {
		err = fmt.Errorf("failed to open uploaded file: %w", err)
		h.logger.Error("File open error", "error", err, "conversionID", conversionID)
		h.updateConversionError(conversionID, err)
		return
	}
	defer srcFile.Close()

	// Get the appropriate converter
	ext := filepath.Ext(fileHeader.Filename())
	if ext == "" {
		err := fmt.Errorf("file has no extension")
		h.logger.Error("File extension error", "error", err, "filename", fileHeader.Filename())
		h.updateConversionError(conversionID, err)
		return
	}

	sourceFormat := ext[1:] // Remove the dot
	h.logger.Info("Getting converter", 
		"sourceFormat", sourceFormat, 
		"targetFormat", targetFormat,
		"conversionID", conversionID,
	)

	converter, err := h.converterFactory.GetConverter(sourceFormat, targetFormat)
	if err != nil {
		err = fmt.Errorf("unsupported conversion from %s to %s: %w", sourceFormat, targetFormat, err)
		h.logger.Error("Converter error", "error", err, "conversionID", conversionID)
		h.updateConversionError(conversionID, err)
		return
	}

	h.logger.Info("Converter found", 
		"converter", fmt.Sprintf("%T", converter),
		"conversionID", conversionID,
	)

	// Prepare conversion options
	options := map[string]interface{}{
		"source_format": sourceFormat,
		"target_format": targetFormat,
	}

	// Convert the file
	h.logger.Info("Starting file conversion", 
		"conversionID", conversionID,
		"sourceFormat", sourceFormat,
		"targetFormat", targetFormat,
	)
	convertedReader, err := converter.Convert(context.Background(), srcFile, options)
	if err != nil {
		err = fmt.Errorf("conversion failed: %w", err)
		h.logger.Error("Conversion error", "error", err, "conversionID", conversionID)
		h.updateConversionError(conversionID, err)
		return
	}
	// Close the reader if it implements io.Closer
	if closer, ok := convertedReader.(io.Closer); ok {
		defer closer.Close()
	}

	// Read the converted data
	h.logger.Info("Reading converted data", "conversionID", conversionID)
	convertedData, err := io.ReadAll(convertedReader)
	if err != nil {
		err = fmt.Errorf("failed to read converted data: %w", err)
		h.logger.Error("Read converted data error", "error", err, "conversionID", conversionID)
		h.updateConversionError(conversionID, err)
		return
	}

	// Generate a unique filename for the converted file
	convertedName := fmt.Sprintf("%s.%s", utils.UUIDv4(), targetFormat)

	h.logger.Info("Saving converted file", 
		"conversionID", conversionID,
		"convertedName", convertedName,
		"size", len(convertedData),
	)

	// Save the converted file to storage
	if _, err := h.storage.Save(context.Background(), convertedName, convertedData); err != nil {
		err = fmt.Errorf("failed to save converted file: %w", err)
		h.logger.Error("Save converted file error", "error", err, "conversionID", conversionID)
		h.updateConversionError(conversionID, err)
		return
	}

	// Update conversion with success status
	h.logger.Info("Conversion completed successfully", "conversionID", conversionID)
	h.updateConversionSuccess(conversionID, convertedName, int64(len(convertedData)))
}


