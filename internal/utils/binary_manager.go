package utils

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
)

// BinaryManager manages binary dependencies
type BinaryManager struct {
	BinDir string
}

// NewBinaryManager creates a new BinaryManager
func NewBinaryManager(binDir string) (*BinaryManager, error) {
	// Create the bin directory if it doesn't exist
	if err := os.MkdirAll(binDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create bin directory: %w", err)
	}

	return &BinaryManager{
		BinDir: binDir,
	}, nil
}

// EnsureBinary downloads and verifies a binary if it doesn't exist or is invalid
func (bm *BinaryManager) EnsureBinary(ctx context.Context, name, url, expectedChecksum string) (string, error) {
	binaryPath := filepath.Join(bm.BinDir, name)

	// Check if the binary exists and is valid
	if bm.isBinaryValid(binaryPath, expectedChecksum) {
		return binaryPath, nil
	}

	// Download the binary
	log.Info().Str("url", url).Msg("Downloading binary")
	if err := bm.downloadFile(ctx, url, binaryPath); err != nil {
		return "", fmt.Errorf("failed to download binary: %w", err)
	}

	// Make the binary executable
	if err := os.Chmod(binaryPath, 0755); err != nil {
		return "", fmt.Errorf("failed to make binary executable: %w", err)
	}

	// Verify the checksum
	if !bm.isBinaryValid(binaryPath, expectedChecksum) {
		return "", fmt.Errorf("checksum verification failed for %s", name)
	}

	return binaryPath, nil
}

// isBinaryValid checks if the binary exists and has the expected checksum
func (bm *BinaryManager) isBinaryValid(path, expectedChecksum string) bool {
	// Check if file exists
	info, err := os.Stat(path)
	if os.IsNotExist(err) || info.IsDir() {
		return false
	}

	// If no checksum is provided, just check if the file exists
	if expectedChecksum == "" {
		return true
	}

	// Calculate file checksum
	file, err := os.Open(path)
	if err != nil {
		return false
	}
	defer file.Close()

	hasher := sha256.New()
	if _, err := io.Copy(hasher, file); err != nil {
		return false
	}

	actualChecksum := hex.EncodeToString(hasher.Sum(nil))
	return strings.EqualFold(actualChecksum, expectedChecksum)
}

// downloadFile downloads a file from a URL to a local path
func (bm *BinaryManager) downloadFile(ctx context.Context, url, destPath string) error {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return err
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Create the destination file
	out, err := os.Create(destPath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the response body to the file
	_, err = io.Copy(out, resp.Body)
	return err
}

// GetBinaryPath returns the full path to a binary in the bin directory
func (bm *BinaryManager) GetBinaryPath(name string) string {
	return filepath.Join(bm.BinDir, name)
}

// Command creates a new command with the binary
func (bm *BinaryManager) Command(ctx context.Context, name string, args ...string) *exec.Cmd {
	binaryPath := bm.GetBinaryPath(name)
	cmd := exec.CommandContext(ctx, binaryPath, args...)
	// Set the PATH to include the bin directory
	cmd.Env = append(os.Environ(), fmt.Sprintf("PATH=%s:%s", bm.BinDir, os.Getenv("PATH")))
	return cmd
}
