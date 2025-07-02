package video

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

// VideoConverter handles video format conversions
type VideoConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	tempDir    string
}

// NewVideoConverter creates a new VideoConverter
func NewVideoConverter(toolManager *tools.ToolManager, tempDir string) *VideoConverter {
	return &VideoConverter{
		BaseConverter: base.NewBaseConverter("video", map[string][]string{
			"mp4":  {"avi", "mov", "mkv", "wmv", "flv", "webm", "3gp", "gif"},
			"avi":  {"mp4", "mov", "mkv", "wmv", "flv", "webm", "3gp"},
			"mov":  {"mp4", "avi", "mkv", "wmv", "flv", "webm", "3gp"},
			"mkv":  {"mp4", "avi", "mov", "wmv", "flv", "webm", "3gp"},
			"wmv":  {"mp4", "avi", "mov", "mkv", "flv", "webm", "3gp"},
			"flv":  {"mp4", "avi", "mov", "mkv", "wmv", "webm", "3gp"},
			"webm": {"mp4", "avi", "mov", "mkv", "wmv", "flv", "3gp"},
			"3gp":  {"mp4", "avi", "mov", "mkv", "wmv", "flv", "webm"},
		}),
		toolManager: toolManager,
		tempDir:     tempDir,
	}
}

// Convert converts a video file from one format to another
func (c *VideoConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
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

	// Add video codec and bitrate options
	videoCodec := c.getVideoCodecForFormat(targetFormat)
	if videoCodec != "" {
		args = append(args, "-c:v", videoCodec)
	}

	// Add audio codec
	audioCodec := c.getAudioCodecForFormat(targetFormat)
	args = append(args, "-c:a", audioCodec)

	// Add video bitrate if specified
	if bitrate, ok := options["video_bitrate"].(string); ok && bitrate != "" {
		args = append(args, "-b:v", bitrate)
	}

	// Add audio bitrate if specified
	if audioBitrate, ok := options["audio_bitrate"].(string); ok && audioBitrate != "" {
		args = append(args, "-b:a", audioBitrate)
	}

	// Add resolution if specified
	if width, ok := options["width"].(int); ok && width > 0 {
		args = append(args, "-vf", fmt.Sprintf("scale=%d:-1", width))
	}

	// Add frame rate if specified
	if fps, ok := options["fps"].(int); ok && fps > 0 {
		args = append(args, "-r", strconv.Itoa(fps))
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
			"failed to convert video file",
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

// getVideoCodecForFormat returns the appropriate video codec for the target format
func (c *VideoConverter) getVideoCodecForFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp4":
		return "libx264"
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
	case "webm":
		return "libvpx"
	case "3gp":
		return "h263"
	case "gif":
		return "gif"
	default:
		return "libx264" // Default to H.264
	}
}

// getAudioCodecForFormat returns the appropriate audio codec for the target format
func (c *VideoConverter) getAudioCodecForFormat(format string) string {
	switch strings.ToLower(format) {
	case "mp4", "mov", "mkv":
		return "aac"
	case "avi":
		return "libmp3lame"
	case "wmv":
		return "wmav2"
	case "flv":
		return "libmp3lame"
	case "webm":
		return "libvorbis"
	case "3gp":
		return "amr_nb"
	case "gif":
		return "" // GIF doesn't have audio
	default:
		return "aac" // Default to AAC
	}
}

// ValidateOptions validates the conversion options
func (c *VideoConverter) ValidateOptions(options map[string]interface{}) error {
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

	// Validate video bitrate if provided
	if bitrate, ok := options["video_bitrate"].(string); ok && bitrate != "" {
		if !strings.HasSuffix(bitrate, "k") {
			return iface.NewConversionError(
				"invalid_option",
				"video_bitrate must end with 'k' (e.g., '2000k')",
				nil,
			)
		}

		numericPart := strings.TrimSuffix(bitrate, "k")
		if _, err := strconv.Atoi(numericPart); err != nil {
			return iface.NewConversionError(
				"invalid_option",
				"video_bitrate must be a number followed by 'k' (e.g., '2000k')",
				err,
			)
		}
	}

	// Validate audio bitrate if provided
	if bitrate, ok := options["audio_bitrate"].(string); ok && bitrate != "" {
		if !strings.HasSuffix(bitrate, "k") {
			return iface.NewConversionError(
				"invalid_option",
				"audio_bitrate must end with 'k' (e.g., '128k')",
				nil,
			)
		}

		numericPart := strings.TrimSuffix(bitrate, "k")
		if _, err := strconv.Atoi(numericPart); err != nil {
			return iface.NewConversionError(
				"invalid_option",
				"audio_bitrate must be a number followed by 'k' (e.g., '128k')",
				err,
			)
		}
	}

	// Validate width if provided
	if width, ok := options["width"].(int); ok && width > 0 {
		if width < 32 || width > 7680 { // 8K resolution
			return iface.NewConversionError(
				"invalid_option",
				"width must be between 32 and 7680 pixels",
				nil,
			)
		}
	}

	// Validate FPS if provided
	if fps, ok := options["fps"].(int); ok && fps > 0 {
		if fps < 1 || fps > 120 {
			return iface.NewConversionError(
				"invalid_option",
				"fps must be between 1 and 120",
				nil,
			)
		}
	}

	return nil
}
