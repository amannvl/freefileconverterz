package storage

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/amannvl/freefileconverterz/pkg/utils"
)

// FileStorage handles file operations
type FileStorage struct {
	basePath    string
	maxFileSize int64
}

// NewFileStorage creates a new FileStorage instance
func NewFileStorage(basePath string, maxFileSize int64) (*FileStorage, error) {
	// Create base directory if it doesn't exist
	if err := os.MkdirAll(basePath, 0755); err != nil {
		return nil, fmt.Errorf("failed to create storage directory: %w", err)
	}

	return &FileStorage{
		basePath:    basePath,
		maxFileSize: maxFileSize,
	}, nil
}

// SaveFile saves an uploaded file to the storage
func (s *FileStorage) SaveFile(fileHeader *multipart.FileHeader) (string, error) {
	// Validate file size
	if fileHeader.Size > s.maxFileSize {
		return "", utils.ErrFileTooLarge
	}

	// Generate a unique filename
	ext := filepath.Ext(fileHeader.Filename)
	hash := md5.Sum([]byte(fmt.Sprintf("%s%d", fileHeader.Filename, time.Now().UnixNano())))
	filename := hex.EncodeToString(hash[:]) + ext
	filePath := filepath.Join(s.basePath, filename)

	// Create the file
	dst, err := os.Create(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %w", err)
	}
	defer dst.Close()

	// Open the uploaded file
	src, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer src.Close()

	// Copy the file contents
	if _, err = io.Copy(dst, src); err != nil {
		os.Remove(filePath) // Clean up if copy fails
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	return filename, nil
}

// GetFilePath returns the full path to a stored file
func (s *FileStorage) GetFilePath(filename string) (string, error) {
	filePath := filepath.Join(s.basePath, filename)

	// Verify the file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", os.ErrNotExist
	}

	return filePath, nil
}

// DeleteFile removes a file from storage
func (s *FileStorage) DeleteFile(filename string) error {
	filePath := filepath.Join(s.basePath, filename)
	return os.Remove(filePath)
}

// Cleanup removes files older than the specified duration
func (s *FileStorage) Cleanup(olderThan time.Duration) error {
	files, err := os.ReadDir(s.basePath)
	if err != nil {
		return fmt.Errorf("failed to read storage directory: %w", err)
	}

	now := time.Now()
	var firstError error

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		info, err := file.Info()
		if err != nil {
			if firstError == nil {
				firstError = err
			}
			continue
		}

		if now.Sub(info.ModTime()) > olderThan {
			if err := os.Remove(filepath.Join(s.basePath, file.Name())); err != nil && firstError == nil {
				firstError = err
			}
		}
	}

	return firstError
}

// GetFileInfo returns information about a stored file
func (s *FileStorage) GetFileInfo(filename string) (os.FileInfo, error) {
	filePath := filepath.Join(s.basePath, filename)
	return os.Stat(filePath)
}

// FileExists checks if a file exists in storage
func (s *FileStorage) FileExists(filename string) bool {
	filePath := filepath.Join(s.basePath, filename)
	_, err := os.Stat(filePath)
	return !os.IsNotExist(err)
}

// GetFileSize returns the size of a file in storage
func (s *FileStorage) GetFileSize(filename string) (int64, error) {
	filePath := filepath.Join(s.basePath, filename)
	info, err := os.Stat(filePath)
	if err != nil {
		return 0, err
	}
	return info.Size(), nil
}

// GetMimeType returns the MIME type of a file
func (s *FileStorage) GetMimeType(filename string) (string, error) {
	filePath := filepath.Join(s.basePath, filename)
	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	// Only the first 512 bytes are used to sniff the content type
	buffer := make([]byte, 512)
	_, err = file.Read(buffer)
	if err != nil {
		return "", err
	}

	return http.DetectContentType(buffer), nil
}

// MoveFile moves a file within the storage
func (s *FileStorage) MoveFile(oldPath, newPath string) error {
	oldPath = filepath.Join(s.basePath, oldPath)
	newPath = filepath.Join(s.basePath, newPath)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(newPath), 0755); err != nil {
		return fmt.Errorf("failed to create directory: %w", err)
	}

	return os.Rename(oldPath, newPath)
} storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/amannvl/freefileconverterz/internal/config"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
)

// FileInfo contains information about a file in storage
type FileInfo struct {
	ID        string
	Name      string
	Size      int64
	CreatedAt time.Time
	UpdatedAt time.Time
}

// Storage defines the interface for storage operations
type Storage interface {
	// Save saves a file to the storage
	Save(ctx context.Context, path string, file interface{}) (string, error)
	// Get retrieves a file from the storage
	Get(ctx context.Context, path string) (io.ReadCloser, error)
	// Delete removes a file from the storage
	Delete(ctx context.Context, path string) error
	// URL returns the public URL for a file
	URL(path string) string
	// Read reads the entire file into memory
	Read(ctx context.Context, path string) ([]byte, error)
	// ReadStream returns a reader for the file
	ReadStream(ctx context.Context, path string) (io.ReadCloser, error)
	// Exists checks if a file exists
	Exists(ctx context.Context, path string) (bool, error)
	// List lists all files in the storage with the given prefix
	List(ctx context.Context, prefix string) ([]FileInfo, error)
	// GetFullPath returns the full filesystem path for a given relative path
	GetFullPath(path string) string
	// Cleanup removes files older than the specified duration
	Cleanup(ctx context.Context, olderThan time.Duration) error
}

// NewStorage creates a new storage instance based on the configuration
func NewStorage(cfg config.StorageConfig) (Storage, error) {
	switch strings.ToLower(cfg.Provider) {
	case "s3":
		return newS3Storage(cfg)
	case "local":
		fallthrough
	default:
		return newLocalStorage(cfg)
	}
}

// LocalStorage implements Storage for local filesystem
type localStorage struct {
	basePath string
	baseURL  string
}

func newLocalStorage(cfg config.StorageConfig) (Storage, error) {
	// Ensure base directory exists
	if err := os.MkdirAll(cfg.UploadDir, 0755); err != nil {
		return nil, err
	}

	return &localStorage{
		basePath: cfg.UploadDir,
		baseURL:  "", // Use relative URLs for local storage
	}, nil
}

func (s *localStorage) Save(ctx context.Context, path string, file interface{}) (string, error) {
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

	return path, nil
}

func (s *localStorage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

// Delete removes a file from the storage
func (s *localStorage) Delete(ctx context.Context, path string) error {
	fullPath := filepath.Join(s.basePath, path)
	return os.Remove(fullPath)
}

// Read reads the entire file into memory
func (s *localStorage) Read(ctx context.Context, path string) ([]byte, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.ReadFile(fullPath)
}

// ReadStream returns a reader for the file
func (s *localStorage) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	fullPath := filepath.Join(s.basePath, path)
	return os.Open(fullPath)
}

// Exists checks if a file exists
func (s *localStorage) Exists(ctx context.Context, path string) (bool, error) {
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

// GetFullPath returns the full filesystem path for a given relative path
func (s *localStorage) GetFullPath(path string) string {
	return filepath.Join(s.basePath, path)
}

// List lists all files in the storage with the given prefix
func (s *localStorage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	fullPath := s.GetFullPath(prefix)
	
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

// Cleanup removes files older than the specified duration
func (s *localStorage) Cleanup(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)
	
	return filepath.Walk(s.basePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		
		// Skip directories and hidden files
		if info.IsDir() || strings.HasPrefix(filepath.Base(path), ".") {
			return nil
		}
		
		// Delete files older than cutoff
		if info.ModTime().Before(cutoff) {
			return os.Remove(path)
		}
		
		return nil
	})
}

func (s *localStorage) URL(path string) string {
	if s.baseURL != "" {
		return strings.TrimSuffix(s.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")
	}
	return "/uploads/" + strings.TrimPrefix(path, "/")
}

// S3Storage implements Storage for AWS S3
type s3Storage struct {
	client  *s3.Client
	bucket  string
	baseURL string
}

// NewS3Storage creates a new S3 storage instance
func NewS3Storage(client *s3.Client, bucket, baseURL string) Storage {
	return &s3Storage{
		client:  client,
		bucket:  bucket,
		baseURL: baseURL,
	}
}

func newS3Storage(cfg config.StorageConfig) (Storage, error) {
	if cfg.S3AccessKeyID == "" || cfg.S3SecretAccessKey == "" || cfg.S3Bucket == "" {
		return nil, errors.New("missing required S3 configuration")
	}

	// Configure AWS client
	creds := credentials.NewStaticCredentialsProvider(
		cfg.S3AccessKeyID,
		cfg.S3SecretAccessKey,
		"",
	)

	client := s3.New(s3.Options{
		Region:      cfg.S3Region,
		Credentials: aws.NewCredentialsCache(creds),
		EndpointResolver: s3.EndpointResolverFunc(func(region string, options s3.EndpointResolverOptions) (aws.Endpoint, error) {
			if cfg.S3Endpoint != "" {
				return aws.Endpoint{
					URL:           cfg.S3Endpoint,
					SigningRegion: region,
				}, nil
			}
			return aws.Endpoint{}, &aws.EndpointNotFoundError{}
		}),
	})

	return &s3Storage{
		client:  client,
		bucket:  cfg.S3Bucket,
		baseURL: cfg.S3Endpoint,
	}, nil
}

func (s *s3Storage) Save(ctx context.Context, path string, file interface{}) (string, error) {
	// Create a reader from the input file
	var reader io.Reader

	switch f := file.(type) {
	case multipart.File:
		// For multipart files, we need to read the content first
		buf := bytes.NewBuffer(nil)
		if _, err := io.Copy(buf, f); err != nil {
			return "", err
		}
		reader = bytes.NewReader(buf.Bytes())
	case io.Reader:
		// For io.Reader, use it directly
		reader = f
	default:
		return "", errors.New("unsupported file type")
	}

	// Upload the file to S3
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.bucket),
		Key:         aws.String(path),
		Body:        reader,
		ContentType: aws.String(getContentType(path)),
	})

	if err != nil {
		return "", err
	}

	return path, nil
}

func (s *s3Storage) Get(ctx context.Context, path string) (io.ReadCloser, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}

	return result.Body, nil
}

func (s *s3Storage) Delete(ctx context.Context, path string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	return err
}

func (s *s3Storage) URL(path string) string {
	if s.baseURL != "" {
		return strings.TrimSuffix(s.baseURL, "/") + "/" + strings.TrimPrefix(path, "/")
	}
	return "https://" + s.bucket + ".s3.amazonaws.com/" + strings.TrimPrefix(path, "/")
}

// Read reads the entire file into memory
func (s *s3Storage) Read(ctx context.Context, path string) ([]byte, error) {
	result, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})
	if err != nil {
		return nil, err
	}
	defer result.Body.Close()

	return io.ReadAll(result.Body)
}

// ReadStream returns a reader for the file (same as Get)
func (s *s3Storage) ReadStream(ctx context.Context, path string) (io.ReadCloser, error) {
	return s.Get(ctx, path)
}

// Exists checks if a file exists in S3
func (s *s3Storage) Exists(ctx context.Context, path string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(path),
	})

	if err != nil {
		var notFound *types.NotFound
		if errors.As(err, &notFound) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

// GetFullPath returns the full S3 path for a given key
func (s *s3Storage) GetFullPath(path string) string {
	return fmt.Sprintf("s3://%s/%s", s.bucket, strings.TrimPrefix(path, "/"))
}

// List lists all files in the storage with the given prefix
func (s *s3Storage) List(ctx context.Context, prefix string) ([]FileInfo, error) {
	var files []FileInfo
	var continuationToken *string

	for {
		// List objects in the bucket with the given prefix
		input := &s3.ListObjectsV2Input{
			Bucket:            aws.String(s.bucket),
			Prefix:            aws.String(prefix),
			ContinuationToken: continuationToken,
		}

		result, err := s.client.ListObjectsV2(ctx, input)
		if err != nil {
			return nil, fmt.Errorf("failed to list objects in bucket %s: %w", s.bucket, err)
		}

		// Convert S3 objects to FileInfo
		for _, obj := range result.Contents {
			files = append(files, FileInfo{
				ID:        aws.ToString(obj.Key),
				Name:      aws.ToString(obj.Key),
				Size:      aws.ToInt64(obj.Size),
				CreatedAt: aws.ToTime(obj.LastModified),
				UpdatedAt: aws.ToTime(obj.LastModified),
			})
		}

		// Check if there are more objects to fetch
		if result.IsTruncated == nil || !*result.IsTruncated {
			break
		}

		if result.NextContinuationToken == nil {
			break
		}
		continuationToken = result.NextContinuationToken
	}

	return files, nil
}

// Cleanup removes files older than the specified duration
func (s *s3Storage) Cleanup(ctx context.Context, olderThan time.Duration) error {
	cutoff := time.Now().Add(-olderThan)

	// List all objects in the bucket
	paginator := s3.NewListObjectsV2Paginator(s.client, &s3.ListObjectsV2Input{
		Bucket: aws.String(s.bucket),
	})

	// Process each page of results
	for paginator.HasMorePages() {
		page, err := paginator.NextPage(ctx)
		if err != nil {
			return fmt.Errorf("failed to list objects: %w", err)
		}

		// Delete objects older than cutoff
		for _, obj := range page.Contents {
			if obj.LastModified.Before(cutoff) {
				_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
					Bucket: aws.String(s.bucket),
					Key:    obj.Key,
				})
				if err != nil {
					return fmt.Errorf("failed to delete object %s: %w", *obj.Key, err)
				}
			}
		}
	}

	return nil
}

// getContentType determines the content type based on file extension
func getContentType(filename string) string {
	ext := strings.ToLower(filepath.Ext(filename))
	switch ext {
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".png":
		return "image/png"
	case ".gif":
		return "image/gif"
	case ".pdf":
		return "application/pdf"
	case ".doc":
		return "application/msword"
	case ".docx":
		return "application/vnd.openxmlformats-officedocument.wordprocessingml.document"
	case ".xls":
		return "application/vnd.ms-excel"
	case ".xlsx":
		return "application/vnd.openxmlformats-officedocument.spreadsheetml.sheet"
	case ".ppt":
		return "application/vnd.ms-powerpoint"
	case ".pptx":
		return "application/vnd.openxmlformats-officedocument.presentationml.presentation"
	case ".txt":
		return "text/plain"
	default:
		return "application/octet-stream"
	}
}
