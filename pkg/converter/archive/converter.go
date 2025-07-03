package archive

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

// ArchiveConverter handles archive format conversions
type ArchiveConverter struct {
	*base.BaseConverter
	toolManager *tools.ToolManager
}

// NewArchiveConverter creates a new ArchiveConverter
func NewArchiveConverter(toolManager *tools.ToolManager, tempDir string) iface.Converter {
	converter := &ArchiveConverter{
		BaseConverter: base.NewBaseConverter(toolManager, tempDir),
		toolManager:   toolManager,
	}

	// Register supported archive formats
	converter.AddSupportedConversion("zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z", "rar")
	converter.AddSupportedConversion("tar", "zip", "tar.gz", "tar.bz2", "tar.xz", "7z")
	converter.AddSupportedConversion("tar.gz", "zip", "tar", "tar.bz2", "tar.xz", "7z")
	converter.AddSupportedConversion("tar.bz2", "zip", "tar", "tar.gz", "tar.xz", "7z")
	converter.AddSupportedConversion("tar.xz", "zip", "tar", "tar.gz", "tar.bz2", "7z")
	converter.AddSupportedConversion("7z", "zip", "tar", "tar.gz", "tar.bz2", "tar.xz")
	converter.AddSupportedConversion("rar", "zip", "tar", "tar.gz", "tar.bz2", "tar.xz", "7z")

	return converter
}

// Convert converts an archive from one format to another
func (c *ArchiveConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Extract source and target formats from file extensions
	sourceFormat := strings.TrimPrefix(filepath.Ext(inputPath), ".")
	if sourceFormat == "" {
		return fmt.Errorf("could not determine source format from file extension")
	}

	targetFormat := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if targetFormat == "" {
		return fmt.Errorf("could not determine target format from file extension")
	}

	// Create a temporary directory for extraction
	extractDir, err := os.MkdirTemp("", "extract-*")
	if err != nil {
		return fmt.Errorf("failed to create temp directory: %w", err)
	}
	defer os.RemoveAll(extractDir)

	// Extract the source archive
	log.Debug().
		Str("source", inputPath).
		Str("target", outputPath).
		Str("source_format", sourceFormat).
		Str("target_format", targetFormat).
		Msg("Extracting source archive")

	if err := c.extractArchive(inputPath, extractDir, sourceFormat); err != nil {
		return fmt.Errorf("failed to extract archive: %w", err)
	}

	// Create output directory if it doesn't exist
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Create the target archive
	log.Debug().
		Str("path", outputPath).
		Str("format", targetFormat).
		Msg("Creating target archive")

	if err := c.createArchive(extractDir, outputPath, targetFormat, nil); err != nil {
		return fmt.Errorf("failed to create archive: %w", err)
	}

	return nil
}

// Cleanup removes temporary files
func (c *ArchiveConverter) Cleanup(files ...string) error {
	for _, file := range files {
		if err := os.RemoveAll(file); err != nil {
			return fmt.Errorf("failed to remove file %s: %w", file, err)
		}
	}
	return nil
}

// SupportsConversion checks if the converter supports the given conversion
func (c *ArchiveConverter) SupportsConversion(sourceFormat, targetFormat string) bool {
	return c.BaseConverter.SupportsConversion(sourceFormat, targetFormat)
}

// extractArchive extracts an archive to the specified directory
func (c *ArchiveConverter) extractArchive(src, dest, format string) error {
	// Create destination directory if it doesn't exist
	if err := os.MkdirAll(dest, 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}

	switch format {
	case "zip":
		return c.extractZip(src, dest)
	case "tar":
		return c.extractTar(src, dest, false)
	case "tar.gz", "tgz":
		return c.extractTar(src, dest, true)
	case "tar.bz2", "tbz2":
		return c.extractTarBz2(src, dest)
	case "tar.xz", "txz":
		return c.extractTarXz(src, dest)
	case "rar":
		return c.extractRar(src, dest)
	case "7z":
		return c.extract7z(src, dest)
	default:
		return fmt.Errorf("unsupported archive format: %s", format)
	}
}

// createArchive creates an archive from the specified directory
func (c *ArchiveConverter) createArchive(src, dest, format string, options map[string]interface{}) error {
	// Create parent directory if it doesn't exist
	if err := os.MkdirAll(filepath.Dir(dest), 0755); err != nil {
		return fmt.Errorf("failed to create parent directory: %w", err)
	}

	switch format {
	case "zip":
		return c.createZip(src, dest, options)
	case "tar":
		return c.createTar(src, dest, false, options)
	case "tar.gz", "tgz":
		return c.createTar(src, dest, true, options)
	case "tar.bz2", "tbz2":
		return c.createTarBz2(src, dest, options)
	case "tar.xz", "txz":
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
