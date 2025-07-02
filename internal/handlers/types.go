package handlers

import (
	"time"
)

// Conversion represents a file conversion
// @Description File conversion details
type Conversion struct {
	ID            string    `json:"id"`
	Status        string    `json:"status"`
	SourceFormat  string    `json:"source_format"`
	TargetFormat  string    `json:"target_format"`
	OriginalName  string    `json:"original_name"`
	ConvertedName string    `json:"converted_name,omitempty"`
	FileSize      int64     `json:"file_size"`
	CreatedAt     time.Time `json:"created_at"`
	CompletedAt   time.Time `json:"completed_at,omitempty"`
	DownloadURL   string    `json:"download_url,omitempty"`
	Error         string    `json:"error,omitempty"`
}
