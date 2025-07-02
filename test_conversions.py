import os
import http.client
import json
import time
import mimetypes
from urllib.parse import urlencode, urlparse

BASE_URL = "http://localhost:3001/api/v1"
TEST_FILES_DIR = "test_files"

def test_health():
    print("Testing health check...")
    conn = http.client.HTTPConnection("localhost", 3001)
    try:
        conn.request("GET", "/api/v1/health", headers={"Accept": "application/json"})
        response = conn.getresponse()
        data = response.read().decode()
        print(f"Status: {response.status}")
        
        # Try to parse as JSON, fall back to raw response if not JSON
        try:
            json_data = json.loads(data)
            print(f"Response: {json.dumps(json_data, indent=2)}")
        except json.JSONDecodeError:
            print(f"Response (not JSON): {data[:500]}...")
            
        print("-" * 50)
    except Exception as e:
        print(f"Error: {str(e)}")
    finally:
        conn.close()

def test_file_conversion(input_file, output_format):
    print(f"Testing {input_file} to {output_format} conversion...")
    
    # Prepare the file upload
    file_path = os.path.join(TEST_FILES_DIR, input_file)
    file_name = os.path.basename(file_path)
    
    # Read file content
    with open(file_path, 'rb') as f:
        file_content = f.read()
    
    # Create multipart form data
    boundary = '----WebKitFormBoundary7MA4YWxkTrZu0gW'
    body = []
    
    # Add file
    body.append(f'--{boundary}')
    body.append(f'Content-Disposition: form-data; name="file"; filename="{file_name}"')
    body.append('Content-Type: application/octet-stream')
    body.append('')
    body.append(file_content.decode('latin-1'))
    
    # Add target format (using 'format' parameter as per the API)
    body.append(f'--{boundary}')
    body.append('Content-Disposition: form-data; name="format"')
    body.append('')
    body.append(output_format)
    
    body.append(f'--{boundary}--')
    body.append('')
    
    # Join body parts with CRLF
    body = '\r\n'.join(body)
    
    # Send request
    conn = http.client.HTTPConnection("localhost", 3001)
    headers = {
        'Content-Type': f'multipart/form-data; boundary={boundary}',
        'Content-Length': str(len(body.encode('latin-1')))
    }
    
    try:
        # Add Accept header to request JSON
        headers['Accept'] = 'application/json'
        conn.request("POST", "/api/v1/convert", body=body.encode('latin-1'), headers=headers)
        response = conn.getresponse()
        response_data = response.read().decode()
        
        if response.status != 200:
            print(f"Error: {response.status} - {response_data}")
            return
        
        try:
            result = json.loads(response_data)
            task_id = result.get('id')
            if not task_id:
                print(f"Unexpected response format: {response_data}")
                return
                
            print(f"Conversion started. Task ID: {task_id}")
            
            # Check status
            max_attempts = 10
            for attempt in range(max_attempts):
                time.sleep(2)  # Wait for 2 seconds between status checks
                
                status_conn = http.client.HTTPConnection("localhost", 3001)
                status_conn.request("GET", f"/api/v1/convert/{task_id}/status")
                status_response = status_conn.getresponse()
                
                if status_response.status != 200:
                    print(f"Error checking status: {status_response.status} - {status_response.read().decode()}")
                    break
                    
                status_data = json.loads(status_response.read().decode())
                status_conn.close()
                
                print(f"Status: {status_data.get('status')} (Attempt {attempt + 1}/{max_attempts})")
                
                if status_data.get('status') == 'completed':
                    download_url = status_data.get('download_url')
                    if download_url:
                        print(f"Conversion successful! Download URL: {download_url}")
                        # Parse the URL to get the path
                        parsed_url = urlparse(download_url)
                        file_path = parsed_url.path
                        
                        # Download the file
                        download_conn = http.client.HTTPConnection("localhost", 3001)
                        download_conn.request("GET", file_path)
                        download_response = download_conn.getresponse()
                        
                        if download_response.status == 200:
                            output_file = os.path.join(TEST_FILES_DIR, f"converted_{task_id}.{output_format}")
                            with open(output_file, 'wb') as f:
                                f.write(download_response.read())
                            print(f"File saved as: {output_file}")
                        else:
                            print(f"Failed to download file: {download_response.status} - {download_response.read().decode()}")
                        
                        download_conn.close()
                    break
                    
                if status_data.get('status') == 'failed':
                    print(f"Conversion failed: {status_data.get('error', 'Unknown error')}")
                    break
                    
            else:
                print(f"Conversion did not complete after {max_attempts} attempts")
                
        except json.JSONDecodeError as e:
            print(f"Failed to parse response as JSON: {e}")
            print(f"Response: {response_data}")
            
    except Exception as e:
        print(f"Error during conversion: {str(e)}")
    finally:
        conn.close()
    
    print("-" * 50)

def main():
    # Test health check
    test_health()
    
    # Test different conversions with available test files
    test_cases = [
        # Image conversions
        ("test.jpg", "png"),
        ("test.png", "jpg"),
        ("test.jpg", "webp"),
        
        # Text file (test basic file handling)
        ("test.txt", "pdf"),
        
        # Add more test cases as needed
    ]
    
    for input_file, output_format in test_cases:
        if os.path.exists(os.path.join(TEST_FILES_DIR, input_file)):
            test_file_conversion(input_file, output_format)
        else:
            print(f"Skipping {input_file} -> {output_format} (file not found)")

if __name__ == "__main__":
    main()
