package models

import (
	"time"
)

// ConversionStatus represents the status of a conversion task
type ConversionStatus string

const (
	StatusPending    ConversionStatus = "pending"
	StatusProcessing ConversionStatus = "processing"
	StatusCompleted  ConversionStatus = "completed"
	StatusFailed     ConversionStatus = "failed"
)

// Conversion represents a file conversion task
type Conversion struct {
	ID           string          `json:"id"`
	UserID       string          `json:"user_id"`
	SourceFileID string          `json:"source_file_id"`
	Status       ConversionStatus `json:"status"`
	SourceFormat string          `json:"source_format"`
	TargetFormat string          `json:"target_format"`
	Options      map[string]interface{} `json:"options"`
	Error        string          `json:"error,omitempty"`
	DownloadURL  string          `json:"download_url,omitempty"`
	CreatedAt    time.Time       `json:"created_at"`
	UpdatedAt    time.Time       `json:"updated_at"`
}

// ConversionRequest represents a conversion request
type ConversionRequest struct {
	SourceFileID string                 `json:"source_file_id" validate:"required"`
	TargetFormat string                 `json:"target_format" validate:"required"`
	Options      map[string]interface{} `json:"options"`
}

// ConversionResponse represents the response for a conversion request
type ConversionResponse struct {
	ID           string    `json:"id"`
	Status       string    `json:"status"`
	DownloadURL  string    `json:"download_url,omitempty"`
	Error        string    `json:"error,omitempty"`
	QueuedAt     time.Time `json:"queued_at"`
	CompletedAt  *time.Time `json:"completed_at,omitempty"`
}
