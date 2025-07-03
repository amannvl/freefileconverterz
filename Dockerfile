# Build stage
FROM golang:1.23.0-alpine AS builder

WORKDIR /app

# Install build dependencies
RUN apk add --no-cache gcc musl-dev

# Download dependencies
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Build the application
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /freefileconverterz ./cmd/api

# Final stage
FROM ubuntu:22.04

# Install runtime dependencies
RUN apt-get update && \
    DEBIAN_FRONTEND=noninteractive apt-get install -y --no-install-recommends \
    ffmpeg \
    imagemagick \
    p7zip-full \
    unrar \
    libreoffice \
    fonts-freefont-ttf \
    ghostscript \
    libsm6 \
    libxrender1 \
    libxext6 \
    libgl1-mesa-glx \
    libglib2.0-0 \
    libxcb-shm0 \
    libxcb-xfixes0 \
    libxcb-randr0 \
    libxcb-shape0 \
    libx11-xcb1 \
    libx11-6 \
    libxext6 \
    libxi6 \
    libxrender1 \
    libxrandr2 \
    libxfixes3 \
    libxcb1 \
    libx11-xcb-dev \
    libxcb1-dev \
    libx11-dev \
    libxcb-keysyms1-dev \
    libxcb-util0-dev \
    libxcb-image0-dev \
    libxcb-shm0-dev \
    libxcb-icccm4-dev \
    libxcb-sync-dev \
    libxcb-xfixes0-dev \
    libxcb-shape0-dev \
    libxcb-randr0-dev \
    libxcb-render-util0-dev \
    libxcb-xinerama0 \
    libxcb-xinerama0-dev \
    libxcb-xkb-dev \
    libxkbcommon-dev \
    libxkbcommon-x11-dev \
    && rm -rf /var/lib/apt/lists/*

# Create non-root user and group with specific UID/GID
RUN groupadd -r -g 999 appgroup && \
    useradd -r -u 999 -g appgroup -m -d /home/appuser -s /bin/bash appuser && \
    mkdir -p /app/uploads && \
    mkdir -p /home/appuser/.config/libreoffice/4/ && \
    chown -R appuser:appgroup /app /home/appuser && \
    chmod 777 /app/uploads && \
    chmod 755 /home/appuser /home/appuser/.config /home/appuser/.config/libreoffice /home/appuser/.config/libreoffice/4

# Set the working directory
WORKDIR /app

# Switch to non-root user
USER appuser

# Set environment variables for LibreOffice
ENV HOME=/home/appuser \
    USER=appuser \
    UID=999 \
    GID=999 \
    LANG=en_US.UTF-8 \
    LANGUAGE=en_US:en \
    LC_ALL=en_US.UTF-8

# Copy the binary from builder
COPY --from=builder /freefileconverterz /app/
COPY --chown=appuser:appgroup ./uploads /app/uploads

# Set environment variables
ENV UPLOAD_PATH=/app/uploads
ENV MAX_UPLOAD_SIZE=104857600
ENV APP_ENV=production \
    PORT=8080 \
    UPLOAD_DIR=/app/uploads

# Create LibreOffice config directory and set permissions
RUN mkdir -p /home/appuser/.config/libreoffice/4/ \
    && chown -R appuser:appgroup /home/appuser/ \
    && chmod -R 755 /home/appuser/

# Switch to non-root user
USER appuser

# Set HOME environment variable
ENV HOME=/home/appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/freefileconverterz"]
