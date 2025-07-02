package document

import (
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// DocumentConverter handles document format conversions
type DocumentConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	tempDir     string
}

// Ensure DocumentConverter implements iface.Converter
var _ iface.Converter = (*DocumentConverter)(nil)

// NewDocumentConverter creates a new DocumentConverter
func NewDocumentConverter(toolManager *tools.ToolManager, tempDir string) *DocumentConverter {
	supportedFormats := map[string][]string{
		"doc":  {"pdf", "docx", "odt", "txt", "rtf"},
		"docx": {"pdf", "doc", "odt", "txt", "rtf"},
		"odt":  {"pdf", "doc", "docx", "txt", "rtf"},
		"pdf":  {"docx", "doc", "odt", "txt", "rtf"},
		"rtf":  {"docx", "doc", "odt", "pdf", "txt"},
		"txt":  {"docx", "doc", "odt", "pdf", "rtf"},
	}

	return &DocumentConverter{
		BaseConverter: base.NewBaseConverter("document", supportedFormats),
		toolManager:   toolManager,
		tempDir:       tempDir,
	}
}

// Convert converts a document from one format to another using LibreOffice
func (c *DocumentConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
	if err := c.ValidateOptions(options); err != nil {
		return nil, err
	}

	sourceFormat, _ := options["source_format"].(string)
	targetFormat, _ := options["target_format"].(string)

	log.Info().
		Str("source_format", sourceFormat).
		Str("target_format", targetFormat).
		Msg("Starting document conversion")

	// Create a temporary file for the input
	tempInput, err := os.CreateTemp(c.tempDir, "input-*."+sourceFormat)
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary input file: %w", err)
	}
	defer os.Remove(tempInput.Name())

	// Write the input to the temporary file
	if _, err := io.Copy(tempInput, input); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempInput.Close()

	// Create a temporary directory for LibreOffice output
	outputDir, err := os.MkdirTemp(c.tempDir, "libreoffice-out-")
	if err != nil {
		return nil, fmt.Errorf("failed to create output directory: %w", err)
	}
	defer os.RemoveAll(outputDir)

	// Determine the output format for LibreOffice
	libreofficeFormat, err := getLibreOfficeFormat(targetFormat)
	if err != nil {
		return nil, err
	}

	// Create and run the LibreOffice command
	cmd, err := c.toolManager.Command(
		ctx,
		"libreoffice",
		"--headless",
		"--convert-to", libreofficeFormat,
		"--outdir", outputDir,
		tempInput.Name(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to create LibreOffice command: %w", err)
	}

	// Run the conversion
	output, err := cmd.CombinedOutput()
	if err != nil {
		log.Error().
			Err(err).
			Str("output", string(output)).
			Msg("LibreOffice conversion failed")
		return nil, iface.NewConversionError(
			"conversion_failed",
			"failed to convert document",
			err,
		)
	}

	// Find the output file
	outputFiles, err := filepath.Glob(filepath.Join(outputDir, "*."+targetFormat))
	if err != nil || len(outputFiles) == 0 {
		return nil, iface.NewConversionError(
			"output_not_found",
			"could not find output file",
			err,
		)
	}

	// Open the output file
	outputFile, err := os.Open(outputFiles[0])
	if err != nil {
		return nil, fmt.Errorf("failed to open output file: %w", err)
	}

	return outputFile, nil
}

// getLibreOfficeFormat returns the LibreOffice format filter for the given file extension
func getLibreOfficeFormat(ext string) (string, error) {
	switch ext {
	case "pdf":
		return "writer_pdf_Export", nil
	case "docx":
		return "MS Word 2007 XML", nil
	case "doc":
		return "MS Word 97", nil
	case "odt":
		return "writer8", nil
	case "rtf":
		return "Rich Text Format", nil
	case "txt":
		return "Text (encoded)", nil
	case "html":
		return "HTML (StarWriter)", nil
	default:
		return "", fmt.Errorf("unsupported output format: %s", ext)
	}
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
