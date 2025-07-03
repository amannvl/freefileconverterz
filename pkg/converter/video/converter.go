package video

import (
	"context"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// VideoConverter handles video format conversions
type VideoConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
}

// NewVideoConverter creates a new VideoConverter
func NewVideoConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &VideoConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported formats and conversions
	converter.AddSupportedConversion("mp4", "avi", "mov", "mkv", "wmv", "flv", "webm", "3gp", "gif")
	converter.AddSupportedConversion("avi", "mp4", "mov", "mkv", "wmv", "flv", "webm", "3gp")
	converter.AddSupportedConversion("mov", "mp4", "avi", "mkv", "wmv", "flv", "webm", "3gp")
	converter.AddSupportedConversion("mkv", "mp4", "avi", "mov", "wmv", "flv", "webm", "3gp")
	converter.AddSupportedConversion("wmv", "mp4", "avi", "mov", "mkv", "flv", "webm", "3gp")
	converter.AddSupportedConversion("flv", "mp4", "avi", "mov", "mkv", "wmv", "webm", "3gp")
	converter.AddSupportedConversion("webm", "mp4", "avi", "mov", "mkv", "wmv", "flv", "3gp")
	converter.AddSupportedConversion("3gp", "mp4", "avi", "mov", "mkv", "wmv", "flv", "webm")

	return converter
}

// Convert converts a video file from one format to another
func (c *VideoConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
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
		Msg("Starting video conversion with FFmpeg")

	// Get FFmpeg path from tool manager
	ffmpegPath, err := c.toolManager.GetFFmpegPath()
	if err != nil {
		return iface.NewConversionError(
			"tool_not_found",
			"FFmpeg not found",
			err,
		)
	}

	// Build FFmpeg command
	args := []string{
		"-y", // Overwrite output file if it exists
		"-i", inputPath,
	}

	// Add video codec for the target format
	videoCodec := c.getVideoCodecForFormat(targetFormat)
	if videoCodec != "" {
		args = append(args, "-c:v", videoCodec)
	}

	// Add audio codec for the target format
	audioCodec := c.getAudioCodecForFormat(targetFormat)
	if audioCodec != "" {
		args = append(args, "-c:a", audioCodec)
	}

	// Add output file
	args = append(args, outputPath)

	// Execute FFmpeg
	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Video conversion failed")

		return iface.NewConversionError(
			"conversion_failed",
			"failed to convert video",
			err,
		)
	}

	log.Info().
		Str("output", outputPath).
		Msg("Video conversion completed successfully")

	return nil
}

// getVideoCodecForFormat returns the appropriate video codec for the target format
func (c *VideoConverter) getVideoCodecForFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp4":
		return "libx264"
	case "webm":
		return "libvpx"
	case "avi":
		return "mpeg4"
	case "mov":
		return "mpeg4"
	case "mkv":
		return "libx264"
	case "wmv":
		return "wmv2"
	case "flv":
		return "flv"
	case "3gp":
		return "mpeg4"
	case "gif":
		return "gif"
	default:
		// Let FFmpeg choose the default codec
		return ""
	}
}

// getAudioCodecForFormat returns the appropriate audio codec for the target format
func (c *VideoConverter) getAudioCodecForFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp4", "m4a":
		return "aac"
	case "webm":
		return "libvorbis"
	case "avi", "mov", "mkv", "wmv", "flv", "3gp":
		return "aac"
	case "gif":
		// GIF doesn't support audio
		return "none"
	default:
		// Let FFmpeg choose the default codec
		return "aac"
	}
}

// Cleanup removes temporary files created during conversion
func (c *VideoConverter) Cleanup(files ...string) error {
	var lastErr error
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to remove temporary file")
			lastErr = err
		}
	}
	return lastErr
}
