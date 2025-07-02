#!/bin/bash

set -euo pipefail

# Colors for output
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
RED='\033[0;31m'
NC='\033[0m' # No Color

# Detect OS and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Map architecture to standard names
case "$ARCH" in
    x86_64)  ARCH="amd64" ;;
    aarch64) ARCH="arm64" ;;
    arm*)    ARCH="arm" ;;
    *)       ARCH="unknown" ;;
esac

# Configuration
TEMP_DIR="/tmp/freefileconverterz"
BIN_DIR="$(pwd)/bin/${OS}/${ARCH}"
MAX_RETRIES=3

# Create necessary directories
mkdir -p "$BIN_DIR"
mkdir -p "$TEMP_DIR"

# Function to print status messages
info() {
    echo -e "${GREEN}[INFO]${NC} $1"
}

warn() {
    echo -e "${YELLOW}[WARN]${NC} $1"
}

error() {
    echo -e "${RED}[ERROR]${NC} $1" >&2
    exit 1
}

# Function to download a file with retries
download_with_retry() {
    local url=$1
    local output=$2
    local retries=$MAX_RETRIES
    
    while [ $retries -gt 0 ]; do
        if curl -L "$url" -o "$output" --progress-bar --connect-timeout 30 --max-time 300; then
            return 0
        fi
        retries=$((retries - 1))
        warn "Download failed, $retries retries remaining..."
        sleep 5
    done
    
    error "Failed to download $url after $MAX_RETRIES attempts"
    return 1
}

# Function to extract archives
extract_archive() {
    local file=$1
    local output_dir=$2
    
    mkdir -p "$output_dir"
    
    case $file in
        *.tar.gz|*.tgz)
            tar -xzf "$file" -C "$output_dir" --strip-components=1
            ;;
        *.tar.xz)
            tar -xJf "$file" -C "$output_dir" --strip-components=1
            ;;
        *.zip)
            unzip -q "$file" -d "$output_dir"
            ;;
        *)
            error "Unsupported archive format: $file"
            ;;
    esac
}

# Function to setup binaries based on OS
setup_binaries() {
    info "Detected OS: $OS, Architecture: $ARCH"
    
    case "$OS" in
        linux)
            setup_linux_binaries
            ;;
        darwin)
            setup_macos_binaries
            ;;
        *)
            error "Unsupported operating system: $OS"
            ;;
    esac
}

# Setup Linux binaries
setup_linux_binaries() {
    info "Setting up Linux binaries..."
    
    # Install system dependencies if needed
    if ! command -v curl &> /dev/null || ! command -v tar &> /dev/null; then
        info "Installing required system packages..."
        if command -v apt-get &> /dev/null; then
            sudo apt-get update && sudo apt-get install -y curl tar xz-utils unzip
        elif command -v dnf &> /dev/null; then
            sudo dnf install -y curl tar xz unzip
        elif command -v yum &> /dev/null; then
            sudo yum install -y curl tar xz unzip
        else
            warn "Could not install system dependencies. Some features might not work."
        fi
    fi
    
    # Download and setup FFmpeg
    if [ ! -f "$BIN_DIR/ffmpeg" ]; then
        info "Downloading FFmpeg..."
        FFMPEG_URL="https://github.com/BtbN/FFmpeg-Builds/releases/download/latest/ffmpeg-master-latest-${OS}64-gpl.tar.xz"
        download_with_retry "$FFMPEG_URL" "$TEMP_DIR/ffmpeg.tar.xz"
        extract_archive "$TEMP_DIR/ffmpeg.tar.xz" "$TEMP_DIR/ffmpeg"
        cp "$TEMP_DIR/ffmpeg/bin/ffmpeg" "$BIN_DIR/"
        cp "$TEMP_DIR/ffmpeg/bin/ffprobe" "$BIN_DIR/"
    fi
    
    # Download and setup ImageMagick
    if [ ! -f "$BIN_DIR/convert" ]; then
        info "Downloading ImageMagick..."
        IMAGEMAGICK_URL="https://download.imagemagick.org/ImageMagick/download/binaries/ImageMagick-${ARCH}-pc-linux-gnu.tar.xz"
        download_with_retry "$IMAGEMAGICK_URL" "$TEMP_DIR/imagemagick.tar.xz"
        extract_archive "$TEMP_DIR/imagemagick.tar.xz" "$TEMP_DIR/imagemagick"
        cp "$TEMP_DIR/imagemagick/convert" "$BIN_DIR/"
        cp "$TEMP_DIR/imagemagick/magick" "$BIN_DIR/"
    fi
    
    # Download and setup 7-Zip
    if [ ! -f "$BIN_DIR/7z" ]; then
        info "Downloading 7-Zip..."
        if [ "$ARCH" = "amd64" ]; then
            SEVENZIP_URL="https://www.7-zip.org/a/7z2301-linux-x64.tar.xz"
        else
            SEVENZIP_URL="https://www.7-zip.org/a/7z2301-linux-arm64.tar.xz"
        fi
        download_with_retry "$SEVENZIP_URL" "$TEMP_DIR/7z.tar.xz"
        extract_archive "$TEMP_DIR/7z.tar.xz" "$TEMP_DIR/7z"
        cp "$TEMP_DIR/7z/7zz" "$BIN_DIR/7z"
    fi
    
    # Download and setup UnRAR
    if [ ! -f "$BIN_DIR/unrar" ]; then
        info "Downloading UnRAR..."
        if [ "$ARCH" = "amd64" ]; then
            UNRAR_URL="https://www.rarlab.com/rar/rarlinux-x64-623.tar.gz"
        else
            UNRAR_URL="https://www.rarlab.com/rar/rarlinux-arm-623.tar.gz"
        fi
        download_with_retry "$UNRAR_URL" "$TEMP_DIR/unrar.tar.gz"
        extract_archive "$TEMP_DIR/unrar.tar.gz" "$TEMP_DIR/unrar"
        cp "$TEMP_DIR/unrar/rar/unrar" "$BIN_DIR/"
    fi
    
    # Download and setup LibreOffice
    if [ ! -f "$BIN_DIR/soffice" ]; then
        info "Downloading LibreOffice..."
        if [ "$ARCH" = "amd64" ]; then
            LO_URL="https://download.documentfoundation.org/libreoffice/stable/7.5.8/linux/x86_64/LibreOffice_7.5.8_Linux_x86-64_rpm_langpack_en-GB.tar.gz"
        else
            LO_URL="https://download.documentfoundation.org/libreoffice/stable/7.5.8/linux/aarch64/LibreOffice_7.5.8_Linux_aarch64_rpm_langpack_en-GB.tar.gz"
        fi
        download_with_retry "$LO_URL" "$TEMP_DIR/lo.tar.gz"
        extract_archive "$TEMP_DIR/lo.tar.gz" "$TEMP_DIR/lo"
        find "$TEMP_DIR/lo" -name "LibreOffice*.AppImage" -exec cp {} "$BIN_DIR/soffice" \;
        chmod +x "$BIN_DIR/soffice"
    fi
}

# Setup macOS binaries
setup_macos_binaries() {
    info "Setting up macOS binaries..."
    
    # Check if Homebrew is installed
    if ! command -v brew &> /dev/null; then
        info "Installing Homebrew..."
        /bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"
        eval "$(/opt/homebrew/bin/brew shellenv)"
    fi
    
    # Install required packages
    info "Installing required packages via Homebrew..."
    brew install ffmpeg imagemagick p7zip unar
    
    # Create symlinks in our bin directory
    mkdir -p "$BIN_DIR"
    ln -sf "$(which ffmpeg)" "$BIN_DIR/ffmpeg"
    ln -sf "$(which convert)" "$BIN_DIR/convert"
    ln -sf "$(which 7z)" "$BIN_DIR/7z"
    ln -sf "$(which unar)" "$BIN_DIR/unar"
    ln -sf "$(which lsar)" "$BIN_DIR/lsar"
    ln -sf "$(which soffice)" "$BIN_DIR/soffice" 2>/dev/null || warn "LibreOffice not found. Please install it manually."
    
    info "macOS binaries have been set up using Homebrew"
}

# Main execution
main() {
    info "Starting binary setup for $OS $ARCH..."
    
    # Setup binaries based on OS
    setup_binaries
    
    # Make all binaries executable
    chmod +x "$BIN_DIR/"*
    
    info "✅ All binaries have been installed to $BIN_DIR"
    info "Add this directory to your PATH or use the full path to the binaries."
    
    # Clean up
    rm -rf "$TEMP_DIR"
}

# Run the main function
main "$@"
