# Build stage
FROM golang:1.21-alpine AS builder

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
FROM alpine:latest

# Install runtime dependencies
RUN apk --no-cache add \
    ffmpeg \
    imagemagick \
    p7zip \
    unar \
    libreoffice \
    ttf-freefont

# Create non-root user
RUN addgroup -S appgroup && adduser -S appuser -G appgroup

# Create necessary directories
RUN mkdir -p /app/uploads && \
    chown -R appuser:appgroup /app

WORKDIR /app

# Copy the binary from builder
COPY --from=builder /freefileconverterz /app/
COPY --chown=appuser:appgroup ./uploads /app/uploads

# Set environment variables
ENV UPLOAD_PATH=/app/uploads
ENV MAX_UPLOAD_SIZE=104857600
ENV APP_ENV=production \
    PORT=8080 \
    UPLOAD_DIR=/app/uploads

# Switch to non-root user
USER appuser

# Expose port
EXPOSE 8080

# Run the application
CMD ["/app/freefileconverterz"]
