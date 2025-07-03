# FreeFileConverterZ

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/amannvl/freefileconverterz)](https://goreportcard.com/report/github.com/amannvl/freefileconverterz)
[![Docker Pulls](https://img.shields.io/docker/pulls/amannvl/freefileconverterz)](https://hub.docker.com/r/amannvl/freefileconverterz)

FreeFileConverterZ is a high-performance, web-based file conversion platform that enables users to convert between a wide variety of document, image, audio, and video formats. Built with Go and React, it offers a modern, responsive interface with a robust backend.

## 🚀 Features

- **Multiple Format Support**: Convert between various document, image, audio, and video formats
- **High Performance**: Built with Go for fast file processing
- **Modern Web Interface**: Responsive React frontend with drag-and-drop support
- **RESTful API**: Programmatic access to conversion services
- **Docker Support**: Easy deployment with Docker and Docker Compose
- **Cross-Platform**: Works on Windows, macOS, and Linux

## 🛠️ Tech Stack

- **Backend**: Go 1.21+
- **Frontend**: React, TypeScript, TailwindCSS
- **API**: RESTful with Gin framework
- **Storage**: Local filesystem or S3-compatible storage
- **Containerization**: Docker
- **Logging**: Zerolog for structured logging

## 📋 Prerequisites

- Go 1.21+
- Node.js 18+ (for frontend development)
- Docker and Docker Compose (for containerized deployment)
- Required system tools (handled automatically in Docker):
  - FFmpeg (for audio/video conversion)
  - ImageMagick (for image processing)
  - LibreOffice (for document conversion)
  - p7zip, unrar (for archive handling)

## 🚀 Quick Start

### Development Setup

1. **Clone the repository**:
   ```bash
   git clone https://github.com/amannvl/freefileconverterz.git
   cd freefileconverterz
   ```

2. **Set up environment variables**:
   ```bash
   cp .env.example .env
   # Edit .env as needed
   ```

3. **Install Go dependencies**:
   ```bash
   go mod download
   ```

4. **Build and run the backend**:
   ```bash
   go run main.go
   ```

5. **Set up the frontend**:
   ```bash
   cd frontend
   npm install
   npm start
   ```

The application will be available at `http://localhost:3000`

### Production Deployment

1. **Using Docker Compose (recommended)**:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d --build
   ```

2. **Using the binary**:
   ```bash
   # Build the application
   make build
   
   # Run the server
   ./bin/freefileconverterz
   ```

## 🔧 Configuration

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment (development, production) | `development` |
| `PORT` | Port to listen on | `8080` |
| `UPLOAD_DIR` | Directory to store uploaded files | `./uploads` |
| `TEMP_DIR` | Directory for temporary files | `./temp` |
| `MAX_UPLOAD_SIZE` | Maximum upload size in bytes | `104857600` (100MB) |
| `JWT_SECRET` | Secret key for JWT authentication | Randomly generated |
| `CORS_ALLOWED_ORIGINS` | Allowed CORS origins | `*` |

### File Storage

The application can be configured to use either local filesystem or S3-compatible storage by setting the appropriate environment variables.

## 📚 API Documentation

The API documentation is available at `/api/docs` when running in development mode.

### Endpoints

- `GET /api/v1/health` - Health check
- `GET /api/v1/formats` - List supported formats
- `POST /api/v1/convert` - Convert a file

## 🧪 Testing

Run the test suite:

```bash
go test -v ./...
```

## 🤝 Contributing

Contributions are welcome! Please read our [Contributing Guidelines](CONTRIBUTING.md) for details.

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
## Project Structure

```
freefileconverterz/
├── cmd/                  # Application entry points
│   └── server/           # Main server application
├── internal/             # Private application code
│   ├── config/           # Configuration management
│   ├── handlers/         # HTTP request handlers
│   ├── middleware/       # HTTP middleware
│   ├── services/         # Business logic
│   └── utils/            # Utility functions
├── pkg/                  # Reusable packages
│   ├── converter/        # File conversion logic
│   │   ├── document/     # Document converters
│   │   ├── image/        # Image converters
│   │   ├── audio/        # Audio converters
│   │   ├── video/        # Video converters
│   │   └── archive/      # Archive converters
│   └── storage/          # File storage abstraction
├── static/               # Static files (CSS, JS, images)
│   ├── css/
│   ├── js/
│   └── img/
├── uploads/              # Temporary file uploads
├── views/                # HTML templates
│   ├── layouts/          # Base templates
│   ├── partials/         # Reusable template components
│   └── *.html            # Page templates
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── Dockerfile           # Docker configuration
├── docker-compose.yml   # Docker Compose configuration
└── README.md           # Project documentation
```

## API Documentation

The API is stateless and requires no authentication for basic file conversion. For advanced features, you may need to include an API key in the request header:

```
X-API-Key: your-api-key
```

### Endpoints

#### Public Endpoints

- `POST /api/v1/convert` - Convert a file
  ```
  Content-Type: multipart/form-data
  
  file: [file to convert]
  format: [target format]
  ```

- `GET /api/v1/status/:id` - Check conversion status
- `GET /download/:id` - Download a converted file

#### Admin Endpoints (Require Authentication)

- `GET /admin/dashboard` - Admin dashboard
- `GET /admin/conversions` - List all conversions
- `GET /admin/users` - List all users
- `GET /admin/settings` - Get system settings
- `PUT /admin/settings` - Update system settings

## Contributing

1. Fork the repository
2. Create a new branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- [Fiber](https://gofiber.io/) - Fast HTTP framework for Go
- [Tailwind CSS](https://tailwindcss.com/) - Utility-first CSS framework
- [Font Awesome](https://fontawesome.com/) - Icon library
- [FFmpeg](https://ffmpeg.org/) - Multimedia framework
- [LibreOffice](https://www.libreoffice.org/) - Office suite
- [7-Zip](https://www.7-zip.org/) - File archiver
- [Calibre](https://calibre-ebook.com/) - E-book management
