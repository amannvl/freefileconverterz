package utils

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"mime"
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
)

// LogLevel represents the log level
type LogLevel int

const (
	// DebugLevel logs everything
	DebugLevel LogLevel = iota
	// InfoLevel logs important information
	InfoLevel
	// ErrorLevel logs only errors
	ErrorLevel
)

// Logger is a simple logger interface
type Logger interface {
	Debug(msg string, keyvals ...interface{})
	Info(msg string, keyvals ...interface{})
	Error(msg string, keyvals ...interface{})
	Fatal(msg string, keyvals ...interface{})
}

// zerologLogger wraps zerolog.Logger to implement our Logger interface
type zerologLogger struct {
	logger zerolog.Logger
}

// NewLogger creates a new logger with the specified level
func NewLogger(level LogLevel) Logger {
	zerolog.TimeFieldFormat = time.RFC3339Nano
	zerolog.SetGlobalLevel(zerolog.InfoLevel)

	var zl zerolog.Logger

	switch level {
	case DebugLevel:
		zerolog.SetGlobalLevel(zerolog.DebugLevel)
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			With().
			Timestamp().
			Caller().
			Logger()
	case InfoLevel:
		zerolog.SetGlobalLevel(zerolog.InfoLevel)
		zl = zerolog.New(zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: "15:04:05"}).
			With().
			Timestamp().
			Logger()
	case ErrorLevel:
		zerolog.SetGlobalLevel(zerolog.ErrorLevel)
		zl = zerolog.New(os.Stderr).
			With().
			Timestamp().
			Logger()
	}

	return &zerologLogger{logger: zl}
}

// Debug logs a debug message
func (l *zerologLogger) Debug(msg string, keyvals ...interface{}) {
	e := l.logger.Debug()
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			e = e.Interface(fmt.Sprint(keyvals[i]), keyvals[i+1])
		}
	}
	e.Msg(msg)
}

// Info logs an info message
func (l *zerologLogger) Info(msg string, keyvals ...interface{}) {
	e := l.logger.Info()
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			e = e.Interface(fmt.Sprint(keyvals[i]), keyvals[i+1])
		}
	}
	e.Msg(msg)
}

// Error logs an error message
func (l *zerologLogger) Error(msg string, keyvals ...interface{}) {
	e := l.logger.Error()
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			e = e.Interface(fmt.Sprint(keyvals[i]), keyvals[i+1])
		}
	}
	e.Msg(msg)
}

// Fatal logs a fatal message and exits
func (l *zerologLogger) Fatal(msg string, keyvals ...interface{}) {
	e := l.logger.Fatal()
	for i := 0; i < len(keyvals); i += 2 {
		if i+1 < len(keyvals) {
			e = e.Interface(fmt.Sprint(keyvals[i]), keyvals[i+1])
		}
	}
	e.Msg(msg)
}

// FileInfo contains information about a file
type FileInfo struct {
	Name     string
	Path     string
	Size     int64
	MimeType string
	Hash     string
}

// GetFileInfo returns information about a file
func GetFileInfo(filePath string) (*FileInfo, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		return nil, fmt.Errorf("failed to get file stats: %w", err)
	}

	hash := md5.New()
	if _, err := io.Copy(hash, file); err != nil {
		return nil, fmt.Errorf("failed to calculate file hash: %w", err)
	}

	mimeType := mime.TypeByExtension(filepath.Ext(filePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}

	return &FileInfo{
		Name:     filepath.Base(filePath),
		Path:     filePath,
		Size:     stat.Size(),
		MimeType: mimeType,
		Hash:     hex.EncodeToString(hash.Sum(nil)),
	}, nil
}
