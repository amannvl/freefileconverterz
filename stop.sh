#!/bin/bash

# Configuration - should match run.sh
CONTAINER_NAME="freefileconverterz"

# Colors for better output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Check if the container is running
if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
    echo -e "${YELLOW}Stopping and removing ${CONTAINER_NAME}...${NC}"
    # Stop the container
    docker stop $CONTAINER_NAME 2>/dev/null
    # Remove the container
    docker rm $CONTAINER_NAME 2>/dev/null
    echo -e "${GREEN}✅ ${CONTAINER_NAME} has been stopped and removed.${NC}
"
    
    # Show remaining containers (for debugging)
    echo -e "${YELLOW}Remaining Docker containers:${NC}"
    docker ps -a --format "{{.Names}} ({{.Status}})" | grep -v "^${CONTAINER_NAME}$" || echo "No other containers found."
else
    echo -e "${YELLOW}Container ${CONTAINER_NAME} is not running.${NC}"
    # Check if the container exists but is stopped
    if docker ps -a --format '{{.Names}}' | grep -q "^${CONTAINER_NAME}$"; then
        echo -e "${YELLOW}Removing stopped container ${CONTAINER_NAME}...${NC}"
        docker rm $CONTAINER_NAME 2>/dev/null
        echo -e "${GREEN}✅ ${CONTAINER_NAME} has been removed.${NC}"
    fi
fi

echo -e "\n${YELLOW}To start the application again, run:${NC} ${GREEN}./run.sh${NC}"
echo -e "${YELLOW}To view all containers, run:${NC} ${GREEN}docker ps -a${NC}"
