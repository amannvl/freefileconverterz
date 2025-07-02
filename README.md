# FreeFileConverterZ

[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![Go Report Card](https://goreportcard.com/badge/github.com/amannvl/freefileconverterz)](https://goreportcard.com/report/github.com/amannvl/freefileconverterz)
[![Docker Pulls](https://img.shields.io/docker/pulls/amannvl/freefileconverterz)](https://hub.docker.com/r/amannvl/freefileconverterz)

FreeFileConverterZ is a comprehensive, web-based file conversion platform that enables users to convert between a wide variety of document, image, audio, video, archive, and specialized file formats. The application prioritizes ease of use, speed, and security while handling file conversions in the cloud.

## 🚀 Features

- **Multiple Format Support**: Convert between 100+ file formats including documents, images, audio, and video
- **Fast Conversions**: Utilize high-performance backend services for quick file processing
- **Stateless Architecture**: No database required - simple and easy to deploy
- **Secure & Private**: Files are automatically deleted after conversion
- **No Registration Required**: Start converting files immediately without creating an account
- **Responsive Design**: Works on desktop, tablet, and mobile devices
- **Docker Ready**: Easy deployment with Docker Compose
- **RESTful API**: Programmatic access to conversion services

## Supported Formats

### Documents
- **Word**: DOC, DOCX, ODT, RTF, TXT
- **PDF**: PDF to Word, Excel, PowerPoint, Images, and more
- **Excel**: XLS, XLSX, CSV, ODS
- **PowerPoint**: PPT, PPTX, ODP
- **E-books**: EPUB, MOBI, AZW, FB2

### Images
- **Raster**: JPG, PNG, GIF, WEBP, BMP, TIFF, HEIC
- **Vector**: SVG, AI, EPS, PDF

### Audio
- MP3, WAV, AAC, FLAC, ALAC, AIFF, WMA, OGG, M4A, OPUS

### Video
- MP4, AVI, MOV, MKV, WMV, FLV, WEBM, 3GP, MTS, M2TS

### Archives
- ZIP, RAR, 7Z, TAR, TAR.GZ, TAR.BZ2, TAR.XZ, ISO

## 🚀 Quick Start

### Prerequisites

- Docker and Docker Compose (for containerized deployment)
- Go 1.21+ and Node.js 18+ (for development)
- For Linux systems, the following system packages are recommended but not required (will be handled automatically):
   # On macOS (using Homebrew)
   brew install unrar p7zip ffmpeg imagemagick libreoffice
   ```

2. **Local binaries** (recommended for production):
   ```bash
   # Download and set up all required binaries in ./bin/linux/amd64/
   make setup-binaries
   
   # The application will automatically use these binaries if they exist
   ```

### Building for Production

```bash
# Build the application
make build

# The binary will be available at ./bin/freefileconverterz
```

## 🐳 Production Deployment

### Prerequisites

- Docker 20.10+
- Docker Compose 2.0+

### Steps

1. Clone the repository:
   ```bash
   git clone https://github.com/amannvl/freefileconverterz.git
   cd freefileconverterz
   ```

2. Build and start the production stack:
   ```bash
   docker-compose -f docker-compose.prod.yml up -d
   ```

3. The application will be available at http://localhost:3000

### Environment Variables

You can customize the application behavior using these environment variables:

| Variable | Description | Default |
|----------|-------------|---------|
| `APP_ENV` | Application environment (development, production) | `production` |
| `PORT` | Port to listen on | `3000` |
| `UPLOAD_DIR` | Directory to store uploaded files | `/app/uploads` |
| `MAX_UPLOAD_SIZE` | Maximum file upload size in bytes | `104857600` (100MB) |
| `FILE_RETENTION` | How long to keep converted files | `1h` |
| `LOG_LEVEL` | Logging level (debug, info, warn, error) | `info` |
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
