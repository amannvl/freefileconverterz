package handlers

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
)

// GetConversionStatus gets the status of a conversion
// @Summary Get conversion status
// @Description Gets the status of a file conversion by ID
// @Tags conversion
// @Produce json
// @Param id path string true "Conversion ID"
// @Success 200 {object} Conversion
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/convert/{id}/status [get]
func (h *Handler) GetConversionStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	h.mu.RLock()
	conversion, exists := h.conversions[id]
	h.mu.RUnlock()

	// If conversion is not found in memory, check if the file exists in storage
	if !exists {
		// Try to find the file in storage
		files, err := h.storage.List(c.Context(), "")
		if err != nil {
			h.logger.Error("Failed to list files in storage", "error", err)
			return h.errorResponse(c, fiber.StatusInternalServerError, "internal_error", "Failed to check conversion status", err)
		}

		// Look for a file with the conversion ID in its name
		for _, file := range files {
			if file.ID == id || strings.HasPrefix(file.Name, id) {
				h.mu.Lock()
				// Create a new conversion entry
				conversion = &Conversion{
					ID:             id,
					Status:         "completed",
					ConvertedName:  file.Name,
					OriginalName:   file.Name,
					FileSize:       file.Size,
					CreatedAt:      file.CreatedAt,
					CompletedAt:    file.UpdatedAt,
				}
				h.conversions[id] = conversion
				h.mu.Unlock()
				exists = true
				break
			}
		}

		if !exists {
			return h.errorResponse(c, fiber.StatusNotFound, "not_found", "Conversion not found", nil)
		}
	}

	// If conversion is complete, update the download URL
	if conversion.Status == "completed" && conversion.DownloadURL == "" {
		downloadURL, err := h.generateDownloadURL(conversion.ConvertedName)
		if err != nil {
			return h.errorResponse(c, fiber.StatusInternalServerError, "internal_error", "Failed to generate download URL", err)
		}
		conversion.DownloadURL = downloadURL
	}

	return c.JSON(conversion)
}

// DownloadFile downloads a converted file
// @Summary Download converted file
// @Description Downloads a converted file by ID
// @Tags conversion
// @Produce octet-stream
// @Param id path string true "Conversion ID"
// @Success 200
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/convert/{id}/download [get]
func (h *Handler) DownloadFile(c *fiber.Ctx) error {
	id := c.Params("id")
	h.mu.RLock()
	conversion, exists := h.conversions[id]
	h.mu.RUnlock()

	var convertedFileName string

	// If conversion is not found in memory or not completed, try to find the file in storage
	if !exists || conversion.Status != "completed" {
		h.logger.Info("Conversion not found in memory, checking storage", "conversionID", id)
		
		// List all files in storage to find one that matches the conversion ID
		files, err := h.storage.List(c.Context(), "")
		if err != nil {
			h.logger.Error("Failed to list files in storage", "error", err)
			return h.errorResponse(c, fiber.StatusInternalServerError, "internal_error", "Failed to check storage", err)
		}

		// Look for a file with the conversion ID in its name
		found := false
		for _, file := range files {
			if strings.HasPrefix(file.Name, id) {
				convertedFileName = file.Name
				found = true
				break
			}
		}

		if !found {
			return h.errorResponse(c, fiber.StatusNotFound, "not_found", "File not found or conversion not complete", nil)
		}

		// If we found the file but not in memory, create a conversion entry
		h.mu.Lock()
		conversion = &Conversion{
			ID:             id,
			Status:         "completed",
			ConvertedName:  convertedFileName,
			OriginalName:   convertedFileName,
			FileSize:       0, // We don't know the size yet
			CompletedAt:    time.Now(),
		}
		h.conversions[id] = conversion
		h.mu.Unlock()
	} else {
		convertedFileName = conversion.ConvertedName
	}

	// Get the file from storage
	file, err := h.storage.Get(c.Context(), convertedFileName)
	if err != nil {
		h.logger.Error("Failed to get file from storage", "error", err, "filename", convertedFileName)
		return h.errorResponse(c, fiber.StatusNotFound, "not_found", "File not found in storage", err)
	}
	defer file.Close()

	// Get file info from storage
	files, err := h.storage.List(c.Context(), "")
	if err != nil {
		h.logger.Error("Failed to list files in storage", "error", err)
		return h.errorResponse(c, fiber.StatusInternalServerError, "internal_error", "Failed to get file info", err)
	}

	// Find the file in the storage to get its size
	var fileSize int64
	for _, f := range files {
		if f.Name == convertedFileName {
			fileSize = f.Size
			break
		}
	}

	// Set appropriate headers for download
	c.Set(fiber.HeaderContentDisposition, fmt.Sprintf("attachment; filename=%s", convertedFileName))
	if fileSize > 0 {
		c.Set(fiber.HeaderContentLength, strconv.FormatInt(fileSize, 10))
	}
	c.Set(fiber.HeaderContentType, "application/octet-stream")
	
	// Stream the file from storage
	return c.Status(fiber.StatusOK).SendStream(file)
}

// ListConversions lists all conversions
// @Summary List conversions
// @Description Lists all file conversions
// @Tags conversion
// @Produce json
// @Success 200 {array} Conversion
// @Router /api/v1/convert [get]
func (h *Handler) ListConversions(c *fiber.Ctx) error {
	h.mu.RLock()
	defer h.mu.RUnlock()

	// Convert map to slice
	conversionList := make([]*Conversion, 0, len(h.conversions))
	for _, conv := range h.conversions {
		conversionList = append(conversionList, conv)
	}

	return c.JSON(conversionList)
}

// GetConversion gets a single conversion by ID
// @Summary Get conversion
// @Description Gets a file conversion by ID
// @Tags conversion
// @Produce json
// @Param id path string true "Conversion ID"
// @Success 200 {object} Conversion
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/convert/{id} [get]
func (h *Handler) GetConversion(c *fiber.Ctx) error {
	id := c.Params("id")
	h.mu.RLock()
	defer h.mu.RUnlock()

	conversion, exists := h.conversions[id]
	if !exists {
		return h.errorResponse(c, fiber.StatusNotFound, "not_found", "Conversion not found", nil)
	}

	return c.JSON(conversion)
}

// DeleteConversion deletes a conversion
// @Summary Delete conversion
// @Description Deletes a file conversion by ID
// @Tags conversion
// @Param id path string true "Conversion ID"
// @Success 204
// @Failure 404 {object} map[string]interface{}
// @Router /api/v1/convert/{id} [delete]
func (h *Handler) DeleteConversion(c *fiber.Ctx) error {
	id := c.Params("id")
	h.mu.Lock()
	defer h.mu.Unlock()

	conversion, exists := h.conversions[id]
	if !exists {
		return h.errorResponse(c, fiber.StatusNotFound, "not_found", "Conversion not found", nil)
	}

	// Delete the file from storage
	if conversion.ConvertedName != "" {
		_ = h.storage.Delete(c.Context(), conversion.ConvertedName) // Best effort delete
	}

	// Remove from map
	delete(h.conversions, id)

	return c.SendStatus(fiber.StatusNoContent)
}

// generateDownloadURL creates a download URL for a converted file
func (h *Handler) generateDownloadURL(filename string) (string, error) {
	if h.storage == nil {
		return "", fmt.Errorf("storage not initialized")
	}

	// Use the storage's URL method to generate the download URL
	// The storage implementation (S3 or local) will handle the URL generation appropriately
	url := h.storage.URL(filename)
	if url == "" {
		return "", fmt.Errorf("failed to generate download URL")
	}

	return url, nil
}
