# Contributing to FreeFileConverterZ

First off, thank you for considering contributing to FreeFileConverterZ! It's people like you that make FreeFileConverterZ such a great tool.

## Code of Conduct

This project and everyone participating in it is governed by our [Code of Conduct](CODE_OF_CONDUCT.md). By participating, you are expected to uphold this code. Please report any unacceptable behavior to the project maintainers.

## How Can I Contribute?

### Reporting Bugs

- **Ensure the bug was not already reported** by searching on GitHub under [Issues](https://github.com/amannvl/freefileconverterz/issues).
- If you're unable to find an open issue addressing the problem, [open a new one](https://github.com/amannvl/freefileconverterz/issues/new). Be sure to include:
  - A clear and descriptive title
  - A description of the expected behavior and the observed behavior
  - Steps to reproduce the issue
  - Any relevant error messages
  - Your environment (OS, browser, etc.)

### Suggesting Enhancements

- Use GitHub Issues to submit enhancement suggestions
- Clearly describe the enhancement and why you believe it would be useful
- Include any relevant screenshots or examples

### Your First Code Contribution

1. **Fork the repository**
2. **Create a new branch** for your feature or bug fix
   ```
   git checkout -b feature/amazing-feature
   ```
3. **Make your changes**
4. **Run tests** to ensure everything works
5. **Commit your changes** with a descriptive message
   ```
   git commit -m "Add amazing feature"
   ```
6. **Push to your fork**
   ```
   git push origin feature/amazing-feature
   ```
7. **Open a Pull Request**

### Code Style

- Follow the existing code style
- Run `go fmt` before committing
- Ensure all tests pass
- Add tests for new features
- Update documentation as needed

## Development Setup

### Prerequisites

- Go 1.21+
- Node.js 18+
- Docker and Docker Compose (optional)
- FFmpeg, ImageMagick, LibreOffice, 7-Zip, and Calibre for file conversion

### Getting Started

1. **Clone the repository**
   ```
   git clone https://github.com/amannvl/freefileconverterz.git
   cd freefileconverterz
   ```

2. **Set up environment variables**
   ```
   cp .env.example .env
   # Edit .env with your configuration
   ```

3. **Install dependencies**
   ```
   # Install Go dependencies
   go mod download
   
   # Install Node.js dependencies
   npm install
   
   # Build frontend assets
   npm run build
   ```

4. **Start the development server**
   ```
   # Using Docker Compose (recommended)
   docker-compose up -d
   
   # Or run locally
   go run cmd/server/main.go
   ```

5. **Visit** http://localhost:3000

## Testing

Run the test suite:

```bash
# Run all tests
go test ./...

# Run tests with coverage
go test -coverprofile=coverage.out ./...

# View coverage report
go tool cover -html=coverage.out
```

## Pull Request Process

1. Ensure any install or build dependencies are removed before the end of the layer when doing a build.
2. Update the README.md with details of changes to the interface, this includes new environment variables, exposed ports, useful file locations, and container parameters.
3. Increase the version numbers in any example files and the README.md to the new version that this Pull Request would represent. The versioning scheme we use is [SemVer](http://semver.org/).
4. You may merge the Pull Request in once you have the sign-off of two other developers, or if you do not have permission to do that, you may request the second reviewer to merge it for you.

## License

By contributing, you agree that your contributions will be licensed under its MIT License.
