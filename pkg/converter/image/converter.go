package image

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// ImageConverter handles image format conversions
type ImageConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
}

// NewImageConverter creates a new ImageConverter
func NewImageConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &ImageConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported formats
	converter.AddSupportedConversion("jpg", "png", "jpeg", "gif", "bmp", "tiff", "webp")
	converter.AddSupportedConversion("jpeg", "png", "jpg", "gif", "bmp", "tiff", "webp")
	converter.AddSupportedConversion("png", "jpg", "jpeg", "gif", "bmp", "tiff", "webp")
	converter.AddSupportedConversion("gif", "png", "jpg", "jpeg", "bmp", "tiff")
	converter.AddSupportedConversion("bmp", "png", "jpg", "jpeg", "gif", "tiff")
	converter.AddSupportedConversion("tiff", "png", "jpg", "jpeg", "gif", "bmp")
	converter.AddSupportedConversion("webp", "png", "jpg", "jpeg")
	converter.AddSupportedConversion("heic", "jpg", "jpeg", "png")
	converter.AddSupportedConversion("heif", "jpg", "jpeg", "png")

	return converter
}

// Convert converts an image from one format to another using ImageMagick
func (c *ImageConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Get the file extension to determine the target format
	extension := filepath.Ext(outputPath)
	if len(extension) == 0 {
		return iface.NewConversionError(
			"invalid_output",
			"output path must have an extension",
			nil,
		)
	}
	targetFormat := extension[1:] // Remove the dot

	// Log the conversion attempt
	log.Info().
		Str("source", inputPath).
		Str("target", outputPath).
		Str("target_format", targetFormat).
		Msg("Starting image conversion with ImageMagick")

	// Get ImageMagick path from tool manager
	convertPath, err := c.toolManager.GetImageMagickPath()
	if err != nil {
		return iface.NewConversionError(
			"tool_not_found",
			"ImageMagick not found",
			err,
		)
	}

	// Use ImageMagick for all conversions for consistency and better format support
	cmd := exec.CommandContext(ctx, convertPath, inputPath, outputPath)

	// TODO: Add support for quality and other options

	// Run the conversion
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Image conversion failed")

		return iface.NewConversionError(
			"conversion_failed",
			"failed to convert image",
			err,
		)
	}

	return nil
}

// Cleanup removes temporary files created during conversion
func (c *ImageConverter) Cleanup(files ...string) error {
	var lastErr error
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to remove temporary file")
			lastErr = err
		}
	}
	return lastErr
}
