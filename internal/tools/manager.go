package tools

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/amannvl/freefileconverterz/internal/utils"
)

type ToolManager struct {
	binManager *utils.BinaryManager
	tempDir    string
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

	// Create the binary manager
	binManager, err := utils.NewBinaryManager(binDir)
	if err != nil {
		return nil, fmt.Errorf("failed to create binary manager: %w", err)
	}

	return &ToolManager{
		binManager: binManager,
		tempDir:    tempDir,
	}, nil
}

func isExecutable(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	// Check if the file is executable
	if runtime.GOOS != "windows" && info.Mode()&0111 == 0 {
		return false
	}

	return true
}

// EnsureTools verifies that all required tools are available
func (tm *ToolManager) EnsureTools() error {
	var missingTools []string

	// Check LibreOffice
	if path, err := tm.GetLibreOfficePath(); err != nil {
		log.Printf("LibreOffice not found: %v", err)
		missingTools = append(missingTools, "LibreOffice (for document conversions)")
	} else if !isExecutable(path) {
		log.Printf("LibreOffice is not executable at %s", path)
		missingTools = append(missingTools, "LibreOffice (not executable)")
	}

	// Check ImageMagick
	if path, err := tm.GetImageMagickPath(); err != nil {
		log.Printf("ImageMagick not found: %v", err)
		missingTools = append(missingTools, "ImageMagick (for image conversions)")
	} else if !isExecutable(path) {
		log.Printf("ImageMagick convert is not executable at %s", path)
		missingTools = append(missingTools, "ImageMagick (not executable)")
	}

	// Check FFmpeg
	if path, err := tm.GetFFmpegPath(); err != nil {
		log.Printf("FFmpeg not found: %v", err)
		missingTools = append(missingTools, "FFmpeg (for audio/video conversions)")
	} else if !isExecutable(path) {
		log.Printf("FFmpeg is not executable at %s", path)
		missingTools = append(missingTools, "FFmpeg (not executable)")
	}

	// Check 7z
	if path, err := tm.Get7zPath(); err != nil {
		log.Printf("7z not found: %v", err)
		missingTools = append(missingTools, "p7zip (for archive operations)")
	} else if !isExecutable(path) {
		log.Printf("7z is not executable at %s", path)
		missingTools = append(missingTools, "7z (not executable)")
	}

	// Check unrar (optional, but log a warning if not found)
	if path, err := tm.GetUnrarPath(); err != nil {
		log.Printf("Warning: unrar not found. RAR archive support will be disabled: %v", err)
	} else if !isExecutable(path) {
		log.Printf("Warning: unrar is not executable at %s. RAR archive support will be disabled", path)
	}

	if len(missingTools) > 0 {
		errMsg := "The following required tools are missing or not executable:\n"
		for _, tool := range missingTools {
			errMsg += fmt.Sprintf("- %s\n", tool)
		}
		errMsg += "\nPlease install the missing tools and try again.\n"
		errMsg += "On Ubuntu/Debian, you can install them with:\n"
		errMsg += "sudo apt-get update && sudo apt-get install -y libreoffice imagemagick ffmpeg p7zip-full unrar"
		
		return fmt.Errorf(errMsg)
	}

	log.Println("All required tools are available and executable")
	return nil
}

// GetBinaryPath returns the full path to a binary in the bin directory
func (bm *BinaryManager) GetBinaryPath(name string) string {
	return filepath.Join(bm.BinDir, name)
}

// GetLibreOfficePath returns the path to the LibreOffice binary
func (tm *ToolManager) GetLibreOfficePath() (string, error) {
	// First try to find soffice in PATH
	if path, err := exec.LookPath("soffice"); err == nil {
		return path, nil
	}
	
	// Try common Linux paths
	commonPaths := []string{
		"/usr/bin/soffice",
		"/usr/local/bin/soffice",
		"/opt/libreoffice/program/soffice",
		"/usr/lib/libreoffice/program/soffice",
	}
	
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	// On macOS, try the standard macOS path
	if runtime.GOOS == "darwin" {
		stdPath := "/Applications/LibreOffice.app/Contents/MacOS/soffice"
		if _, err := os.Stat(stdPath); err == nil {
			return stdPath, nil
		}
	}

	// For other platforms or if not found on macOS, use the bundled binary
	return "", fmt.Errorf("LibreOffice not found. Please install LibreOffice or provide the path to the binary")
}

// GetImageMagickPath returns the path to the ImageMagick convert binary
func (tm *ToolManager) GetImageMagickPath() (string, error) {
	// First try to find convert in PATH
	if path, err := exec.LookPath("convert"); err == nil {
		return path, nil
	}
	
	// Try common paths
	commonPaths := []string{
		"/usr/bin/convert",
		"/usr/local/bin/convert",
		"/opt/ImageMagick/bin/convert",
	}
	
	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}
	
	// Fall back to the bundled binary if available
	bundledPath := tm.binManager.GetBinaryPath("convert")
	if _, err := os.Stat(bundledPath); err == nil {
		return bundledPath, nil
	}
	
	return "", fmt.Errorf("ImageMagick not found. Please install ImageMagick or provide the path to the convert binary")
}

// GetFFmpegPath returns the path to the FFmpeg binary
func (tm *ToolManager) GetFFmpegPath() (string, error) {
	// First try to find ffmpeg in PATH
	if path, err := exec.LookPath("ffmpeg"); err == nil {
		return path, nil
	}

	// Try common paths
	commonPaths := []string{
		"/usr/bin/ffmpeg",
		"/usr/local/bin/ffmpeg",
		"/opt/ffmpeg/bin/ffmpeg",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Fall back to the bundled binary if available
	bundledPath := tm.binManager.GetBinaryPath("ffmpeg")
	if _, err := os.Stat(bundledPath); err == nil {
		return bundledPath, nil
	}

	return "", fmt.Errorf("FFmpeg not found. Please install FFmpeg or provide the path to the ffmpeg binary")
}

// Get7zPath returns the path to the 7z binary
func (tm *ToolManager) Get7zPath() (string, error) {
	// First try to find 7z in PATH
	if path, err := exec.LookPath("7z"); err == nil {
		return path, nil
	}

	// Try common paths
	commonPaths := []string{
		"/usr/bin/7z",
		"/usr/local/bin/7z",
		"/opt/7z/bin/7z",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Fall back to the bundled binary if available
	bundledPath := tm.binManager.GetBinaryPath("7z")
	if _, err := os.Stat(bundledPath); err == nil {
		return bundledPath, nil
	}

	return "", fmt.Errorf("7z not found. Please install p7zip or provide the path to the 7z binary")
}

// GetUnrarPath returns the path to the unrar binary
func (tm *ToolManager) GetUnrarPath() (string, error) {
	// First try to find unrar in PATH
	if path, err := exec.LookPath("unrar"); err == nil {
		return path, nil
	}

	// Try common paths
	commonPaths := []string{
		"/usr/bin/unrar",
		"/usr/local/bin/unrar",
		"/opt/unrar/unrar",
	}

	for _, path := range commonPaths {
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}
	}

	// Fall back to the bundled binary if available
	bundledPath := tm.binManager.GetBinaryPath("unrar")
	if _, err := os.Stat(bundledPath); err == nil {
		return bundledPath, nil
	}

	return "", fmt.Errorf("unrar not found. Please install unrar or provide the path to the unrar binary")
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
				log.Printf("ERROR: Failed to remove temp file %s: %v", path, err)
			}
		}
	}

	return nil
}
