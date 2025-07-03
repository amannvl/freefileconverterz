#!/bin/bash
set -e

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}🚀 Starting FreeFileConverterZ development setup...${NC}"

# Check for required tools
check_command() {
    if ! command -v $1 &> /dev/null; then
        echo -e "${YELLOW}⚠️  $1 is not installed. Please install it and try again.${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓ $1 is installed${NC}"
}

# Check prerequisites
echo -e "\n${YELLOW}🔍 Checking prerequisites...${NC}"
check_command go
check_command node
check_command npm
check_command docker
check_command docker-compose

# Set up environment variables
echo -e "\n${YELLOW}⚙️  Setting up environment...${NC}"
if [ ! -f .env ]; then
    echo -e "${YELLOW}📄 Creating .env file from example...${NC}"
    cp .env.example .env
    echo -e "${GREEN}✓ Created .env file${NC}"
else
    echo -e "${GREEN}✓ .env file already exists${NC}"
fi

# Install Go dependencies
echo -e "\n${YELLOW}📦 Installing Go dependencies...${NC}"
go mod download

# Install Node.js dependencies
echo -e "\n${YELLOW}📦 Installing Node.js dependencies...${NC}"
cd frontend
npm install
cd ..

# Build the frontend
echo -e "\n${YELLOW}🏗️  Building frontend...${NC}"
cd frontend
npm run build
cd ..

# Create required directories
echo -e "\n${YELLOW}📁 Creating required directories...${NC}"
mkdir -p uploads temp

# Set permissions
echo -e "\n${YELLOW}🔒 Setting up permissions...${NC}"
chmod -R 755 uploads temp

# Build the application
echo -e "\n${YELLOW}🔨 Building the application...${NC}
make build

echo -e "\n${GREEN}✨ Setup complete! You can now start the application with:${NC}"
echo -e "${GREEN}   make run${NC} - Start the application in development mode"
echo -e "${GREEN}   docker-compose up -d${NC} - Start with Docker"
echo -e "\n${YELLOW}Happy coding! 🎉${NC}"
