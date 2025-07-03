#!/bin/bash

# Set up test environment
export GO111MODULE=on

echo "=== Generating test files ==="
chmod +x ./scripts/generate_test_files.sh
./scripts/generate_test_files.sh

echo -e "\n=== Running conversion tests ==="
go test -v ./pkg/converter/ -run TestDocumentConversions

# Clean up test files if needed
# rm -rf testdata/
