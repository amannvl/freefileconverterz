package models

import (
	"mime/multipart"
	"time"
)

// File represents a file in the system
type File struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	Name        string    `json:"name"`
	Size        int64     `json:"size"`
	ContentType string    `json:"content_type"`
	Path        string    `json:"path"`
	Hash        string    `json:"hash"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// FileUpload represents a file upload request
type FileUpload struct {
	File       *multipart.FileHeader `form:"file" validate:"required"`
	TargetType string               `form:"target_type" validate:"required,oneof=pdf docx jpg png mp3 mp4 zip"`
	Options    map[string]interface{} `form:"options"`
}
