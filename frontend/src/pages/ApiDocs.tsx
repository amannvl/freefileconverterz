const ApiDocs = () => {
  return (
    <div className="prose max-w-5xl mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">FreeFileConverterZ API Documentation</h1>
      
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Introduction</h2>
        <p className="text-gray-600 mb-6">
          The FreeFileConverterZ API provides programmatic access to file conversion capabilities.
          This documentation covers all available endpoints, request/response formats, and authentication.
        </p>
        
        <div className="bg-blue-50 border-l-4 border-blue-400 p-4 mb-4">
          <p className="text-blue-700">
            <strong>Base URL:</strong> <code>http://localhost:4000/api/v1</code>
          </p>
          <p className="text-blue-700 mt-2">
            <strong>Note:</strong> In production, replace <code>http://localhost:4000</code> with your actual domain.
          </p>
        </div>
      </div>

      {/* Authentication */}
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Authentication</h2>
        <p className="text-gray-600 mb-4">
          Most endpoints require authentication using JWT tokens. Include the token in the <code>Authorization</code> header:
        </p>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4">
          Authorization: Bearer YOUR_JWT_TOKEN
        </div>
        
        <h3 className="text-xl font-semibold text-gray-800 mb-3">Register</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4">
          <div className="text-purple-600">POST /register</div>
        </div>
        
        <h4 className="font-medium text-gray-800 mt-4 mb-2">Request Body</h4>
        <pre className="bg-gray-100 p-4 rounded-md text-sm mb-4 overflow-x-auto">
{`{
  "email": "user@example.com",
  "password": "securepassword123",
  "name": "John Doe"
}`}
        </pre>
        
        <h4 className="font-medium text-gray-800 mt-4 mb-2">Response (201 Created)</h4>
        <pre className="bg-gray-100 p-4 rounded-md text-sm mb-6 overflow-x-auto">
{`{
  "success": true,
  "message": "User registered successfully",
  "data": {
    "id": "123",
    "email": "user@example.com",
    "name": "John Doe"
  }
}`}
        </pre>
        
        <h3 className="text-xl font-semibold text-gray-800 mb-3">Login</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4">
          <div className="text-purple-600">POST /login</div>
        </div>
        
        <h4 className="font-medium text-gray-800 mt-4 mb-2">Request Body</h4>
        <pre className="bg-gray-100 p-4 rounded-md text-sm mb-4 overflow-x-auto">
{`{
  "email": "user@example.com",
  "password": "securepassword123"
}`}
        </pre>
        
        <h4 className="font-medium text-gray-800 mt-4 mb-2">Response (200 OK)</h4>
        <pre className="bg-gray-100 p-4 rounded-md text-sm mb-6 overflow-x-auto">
{`{
  "success": true,
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "refresh_token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."
}`}
        </pre>
      </div>

      {/* File Conversion */}
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-6">File Conversion</h2>
        
        <div className="mb-8">
          <div className="flex items-center mb-2">
            <span className="bg-green-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">POST</span>
            <h3 className="text-xl font-semibold text-gray-800">/convert</h3>
            <span className="ml-2 bg-yellow-100 text-yellow-800 text-xs font-medium px-2.5 py-0.5 rounded">Requires Auth</span>
          </div>
          
          <p className="text-gray-600 mb-4">
            Convert a file to a different format. Supported formats include:
          </p>
          <ul className="list-disc pl-6 mb-4 text-gray-600 space-y-1">
            <li><strong>Documents:</strong> PDF, DOCX, TXT, RTF, ODT</li>
            <li><strong>Images:</strong> JPG, PNG, GIF, WEBP, TIFF</li>
            <li><strong>Presentations:</strong> PPTX, ODP</li>
            <li><strong>Spreadsheets:</strong> XLSX, CSV, ODS</li>
          </ul>

          <h4 className="font-medium text-gray-800 mt-4 mb-2">Request Headers</h4>
          <div className="border rounded-md overflow-hidden mb-4">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Value</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Required</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">Content-Type</td>
                  <td className="px-4 py-2 text-sm text-gray-500">multipart/form-data</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Yes</td>
                </tr>
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">Authorization</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Bearer &lt;token&gt;</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Yes</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h4 className="font-medium text-gray-800 mt-4 mb-2">Request Body (Form Data)</h4>
          <div className="border rounded-md overflow-hidden mb-4">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Required</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">file</td>
                  <td className="px-4 py-2 text-sm text-gray-500">File</td>
                  <td className="px-4 py-2 text-sm text-gray-500">The file to convert (max 50MB)</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Yes</td>
                </tr>
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">format</td>
                  <td className="px-4 py-2 text-sm text-gray-500">String</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Target format (e.g., 'docx', 'pdf', 'jpg')</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Yes</td>
                </tr>
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">options</td>
                  <td className="px-4 py-2 text-sm text-gray-500">JSON String</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Conversion options (format-specific)</td>
                  <td className="px-4 py-2 text-sm text-gray-500">No</td>
                </tr>
              </tbody>
            </table>
          </div>

          <h4 className="font-medium text-gray-800 mt-4 mb-2">Success Response (202 Accepted)</h4>
          <p className="text-gray-600 text-sm mb-2">
            The conversion request has been accepted and is being processed.
          </p>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`{
  "success": true,
  "message": "Conversion started",
  "data": {
    "id": "conversion_123",
    "status": "processing",
    "originalName": "document.pdf",
    "targetFormat": "docx",
    "createdAt": "2023-10-01T12:00:00Z",
    "statusUrl": "/api/v1/convert/conversion_123/status",
    "downloadUrl": "/api/v1/convert/conversion_123/download"
  }
}`}</pre>
          </div>

          <h4 className="font-medium text-gray-800 mt-4 mb-2">Error Responses</h4>
          <div className="space-y-4 mb-6">
            <div>
              <p className="text-gray-600 text-sm mb-1"><strong>400 Bad Request</strong> - Invalid request</p>
              <div className="bg-gray-100 p-4 rounded-md font-mono text-sm overflow-x-auto">
                <pre>{`{
  "success": false,
  "error": {
    "code": "invalid_request",
    "message": "Missing required field: file"
  }
}`}</pre>
              </div>
            </div>
            
            <div>
              <p className="text-gray-600 text-sm mb-1"><strong>415 Unsupported Media Type</strong> - Unsupported file format</p>
              <div className="bg-gray-100 p-4 rounded-md font-mono text-sm overflow-x-auto">
                <pre>{`{
  "success": false,
  "error": {
    "code": "unsupported_format",
    "message": "Unsupported file format: .exe"
  }
}`}</pre>
              </div>
            </div>

            <div>
              <p className="text-gray-600 text-sm mb-1"><strong>429 Too Many Requests</strong> - Rate limit exceeded</p>
              <div className="bg-gray-100 p-4 rounded-md font-mono text-sm overflow-x-auto">
                <pre>{`{
  "success": false,
  "error": {
    "code": "rate_limit_exceeded",
    "message": "Too many requests, please try again later",
    "retryAfter": 60
  }
}`}</pre>
              </div>
            </div>
          </div>

          <h4 className="font-medium text-gray-800 mt-6 mb-2">Example cURL</h4>
          <div className="bg-gray-900 text-green-400 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`# Start a conversion
curl -X POST \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -F "file=@document.pdf" \
  -F "format=docx" \
  -F "options={\"quality\":90}" \
  http://localhost:4000/api/v1/convert`}</pre>
          </div>
        </div>

        {/* Check Conversion Status */}
        <div className="mb-8 pt-6 border-t border-gray-200">
          <div className="flex items-center mb-2">
            <span className="bg-blue-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">GET</span>
            <h3 className="text-xl font-semibold text-gray-800">/convert/&#123;id&#125;/status</h3>
          </div>
          
          <p className="text-gray-600 mb-4">
            Check the status of a file conversion. The <code>id</code> is returned in the conversion response.
          </p>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Path Parameters</h4>
          <div className="border rounded-md overflow-hidden mb-4">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">id</td>
                  <td className="px-4 py-2 text-sm text-gray-500">String</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Conversion ID</td>
                </tr>
              </tbody>
            </table>
          </div>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Success Response (200 OK)</h4>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`{
  "success": true,
  "data": {
    "id": "conversion_123",
    "status": "completed",
    "originalName": "document.pdf",
    "convertedName": "document.docx",
    "targetFormat": "docx",
    "fileSize": 12345,
    "createdAt": "2023-10-01T12:00:00Z",
    "completedAt": "2023-10-01T12:00:30Z",
    "downloadUrl": "/api/v1/convert/conversion_123/download"
  }
}`}</pre>
          </div>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Status Values</h4>
          <ul className="list-disc pl-6 mb-4 text-gray-600 space-y-1 text-sm">
            <li><code>queued</code> - Waiting to be processed</li>
            <li><code>processing</code> - Currently being converted</li>
            <li><code>completed</code> - Successfully converted</li>
            <li><code>failed</code> - Conversion failed (check error message)</li>
          </ul>
        </div>
        
        {/* Download Converted File */}
        <div className="mb-8 pt-6 border-t border-gray-200">
          <div className="flex items-center mb-2">
            <span className="bg-purple-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">GET</span>
            <h3 className="text-xl font-semibold text-gray-800">/convert/&#123;id&#125;/download</h3>
          </div>
          
          <p className="text-gray-600 mb-4">
            Download a converted file. The file will be served as a download with appropriate headers.
          </p>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Example cURL</h4>
          <div className="bg-gray-900 text-green-400 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`# Download the converted file
curl -OJ -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  http://localhost:4000/api/v1/convert/conversion_123/download`}</pre>
          </div>
          
          <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4 mb-6">
            <p className="text-yellow-700 text-sm">
              <strong>Note:</strong> The <code>-O</code> flag saves the file with its original name, 
              and <code>-J</code> honors the Content-Disposition header to set the correct filename.
            </p>
          </div>
        </div>
      </div>

      {/* User Management */}
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-6">User Management</h2>
        
        {/* Get Current User */}
        <div className="mb-8">
          <div className="flex items-center mb-2">
            <span className="bg-blue-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">GET</span>
            <h3 className="text-xl font-semibold text-gray-800">/users/me</h3>
            <span className="ml-2 bg-yellow-100 text-yellow-800 text-xs font-medium px-2.5 py-0.5 rounded">Requires Auth</span>
          </div>
          
          <p className="text-gray-600 mb-4">
            Get the currently authenticated user's profile information.
          </p>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Success Response (200 OK)</h4>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`{
  "success": true,
  "data": {
    "id": "user_123",
    "email": "user@example.com",
    "name": "John Doe",
    "createdAt": "2023-01-01T00:00:00Z",
    "lastLogin": "2023-10-01T12:00:00Z"
  }
}`}</pre>
          </div>
        </div>
        
        {/* Update User */}
        <div className="mb-8 pt-6 border-t border-gray-200">
          <div className="flex items-center mb-2">
            <span className="bg-yellow-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">PATCH</span>
            <h3 className="text-xl font-semibold text-gray-800">/users/me</h3>
            <span className="ml-2 bg-yellow-100 text-yellow-800 text-xs font-medium px-2.5 py-0.5 rounded">Requires Auth</span>
          </div>
          
          <p className="text-gray-600 mb-4">
            Update the currently authenticated user's profile information.
          </p>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Request Body</h4>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4 overflow-x-auto">
            <pre>{`{
  "name": "Updated Name",
  "currentPassword": "currentpassword123",
  "newPassword": "newsecurepassword456"
}`}</pre>
          </div>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Success Response (200 OK)</h4>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`{
  "success": true,
  "message": "Profile updated successfully"
}`}</pre>
          </div>
        </div>
      </div>

      {/* Admin Endpoints */}
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-6">Admin Endpoints</h2>
        <p className="text-gray-600 mb-4">
          These endpoints are only accessible to users with admin privileges.
        </p>
        
        {/* List All Users */}
        <div className="mb-8">
          <div className="flex items-center mb-2">
            <span className="bg-blue-500 text-white text-xs font-medium px-2.5 py-0.5 rounded mr-2">GET</span>
            <h3 className="text-xl font-semibold text-gray-800">/admin/users</h3>
            <span className="ml-2 bg-red-100 text-red-800 text-xs font-medium px-2.5 py-0.5 rounded">Admin Only</span>
          </div>
          
          <p className="text-gray-600 mb-4">
            List all users (paginated).
          </p>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Query Parameters</h4>
          <div className="border rounded-md overflow-hidden mb-4">
            <table className="min-w-full divide-y divide-gray-200">
              <thead className="bg-gray-50">
                <tr>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Name</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Type</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Default</th>
                  <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
                </tr>
              </thead>
              <tbody className="bg-white divide-y divide-gray-200">
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">page</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Integer</td>
                  <td className="px-4 py-2 text-sm text-gray-500">1</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Page number</td>
                </tr>
                <tr>
                  <td className="px-4 py-2 text-sm text-gray-900">limit</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Integer</td>
                  <td className="px-4 py-2 text-sm text-gray-500">20</td>
                  <td className="px-4 py-2 text-sm text-gray-500">Items per page (max 100)</td>
                </tr>
              </tbody>
            </table>
          </div>
          
          <h4 className="font-medium text-gray-800 mt-4 mb-2">Success Response (200 OK)</h4>
          <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
            <pre>{`{
  "success": true,
  "data": {
    "items": [
      {
        "id": "user_123",
        "email": "user@example.com",
        "name": "John Doe",
        "role": "user",
        "createdAt": "2023-01-01T00:00:00Z",
        "lastLogin": "2023-10-01T12:00:00Z"
      }
    ],
    "pagination": {
      "total": 1,
      "page": 1,
      "limit": 20,
      "totalPages": 1
    }
  }
}`}</pre>
          </div>
        </div>
      </div>

      {/* Rate Limiting & Errors */}
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Rate Limiting</h2>
        <div className="space-y-4">
          <div>
            <p className="text-gray-600 mb-2">
              The API is rate limited to prevent abuse and ensure fair usage. The current limits are:
            </p>
            <ul className="list-disc pl-6 text-gray-600 space-y-1">
              <li><strong>100 requests per minute</strong> per IP address for public endpoints</li>
              <li><strong>1000 requests per minute</strong> per authenticated user</li>
              <li><strong>10 concurrent conversions</strong> per user</li>
            </ul>
          </div>
          
          <div className="bg-yellow-50 border-l-4 border-yellow-400 p-4">
            <p className="text-yellow-700">
              <strong>Note:</strong> When rate limited, the API will return a <code>429 Too Many Requests</code> 
              response with a <code>Retry-After</code> header indicating how long to wait before making another request.
            </p>
          </div>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Error Handling</h2>
        <p className="text-gray-600 mb-4">
          The API uses standard HTTP status codes to indicate success or failure of an API request.
          In case of an error, the response will include a JSON object with an <code>error</code> field 
          containing details about what went wrong.
        </p>
        
        <h3 className="text-xl font-semibold text-gray-800 mt-6 mb-3">Common Error Codes</h3>
        <div className="border rounded-md overflow-hidden mb-6">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Code</th>
                <th className="px-4 py-2 text-left text-xs font-medium text-gray-500 uppercase">Description</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">400</td>
                <td className="px-4 py-2 text-sm text-gray-600">Bad Request - Invalid request parameters</td>
              </tr>
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">401</td>
                <td className="px-4 py-2 text-sm text-gray-600">Unauthorized - Authentication required</td>
              </tr>
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">403</td>
                <td className="px-4 py-2 text-sm text-gray-600">Forbidden - Insufficient permissions</td>
              </tr>
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">404</td>
                <td className="px-4 py-2 text-sm text-gray-600">Not Found - Resource not found</td>
              </tr>
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">429</td>
                <td className="px-4 py-2 text-sm text-gray-600">Too Many Requests - Rate limit exceeded</td>
              </tr>
              <tr>
                <td className="px-4 py-2 text-sm font-medium text-gray-900">500</td>
                <td className="px-4 py-2 text-sm text-gray-600">Internal Server Error - Something went wrong on our end</td>
              </tr>
            </tbody>
          </table>
        </div>
        
        <h3 className="text-xl font-semibold text-gray-800 mt-6 mb-3">Example Error Response</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4 overflow-x-auto">
          <pre>{`{
  "success": false,
  "error": {
    "code": "invalid_credentials",
    "message": "Invalid email or password",
    "details": {
      "field": "password",
      "suggestion": "Check your password and try again"
    },
    "requestId": "req_1234567890"
  }
}`}</pre>
        </div>
      </div>
    </div>
  );
};

export default ApiDocs;
