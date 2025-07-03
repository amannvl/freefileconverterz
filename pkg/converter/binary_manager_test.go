package converter

import (
	"os"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestBinaryManager(t *testing.T) {
	// Create a temporary directory for test binaries
	tempDir := t.TempDir()
	binDir := filepath.Join(tempDir, "bin")
	os.MkdirAll(filepath.Join(binDir, runtime.GOOS, runtime.GOARCH), 0755)

	tests := []struct {
		name        string
		setup       func()
		binaryName  string
		expectError bool
	}{
		{
			name: "Binary found in bin directory",
			setup: func() {
				testBin := filepath.Join(binDir, runtime.GOOS, runtime.GOARCH, "testbin")
				if runtime.GOOS == "windows" {
					testBin += ".exe"
				}
				f, err := os.Create(testBin)
				require.NoError(t, err)
				f.Close()
				os.Chmod(testBin, 0755)
			},
			binaryName:  "testbin",
			expectError: false,
		},
		{
			name:        "Binary found in PATH",
			setup:       func() {},
			binaryName:  "go", // Should be in PATH
			expectError: false,
		},
		{
			name:        "Binary not found",
			setup:       func() {},
			binaryName:  "nonexistentbinary123",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup test case
			if tt.setup != nil {
				tt.setup()
			}

			// Create binary manager
			bm := NewBinaryManager(binDir)

			// Test GetBinaryPath
			path, err := bm.GetBinaryPath(tt.binaryName)
			if tt.expectError {
				assert.Error(t, err)
				assert.Empty(t, path)
			} else {
				assert.NoError(t, err)
				assert.NotEmpty(t, path)

				// Test IsAvailable
				assert.True(t, bm.IsAvailable(tt.binaryName))

				// Test Execute (basic test, just check it doesn't panic)
				if tt.binaryName == "go" {
					// Only run this test for the 'go' binary which we know exists
					output, err := bm.Execute(tt.binaryName, "version")
					assert.NoError(t, err)
					assert.Contains(t, string(output), "go version")
				}
			}

			// Test IsAvailable for non-existent binary
			assert.False(t, bm.IsAvailable("nonexistentbinary123"))
		})
	}
}
