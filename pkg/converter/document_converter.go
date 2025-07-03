package converter

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"strconv"
	"strings"
	"sync"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// DocumentConverter handles document format conversions using LibreOffice
type DocumentConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	mu          sync.RWMutex
}

// NewDocumentConverter creates a new DocumentConverter
func NewDocumentConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &DocumentConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported formats and conversions
	// Text documents
	converter.AddSupportedConversion("doc", "pdf", "docx", "odt", "rtf", "txt", "html")
	converter.AddSupportedConversion("docx", "pdf", "doc", "odt", "rtf", "txt", "html")
	converter.AddSupportedConversion("odt", "pdf", "doc", "docx", "rtf", "txt", "html")
	converter.AddSupportedConversion("rtf", "pdf", "doc", "docx", "odt", "txt", "html")
	converter.AddSupportedConversion("txt", "pdf", "doc", "docx", "odt", "rtf", "html")
	converter.AddSupportedConversion("html", "pdf", "doc", "docx", "odt", "rtf", "txt")

	// Spreadsheets
	converter.AddSupportedConversion("xls", "pdf", "xlsx", "ods", "csv")
	converter.AddSupportedConversion("xlsx", "pdf", "xls", "ods", "csv")
	converter.AddSupportedConversion("ods", "pdf", "xls", "xlsx", "csv")
	converter.AddSupportedConversion("csv", "xls", "xlsx", "ods")

	// Presentations
	converter.AddSupportedConversion("ppt", "pdf", "pptx", "odp")
	converter.AddSupportedConversion("pptx", "pdf", "ppt", "odp")
	converter.AddSupportedConversion("odp", "pdf", "ppt", "pptx")

	return converter
}

// getLibreOfficeFormat maps file extensions to LibreOffice filter names
func (c *DocumentConverter) getLibreOfficeFormat(ext string) (string, error) {
	formatMap := map[string]string{
		// Document formats
		"doc":  "MS Word 97",
		"docx": "MS Word 2007 XML",
		"odt":  "writer8",
		"rtf":  "Rich Text Format",
		"txt":  "Text (encoded):UTF8",
		"html": "HTML (StarWriter)",

		// Spreadsheet formats
		"xls":  "MS Excel 97",
		"xlsx": "Calc MS Excel 2007 XML",
		"ods":  "calc8",
		"csv":  "Text - txt - csv (StarCalc)",

		// Presentation formats
		"ppt":  "MS PowerPoint 97",
		"pptx": "MS PowerPoint 2007 XML",
		"odp":  "impress8",

		// PDF
		"pdf": "writer_pdf_Export",
	}

	ext = strings.ToLower(ext)
	if format, exists := formatMap[ext]; exists {
		return format, nil
	}

	return "", fmt.Errorf("unsupported document format: %s", ext)
}

// Convert converts a document from one format to another using LibreOffice
func (c *DocumentConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Implementation of the Convert method
	// This is a placeholder - the actual implementation would go here
	return nil
}

// killLibreOffice kills any running LibreOffice processes
func killLibreOffice() error {
	cmd := exec.Command("pkill", "-f", "soffice.bin")
	if err := cmd.Run(); err != nil {
		// Ignore error if no processes were found to kill
		if exitErr, ok := err.(*exec.ExitError); ok && exitErr.ExitCode() == 1 {
			return nil
		}
		return fmt.Errorf("failed to kill LibreOffice processes: %w", err)
	}
	return nil
}

// setFilePermissions sets the correct permissions for the output file
func setFilePermissions(path string) error {
	// Set read/write permissions for owner and group
	if err := os.Chmod(path, 0664); err != nil {
		return fmt.Errorf("failed to set file permissions: %w", err)
	}

	// Get the current user info
	currentUser, err := user.Current()
	if err != nil {
		return fmt.Errorf("failed to get current user: %w", err)
	}

	// Get the group ID for the appuser group
	group, err := user.LookupGroup("appuser")
	if err != nil {
		// If we can't find the appuser group, just log a warning and continue
		log.Warn().Err(err).Msg("Failed to find appuser group, using current group")
		group = &user.Group{Gid: currentUser.Gid}
	}

	// Convert group ID to int
	gid, err := strconv.Atoi(group.Gid)
	if err != nil {
		return fmt.Errorf("invalid group ID: %w", err)
	}

	// Set the group ownership
	if err := os.Chown(path, -1, gid); err != nil {
		return fmt.Errorf("failed to set file group: %w", err)
	}

	return nil
}
