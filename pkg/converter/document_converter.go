package converter

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/amannvl/freefileconverterz/pkg/converter/base"
	"github.com/amannvl/freefileconverterz/pkg/converter/iface"
	"github.com/rs/zerolog/log"
)

// DocumentConverter handles document format conversions using LibreOffice
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

// Convert converts a document from one format to another using LibreOffice
func (c *DocumentConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Log the start of the conversion with full paths
	log.Info().
		Str("input_path", inputPath).
		Str("output_path", outputPath).
		Msg("Starting document conversion")

	// Verify input file exists and is readable
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Error().
			Err(err).
			Str("input_path", inputPath).
			Msg("Input file does not exist")
		return fmt.Errorf("input file does not exist: %s", inputPath)
	}

	// Get output directory and ensure it exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Error().
			Err(err).
			Str("output_dir", outputDir).
			Msg("Failed to create output directory")
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Kill any existing LibreOffice processes to avoid conflicts
	if err := exec.Command("pkill", "-f", "soffice.bin").Run(); err != nil {
		// Ignore error if no processes were found to kill
		if exitErr, ok := err.(*exec.ExitError); !ok || exitErr.ExitCode() != 1 {
			log.Warn().
				Err(err).
				Msg("Failed to kill existing LibreOffice processes")
		}
	}
	// Log the start of the conversion with full paths
	log.Info().
		Str("input_path", inputPath).
		Str("output_path", outputPath).
		Msg("Starting document conversion")

	// Verify input file exists and is readable
	inputInfo, err := os.Stat(inputPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("input_path", inputPath).
			Msg("Input file does not exist or is not accessible")
		return fmt.Errorf("input file does not exist or is not accessible: %w", err)
	}

	// Get output directory and ensure it exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		log.Error().
			Err(err).
			Str("output_dir", outputDir).
			Msg("Failed to create output directory")
		return fmt.Errorf("failed to create output directory: %w", err)
	}

	// Get the target format from the output file extension
	ext := filepath.Ext(outputPath)
	if ext == "" {
		err := fmt.Errorf("output file must have an extension")
		log.Error().
			Err(err).
			Str("output_path", outputPath).
			Msg("Output file has no extension")
		return err
	}
	targetFormat := strings.TrimPrefix(ext, ".")

	// Map file extensions to LibreOffice filter names
	format, err := getLibreOfficeFormat(targetFormat)
	if err != nil {
		log.Error().
			Err(err).
			Str("target_format", targetFormat).
			Msg("Unsupported output format")
		return fmt.Errorf("unsupported output format: %s: %w", targetFormat, err)
	}
	// Add a timeout to the context if not already set
	if _, hasDeadline := ctx.Deadline(); !hasDeadline {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, 5*time.Minute)
		defer cancel()
	}
	var (
		err       error
		inputInfo os.FileInfo
	)

	// Log the start of the conversion with full paths
	log.Info().
		Str("input_path", inputPath).
		Str("output_path", outputPath).
		Str("working_dir", c.tempDir).
		Msg("Starting document conversion")

	// Log environment variables for debugging
	envVars := os.Environ()
	envVarsStr := ""
	for _, v := range envVars {
		envVarsStr += v + "\n"
	}

	log.Debug().
		Str("environment", envVarsStr).
		Msg("Environment variables")

	// Get the file extension to determine the target format
	extension := strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if extension == "" {
		err := fmt.Errorf("output path has no extension")
		log.Error().Err(err).Str("path", outputPath).Msg("Invalid output path")
		return iface.NewConversionError(
			"invalid_output",
			err.Error(),
			err,
		)
	}
	targetFormat := strings.ToLower(extension)

	// Verify input file exists and is readable
	inputInfo, err = os.Stat(inputPath)
	if err != nil {
		log.Error().
			Err(err).
			Str("path", inputPath).
			Msg("Failed to stat input file")
		return iface.NewConversionError(
			"input_not_found",
			fmt.Sprintf("failed to access input file: %v", err),
			err,
		)
	}

	if inputInfo.IsDir() {
		err := fmt.Errorf("input path is a directory")
		log.Error().Err(err).Str("path", inputPath).Msg("Input is a directory")
		return iface.NewConversionError("invalid_input", err.Error(), err)
	}

	log.Debug().
		Str("input_path", inputPath).
		Int64("size", inputInfo.Size()).
		Stringer("mode", inputInfo.Mode()).
		Msg("Input file details")

	// Get the output directory and ensure it exists
	outputDir := filepath.Dir(outputPath)
	log.Debug().
		Str("output_dir", outputDir).
		Str("absolute_output_dir", outputDir).
		Msg("Ensuring output directory exists")

	// Create the output directory with 0777 permissions
	if err := os.MkdirAll(outputDir, 0777); err != nil {
		log.Error().
			Err(err).
			Str("path", outputDir).
			Msg("Failed to create output directory")
		return iface.NewConversionError(
			"create_output_dir_failed",
			fmt.Sprintf("failed to create output directory: %v", err),
			err,
		)
	}

	// Ensure the output directory has the correct permissions
	if err := os.Chmod(outputDir, 0777); err != nil {
		log.Warn().
			Err(err).
			Str("path", outputDir).
			Msg("Failed to set permissions on output directory")
	}

	if mkdirErr := os.MkdirAll(outputDir, 0755); mkdirErr != nil {
		log.Error().
			Err(mkdirErr).
			Str("path", outputDir).
			Msg("Failed to create output directory")
		return iface.NewConversionError(
			"create_output_dir_failed",
			fmt.Sprintf("failed to create output directory: %v", mkdirErr),
			mkdirErr,
		)
	}

	// Log directory permissions
	if dirInfo, statErr := os.Stat(outputDir); statErr != nil {
		log.Warn().
			Err(statErr).
			Str("path", outputDir).
			Msg("Failed to stat output directory")
	} else {
		log.Debug().
			Str("path", outputDir).
			Stringer("mode", dirInfo.Mode()).
			Msg("Output directory info")
	}

	// Ensure output directory is writable
	log.Debug().
		Str("path", outputDir).
		Msg("Setting directory permissions to 0755")

	if chmodErr := os.Chmod(outputDir, 0755); chmodErr != nil {
		log.Error().
			Err(chmodErr).
			Str("path", outputDir).
			Msg("Failed to set permissions on output directory")
	} else {
		log.Debug().
			Str("path", outputDir).
			Msg("Successfully set directory permissions")
	}

	// Log the conversion attempt with all relevant details
	log.Info().
		Str("source", inputPath).
		Str("target", outputPath).
		Str("target_format", targetFormat).
		Str("output_dir", outputDir).
		Int64("input_size", inputInfo.Size()).
		Stringer("input_mode", inputInfo.Mode()).
		Msg("Starting document conversion with LibreOffice")

	// Log environment variables for debugging
	envVars := os.Environ()
	log.Debug().
		Strs("environment", envVars).
		Msg("Environment variables")

	// Check if we can execute LibreOffice
	libreOfficePath, err := c.toolManager.GetLibreOfficePath()
	if err != nil {
		log.Error().
			Err(err).
			Msg("LibreOffice not found")
		return iface.NewConversionError(
			"tool_not_found",
			"LibreOffice is not installed or not in PATH",
			err,
		)
	}

	// Verify LibreOffice is executable
	if _, err := os.Stat(libreOfficePath); os.IsNotExist(err) {
		log.Error().
			Str("path", libreOfficePath).
			Msg("LibreOffice binary not found at the specified path")
		return iface.NewConversionError(
			"tool_not_found",
			fmt.Sprintf("LibreOffice binary not found at %s", libreOfficePath),
			err,
		)
	}

	// Check if the input file exists and is readable
	if _, err := os.Stat(inputPath); os.IsNotExist(err) {
		log.Error().
			Str("path", inputPath).
			Msg("Input file does not exist")
		return iface.NewConversionError(
			"input_not_found",
			fmt.Sprintf("Input file does not exist: %s", inputPath),
			err,
		)
	}

	// Check if LibreOffice is installed
	if _, err := c.toolManager.GetLibreOfficePath(); err != nil {
		return iface.NewConversionError(
			"tool_not_found",
			"LibreOffice is required for document conversion",
			err,
		)
	}

	// Get the output directory and ensure it exists
	outputDir := filepath.Dir(outputPath)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		return iface.NewConversionError(
			"create_output_dir_failed",
			fmt.Sprintf("failed to create output directory: %v", err),
			err,
		)
	}

	// Get the LibreOffice format for the target format
	format, err := getLibreOfficeFormat(targetFormat)
	if err != nil {
		return iface.NewConversionError(
			"unsupported_format",
			fmt.Sprintf("unsupported target format: %s", targetFormat),
			err,
		)
	}

	// Log the format being used for conversion
	log.Debug().
		Str("format", format).
		Str("output_dir", outputDir).
		Str("input_path", inputPath).
		Msg("Preparing LibreOffice conversion command")

	// Log the current working directory
	if wd, err := os.Getwd(); err == nil {
		log.Info().
			Str("current_working_directory", wd).
			Msg("Current working directory")
	}

	// Log input file info
	log.Info().
		Str("input_path", inputPath).
		Int64("size", inputInfo.Size()).
		Stringer("mode", inputInfo.Mode()).
		Msg("Input file info")

	// Log output directory info
	if outputDirInfo, err := os.Stat(outputDir); err == nil {
		log.Info().
			Str("output_dir", outputDir).
			Stringer("mode", outputDirInfo.Mode()).
			Msg("Output directory info")
	} else {
		log.Error().
			Err(err).
			Str("output_dir", outputDir).
			Msg("Failed to get output directory info")
	}

	// Kill any existing LibreOffice processes
	if err := killLibreOffice(); err != nil {
		log.Warn().
			Err(err).
			Msg("Failed to kill existing LibreOffice processes")
	}

	// Build the command to run our conversion script
	cmd := exec.CommandContext(
		ctx,
		"/app/convert.sh",
		inputPath,
		outputDir,
		filepath.Base(outputPath),
	)

	// Log the exact command being executed
	log.Info().
		Str("command", cmd.String()).
		Str("working_dir", cmd.Dir).
		Str("input_path", inputPath).
		Str("output_path", outputPath).
		Str("output_dir", outputDir).
		Str("format", format).
		Msg("Executing LibreOffice command")

	// Set the working directory to the output directory to ensure files are written correctly
	cmd.Dir = outputDir
	
	// Log the current directory and environment
	if wd, err := os.Getwd(); err == nil {
		log.Info().
			Str("current_working_directory", wd).
			Msg("Current working directory")
	}

	// Set environment variables for the command
	cmd.Env = append(os.Environ(),
		"HOME=/home/appuser",
		"USER=appuser",
		"USERNAME=appuser",
		"LOGNAME=appuser",
	)
	
	// Log environment variables for debugging
	envVarsStr := ""
	for _, envVar := range cmd.Env {
		envVarsStr += envVar + "\n"
	}
	log.Debug().
		Str("environment", envVarsStr).
		Msg("Command environment variables")

	// Log the command being executed with full context
	log.Debug().
		Str("command", cmd.String()).
		Str("working_dir", cmd.Dir).
		Str("output_dir", outputDir).
		Str("input_path", inputPath).
		Str("output_path", outputPath).
		Str("format", format).
		Strs("environment", cmd.Env).
		Msg("Executing LibreOffice conversion command")

	// Log the current directory and its contents
	var (
		cwd     string
		cwdErr  error
	)
	cwd, cwdErr = os.Getwd()
	if cwdErr != nil {
		log.Error().
			Err(cwdErr).
			Msg("Failed to get current working directory")
	} else {
		log.Debug().
			Str("cwd", cwd).
			Msg("Current working directory")
	}

	// Log environment variables for debugging
	envVars := make([]string, 0, len(cmd.Env))
	for _, e := range cmd.Env {
		envVars = append(envVars, e)
	}
	log.Debug().
		Strs("environment", envVars).
		Msg("Command environment variables")

	// Log environment variables for debugging
	envVars = os.Environ()
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "HOME=") || 
		   strings.HasPrefix(envVar, "USER=") ||
		   strings.HasPrefix(envVar, "PATH=") ||
		   strings.HasPrefix(envVar, "LANG=") {
			log.Debug().Str("env", envVar).Msg("Environment variable")
		}
	}

	// Log current working directory
	if wd, err := os.Getwd(); err == nil {
		log.Debug().Str("cwd", wd).Msg("Current working directory")
	}

	// Log input file details
	if fi, err := os.Stat(inputPath); err == nil {
		log.Debug().
			Str("input_file", inputPath).
			Int64("size", fi.Size()).
			Stringer("mode", fi.Mode()).
			Msg("Input file details")
	} else {
		log.Error().Err(err).Str("path", inputPath).Msg("Failed to stat input file")
	}

	// Log output directory permissions
	if fi, err := os.Stat(outputDir); err == nil {
		log.Debug().
			Str("output_dir", outputDir).
			Stringer("mode", fi.Mode()).
			Stringer("perm", fi.Mode().Perm()).
			Msg("Output directory details")
	} else {
		log.Error().Err(err).Str("path", outputDir).Msg("Failed to stat output directory")
	}

	// The command completed successfully, but we still need to check the output file
	outputPath = filepath.Clean(outputPath)

	// Log the output directory and its contents
	log.Debug().
		Str("directory", outputDir).
		Msg("Output directory")

	// Check if the output directory exists and is accessible
	dirInfo, statErr := os.Stat(outputDir)
	if os.IsNotExist(statErr) {
		log.Error().
			Err(statErr).
			Str("path", outputDir).
			Msg("Output directory does not exist")
	} else if statErr != nil {
		log.Error().
			Err(statErr).
			Str("path", outputDir).
			Msg("Failed to stat output directory")
	} else {
		log.Debug().
			Str("path", outputDir).
			Stringer("mode", dirInfo.Mode()).
			Msg("Output directory info")

		// List directory contents for debugging
		files, readDirErr := os.ReadDir(outputDir)
		if readDirErr != nil {
			log.Error().
				Err(readDirErr).
				Str("path", outputDir).
				Msg("Failed to read output directory")
		} else {
			fileList := make([]string, 0, len(files))
			for _, file := range files {
				fileInfo, _ := file.Info()
				fileList = append(fileList, fmt.Sprintf("%s (%d bytes, mode: %s)", 
					file.Name(), fileInfo.Size(), fileInfo.Mode()))
			}
			log.Debug().
				Strs("files", fileList).
				Msg("Output directory contents")
		}
	
	// Run the command
	log.Info().
		Str("command", cmd.String()).
		Str("format", format).
		Str("output_dir", outputDir).
		Msg("Running LibreOffice command")

	// Capture command output
	var stdoutBuf, stderrBuf bytes.Buffer
	cmd.Stdout = &stdoutBuf
	cmd.Stderr = &stderrBuf

	// Log environment variables for debugging
	log.Info().
		Str("PATH", os.Getenv("PATH")).
		Str("HOME", os.Getenv("HOME")).
		Msg("Environment variables")

	// Log the current working directory
	if wd, err := os.Getwd(); err == nil {
		log.Info().
			Str("current_working_directory", wd).
			Msg("Current working directory")
			Msg("LibreOffice command start error")
		
		return iface.NewConversionError(
			"conversion_failed",
			fmt.Sprintf("failed to start LibreOffice: %v, stderr: %s", err, stderrBuf.String()),
			err,
		)
	}
	
	// Log the process ID
	log.Info().
		Int("pid", cmd.Process.Pid).
		Msg("LibreOffice process started")

	// Wait for the command to complete with a timeout
	done := make(chan error, 1)
	go func() {
		err := cmd.Wait()
		log.Info().
			Int("pid", cmd.Process.Pid).
			Err(err).
			Msg("LibreOffice process completed")
		done <- err
	}()

	select {
	case <-ctx.Done():
		// Context was cancelled or timed out
		err := ctx.Err()
		log.Error().
			Err(err).
			Int("pid", cmd.Process.Pid).
			Msg("Conversion context cancelled or timed out")
		
		if cmd.Process != nil {
			log.Warn().
				Int("pid", cmd.Process.Pid).
				Msg("Killing LibreOffice process due to timeout")
			
			if killErr := cmd.Process.Kill(); killErr != nil {
				log.Error().
					Err(killErr).
					Int("pid", cmd.Process.Pid).
					Msg("Failed to kill LibreOffice process after timeout")
			}
		}
		
		// Log the output buffers
		log.Error().
			Str("stdout", stdoutBuf.String()).
			Str("stderr", stderrBuf.String()).
			Msg("Command output before timeout")
		
		return iface.NewConversionError(
			"conversion_timeout",
			fmt.Sprintf("document conversion timed out: %v", err),
			err,
		)

	case err := <-done:
		output := stdoutBuf.String()
		stderr := stderrBuf.String()

		// Log command output regardless of success/failure
		logLevel := log.Info()
		if err != nil {
			logLevel = log.Error()
		}

		logLevel.
			Str("command", cmd.String()).
			Str("working_dir", cmd.Dir).
			Str("stdout", output).
			Str("stderr", stderr).
			Str("error", fmt.Sprintf("%v", err)).
			Str("error_type", fmt.Sprintf("%T", err)).
			Msg("LibreOffice command completed")

		if err != nil {
			log.Error().
				Err(err).
				Str("command", cmd.String()).
				Str("stdout", output).
				Str("stderr", stderr).
				Str("error_type", fmt.Sprintf("%T", err)).
				Msg("Document conversion failed")

			if exitErr, ok := err.(*exec.ExitError); ok {
				log.Error().
					Int("exit_code", exitErr.ExitCode()).
					Str("stderr", stderr).
					Msg("LibreOffice command exited with error")
			}

			return iface.NewConversionError(
				"conversion_failed",
				fmt.Sprintf("failed to convert document: %v, stderr: %s", err, stderr),
				err,
			)
		}

		log.Info().
			Str("output", outputPath).
			Msg("Document conversion completed successfully")

		// Check if the output file was created
	outputFile := filepath.Join(outputDir, strings.TrimSuffix(filepath.Base(inputPath), ".pdf") + ".docx")
	if _, err := os.Stat(outputFile); os.IsNotExist(err) {
		// File doesn't exist, check if it was created with a different name
		files, err := os.ReadDir(outputDir)
		if err != nil {
				
				if err := os.Rename(actualPath, outputPath); err != nil {
					log.Error().
						Err(err).
						Str("source", actualPath).
						Str("target", outputPath).
						Msg("Failed to rename output file")
					
					return iface.NewConversionError(
						"rename_failed",
						fmt.Sprintf("failed to rename output file: %v", err),
						err,
					)
				}
				found = true
				break
			}
		}

		if !found {
			log.Error().
				Str("expected_file", outputPath).
				Int("files_found", len(files)).
				Msg("Converted file not found in output directory")
			
			return iface.NewConversionError(
				"output_not_found",
				"failed to find converted file in output directory",
				nil,
			)
		}
	}

	log.Info().
		Str("output", outputPath).
		Msg("Document conversion completed successfully")

	return nil
}

// getLibreOfficeFormat maps file extensions to LibreOffice filter names
func getLibreOfficeFormat(ext string) (string, error) {
	switch strings.ToLower(ext) {
	case "docx":
		// Using empty string as the format to let LibreOffice use the default
		return "docx", nil
	case "pdf":
		return "pdf", nil
	case "odt":
		return "odt", nil
	case "txt":
		return "txt:Text (encoded):UTF8", nil
	case "html":
		return "html:XHTML Writer File:UTF8", nil
	case "csv":
		return "csv:Text - txt - csv (StarCalc)", nil
	
	// Presentations
	case "ppt":
		return "ppt:MS PowerPoint 97", nil
	case "pptx":
		return "pptx:MS PowerPoint 2007 XML", nil
	case "odp":
		return "odp:impress8", nil
	
	default:
		return "", fmt.Errorf("unsupported document format: %s", ext)
	}
}
