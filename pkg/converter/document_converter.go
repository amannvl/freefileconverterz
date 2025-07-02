package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
)

// DocumentConverter handles document format conversions using LibreOffice
type DocumentConverter struct {
	binMgr *BinaryManager
}

// NewDocumentConverter creates a new DocumentConverter
func NewDocumentConverter() *DocumentConverter {
	return &DocumentConverter{
		binMgr: NewBinaryManager(),
	}
}

// SupportedFormats returns the supported document formats and conversions
func (dc *DocumentConverter) SupportedFormats() map[string][]string {
	return map[string][]string{
		// Text documents
		"doc":  {"pdf", "docx", "odt", "rtf", "txt", "html"},
		"docx": {"pdf", "doc", "odt", "rtf", "txt", "html"},
		"odt":  {"pdf", "doc", "docx", "rtf", "txt", "html"},
		"rtf":  {"pdf", "doc", "docx", "odt", "txt", "html"},
		"txt":  {"pdf", "doc", "docx", "odt", "rtf", "html"},
		"html": {"pdf", "doc", "docx", "odt", "rtf", "txt"},
		
		// Spreadsheets
		"xls":  {"pdf", "xlsx", "ods", "csv"},
		"xlsx": {"pdf", "xls", "ods", "csv"},
		"ods":  {"pdf", "xls", "xlsx", "csv"},
		"csv":  {"xls", "xlsx", "ods"},
		
		// Presentations
		"ppt":  {"pdf", "pptx", "odp"},
		"pptx": {"pdf", "ppt", "odp"},
		"odp":  {"pdf", "ppt", "pptx"},
	}
}

// Convert converts a document from one format to another using LibreOffice
func (dc *DocumentConverter) Convert(ctx context.Context, inputPath, outputPath string, options map[string]interface{}) error {
	// Check if soffice (LibreOffice) is available
	if !dc.binMgr.BinaryExists("soffice") {
		return fmt.Errorf("LibreOffice (soffice) is required for document conversion")
	}

	// Create a temporary directory for LibreOffice output
	tempDir, err := os.MkdirTemp("", "libreoffice-convert-")
	if err != nil {
		return fmt.Errorf("failed to create temporary directory: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Get the output format from the file extension
	extension := filepath.Ext(outputPath)
	if len(extension) == 0 {
		return fmt.Errorf("output file must have an extension")
	}
	outputFormat := extension[1:] // Remove the dot

	// Build the LibreOffice command
	cmdArgs := []string{
		"--headless",
		"--convert-to", outputFormat,
		"--outdir", tempDir,
	}

	// Add input file path
	cmdArgs = append(cmdArgs, inputPath)

	// Run LibreOffice
	cmd := exec.CommandContext(ctx, "soffice", cmdArgs...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("LibreOffice conversion failed: %v", err)
	}

	// Find the converted file
	convertedFile := filepath.Join(tempDir, filepath.Base(inputPath))
	ext := filepath.Ext(convertedFile)
	convertedFile = convertedFile[:len(convertedFile)-len(ext)] + "." + outputFormat

	// Check if the converted file exists
	if _, err := os.Stat(convertedFile); os.IsNotExist(err) {
		return fmt.Errorf("converted file not found: %s", convertedFile)
	}

	// Move the converted file to the desired output path
	if err := os.Rename(convertedFile, outputPath); err != nil {
		// If rename fails (cross-device), copy the file
		return copyFile(convertedFile, outputPath)
	}

	return nil
}

// copyFile copies a file from src to dst
func copyFile(src, dst string) error {
	input, err := os.ReadFile(src)
	if err != nil {
		return fmt.Errorf("failed to read source file: %v", err)
	}

	if err := os.WriteFile(dst, input, 0644); err != nil {
		return fmt.Errorf("failed to write destination file: %v", err)
	}

	return nil
}
