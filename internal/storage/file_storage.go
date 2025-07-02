package storage

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"
)

// FileStorage implements file storage operations for the converter service
type FileStorage struct {
	basePath string
}

// NewFileStorage creates a new FileStorage instance that implements the Storage interface
func NewFileStorage(basePath string) (Storage, error) {
	// Create the base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, err
	}

	return &FileStorage{
		basePath: basePath,
	}, nil
}

// Save saves a file to the storage
func (s *FileStorage) Save(ctx context.Context, path string, file interface{}) (string, error) {
	fullPath := filepath.Join(s.basePath, path)

	// Create the directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(fullPath), 0755); err != nil {
		return "", err
	}

	// Handle different file types
	switch f := file.(type) {
	case []byte:
		// Write bytes directly
		if err := os.WriteFile(fullPath, f, 0644); err != nil {
			return "", err
		}
	case io.Reader:
		// Copy from reader to file
		dst, err := os.Create(fullPath)
		if err != nil {
			return "", err
		}
		defer dst.Close()

		if _, err := io.Copy(dst, f); err != nil {
			os.Remove(fullPath) // Clean up if copy fails
			return "", err
		}
	case *multipart.FileHeader:
		// Handle multipart file uploads
		src, err := f.Open()
		if err != nil {
			return "", err
		}
		defer src.Close()

		dst, err := os.Create(fullPath)
		if err != nil {
			return "", err
		}
		defer dst.Close()

		if _, err = io.Copy(dst, src); err != nil {
			os.Remove(fullPath) // Clean up if copy fails
			return "", err
		}
	default:
		return "", errors.New("unsupported file type")
	}

	return fullPath, nil
}

// Get retrieves a file from the storage
func (s *FileStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

// Read reads a file from the storage and returns its content as bytes
func (s *FileStorage) Read(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

// ReadStream returns a reader for the file
func (s *FileStorage) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

// Delete deletes a file from the storage
func (s *FileStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// Exists checks if a file exists in the storage
func (s *FileStorage) Exists(ctx context.Context, path string) (bool, error) {
	fullPath := filepath.Join(s.basePath, path)
	_, err := os.Stat(fullPath)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// GetFullPath returns the full filesystem path for a given storage path
func (s *FileStorage) GetFullPath(path string) string {
	return filepath.Join(s.basePath, path)
}

// CreateTempFile creates a temporary file in the storage
func (s *FileStorage) CreateTempFile(prefix, suffix string) (*os.File, error) {
	tempDir := filepath.Join(s.basePath, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, err
	}
	return os.CreateTemp(tempDir, prefix+"*"+suffix)
}

// CreateTempDir creates a temporary directory in the storage
func (s *FileStorage) CreateTempDir(prefix string) (string, error) {
	tempDir := filepath.Join(s.basePath, "temp")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return "", err
	}
	return os.MkdirTemp(tempDir, prefix+"*")
}

// Cleanup removes files older than the specified duration
func (s *FileStorage) Cleanup(ctx context.Context, olderThan time.Duration) error {
	tempDir := filepath.Join(s.basePath, "temp")
	
	// If temp directory doesn't exist, nothing to clean up
	if _, err := os.Stat(tempDir); os.IsNotExist(err) {
		return nil
	}

	// Walk through the temp directory
	return filepath.Walk(tempDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		// Skip the temp directory itself and any subdirectories
		if path == tempDir || info.IsDir() {
			return nil
		}

		// Skip files that are not older than the specified duration
		if time.Since(info.ModTime()) < olderThan {
			return nil
		}

		// Delete the file
		return os.Remove(path)
	})
}

// URL returns the URL for the given path
func (s *FileStorage) URL(path string) string {
	// In a real implementation, this would return a full URL to the file
	// For local storage, we just return the path
	return path
}

// List lists all files in the storage with the given prefix
func (s *FileStorage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	fullPath := filepath.Join(s.basePath, prefix)
	
	// Read the directory
	entries, err := os.ReadDir(fullPath)
	if err != nil {
		if os.IsNotExist(err) {
			return []FileInfo{}, nil
		}
		return nil, fmt.Errorf("failed to read directory %s: %w", fullPath, err)
	}

	var files []FileInfo
	for _, entry := range entries {
		if entry.IsDir() {
			continue
		}

		info, err := entry.Info()
		if err != nil {
			continue
		}

		files = append(files, FileInfo{
			ID:        entry.Name(),
			Name:      entry.Name(),
			Size:      info.Size(),
			CreatedAt: info.ModTime(),
			UpdatedAt: info.ModTime(),
		})
	}

	return files, nil
}
