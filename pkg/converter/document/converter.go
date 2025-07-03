package document

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/rs/zerolog/log"
	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
)

// DocumentConverter handles document format conversions
type DocumentConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
}

// NewDocumentConverter creates a new DocumentConverter
func NewDocumentConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &DocumentConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported formats and conversions
	converter.AddSupportedConversion("doc", "pdf", "docx", "odt", "txt", "rtf")
	converter.AddSupportedConversion("docx", "pdf", "doc", "odt", "txt", "rtf")
	converter.AddSupportedConversion("odt", "pdf", "doc", "docx", "txt", "rtf")
	converter.AddSupportedConversion("pdf", "docx", "doc", "odt", "txt", "rtf")
	converter.AddSupportedConversion("rtf", "docx", "doc", "odt", "pdf", "txt")
	converter.AddSupportedConversion("txt", "docx", "doc", "odt", "pdf", "rtf")

	return converter
}

// Convert converts a document from one format to another using LibreOffice
func (c *DocumentConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Extract source and target formats from file extensions
	sourceFormat := strings.TrimPrefix(filepath.Ext(inputPath), ".")
	if sourceFormat == "" {
		return fmt.Errorf("could not determine source format from file extension")
	}

	targetFormat := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if targetFormat == "" {
		return fmt.Errorf("could not determine target format from file extension")
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create a temporary directory for the conversion
	tempDir, err := c.getOutputDir()
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(tempDir)

	// Convert the document
	if err := c.convertWithLibreOffice(inputPath, tempDir, targetFormat); err != nil {
		return fmt.Errorf("conversion failed: %w", err)
	}

	// Find the converted file
	files, err := filepath.Glob(filepath.Join(tempDir, "*"))
	if err != nil || len(files) == 0 {
		return fmt.Errorf("failed to find converted file in output directory")
	}

	// Move the converted file to the final location
	if err := os.Rename(files[0], outputPath); err != nil {
		return fmt.Errorf("failed to move output file: %w", err)
	}

	return nil
}

func (c *DocumentConverter) convertWithLibreOffice(inputPath, outputDir, targetFormat string) error {
	// Determine the output format for LibreOffice
	libreofficeFormat, err := getLibreOfficeFormat(targetFormat)
	if err != nil {
		return fmt.Errorf("unsupported target format: %w", err)
	}

	// Create output directory if it doesn't exist
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create and run the LibreOffice command
	cmd := exec.Command(
		"libreoffice",
		"--headless",
		"--convert-to", libreofficeFormat,
		"--outdir", outputDir,
		inputPath,
	)

	// Run the command with output capture
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("LibreOffice conversion failed")
		return fmt.Errorf("libreoffice conversion failed: %w\nOutput: %s", err, string(output))
	}

	return nil
}

// getLibreOfficeFormat returns the LibreOffice format filter for the given file extension
func getLibreOfficeFormat(ext string) (string, error) {
	switch strings.ToLower(ext) {
	case "doc":
		return "MS Word 97", nil
	case "docx":
		return "MS Word 2007 XML", nil
	case "odt":
		return "writer8", nil
	case "pdf":
		return "writer_pdf_Export", nil
	case "rtf":
		return "Rich Text Format", nil
	case "txt":
		return "Text (encoded):UTF8", nil
	default:
		return "", fmt.Errorf("unsupported document format: %s", ext)
	}
}

// getOutputDir creates and returns a temporary directory for the conversion
func (c *DocumentConverter) getOutputDir() (string, error) {
	tempDir, err := os.MkdirTemp("", "doc_convert_*")
	if err != nil {
		log.Error().Err(err).Msg("Failed to create temp directory")
		return "", fmt.Errorf("failed to create temp directory: %w", err)
	}
	return tempDir, nil
}

// ValidateOptions validates the conversion options
func (c *DocumentConverter) ValidateOptions(options map[string]interface{}) error {
	// Check required options
	if _, ok := options["source_format"]; !ok {
		return iface.NewConversionError("missing_option", "source_format is required", nil)
	}

	if _, ok := options["target_format"]; !ok {
		return iface.NewConversionError("missing_option", "target_format is required", nil)
	}

	sourceFormat, ok := options["source_format"].(string)
	if !ok {
		return iface.NewConversionError("invalid_option", "source_format must be a string", nil)
	}

	targetFormat, ok := options["target_format"].(string)
	if !ok {
		return iface.NewConversionError("invalid_option", "target_format must be a string", nil)
	}

	if !c.SupportsConversion(sourceFormat, targetFormat) {
		return iface.NewConversionError("unsupported_conversion", 
			fmt.Sprintf("conversion from %s to %s is not supported", sourceFormat, targetFormat), nil)
	}

	// Check if required external tools are installed
	if _, err := exec.LookPath("libreoffice"); err != nil {
		return iface.NewConversionError("dependency_error", "LibreOffice is required for document conversion", err)
	}

	return nil
}

// Cleanup removes temporary files created during conversion
func (c *DocumentConverter) Cleanup(files ...string) error {
	var lastErr error
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			log.Error().Err(err).Str("file", file).Msg("Failed to remove temporary file")
			lastErr = err
		}
	}
	return lastErr
}
