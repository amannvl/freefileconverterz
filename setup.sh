#!/bin/bash

# FreeFileConverterZ Setup Script
# This script helps set up the development and production environment

set -e

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Function to print section headers
section() {
    echo -e "\n${YELLOW}==> $1${NC}"
}

# Function to check if a command exists
command_exists() {
    command -v "$1" &> /dev/null
}

# Check if running as root
if [ "$(id -u)" -eq 0 ]; then
    echo -e "${RED}Error: This script should not be run as root.${NC}" >&2
    exit 1
fi

section "🚀 FreeFileConverterZ Setup"
echo "This script will help you set up FreeFileConverterZ for development or production."

# Check if Docker is installed
if ! command_exists docker; then
    echo -e "${RED}Error: Docker is not installed. Please install Docker first.${NC}" >&2
    exit 1
fi

# Check if Docker Compose is installed
if ! command_exists docker-compose; then
    echo -e "${RED}Error: Docker Compose is not installed. Please install Docker Compose.${NC}" >&2
    exit 1
fi

# Ask for environment type
read -p "Are you setting up for development or production? [dev/prod] " ENV_TYPE
ENV_TYPE=${ENV_TYPE:-dev}

if [ "$ENV_TYPE" != "dev" ] && [ "$ENV_TYPE" != "prod" ]; then
    echo -e "${RED}Error: Invalid environment type. Please choose 'dev' or 'prod'.${NC}"
    exit 1
fi

section "📂 Creating necessary directories"
mkdir -p uploads
mkdir -p bin

if [ "$ENV_TYPE" = "dev" ]; then
    # Development setup
    section "🔧 Setting up development environment"
    
    # Check for Go
    if ! command_exists go; then
        echo -e "${RED}Error: Go is not installed. Please install Go 1.21 or later.${NC}" >&2
        exit 1
    fi
    
    # Check for Node.js and npm
    if ! command_exists node || ! command_exists npm; then
        echo -e "${RED}Error: Node.js and npm are required for development.${NC}" >&2
        exit 1
    fi
    
    # Install frontend dependencies
    section "📦 Installing frontend dependencies"
    cd frontend
    npm install
    
    # Build frontend assets
    section "🔨 Building frontend assets"
    npm run build
    cd ..
    
    # Install Go dependencies
    section "📦 Installing Go dependencies"
    go mod download
    
    # Install development tools
    section "🛠️  Installing development tools"
    go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
    
    # Set up environment variables
    if [ ! -f .env ]; then
        cp .env.example .env
        echo -e "${GREEN}✓ Created .env file${NC}"
        echo -e "${YELLOW}ℹ️  Please update the .env file with your configuration${NC}"
    fi
    
    # Start the application services
    section "🚀 Starting application services"
    docker-compose up -d
    
    echo -e "\n${GREEN}✅ Development setup complete!${NC}"
    echo -e "Start the development server with: ${YELLOW}go run cmd/server/main.go${NC}"
    echo -e "Frontend dev server: ${YELLOW}cd frontend && npm run dev${NC}"
    
else
    # Production setup
    section "🚀 Setting up production environment"
    
    # Check for .env file
    if [ ! -f .env ]; then
        cp .env.example .env
        echo -e "${GREEN}✓ Created .env file${NC}"
        echo -e "${RED}⚠️  Please update the .env file with production values before continuing!${NC}"
        exit 1
    fi
    
    # Build Docker images
    section "🐳 Building Docker images"
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml build
    
    # Start services
    section "🚀 Starting services"
    docker-compose -f docker-compose.yml -f docker-compose.prod.yml up -d
    
    echo -e "\n${GREEN}✅ Production setup complete!${NC}"
    echo -e "Application is running at ${YELLOW}http://localhost:3000${NC}"
    echo -e "View logs with: ${YELLOW}docker-compose logs -f${NC}"
fi

echo -e "\n${GREEN}✨ Setup completed successfully!${NC}"

echo "✨ Setup complete!"
echo "Next steps:"
echo "1. Update the .env file with your configuration"
echo "2. Run 'make docker-run' to start the development environment"
echo "3. Visit http://localhost:3000 in your browser"
echo ""
echo "All dependencies and binaries have been automatically downloaded and set up."
echo "You can now start using FreeFileConverterZ."
echo ""
echo "Happy coding! 🚀"
