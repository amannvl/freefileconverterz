package archive

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
)

// ArchiveConverter handles archive format conversions
type ArchiveConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
	tempDir    string
}

// NewArchiveConverter creates a new ArchiveConverter
func NewArchiveConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	return &ArchiveConverter{
		BaseConverter: base.NewBaseConverter("archive", getSupportedFormats()),
		toolManager:   toolManager,
		tempDir:       tempDir,
	}
}



// getSupportedFormats returns the supported archive formats and their conversions
func getSupportedFormats() map[string][]string {
	return map[string][]string{
		"zip":    {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"tar":    {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"tar.gz": {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"tar.bz2":{"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"tar.xz": {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"7z":     {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
		"rar":    {"zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z"},
	}
}

// ValidateOptions validates the conversion options
func (c *ArchiveConverter) ValidateOptions(options map[string]interface{}) error {
	// No specific validation needed for archive conversion
	return nil
}



// Convert converts an archive from one format to another
func (c *ArchiveConverter) Convert(ctx context.Context, input io.Reader, options map[string]interface{}) (io.Reader, error) {
	// Create a temporary file for the input
	tempInput, err := os.CreateTemp(c.tempDir, "input-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary file: %w", err)
	}
	defer os.Remove(tempInput.Name())

	// Write the input to the temporary file
	if _, err := io.Copy(tempInput, input); err != nil {
		return nil, fmt.Errorf("failed to write to temporary file: %w", err)
	}
	tempInput.Close()

	// Create a temporary file for the output
	tempOutput, err := os.CreateTemp(c.tempDir, "output-*.tmp")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary output file: %w", err)
	}
	tempOutput.Close()
	defer os.Remove(tempOutput.Name())

	// Get the source and target formats from options
	sourceFormat, ok := options["source_format"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid source_format in options")
	}

	targetFormat, ok := options["target_format"].(string)
	if !ok {
		return nil, fmt.Errorf("missing or invalid target_format in options")
	}

	// Create a temporary directory for extraction
	extractDir, err := os.MkdirTemp(c.tempDir, "extract-")
	if err != nil {
		return nil, fmt.Errorf("failed to create temporary directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract the input archive
	if err := c.extractArchive(tempInput.Name(), extractDir, sourceFormat); err != nil {
		return nil, fmt.Errorf("failed to extract archive: %w", err)
	}

	// Create the output archive
	if err := c.createArchive(extractDir, tempOutput.Name(), targetFormat, options); err != nil {
		return nil, fmt.Errorf("failed to create archive: %w", err)
	}

	// Read the output file
	output, err := os.ReadFile(tempOutput.Name())
	if err != nil {
		return nil, fmt.Errorf("failed to read output file: %w", err)
	}

	return bytes.NewReader(output), nil
}

// extractArchive extracts an archive to the specified directory
func (c *ArchiveConverter) extractArchive(src, dest, format string) error {
	switch format {
	case "zip":
		return c.extractZip(src, dest)
	case "tar":
		return c.extractTar(src, dest, false)
	case "tar.gz", "tgz":
		return c.extractTar(src, dest, true)
	case "tar.bz2":
		return c.extractTarBz2(src, dest)
	case "tar.xz":
		return c.extractTarXz(src, dest)
	case "7z":
		return c.extract7z(src, dest)
	case "rar":
		return c.extractRar(src, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// createArchive creates an archive from the specified directory
func (c *ArchiveConverter) createArchive(src, dest, format string, options map[string]interface{}) error {
	switch format {
	case "zip":
		return c.createZip(src, dest, options)
	case "tar":
		return c.createTar(src, dest, false, options)
	case "tar.gz", "tgz":
		return c.createTar(src, dest, true, options)
	case "tar.bz2":
		return c.createTarBz2(src, dest, options)
	case "tar.xz":
		return c.createTarXz(src, dest, options)
	case "7z":
		return c.create7z(src, dest, options)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// extractZip extracts a ZIP archive
func (c *ArchiveConverter) extractZip(src, dest string) error {
	return exec.Command("unzip", "-o", src, "-d", dest).Run()
}

// extractTar extracts a TAR archive, optionally with gzip compression
func (c *ArchiveConverter) extractTar(src, dest string, gzipped bool) error {
	args := []string{"-xf", src, "-C", dest}
	if gzipped {
		args = append([]string{"-z"}, args...)
	}
	return exec.Command("tar", args...).Run()
}

// extractTarBz2 extracts a BZIP2 compressed TAR archive
func (c *ArchiveConverter) extractTarBz2(src, dest string) error {
	return exec.Command("tar", "-xjf", src, "-C", dest).Run()
}

// extractTarXz extracts a XZ compressed TAR archive
func (c *ArchiveConverter) extractTarXz(src, dest string) error {
	return exec.Command("tar", "-xJf", src, "-C", dest).Run()
}

// extract7z extracts a 7-Zip archive
func (c *ArchiveConverter) extract7z(src, dest string) error {
	return exec.Command("7z", "x", "-o"+dest, src).Run()
}

// extractRar extracts a RAR archive
func (c *ArchiveConverter) extractRar(src, dest string) error {
	return exec.Command("unrar", "x", "-o+", src, dest).Run()
}

// createZip creates a ZIP archive
func (c *ArchiveConverter) createZip(src, dest string, options map[string]interface{}) error {
	return exec.Command("zip", "-r", dest, ".").Run()
}

// createTar creates a TAR archive, optionally with gzip compression
func (c *ArchiveConverter) createTar(src, dest string, gzipped bool, options map[string]interface{}) error {
	args := []string{"-cf", dest, "-C", filepath.Dir(src), filepath.Base(src)}
	if gzipped {
		args = append([]string{"-z"}, args...)
	}
	return exec.Command("tar", args...).Run()
}

// createTarBz2 creates a BZIP2 compressed TAR archive
func (c *ArchiveConverter) createTarBz2(src, dest string, options map[string]interface{}) error {
	return exec.Command("tar", "-cjf", dest, "-C", filepath.Dir(src), filepath.Base(src)).Run()
}

// createTarXz creates a XZ compressed TAR archive
func (c *ArchiveConverter) createTarXz(src, dest string, options map[string]interface{}) error {
	return exec.Command("tar", "-cJf", dest, "-C", filepath.Dir(src), filepath.Base(src)).Run()
}

// create7z creates a 7-Zip archive
func (c *ArchiveConverter) create7z(src, dest string, options map[string]interface{}) error {
	return exec.Command("7z", "a", dest, src).Run()
}
