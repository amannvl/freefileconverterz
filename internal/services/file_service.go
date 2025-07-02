package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/amannvl/freefileconverterz/internal/storage"
)

// FileService handles file operations
type FileService struct {
	storage storage.Storage
	config *config.StorageConfig
}

// NewFileService creates a new FileService instance
func NewFileService(storage storage.Storage, cfg *config.StorageConfig) *FileService {
	return &FileService{
		storage: storage,
		config:  cfg,
	}
}

// UploadFile uploads a file to the storage
func (s *FileService) UploadFile(ctx context.Context, fileHeader *multipart.FileHeader, userID string) (string, error) {
	// Validate file size
	if fileHeader.Size > s.config.MaxUploadSize {
		return "", errors.New("file size exceeds maximum allowed size")
	}

	// Validate file type
	if !s.isFileTypeAllowed(fileHeader.Filename) {
		return "", errors.New("file type not allowed")
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	filename := generateUniqueFilename(userID, ext)

	// Save the file
	path, err := s.storage.Save(ctx, filename, file)
	if err != nil {
		return "", err
	}

	return path, nil
}

// isFileTypeAllowed checks if the file type is allowed
func (s *FileService) isFileTypeAllowed(filename string) bool {	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) == 0 {
		return false
	}

	// Remove the leading dot
	ext = ext[1:]

	for _, allowedType := range s.config.AllowedFileTypes {
		// Handle wildcard patterns like "image/*"
		if strings.Contains(allowedType, "/*") {
			prefix := strings.TrimSuffix(allowedType, "/*")
			if strings.HasPrefix(ext, prefix) {
				return true
			}
		}

		// Handle exact matches
		if strings.EqualFold(ext, allowedType) {
			return true
		}

		// Handle extensions with leading dot
		if strings.HasPrefix(allowedType, ".") && strings.EqualFold(ext, allowedType[1:]) {
			return true
		}
	}

	return false
}

// generateUniqueFilename generates a unique filename
func generateUniqueFilename(userID, ext string) string {
	timestamp := time.Now().UnixNano()
	return filepath.Join(userID, "files", fmt.Sprintf("%d%s", timestamp, ext))
}

// GetFile retrieves a file from storage
func (s *FileService) GetFile(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.storage.Get(ctx, path)
}

// DeleteFile removes a file from storage
func (s *FileService) DeleteFile(ctx context.Context, path string) error {
	return s.storage.Delete(ctx, path)
}

// GetFileURL returns the public URL for a file
func (s *FileService) GetFileURL(path string) string {
	return s.storage.URL(path)
}
