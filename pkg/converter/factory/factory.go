package factory

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/archive"
	"github.com/amannvl/freefileconverterz/pkg/converter/audio"
	"github.com/amannvl/freefileconverterz/pkg/converter/document"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/amannvl/freefileconverterz/pkg/converter/image"
	"github.com/amannvl/freefileconverterz/pkg/converter/video"
)

// ConverterType represents the type of converter
type ConverterType string

const (
	// DocumentConverterType handles document format conversions
	DocumentConverterType ConverterType = "document"
	// ImageConverterType handles image format conversions
	ImageConverterType ConverterType = "image"
	// AudioConverterType handles audio format conversions
	AudioConverterType ConverterType = "audio"
	// VideoConverterType handles video format conversions
	VideoConverterType ConverterType = "video"
	// ArchiveConverterType handles archive format conversions
	ArchiveConverterType ConverterType = "archive"
)

// ConverterFactory creates and manages converters
type ConverterFactory struct {
	tempDir     string
	toolManager *tools.ToolManager
}

// NewConverterFactory creates a new ConverterFactory
func NewConverterFactory(tempDir string, toolManager *tools.ToolManager) (*ConverterFactory, error) {
	// Create the temp directory if it doesn't exist
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	return &ConverterFactory{
		tempDir:     tempDir,
		toolManager: toolManager,
	}, nil
}

// GetConverter returns the appropriate converter for the given source and target formats
func (f *ConverterFactory) GetConverter(sourceFormat, targetFormat string) (iface.Converter, error) {
	// Determine the converter type based on the source format
	converterType := f.determineConverterType(sourceFormat)

	// Create the appropriate converter
	switch converterType {
	case DocumentConverterType:
		docConv := document.NewDocumentConverter(f.toolManager, f.tempDir)
		if docConv.SupportsConversion(sourceFormat, targetFormat) {
			return docConv, nil
		}
	case ImageConverterType:
		imgConv := image.NewImageConverter(f.toolManager, f.tempDir)
		if imgConv.SupportsConversion(sourceFormat, targetFormat) {
			return imgConv, nil
		}
	case AudioConverterType:
		audioConv := audio.NewAudioConverter(f.toolManager, f.tempDir)
		if audioConv.SupportsConversion(sourceFormat, targetFormat) {
			return audioConv, nil
		}
	case VideoConverterType:
		videoConv := video.NewVideoConverter(f.toolManager, f.tempDir)
		if videoConv.SupportsConversion(sourceFormat, targetFormat) {
			return videoConv, nil
		}
	case ArchiveConverterType:
		archiveConv := archive.NewArchiveConverter(f.toolManager, f.tempDir)
		if archiveConv.SupportsConversion(sourceFormat, targetFormat) {
			return archiveConv, nil
		}
	}

	return nil, fmt.Errorf("no converter found for %s to %s conversion", sourceFormat, targetFormat)
}

// determineConverterType determines the converter type based on the file extension
func (f *ConverterFactory) determineConverterType(filename string) ConverterType {
	ext := strings.ToLower(strings.TrimPrefix(filepath.Ext(filename), "."))

	// Document formats
	docFormats := map[string]bool{
		"pdf": true, "doc": true, "docx": true, "odt": true, "rtf": true,
		"txt": true, "odp": true, "ods": true, "xls": true, "xlsx": true,
		"ppt": true, "pptx": true, "epub": true, "mobi": true, "azw": true,
		"fb2": true, "lit": true,
	}

	if docFormats[ext] {
		return DocumentConverterType
	}

	// Image formats
	imgFormats := map[string]bool{
		"jpg": true, "jpeg": true, "png": true, "gif": true, "bmp": true,
		"tiff": true, "webp": true, "heic": true, "heif": true,
	}

	if imgFormats[ext] {
		return ImageConverterType
	}

	// Audio formats
	audioFormats := map[string]bool{
		"mp3": true, "wav": true, "aac": true, "flac": true, "ogg": true,
		"m4a": true, "wma": true,
	}

	if audioFormats[ext] {
		return AudioConverterType
	}

	// Video formats
	videoFormats := map[string]bool{
		"mp4": true, "avi": true, "mov": true, "mkv": true, "wmv": true,
		"flv": true, "webm": true, "3gp": true, "m4v": true, "mpg": true,
		"mpeg": true,
	}

	if videoFormats[ext] {
		return VideoConverterType
	}

	// Archive formats
	archiveFormats := map[string]bool{
		"zip": true, "tar": true, "gz": true, "bz2": true, "xz": true,
		"7z": true, "rar": true, "tgz": true, "tbz2": true, "txz": true,
	}

	if archiveFormats[ext] {
		return ArchiveConverterType
	}

	// Default to document converter for unknown formats
	return DocumentConverterType
}

// Cleanup removes all temporary files created by the factory
func (f *ConverterFactory) Cleanup() error {
	return os.RemoveAll(f.tempDir)
}

// CreateTempFile creates a temporary file in the factory's temp directory
func (f *ConverterFactory) CreateTempFile(prefix, suffix string) (*os.File, error) {
	return os.CreateTemp(f.tempDir, prefix+"*"+suffix)
}

// CreateTempDir creates a temporary directory in the factory's temp directory
func (f *ConverterFactory) CreateTempDir(prefix string) (string, error) {
	return os.MkdirTemp(f.tempDir, prefix+"*")
}
