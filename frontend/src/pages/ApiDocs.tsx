const ApiDocs = () => {
  return (
    <div className="prose max-w-4xl mx-auto py-8 px-4">
      <h1 className="text-3xl font-bold text-gray-900 mb-8">API Documentation</h1>
      
      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Introduction</h2>
        <p className="text-gray-600 mb-6">
          The FreeFileConverterZ API allows you to convert files between different formats programmatically.
          Below you'll find the documentation for the available endpoints.
        </p>
      </div>

      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Base URL</h2>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4">
          http://localhost:4000/api/v1
        </div>
        <p className="text-gray-600 mb-4">
          Note: In production, replace <code>http://localhost:4000</code> with your actual domain.
        </p>
      </div>

      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Convert File</h2>
        <p className="text-gray-600 mb-4">
          Convert a file to a different format.
        </p>
        
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4 overflow-x-auto">
          <div className="text-purple-600">POST</div>
          <div>/convert</div>
        </div>

        <h3 className="text-lg font-medium text-gray-800 mt-6 mb-2">Request Headers</h3>
        <div className="border rounded-md overflow-hidden mb-6">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Required</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">Content-Type</td>
                <td className="px-6 py-4 text-sm text-gray-500">multipart/form-data</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Yes</td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3 className="text-lg font-medium text-gray-800 mt-6 mb-2">Request Body (Form Data)</h3>
        <div className="border rounded-md overflow-hidden mb-6">
          <table className="min-w-full divide-y divide-gray-200">
            <thead className="bg-gray-50">
              <tr>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Name</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Type</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Description</th>
                <th className="px-6 py-3 text-left text-xs font-medium text-gray-500 uppercase tracking-wider">Required</th>
              </tr>
            </thead>
            <tbody className="bg-white divide-y divide-gray-200">
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">file</td>
                <td className="px-6 py-4 text-sm text-gray-500">File</td>
                <td className="px-6 py-4 text-sm text-gray-500">The file to convert</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Yes</td>
              </tr>
              <tr>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-900">format</td>
                <td className="px-6 py-4 text-sm text-gray-500">String</td>
                <td className="px-6 py-4 text-sm text-gray-500">The target format (e.g., 'docx', 'pdf')</td>
                <td className="px-6 py-4 whitespace-nowrap text-sm text-gray-500">Yes</td>
              </tr>
            </tbody>
          </table>
        </div>

        <h3 className="text-lg font-medium text-gray-800 mt-6 mb-2">Success Response (200 OK)</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
          <pre>{`{
  "success": true,
  "message": "File converted successfully",
  "data": {
    "originalName": "example.pdf",
    "convertedName": "example.docx",
    "size": 12345,
    "downloadUrl": "/download/example.docx"
  }
}`}</pre>
        </div>

        <h3 className="text-lg font-medium text-gray-800 mt-6 mb-2">Error Response (400 Bad Request)</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-6 overflow-x-auto">
          <pre>{`{
  "success": false,
  "error": {
    "code": "invalid_file_format",
    "message": "Unsupported file format"
  }
}`}</pre>
        </div>

        <h3 className="text-lg font-medium text-gray-800 mt-6 mb-2">Example cURL</h3>
        <div className="bg-gray-100 p-4 rounded-md font-mono text-sm mb-4 overflow-x-auto">
          <pre>{`curl -X POST \
  -F "file=@document.pdf" \
  -F "format=docx" \
  http://localhost:4000/api/v1/convert -o output.docx`}</pre>
        </div>
      </div>

      <div className="bg-white shadow rounded-lg p-6 mb-8">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Rate Limiting</h2>
        <p className="text-gray-600 mb-4">
          The API is rate limited to 60 requests per minute per IP address.
        </p>
      </div>

      <div className="bg-white shadow rounded-lg p-6">
        <h2 className="text-2xl font-semibold text-gray-800 mb-4">Authentication</h2>
        <p className="text-gray-600 mb-4">
          Currently, the API does not require authentication, but this may change in the future.
        </p>
      </div>
    </div>
  );
};

export default ApiDocs;
