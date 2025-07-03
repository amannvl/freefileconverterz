#!/bin/bash

# Get input and output paths from arguments
input_path="$1"
output_dir="$2"
output_file="$3"

# Change to the directory containing the input file
cd "$(dirname "$input_path")"

# Run the conversion
libreoffice --headless --convert-to docx --outdir "$output_dir" "$(basename "$input_path")"

# Check if the output file was created
output_path="${output_dir}/$(basename "${input_path%.*}").docx"
if [ -f "$output_path" ]; then
    # If the output filename doesn't match the expected name, rename it
    if [ "$output_path" != "$output_dir/$output_file" ]; then
        mv "$output_path" "$output_dir/$output_file"
    fi
    exit 0
else
    echo "Conversion failed - output file not found"
    exit 1
fi
