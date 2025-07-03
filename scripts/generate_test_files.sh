#!/bin/bash

# Create test data directory if it doesn't exist
mkdir -p testdata

# Create a simple text file
echo "This is a test document for FreeFileConverterZ" > testdata/sample.txt

# Create a simple DOCX file (requires pandoc)
if command -v pandoc &> /dev/null; then
    echo "Creating sample.docx..."
    echo "# Test Document\n\nThis is a test document for FreeFileConverterZ" > testdata/sample.md
    pandoc testdata/sample.md -o testdata/sample.docx
    rm testdata/sample.md
else
    echo "pandoc not found, skipping DOCX generation"
fi

echo "Test files generated in testdata/ directory"
