package handlers

import (
	"crypto/rand"
	"encoding/hex"
	"path/filepath"
	"strings"
)

// generateID creates a new unique ID
func generateID() string {
	b := make([]byte, 16)
	_, _ = rand.Read(b)
	return hex.EncodeToString(b)
}

// getFileExtension returns the file extension without the dot
func getFileExtension(filename string) string {
	ext := filepath.Ext(filename)
	if ext == "" {
		return ""
	}
	return strings.TrimPrefix(ext, ".")
}
