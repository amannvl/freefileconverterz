#!/bin/bash
set -e

# Configuration
FRONTEND_PORT=9002
BACKEND_PORT=9001
CONTAINER_NAME="freefileconverterz"
DOCKER_IMAGE="freefileconverterz"
DOCKERFILE="Dockerfile.combined"

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if Docker is running
if ! docker info > /dev/null 2>&1; then
    echo "Error: Docker is not running. Please start Docker and try again."
    exit 1
fi

# Function to check if a port is in use
port_in_use() {
    local port=$1
    if lsof -i :"$port" > /dev/null 2>&1; then
        echo "Port $port is already in use. Please free the port and try again."
        exit 1
    fi
}

# Check if ports are available
port_in_use $FRONTEND_PORT
port_in_use $BACKEND_PORT

# Function to stop and remove existing containers
stop_existing_containers() {
    local container_name=$1
    if docker ps -a --format '{{.Names}}' | grep -q "^${container_name}$"; then
        echo -e "${YELLOW}Stopping and removing existing '${container_name}' container...${NC}"
        docker stop "$container_name" 2>/dev/null || true
        docker rm -f "$container_name" 2>/dev/null || true
    fi
}

# Stop and remove any existing containers
stop_existing_containers "$CONTAINER_NAME"

# Clean up any dangling containers from previous runs
echo -e "${YELLOW}Cleaning up any dangling containers...${NC}"
docker ps -a --filter "name=${CONTAINER_NAME}" --format "{{.ID}}" | xargs -r docker rm -f 2>/dev/null || true

# Build the Docker image
echo -e "${YELLOW}Building the Docker image (this may take a few minutes)...${NC}"
if ! docker build -t $DOCKER_IMAGE -f $DOCKERFILE .; then
    echo -e "❌ ${YELLOW}Failed to build Docker image. Please check the error messages above.${NC}"
    exit 1
fi

# Create necessary directories if they don't exist
echo -e "${YELLOW}Setting up directories...${NC}"
mkdir -p uploads temp
chmod -R 777 uploads temp 2>/dev/null || true

# Run the container
echo -e "\n${YELLOW}Starting the application...${NC}"
docker run -d \
  --name $CONTAINER_NAME \
  -p $FRONTEND_PORT:80 \
  -p $BACKEND_PORT:8080 \
  -v "$(pwd)/uploads:/app/uploads" \
  -v "$(pwd)/temp:/app/temp" \
  --restart unless-stopped \
  $DOCKER_IMAGE

# Wait a moment for the container to start
sleep 2

# Verify container is running
if ! docker ps --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "\n❌ ${YELLOW}Container failed to start. Check the logs with:${NC} docker logs $CONTAINER_NAME"
    exit 1
fi

echo -e "\n${GREEN}✅ Application is running!${NC}"
echo -e "\n${GREEN}Access the application at:${NC}"
echo -e "- Frontend: ${GREEN}http://localhost:$FRONTEND_PORT${NC}"
echo -e "- Backend API: ${GREEN}http://localhost:$BACKEND_PORT${NC}"
echo -e "- Health Check: ${GREEN}http://localhost:$BACKEND_PORT/health${NC}"

echo -e "\n${YELLOW}Useful commands:${NC}"
echo -e "- View logs: ${GREEN}docker logs -f $CONTAINER_NAME${NC}"
echo -e "- Stop application: ${GREEN}docker stop $CONTAINER_NAME${NC}"
echo -e "- Remove container: ${GREEN}docker rm -f $CONTAINER_NAME${NC}"
echo -e "- Rebuild and restart: ${GREEN}$0${NC}"

# Show initial logs
echo -e "\n${YELLOW}Showing initial logs (Ctrl+C to exit logs):${NC}"
docker logs -f --tail=20 $CONTAINER_NAME
