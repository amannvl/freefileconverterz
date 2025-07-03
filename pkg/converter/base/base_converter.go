package base

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/rs/zerolog/log"
)

// BaseConverter provides common functionality for all converters
type BaseConverter struct {
	tempDir     string
	toolManager *tools.ToolManager
	supportedFormats map[string][]string // map[sourceFormat][]targetFormat
}

// NewBaseConverter creates a new BaseConverter
func NewBaseConverter(toolManager *tools.ToolManager, tempDir string) *BaseConverter {
	return &BaseConverter{
		tempDir:     tempDir,
		toolManager: toolManager,
		supportedFormats: make(map[string][]string),
	}
}

// SupportsConversion checks if the converter supports the given conversion
func (c *BaseConverter) SupportsConversion(sourceFormat, targetFormat string) bool {
	if targets, exists := c.supportedFormats[sourceFormat]; exists {
		for _, t := range targets {
			if t == targetFormat {
				return true
			}
		}
	}
	return false
}

// AddSupportedConversion adds a supported conversion
func (c *BaseConverter) AddSupportedConversion(sourceFormat string, targetFormats ...string) {
	c.supportedFormats[sourceFormat] = append(c.supportedFormats[sourceFormat], targetFormats...)
}

// CreateTempFile creates a temporary file in the converter's temp directory
func (c *BaseConverter) CreateTempFile(prefix, suffix string) (*os.File, error) {
	return os.CreateTemp(c.tempDir, prefix+"*"+suffix)
}

// CreateTempDir creates a temporary directory in the converter's temp directory
func (c *BaseConverter) CreateTempDir(prefix string) (string, error) {
	tempDir, err := os.MkdirTemp(c.tempDir, prefix+"*")
	if err != nil {
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	return tempDir, nil
}

// Cleanup removes temporary files and directories
func (c *BaseConverter) Cleanup(files ...string) error {
	for _, file := range files {
		if file == "" {
			continue
		}
		if err := os.RemoveAll(file); err != nil {
			log.Warn().Err(err).Str("file", file).Msg("Failed to remove temporary file")
			return err
		}
	}
	return nil
}

// GetOutputFilename generates an output filename with the target extension
// It preserves only the original filename without any hashes or timestamps
func (c *BaseConverter) GetOutputFilename(inputPath, targetExt string) string {
	// Get the base filename without path
	base := filepath.Base(inputPath)
	
	// Remove the extension if it exists
	ext := filepath.Ext(base)
	if ext != "" {
		base = base[:len(base)-len(ext)]
	}
	
	// Remove any hash or timestamp that might have been added
	// This assumes the original filename is before the first underscore
	if idx := filepath.Ext(base); idx != "" {
		base = base[:len(base)-len(idx)]
	}
	
	// Add the target extension
	return base + "." + targetExt
}

// Convert is a placeholder that should be implemented by specific converters
func (c *BaseConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	return fmt.Errorf("convert method not implemented")
}
