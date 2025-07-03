#!/bin/bash

# Stop and remove any existing containers
docker rm -f freefileconverterz 2>/dev/null

# Build the Docker image
echo "Building the Docker image..."
docker build -t freefileconverterz -f Dockerfile.combined .

# Run the container
echo "Starting the application..."
docker run -d \
  --name freefileconverterz \
  -p 80:80 \
  -p 8080:8080 \
  -v $(pwd)/uploads:/app/uploads \
  -v $(pwd)/temp:/app/temp \
  --restart unless-stopped \
  freefileconverterz

echo "\nApplication is running!"
echo "- Frontend: http://localhost"
echo "- Backend API: http://localhost:8080"
echo "- API Documentation: http://localhost/api-docs"

echo "\nTo view logs, run: docker logs -f freefileconverterz"
echo "To stop the application: docker stop freefileconverterz"
