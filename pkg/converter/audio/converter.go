package audio

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// Ensure AudioConverter implements iface.Converter
var _ iface.Converter = (*AudioConverter)(nil)

// AudioConverter handles audio format conversions
type AudioConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	tempDir    string
}

// NewAudioConverter creates a new AudioConverter
func NewAudioConverter(toolManager *tools.ToolManager, tempDir string) *AudioConverter {
	return &AudioConverter{
		BaseConverter: base.NewBaseConverter("audio", map[string][]string{
			"mp3":  {"wav", "aac", "flac", "ogg", "m4a", "wma"},
			"wav":  {"mp3", "aac", "flac", "ogg", "m4a", "wma"},
			"aac":  {"mp3", "wav", "flac", "ogg", "m4a", "wma"},
			"flac": {"mp3", "wav", "aac", "ogg", "m4a", "wma"},
			"ogg":  {"mp3", "wav", "aac", "flac", "m4a", "wma"},
			"m4a":  {"mp3", "wav", "aac", "flac", "ogg", "wma"},
			"wma":  {"mp3", "wav", "aac", "flac", "ogg", "m4a"},
		}),
		toolManager: toolManager,
		tempDir:     tempDir,
	}
}

// Convert converts an audio file from one format to another
func (c *AudioConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
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
	outputFile := filepath.Join(c.tempDir, "output."+targetFormat)
	defer os.Remove(outputFile)

	// Build FFmpeg command
	args := []string{
		"-y", // Overwrite output file if it exists
		"-i", tempInput.Name(),
	}

	// Add audio codec and bitrate options
	codec := c.getCodecForFormat(targetFormat)
	if codec != "" {
		args = append(args, "-c:a", codec)
	}

	// Add bitrate if specified
	if bitrate, ok := options["bitrate"].(string); ok && bitrate != "" {
		args = append(args, "-b:a", bitrate)
	}

	// Add sample rate if specified
	if sampleRate, ok := options["sample_rate"].(int); ok && sampleRate > 0 {
		args = append(args, "-ar", strconv.Itoa(sampleRate))
	}

	// Add channels if specified
	if channels, ok := options["channels"].(int); ok && channels > 0 {
		args = append(args, "-ac", strconv.Itoa(channels))
	}

	// Add output file
	args = append(args, outputFile)

	// Get FFmpeg path from tool manager
	ffmpegPath, err := c.toolManager.GetFFmpegPath()
	if err != nil {
		return nil, iface.NewConversionError(
			"tool_not_found",
			"FFmpeg not found",
			err,
		)
	}

	// Run FFmpeg
	cmd := exec.CommandContext(ctx, ffmpegPath, args...)

	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("ffmpeg command failed")
		return nil, iface.NewConversionError(
			"conversion_failed",
			"failed to convert audio file",
			err,
		)
	}

	// Open the output file for reading
	result, err := os.Open(outputFile)
	if err != nil {
		return nil, iface.NewConversionError(
			"io_error",
			"failed to open output file",
			err,
		)
	}

	return result, nil
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
		return ""
	}
}

// ValidateOptions validates the conversion options
func (c *AudioConverter) ValidateOptions(options map[string]interface{}) error {
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

	// Validate bitrate if provided
	if bitrate, ok := options["bitrate"].(string); ok && bitrate != "" {
		// Check if bitrate is in the correct format (e.g., "128k", "192k", "320k")
		if !strings.HasSuffix(bitrate, "k") {
			return iface.NewConversionError(
				"invalid_option",
				"bitrate must end with 'k' (e.g., '128k')",
				nil,
			)
		}

		// Parse the numeric part of the bitrate
		numericPart := strings.TrimSuffix(bitrate, "k")
		if _, err := strconv.Atoi(numericPart); err != nil {
			return iface.NewConversionError(
				"invalid_option",
				"bitrate must be a number followed by 'k' (e.g., '128k')",
				err,
			)
		}
	}

	// Validate sample rate if provided
	if sampleRate, ok := options["sample_rate"].(int); ok && sampleRate > 0 {
		// Common sample rates: 8000, 11025, 16000, 22050, 32000, 44100, 48000, 88200, 96000, 176400, 192000
		validRates := map[int]bool{
			8000:   true,
			11025:  true,
			16000:  true,
			22050:  true,
			32000:  true,
			44100:  true,
			48000:  true,
			88200:  true,
			96000:  true,
			176400: true,
			192000: true,
		}

		if !validRates[sampleRate] {
			return iface.NewConversionError(
				"invalid_option",
				"unsupported sample rate. Common rates: 8000, 11025, 16000, 22050, 32000, 44100, 48000, 88200, 96000, 176400, 192000",
				nil,
			)
		}
	}

	// Validate channels if provided
	if channels, ok := options["channels"].(int); ok && channels > 0 {
		if channels < 1 || channels > 8 {
			return iface.NewConversionError(
				"invalid_option",
				"number of channels must be between 1 and 8",
				nil,
			)
		}
	}

	return nil
}
