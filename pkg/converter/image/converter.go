package image

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/rs/zerolog/log"
)

// ImageConverter handles image format conversions
type ImageConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	tempDir    string
}

// NewImageConverter creates a new ImageConverter
func NewImageConverter(toolManager *tools.ToolManager, tempDir string) *ImageConverter {
	return &ImageConverter{
		BaseConverter: base.NewBaseConverter("image", map[string][]string{
			"jpg":  {"png", "jpeg", "gif", "bmp", "tiff", "webp"},
			"jpeg": {"png", "jpg", "gif", "bmp", "tiff", "webp"},
			"png":  {"jpg", "jpeg", "gif", "bmp", "tiff", "webp"},
			"gif":  {"png", "jpg", "jpeg", "bmp", "tiff"},
			"bmp":  {"png", "jpg", "jpeg", "gif", "tiff"},
			"tiff": {"png", "jpg", "jpeg", "gif", "bmp"},
			"webp": {"png", "jpg", "jpeg"},
			"heic": {"jpg", "jpeg", "png"},
			"heif": {"jpg", "jpeg", "png"},
		}),
		toolManager: toolManager,
		tempDir:     tempDir,
	}
}

// Convert converts an image from one format to another using ImageMagick for better format support
func (c *ImageConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
	sourceFormat, _ := options["source_format"].(string)
	targetFormat, _ := options["target_format"].(string)

	// Create a temporary file for the input
	tempInput, err := os.CreateTemp(c.tempDir, "input-*."+sourceFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary input file: %w", err)
	}
	defer os.Remove(tempInput.Name())

	// Write the input to the temporary file
	if _, err := io.Copy(tempInput, input); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempInput.Close()

	// Create a temporary file for the output
	outputExt := "." + targetFormat
	tempOutput, err := os.CreateTemp(c.tempDir, "output-*"+outputExt)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary output file: %w", err)
	}
	tempOutput.Close()
	defer os.Remove(tempOutput.Name())

	// Log the conversion attempt
	log.Info().
		Str("source", tempInput.Name()).
		Str("target_format", targetFormat).
		Msg("Starting image conversion with ImageMagick")

	// Get ImageMagick path from tool manager
	convertPath, err := c.toolManager.GetImageMagickPath()
	if err != nil {
		return nil, iface.NewConversionError(
			"tool_not_found",
			"ImageMagick not found",
			err,
		)
	}

	// Use ImageMagick for all conversions for consistency and better format support
	cmd := exec.CommandContext(ctx, convertPath, tempInput.Name())

	// Add quality option if provided
	if q, ok := options["quality"].(int); ok && q > 0 && q <= 100 {
		cmd.Args = append(cmd.Args, "-quality", fmt.Sprintf("%d", q))
	}

	// Add output file
	cmd.Args = append(cmd.Args, tempOutput.Name())

	// Run the conversion
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Image conversion failed")
		return nil, iface.NewConversionError(
			"conversion_failed",
			"failed to convert image",
			err,
		)
	}

	// Return the converted file
	convertedFile, err := os.Open(tempOutput.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to open converted file: %w", err)
	}

	return convertedFile, nil
}

// convertHeic converts HEIC/HEIF images to other formats using ImageMagick
// This is now a compatibility function that uses the main Convert function
func (c *ImageConverter) convertHeic(input io.Reader, targetFormat string) (io.Reader, error) {
	// Since we're using the system's convert command for all image conversions now,
	options := map[string]interface{}{
		"source_format": "heic",
		"target_format": targetFormat,
	}
	return c.Convert(context.Background(), input, options)
}

// ValidateOptions validates the conversion options
func (c *ImageConverter) ValidateOptions(options map[string]interface{}) error {
	// Check required options
	if _, ok := options["source_format"]; !ok {
		return iface.NewConversionError("missing_option", "source_format is required", nil)
	}

	if _, ok := options["target_format"]; !ok {
		return iface.NewConversionError("missing_option", "target_format is required", nil)
	}

	sourceFormat, ok := options["source_format"].(string)
	if !ok {
		return iface.NewConversionError("invalid_option", "source_format must be a string", nil)
	}

	targetFormat, ok := options["target_format"].(string)
	if !ok {
		return iface.NewConversionError("invalid_option", "target_format must be a string", nil)
	}

	if !c.SupportsConversion(sourceFormat, targetFormat) {
		return iface.NewConversionError(
			"unsupported_conversion",
			fmt.Sprintf("conversion from %s to %s is not supported", sourceFormat, targetFormat),
			nil,
		)
	}

	// Validate quality option if provided
	if quality, ok := options["quality"]; ok {
		if q, ok := quality.(int); !ok || q < 1 || q > 100 {
			return iface.NewConversionError(
				"invalid_option",
				"quality must be an integer between 1 and 100",
				nil,
			)
		}
	}

	// Validate compression level if provided
	if level, ok := options["compression_level"]; ok {
		if l, ok := level.(int); !ok || l < 0 || l > 9 {
			return iface.NewConversionError(
				"invalid_option",
				"compression_level must be an integer between 0 and 9",
				nil,
			)
		}
	}

	return nil
}
