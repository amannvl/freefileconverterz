.PHONY: build run test clean setup-binaries test-binaries

# Build the application
build:
	@echo "Building application..."
	@go build -o bin/freefileconverterz ./cmd/api

# Run the application
run: build
	@echo "Starting application..."
	@./bin/freefileconverterz

# Run tests
test:
	@echo "Running tests..."
	@go test -v ./...

# Clean build artifacts
clean:
	@echo "Cleaning..."
	@rm -rf bin/
	@rm -rf cover.out

# Setup required binaries
setup-binaries:
	@echo "Setting up required binaries..."
	@./scripts/setup_binaries.sh

# Test if required binaries are available
test-binaries:
	@echo "Testing required binaries..."
	@./scripts/test_binaries.sh

# Run database migrations
migrate:
	@echo "Running migrations..."
	go run migrations/

# Validate dependencies and build
validate:
	@echo "🔍 Validating dependencies..."
	@(go version && echo "✓ Go is installed" || (echo "✗ Go is not installed" && exit 1))
	@(which docker >/dev/null 2>&1 && docker --version && echo "✓ Docker is installed" || (echo "✗ Docker is not installed" && exit 1))
	@(which docker-compose >/dev/null 2>&1 && docker-compose --version && echo "✓ Docker Compose is installed" || (echo "✗ Docker Compose is not installed" && exit 1))
	@(which node >/dev/null 2>&1 && node --version && echo "✓ Node.js is installed" || (echo "✗ Node.js is not installed" && exit 1))
	@(which npm >/dev/null 2>&1 && npm --version && echo "✓ npm is installed" || (echo "✗ npm is not installed" && exit 1))
	@echo "✅ All dependencies are installed"
	@echo "\n🔍 Validating Go modules..."
	@go mod verify
	@echo "\n📦 Downloading dependencies..."
	@go mod download
	@echo "\n🔨 Building the project..."
	@$(MAKE) build
	@echo "\n✅ Build completed successfully! Run 'make run-local' to start the application."

# Setup local binaries
setup-binaries:
	@echo "📦 Setting up local binaries..."
	@chmod +x scripts/setup_binaries.sh
	@./scripts/setup_binaries.sh

# Install system dependencies (fallback)
install-system-deps:
	@echo "📦 Installing required system dependencies..."
	@echo "This target is deprecated. Please use 'make setup-binaries' instead."
	@echo "Falling back to system package manager..."
	@if [ -f /etc/os-release ]; then \
		. /etc/os-release; \
		if [ "$$ID" = "ubuntu" ] || [ "$$ID" = "debian" ]; then \
			echo "Installing dependencies using apt..."; \
			sudo apt-get update && sudo apt-get install -y unrar p7zip-full ffmpeg imagemagick libreoffice; \
		elif [ "$$ID" = "fedora" ]; then \
			echo "Installing dependencies using dnf..."; \
			sudo dnf install -y unrar p7zip ffmpeg ImageMagick libreoffice; \
		elif [ "$$ID" = "centos" ] || [ "$$ID" = "rhel" ]; then \
			echo "Installing dependencies using yum..."; \
			sudo yum install -y unrar p7zip ffmpeg ImageMagick libreoffice; \
		else \
			echo "⚠️  Unsupported Linux distribution. Please use 'make setup-binaries' instead."; \
			exit 1; \
		fi; \
	elif [ "$$(uname)" = "Darwin" ]; then \
		echo "Detected macOS. Please use 'make setup-binaries' which will use Homebrew."; \
		exit 1; \
	else \
		echo "❌ Unsupported operating system. Please use 'make setup-binaries'."; \
		exit 1; \
	fi

# Install dependencies (alias for backward compatibility)
install-deps: setup-binaries

# Run the application locally with both backend and frontend
run-local: setup-binaries
	@echo "🚀 Starting FreeFileConverterZ locally..."
	@echo "\n🌐 Backend server starting on http://localhost:8080"
	@echo "📱 Frontend will be available at http://localhost:3000"
	@echo "\n🔄 Setting up environment..."
	
	# Create necessary directories
	@mkdir -p /tmp/freefileconverterz/temp/bin
	@mkdir -p /tmp/freefileconverterz/uploads
	
	# Set up binary paths
	@if [ -f "$(PWD)/bin/linux/amd64/unrar" ]; then \
		echo "🔧 Using local binaries from $(PWD)/bin/linux/amd64"; \
		export PATH="$(PWD)/bin/linux/amd64:$$PATH"; \
		export LD_LIBRARY_PATH="$(PWD)/bin/linux/amd64:$$LD_LIBRARY_PATH"; \
	fi
	
	# Start backend server
	@echo "\n🔄 Starting backend server..."
	@(cd cmd/server && go run main.go &) \
	 && sleep 3 \
	 && (cd frontend && npm run dev) \
	 || (echo "\n❌ Failed to start the application. Make sure to run 'make validate' first." && exit 1)
	
	# Cleanup on exit
	@trap 'pkill -f "go run main.go"; exit' INT TERM EXIT
	@wait

# Lint the code
lint:
	@echo "Linting..."
	golangci-lint run

# Build Docker image
docker-build:
	@echo "Building Docker image..."
	docker-compose build

# Start Docker containers
docker-run:
	@echo "Starting Docker containers..."
	docker-compose up -d

# Stop Docker containers
docker-stop:
	@echo "Stopping Docker containers..."
	docker-compose down

# Install development dependencies
setup:
	@echo "Installing development dependencies..."
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger documentation
swagger:
	@echo "Generating Swagger documentation..."
	swag init -g cmd/server/main.go
