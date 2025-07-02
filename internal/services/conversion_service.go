package services

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/amannvl/freefileconverterz/internal/config"
)

// ConversionService handles file conversion operations
type ConversionService struct {
	fileService    *FileService
	tempDir        string
	libreOfficePath string
}

// NewConversionService creates a new ConversionService instance
func NewConversionService(fileService *FileService, cfg *config.StorageConfig, libreOfficePath string) (*ConversionService, error) {
	// Create temp directory if it doesn't exist
	tempDir := filepath.Join(os.TempDir(), "freefileconverterz")
	if err := os.MkdirAll(tempDir, 0755); err != nil {
		return nil, fmt.Errorf("failed to create temp directory: %w", err)
	}

	// If no custom path provided, use "libreoffice" as default
	if libreOfficePath == "" {
		libreOfficePath = "libreoffice"
	}

	return &ConversionService{
		fileService:    fileService,
		tempDir:        tempDir,
		libreOfficePath: libreOfficePath,
	}, nil
}

// ConvertFile converts a file from one format to another
func (s *ConversionService) ConvertFile(
	ctx context.Context,
	fileHeader *multipart.FileHeader,
	targetFormat string,
	userID string,
) (string, error) {
	// Upload the original file
	_, err := s.fileService.UploadFile(ctx, fileHeader, userID)
	if err != nil {
		return "", fmt.Errorf("failed to upload file: %w", err)
	}

	// Generate output filename
	ext := filepath.Ext(fileHeader.Filename)
	outputFilename := strings.TrimSuffix(fileHeader.Filename, ext) + "." + targetFormat

	// Create a temporary directory for this conversion
	conversionDir := filepath.Join(s.tempDir, fmt.Sprintf("%s_%d", userID, time.Now().UnixNano()))
	if err := os.MkdirAll(conversionDir, 0755); err != nil {
		return "", fmt.Errorf("failed to create conversion directory: %w", err)
	}
	defer os.RemoveAll(conversionDir)

	// Download the file to temp directory
	sourceFile := filepath.Join(conversionDir, fileHeader.Filename)
	destFile := filepath.Join(conversionDir, outputFilename)

	// Get the file from storage
	// Note: We're using the uploaded file directly from the fileHeader
	file, err := fileHeader.Open()
	if err != nil {
		return "", fmt.Errorf("failed to open uploaded file: %w", err)
	}
	defer file.Close()

	// Save to temp file
	out, err := os.Create(sourceFile)
	if err != nil {
		return "", fmt.Errorf("failed to create temp file: %w", err)
	}
	defer out.Close()

	if _, err = io.Copy(out, file); err != nil {
		return "", fmt.Errorf("failed to save file: %w", err)
	}

	// Convert the file based on file type
	switch strings.ToLower(targetFormat) {
	case "pdf":
		err = s.convertToPDF(sourceFile, destFile)
	case "docx":
		err = s.convertToDocx(sourceFile, destFile)
	case "jpg", "jpeg", "png":
		err = s.convertToImage(sourceFile, destFile)
	default:
		return "", errors.New("unsupported target format")
	}

	if err != nil {
		return "", fmt.Errorf("conversion failed: %w", err)
	}

	// Upload the converted file
	convertedFile, err := os.Open(destFile)
	if err != nil {
		return "", fmt.Errorf("failed to open converted file: %w", err)
	}
	defer convertedFile.Close()

	// Save the converted file
	destPath := filepath.Join(userID, "converted", outputFilename)
	path, err := s.fileService.storage.Save(ctx, destPath, convertedFile)
	if err != nil {
		return "", fmt.Errorf("failed to save converted file: %w", err)
	}

	return path, nil
}

// convertToPDF converts a file to PDF using LibreOffice
func (s *ConversionService) convertToPDF(sourceFile, destFile string) error {
	log.Printf("Starting PDF conversion: %s -> %s", sourceFile, destFile)
	
	// Check if LibreOffice is installed
	if _, err := exec.LookPath(s.libreOfficePath); err != nil {
		errMsg := fmt.Sprintf("LibreOffice is required for PDF conversion (%s not found): %v", s.libreOfficePath, err)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Ensure source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("Source file does not exist: %s", sourceFile)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(destFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errMsg := fmt.Sprintf("Failed to create output directory %s: %v", outputDir, err)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	log.Printf("Running command: %s --headless --convert-to pdf --outdir %s %s", 
		s.libreOfficePath, outputDir, sourceFile)

	cmd := exec.Command(
		s.libreOfficePath,
		"--headless",
		"--convert-to", "pdf",
		"--outdir", outputDir,
		sourceFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("PDF conversion failed: %s", string(output))
	}

	return nil
}

// convertToDocx converts a file to DOCX using LibreOffice
func (s *ConversionService) convertToDocx(sourceFile, destFile string) error {
	log.Printf("Starting DOCX conversion: %s -> %s", sourceFile, destFile)
	
	// Check if LibreOffice is installed
	if _, err := exec.LookPath(s.libreOfficePath); err != nil {
		errMsg := fmt.Sprintf("LibreOffice is required for DOCX conversion (%s not found): %v", s.libreOfficePath, err)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Ensure source file exists
	if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
		errMsg := fmt.Sprintf("Source file does not exist: %s", sourceFile)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	// Ensure output directory exists
	outputDir := filepath.Dir(destFile)
	if err := os.MkdirAll(outputDir, 0755); err != nil {
		errMsg := fmt.Sprintf("Failed to create output directory %s: %v", outputDir, err)
		log.Printf("Error: %s", errMsg)
		return fmt.Errorf(errMsg)
	}

	log.Printf("Running command: %s --headless --convert-to docx --outdir %s %s", 
		s.libreOfficePath, outputDir, sourceFile)

	cmd := exec.Command(
		s.libreOfficePath,
		"--headless",
		"--convert-to", "docx",
		"--outdir", outputDir,
		sourceFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("DOCX conversion failed: %s", string(output))
	}

	return nil
}

// convertToImage converts a file to an image using ImageMagick
func (s *ConversionService) convertToImage(sourceFile, destFile string) error {
	// Check if ImageMagick is installed
	if _, err := exec.LookPath("convert"); err != nil {
		return errors.New("ImageMagick is required for image conversion")
	}

	cmd := exec.Command(
		"convert",
		sourceFile + "[0]", // First page for multi-page documents
		destFile,
	)

	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("image conversion failed: %s", string(output))
	}

	return nil
}
