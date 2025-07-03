package audio

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

// AudioConverter handles audio format conversions
type AudioConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
}

// NewAudioConverter creates a new AudioConverter
func NewAudioConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &AudioConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported formats
	converter.AddSupportedConversion("mp3", "wav", "aac", "flac", "ogg", "m4a", "wma")
	converter.AddSupportedConversion("wav", "mp3", "aac", "flac", "ogg", "m4a", "wma")
	converter.AddSupportedConversion("aac", "mp3", "wav", "flac", "ogg", "m4a", "wma")
	converter.AddSupportedConversion("flac", "mp3", "wav", "aac", "ogg", "m4a", "wma")
	converter.AddSupportedConversion("ogg", "mp3", "wav", "aac", "flac", "m4a", "wma")
	converter.AddSupportedConversion("m4a", "mp3", "wav", "aac", "flac", "ogg", "wma")
	converter.AddSupportedConversion("wma", "mp3", "wav", "aac", "flac", "ogg", "m4a")

	return converter
}

// Convert converts an audio file from one format to another using FFmpeg
func (c *AudioConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
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
		Msg("Starting audio conversion with FFmpeg")

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

	// Add audio codec for the target format
	codec := c.getCodecForFormat(targetFormat)
	if codec != "" {
		args = append(args, "-c:a", codec)
	}

	// Add output file
	args = append(args, outputPath)

	// Run FFmpeg
	cmd := exec.CommandContext(ctx, ffmpegPath, args...)
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("Audio conversion failed")

		return iface.NewConversionError(
			"conversion_failed",
			"failed to convert audio",
			err,
		)
	}

	return nil
}

// getCodecForFormat returns the appropriate audio codec for the target format
func (c *AudioConverter) getCodecForFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp3":
		return "libmp3lame"
	case "aac":
		return "aac"
	case "flac":
		return "flac"
	case "ogg":
		return "libvorbis"
	case "m4a":
		return "aac"
	case "wma":
		return "wmav2"
	case "wav":
		return "pcm_s16le"
	default:
		// Let FFmpeg choose the default codec for other formats
		return ""
	}
}

// Cleanup removes temporary files created during conversion
func (c *AudioConverter) Cleanup(files ...string) error {
	var lastErr error
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to remove temporary file")
			lastErr = err
		}
	}
	return lastErr
}
