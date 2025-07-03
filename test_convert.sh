#!/bin/bash

# Copy test file to container
docker cp test.pdf freefileconverterz-app-1:/app/test.pdf

# Run conversion inside container
docker exec freefileconverterz-app-1 bash -c "cd /app && libreoffice --headless --convert-to docx --outdir /app/uploads /app/test.pdf"

# Check output
docker exec freefileconverterz-app-1 ls -la /app/uploads/
