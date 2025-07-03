package converter

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// BinaryManager handles the management of external binaries required for conversions
type BinaryManager struct {
	binDir string
}

// NewBinaryManager creates a new instance of BinaryManager
func NewBinaryManager(binDir string) *BinaryManager {
	return &BinaryManager{
		binDir: binDir,
	}
}

// GetBinaryPath returns the full path to a binary for the current platform
func (bm *BinaryManager) GetBinaryPath(name string) (string, error) {
	// Construct the platform-specific binary name
	binaryName := name
	if runtime.GOOS == "windows" {
		binaryName = name + ".exe"
	}

	// Check if the binary exists in the bin directory
	binaryPath := filepath.Join(bm.binDir, runtime.GOOS, runtime.GOARCH, binaryName)
	if _, err := os.Stat(binaryPath); os.IsNotExist(err) {
		// Binary not found, try to find it in PATH
		path, err := exec.LookPath(name)
		if err != nil {
			return "", fmt.Errorf("binary %s not found in %s or in PATH", name, binaryPath)
		}
		return path, nil
	}

	return binaryPath, nil
}

// IsAvailable checks if a binary is available
func (bm *BinaryManager) IsAvailable(name string) bool {
	_, err := bm.GetBinaryPath(name)
	return err == nil
}

// Execute runs a binary with the given arguments
func (bm *BinaryManager) Execute(name string, args ...string) ([]byte, error) {
	path, err := bm.GetBinaryPath(name)
	if err != nil {
		return nil, err
	}

	cmd := exec.Command(path, args...)
	return cmd.CombinedOutput()
}
