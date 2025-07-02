package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
)

// ImageConverter handles image format conversions using ImageMagick
type ImageConverter struct {
	binMgr *BinaryManager
}

// NewImageConverter creates a new ImageConverter
func NewImageConverter() *ImageConverter {
	return &ImageConverter{
		binMgr: NewBinaryManager(),
	}
}

// SupportedFormats returns the supported image formats and conversions
func (ic *ImageConverter) SupportedFormats() map[string][]string {
	return map[string][]string{
		"jpg":  {"png", "jpeg", "gif", "bmp", "tiff", "webp"},
		"jpeg": {"png", "jpg", "gif", "bmp", "tiff", "webp"},
		"png":  {"jpg", "jpeg", "gif", "bmp", "tiff", "webp"},
		"gif":  {"png", "jpg", "jpeg", "bmp", "webp"},
		"bmp":  {"png", "jpg", "jpeg", "tiff"},
		"tiff": {"png", "jpg", "jpeg", "bmp"},
		"webp": {"png", "jpg", "jpeg", "gif"},
	}
}

// Convert converts an image from one format to another using ImageMagick
func (ic *ImageConverter) Convert(ctx context.Context, inputPath, outputPath string, options map[string]interface{}) error {
	// Check if convert (ImageMagick) is available
	if !ic.binMgr.BinaryExists("convert") {
		return fmt.Errorf("ImageMagick (convert) is required for image conversion")
	}

	// Build the convert command
	args := []string{inputPath}

	// Apply options if provided
	if quality, ok := options["quality"].(int); ok && quality > 0 && quality <= 100 {
		args = append(args, "-quality", strconv.Itoa(quality))
	}

	if width, ok := options["width"].(int); ok && width > 0 {
		args = append(args, "-resize", fmt.Sprintf("%d", width))
	}

	if height, ok := options["height"].(int); ok && height > 0 {
		if len(args) > 1 && args[len(args)-2] == "-resize" {
			// Append to existing resize option
			args[len(args)-1] = fmt.Sprintf("%sx%d", args[len(args)-1], height)
		} else {
			args = append(args, "-resize", fmt.Sprintf("x%d", height))
		}
	}

	// Add output file path
	args = append(args, outputPath)

	// Run the convert command
	cmd := exec.CommandContext(ctx, "convert", args...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("ImageMagick conversion failed: %v", err)
	}

	// Verify the output file was created
	if _, err := os.Stat(outputPath); os.IsNotExist(err) {
		return fmt.Errorf("output file was not created: %s", outputPath)
	}

	return nil
}
