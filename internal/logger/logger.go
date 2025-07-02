package logger

import (
	"io"
	"log/slog"
	"os"
	"path/filepath"
)

// SetupLogger initializes and returns a new logger handler based on the provided configuration
func SetupLogger(level, format, logFile string) slog.Handler {
	// Set log level
	var logLevel slog.Level
	switch level {
	case "debug":
		logLevel = slog.LevelDebug
	case "info":
		logLevel = slog.LevelInfo
	case "warn":
		logLevel = slog.LevelWarn
	case "error":
		logLevel = slog.LevelError
	default:
		logLevel = slog.LevelInfo
	}

	// Create output writer
	var output io.Writer = os.Stdout
	if logFile != "" {
		// Ensure directory exists
		dir := filepath.Dir(logFile)
		if err := os.MkdirAll(dir, 0755); err != nil {
			slog.Error("Failed to create log directory, using stdout", "error", err)
		} else {
			file, err := os.OpenFile(logFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
			if err != nil {
				slog.Error("Failed to open log file, using stdout", "error", err, "file", logFile)
			} else {
				output = file
			}
		}
	}

	// Create handler based on format
	var handler slog.Handler
	switch format {
	case "json":
		handler = slog.NewJSONHandler(output, &slog.HandlerOptions{
			Level: logLevel,
		})
	case "text":
		fallthrough
	default:
		handler = slog.NewTextHandler(output, &slog.HandlerOptions{
			Level: logLevel,
		})
	}

	return handler
}

// NewLogger creates a new logger with the given configuration
func NewLogger(level, format, logFile string) *slog.Logger {
	return slog.New(SetupLogger(level, format, logFile))
}

// WithRequestID adds a request ID to the logger
func WithRequestID(logger *slog.Logger, requestID string) *slog.Logger {
	return logger.With("request_id", requestID)
}

// WithError adds an error to the logger
func WithError(logger *slog.Logger, err error) *slog.Logger {
	return logger.With("error", err)
}

// With adds key-value pairs to the logger
func With(logger *slog.Logger, args ...interface{}) *slog.Logger {
	return logger.With(args...)
}
