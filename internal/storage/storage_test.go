package storage

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/amannvl/freefileconverterz/pkg/utils"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFileStorage(t *testing.T) {
	// Setup
	tempDir, err := ioutil.TempDir("", "storage-test")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := utils.NewLogger(utils.InfoLevel)
	storage, err := NewFileStorage(tempDir, 1024*1024, logger)
	require.NoError(t, err)

	t.Run("Save and Read File", func(t *testing.T) {
		// Create a test file
		testContent := []byte("test content")
		testFile := filepath.Join(tempDir, "test.txt")
		err := ioutil.WriteFile(testFile, testContent, 0644)
		require.NoError(t, err)

		// Test GetFilePath and reading the file
		filePath, err := storage.GetFilePath("test.txt")
		assert.NoError(t, err)

		content, err := ioutil.ReadFile(filePath)
		assert.NoError(t, err)
		assert.Equal(t, testContent, content)
	})

	t.Run("File Operations", func(t *testing.T) {
		// Test FileExists
		exists := storage.FileExists("test.txt")
		assert.True(t, exists)

		// Test GetFileSize
		size, err := storage.GetFileSize("test.txt")
		assert.NoError(t, err)
		assert.Greater(t, size, int64(0))

		// Test GetMimeType
		mimeType, err := storage.GetMimeType("test.txt")
		assert.NoError(t, err)
		assert.Contains(t, mimeType, "text/plain")
	})

	t.Run("Move File", func(t *testing.T) {
		err := storage.MoveFile("test.txt", "moved/test.txt")
		assert.NoError(t, err)

		// Verify the file was moved
		exists := storage.FileExists("moved/test.txt")
		assert.True(t, exists)

		exists = storage.FileExists("test.txt")
		assert.False(t, exists)
	})

	t.Run("Cleanup", func(t *testing.T) {
		err := storage.Cleanup(0) // Cleanup all files
		assert.NoError(t, err)

	})
}

// Test helper functions
func createTestFile(path string, size int) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	// Write some data to the file
	data := make([]byte, size)
	for i := 0; i < size; i++ {
		data[i] = byte(i % 256)
	}

	_, err = file.Write(data)
	return err
}
