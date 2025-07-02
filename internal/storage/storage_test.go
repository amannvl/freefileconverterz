package storage

import (
	"bytes"
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

		// Test GetFile
		file, err := storage.GetFile("test.txt")
		assert.NoError(t, err)
		defer file.Close()

		content, err := ioutil.ReadAll(file)
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

		// Verify files were deleted
		exists := storage.FileExists("moved/test.txt")
		assert.False(t, exists)
	})
}

func TestFileStorage_EdgeCases(t *testing.T) {
	// Setup
	tempDir, err := ioutil.TempDir("", "storage-test-edge")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	logger := utils.NewLogger(utils.InfoLevel)
	storage, err := NewFileStorage(tempDir, 1024, logger) // 1KB max size
	require.NoError(t, err)

	t.Run("File Not Found", func(t *testing.T) {
		_, err := storage.GetFile("nonexistent.txt")
		assert.Error(t, err)
		assert.True(t, os.IsNotExist(err))

		size, err := storage.GetFileSize("nonexistent.txt")
		assert.Error(t, err)
		assert.Zero(t, size)

		mimeType, err := storage.GetMimeType("nonexistent.txt")
		assert.Error(t, err)
		assert.Empty(t, mimeType)
	})

	t.Run("Invalid File Operations", func(t *testing.T) {
		err := storage.MoveFile("nonexistent.txt", "new.txt")
		assert.Error(t, err)

		err = storage.DeleteFile("nonexistent.txt")
		assert.NoError(t, err) // Deleting non-existent file should not error
	})

	t.Run("File Size Limit", func(t *testing.T) {
		// Create a file larger than the limit
		largeContent := bytes.Repeat([]byte("a"), 1025) // 1KB + 1 byte
		largeFile := filepath.Join(tempDir, "large.txt")
		err := ioutil.WriteFile(largeFile, largeContent, 0644)
		require.NoError(t, err)

		// Test SaveFile with a file larger than the limit
		_, err = storage.SaveFile(&mockFileHeader{
			filename: "large.txt",
			size:     int64(len(largeContent)),
			content:  largeContent,
		})
		assert.Error(t, err)
	})
}

// mockFileHeader implements a simple multipart.FileHeader for testing
type mockFileHeader struct {
	filename string
	size    int64
	content []byte
}

func (m *mockFileHeader) Open() (multipart.File, error) {
	return ioutil.NopCloser(bytes.NewReader(m.content)), nil
}

func (m *mockFileHeader) Filename() string {
	return m.filename
}

func (m *mockFileHeader) Size() int64 {
	return m.size
}
