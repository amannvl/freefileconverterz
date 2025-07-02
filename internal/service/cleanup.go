package service

import (
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

// CleanupService handles periodic cleanup of old files
type CleanupService struct {
	logger      *slog.Logger
	uploadDir   string
	maxFileAge  time.Duration
	ticker      *time.Ticker
	done        chan bool
}

// NewCleanupService creates a new CleanupService
func NewCleanupService(logger *slog.Logger, uploadDir string, cleanupInterval, maxFileAge time.Duration) *CleanupService {
	return &CleanupService{
		logger:     logger,
		uploadDir:  uploadDir,
		maxFileAge: maxFileAge,
		ticker:     time.NewTicker(cleanupInterval),
		done:       make(chan bool),
	}
}

// Start begins the cleanup process
func (s *CleanupService) Start() {
	go func() {
		s.logger.Info("Starting cleanup service", "interval", s.ticker, "maxFileAge", s.maxFileAge)
		s.cleanup() // Run immediately on start

		for {
			select {
			case <-s.ticker.C:
				s.cleanup()
			case <-s.done:
				s.ticker.Stop()
				return
			}
		}
	}()
}

// Stop stops the cleanup service
func (s *CleanupService) Stop() {
	s.logger.Info("Stopping cleanup service")
	s.done <- true
}

// cleanup removes files older than maxFileAge
func (s *CleanupService) cleanup() {
	startTime := time.Now()
	s.logger.Info("Starting file cleanup")

	files, err := os.ReadDir(s.uploadDir)
	if err != nil {
		s.logger.Error("Failed to read upload directory", "error", err)
		return
	}

	now := time.Now()
	deletedCount := 0
	errCount := 0

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		filePath := filepath.Join(s.uploadDir, file.Name())
		fileInfo, err := file.Info()
		if err != nil {
			s.logger.Error("Failed to get file info", "file", file.Name(), "error", err)
			errCount++
			continue
		}

		// Skip files that are too new
		if now.Sub(fileInfo.ModTime()) < s.maxFileAge {
			continue
		}

		// Delete the file
		if err := os.Remove(filePath); err != nil {
			s.logger.Error("Failed to delete file", "file", file.Name(), "error", err)
			errCount++
		} else {
			deletedCount++
		}
	}

	duration := time.Since(startTime)
	s.logger.Info("File cleanup completed",
		"deleted", deletedCount,
		"errors", errCount,
		"duration", duration,
	)
}
