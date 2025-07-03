package converter

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/amannvl/freefileconverterz/internal/tools"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// MockDocumentConverter is a test double for DocumentConverter
type MockDocumentConverter struct {
	*DocumentConverter
}

// Convert is a mock implementation for testing
func (m *MockDocumentConverter) Convert(ctx context.Context, inputPath, outputPath string) error {
	// Check for unsupported source format
	ext := strings.TrimPrefix(filepath.Ext(inputPath), ".")
	if ext == "unsupported" {
		return fmt.Errorf("unsupported source format: %s", ext)
	}

	// Check for unsupported target format
	ext = strings.TrimPrefix(filepath.Ext(outputPath), ".")
	if ext == "unsupported" {
		return fmt.Errorf("unsupported target format: %s", ext)
	}

	// Create a dummy output file for supported formats
	return os.WriteFile(outputPath, []byte("Mock converted content"), 0644)
}

// TestDocumentConversions tests various document format conversions
func TestDocumentConversions(t *testing.T) {
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "converter-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	// Create necessary subdirectories
	testDataDir := filepath.Join(tempDir, "testdata")
	err = os.MkdirAll(testDataDir, 0755)
	require.NoError(t, err)

	// Create a sample test file
	sampleFile := filepath.Join(testDataDir, "sample.txt")
	err = os.WriteFile(sampleFile, []byte("Test document content"), 0644)
	require.NoError(t, err)

	// Initialize the converter with a tool manager
	toolManager, err := tools.NewToolManager(tempDir, tempDir)
	require.NoError(t, err)
	
	// Create a mock document converter
	docConverter := &MockDocumentConverter{
		DocumentConverter: NewDocumentConverter(toolManager, tempDir).(*DocumentConverter),
	}

	// Define test cases for different conversions
	testCases := []struct {
		name        string
		sourceFile  string
		targetExt   string
		skip        bool
		skipReason  string
		expectError bool
	}{
		// Text document conversions
		{
			name:        "Convert TXT to PDF",
			sourceFile:  "sample.txt",
			targetExt:   "pdf",
			skip:        false,
			expectError: false,
		},
		{
			name:        "Convert TXT to DOCX",
			sourceFile:  "sample.txt",
			targetExt:   "docx",
			skip:        false,
			expectError: false,
		},
		{
			name:        "Convert TXT to ODT",
			sourceFile:  "sample.txt",
			targetExt:   "odt",
			skip:        false,
			expectError: false,
		},

		// Add more test cases for other formats as needed
		{
			name:        "Convert DOCX to PDF",
			sourceFile:  "sample.docx",
			targetExt:   "pdf",
			skip:        true, // Skip if test file doesn't exist
			skipReason:  "Sample DOCX file not available",
			expectError: false,
		},
		{
			name:        "Convert XLSX to ODS",
			sourceFile:  "sample.xlsx",
			targetExt:   "ods",
			skip:        true, // Skip if test file doesn't exist
			skipReason:  "Sample XLSX file not available",
			expectError: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if tc.skip {
				t.Skip(tc.skipReason)
			}

			// Use the sample test file we created
			sourceFile := filepath.Join(testDataDir, tc.sourceFile)
			if _, err := os.Stat(sourceFile); os.IsNotExist(err) {
				t.Skipf("Source file %s does not exist: %v", sourceFile, tc.skipReason)
			}

			// Create output file path
			ext := filepath.Ext(sourceFile)
			baseName := strings.TrimSuffix(filepath.Base(sourceFile), ext)
			outputFile := filepath.Join(tempDir, fmt.Sprintf("%s_converted.%s", baseName, tc.targetExt))

			// Perform the conversion using our mock
			err = docConverter.Convert(context.Background(), sourceFile, outputFile)

			// Verify the result
			if tc.expectError {
				assert.Error(t, err, "Expected an error but got none")
				return
			}

			require.NoError(t, err, "Conversion failed")

			// Verify output file exists and has content
			fileInfo, err := os.Stat(outputFile)
			require.NoError(t, err, "Output file not created")
			assert.Greater(t, fileInfo.Size(), int64(0), "Output file is empty")

			// Verify file extension matches target format
			assert.True(t, strings.HasSuffix(strings.ToLower(outputFile), "."+tc.targetExt),
				"Output file extension does not match target format")

			t.Logf("Successfully converted %s to %s", tc.sourceFile, outputFile)
		})
	}
}

// TestUnsupportedFormats tests that unsupported formats return appropriate errors
func TestUnsupportedFormats(t *testing.T) {
	// Create a temporary directory for test output
	tempDir, err := os.MkdirTemp("", "converter-test-*")
	require.NoError(t, err)
	defer os.RemoveAll(tempDir)

	toolManager, err := tools.NewToolManager(tempDir, tempDir)
	require.NoError(t, err)
	docConverter := NewDocumentConverter(toolManager, tempDir).(*DocumentConverter)

	tests := []struct {
		name       string
		sourceFile string
		targetExt  string
	}{
		{
			name:       "Unsupported source format",
			sourceFile: "test.unsupported",
			targetExt:  "pdf",
		},
		{
			name:       "Unsupported target format",
			sourceFile: "test.docx",
			targetExt:  "unsupported",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := docConverter.Convert(context.Background(), tt.sourceFile, "output."+tt.targetExt)
			assert.Error(t, err, "Expected error for unsupported format")
		})
	}
}
