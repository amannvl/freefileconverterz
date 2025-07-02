package tools

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/amannvl/freefileconverterz/internal/utils"
	"github.com/rs/zerolog/log"
)

type ToolManager struct {
	binManager *utils.BinaryManager
	tempDir string
}

// BinaryManager handles downloading and managing binary dependencies
type BinaryManager struct {
	BinDir string
}

func NewToolManager(binDir, tempDir string) (*ToolManager, error) {
	// Create bin and temp directories if they don't exist
	for _, dir := range []string{binDir, tempDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", dir, err)
		}
	}

	return &ToolManager{
		binManager: &utils.BinaryManager{BinDir: binDir},
		tempDir:    tempDir,
	}, nil
}

// EnsureTools ensures that all required tools are available
func (tm *ToolManager) EnsureTools(ctx context.Context) error {
	// List of tools to check
	toolsToCheck := []struct {
		name    string
		getPath func() (string, error)
	}{
		{"LibreOffice", tm.GetLibreOfficePath},
		{"ImageMagick", tm.GetImageMagickPath},
		{"FFmpeg", tm.GetFFmpegPath},
		{"7z", tm.Get7zPath},
		{"unrar", tm.GetUnrarPath},
	}

	// Check each tool
	for _, tool := range toolsToCheck {
		path, err := tool.getPath()
		if err != nil {
			return fmt.Errorf("failed to find %s: %w", tool.name, err)
		}

		// Check if the tool exists and is executable
		info, err := os.Stat(path)
		if err != nil {
			return fmt.Errorf("%s not found at %s: %w", tool.name, path, err)
		}

		// Check if the file is executable
		if runtime.GOOS != "windows" && info.Mode()&0111 == 0 {
			return fmt.Errorf("%s at %s is not executable", tool.name, path)
		}

		log.Info().
			Str("tool", tool.name).
			Str("path", path).
			Msg("Verified tool is available")
	}

	return nil
}

// GetBinaryPath returns the full path to a binary in the bin directory
func (bm *BinaryManager) GetBinaryPath(name string) string {
	return filepath.Join(bm.BinDir, name)
}

// GetLibreOfficePath returns the path to the LibreOffice binary
func (tm *ToolManager) GetLibreOfficePath() (string, error) {
	// On macOS, we'll use the system LibreOffice if available
	if runtime.GOOS == "darwin" {
		if path, err := exec.LookPath("soffice"); err == nil {
			return path, nil
		}
		// Try the standard macOS path
		stdPath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		if _, err := os.Stat(stdPath); err == nil {
			return stdPath, nil
		}
	}

	// For other platforms or if not found on macOS, use the bundled binary
	return tm.binManager.GetBinaryPath("libreoffice"), nil
}

// GetImageMagickPath returns the path to the ImageMagick convert binary
func (tm *ToolManager) GetImageMagickPath() (string, error) {
	// First try to find in system PATH
	if path, err := exec.LookPath("convert"); err == nil {
		return path, nil
	}

	// If not found, use the bundled binary
	return tm.binManager.GetBinaryPath("convert"), nil
}

// GetFFmpegPath returns the path to the FFmpeg binary
func (tm *ToolManager) GetFFmpegPath() (string, error) {
	// First try to find in system PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// If not found, use the bundled binary
	return tm.binManager.GetBinaryPath("ffmpeg"), nil
}

// Get7zPath returns the path to the 7z binary
func (tm *ToolManager) Get7zPath() (string, error) {
	// First try to find in system PATH
	if path, err := exec.LookPath("7z"); err == nil {
		return path, nil
	}

	// If not found, use the bundled binary
	return tm.binManager.GetBinaryPath("7z"), nil
}

// GetUnrarPath returns the path to the unrar binary
func (tm *ToolManager) GetUnrarPath() (string, error) {
	// First try to find in system PATH
	if path, err := exec.LookPath("unrar"); err == nil {
		return path, nil
	}

	// If not found, use the bundled binary
	return tm.binManager.GetBinaryPath("unrar"), nil
}

// Command creates a new command with the tool
func (tm *ToolManager) Command(ctx context.Context, tool string, args ...string) (*exec.Cmd, error) {
	var path string
	var err error

	switch tool {
	case "libreoffice":
		path, err = tm.GetLibreOfficePath()
	case "imagemagick":
		path, err = tm.GetImageMagickPath()
	case "ffmpeg":
		path, err = tm.GetFFmpegPath()
	case "7z":
		path, err = tm.Get7zPath()
	case "unrar":
		path, err = tm.GetUnrarPath()
	default:
		return nil, fmt.Errorf("unknown tool: %s", tool)
	}

	if err != nil {
		return nil, err
	}

	cmd := exec.CommandContext(ctx, path, args...)
	// Set the PATH to include our bin directory
	cmd.Env = append(os.Environ(), 
		fmt.Sprintf("PATH=%s:%s", filepath.Dir(path), os.Getenv("PATH")),
	)

	return cmd, nil
}

// Cleanup removes temporary files
func (tm *ToolManager) Cleanup() error {
	// Remove all files in temp directory
	entries, err := os.ReadDir(tm.tempDir)
	if err != nil {
		return fmt.Errorf("failed to read temp directory: %w", err)
	}

	for _, entry := range entries {
		if !entry.IsDir() {
			path := filepath.Join(tm.tempDir, entry.Name())
			if err := os.Remove(path); err != nil {
				log.Error().Err(err).Str("path", path).Msg("Failed to remove temp file")
			}
		}
	}

	return nil
}
