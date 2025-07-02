#!/bin/bash

# Configuration
BASE_URL="http://localhost:3000/api/v1"
TEST_FILE="/tmp/test_files/test_image.jpg"
OUTPUT_FORMAT="png"

# Check if test file exists
if [ ! -f "$TEST_FILE" ]; then
    echo "Error: Test file not found at $TEST_FILE"
    exit 1
fi

# Function to make API requests
make_request() {
    local method=$1
    local endpoint=$2
    local data=${3:-}
    
    echo -e "\n[$(date +'%Y-%m-%d %H:%M:%S')] $method $endpoint"
    
    if [ -z "$data" ]; then
        curl -s -X $method "$BASE_URL$endpoint"
    else
        curl -s -X $method -F "file=@$data" "$BASE_URL$endpoint?format=$OUTPUT_FORMAT&conversionType=image"
    fi
}

# Test 1: Upload file
echo "=== Testing File Upload ==="
response=$(make_request "POST" "/convert" "$TEST_FILE")
echo "Upload Response: $response"

# Extract conversion ID from the last line of the response
conversion_id=$(echo "$response" | tail -n 1 | jq -r '.id' 2>/dev/null)
if [ -z "$conversion_id" ] || [ "$conversion_id" = "null" ]; then
    echo "Error: Failed to get conversion ID from response"
    echo "Raw response: $response"
    exit 1
fi

echo "Conversion ID: $conversion_id"

# Test 2: Check status
echo -e "\n=== Testing Status Check ==="
max_attempts=10
attempt=1
status=""

while [ $attempt -le $max_attempts ]; do
    echo "Attempt $attempt of $max_attempts - Checking status..."
    response=$(make_request "GET" "/api/v1/convert/$conversion_id/status")
    echo "Status Response: $response"
    
    # Extract the JSON part from the response
    json_response=$(echo "$response" | tail -n 1)
    status=$(echo "$json_response" | jq -r '.status' 2>/dev/null)
    download_url=$(echo "$json_response" | jq -r '.download_url' 2>/dev/null)
    
    if [ -z "$status" ]; then
        echo "Error: Invalid status in response"
        echo "Raw response: $response"
        exit 1
    fi
    
    if [ "$status" = "completed" ] && [ "$download_url" != "null" ]; then
        echo "Conversion completed successfully!"
        break
    elif [ "$status" = "failed" ]; then
        echo "Error: Conversion failed"
        exit 1
    fi
    
    if [ $attempt -eq $max_attempts ]; then
        echo "Error: Max attempts reached, conversion not completed"
        exit 1
    fi
    
    sleep 2
    ((attempt++))
done

# Test 3: Download file
echo -e "\n=== Testing File Download ==="
if [ -n "$conversion_id" ]; then
    download_url="/api/v1/convert/$conversion_id/download"
    echo "Downloading from: $download_url"
    output_file="/tmp/test_files/converted.$OUTPUT_FORMAT"
    curl -s -o "$output_file" "$download_url"
    
    if [ -f "$output_file" ]; then
        file_size=$(wc -c < "$output_file")
        echo "Download successful! File size: $file_size bytes"
        echo "Output file: $output_file"
    else
        echo "Error: Failed to download file"
        exit 1
    fi
else
    echo "Error: No download URL available"
    exit 1
fi

echo -e "\n=== All tests completed successfully! ==="
