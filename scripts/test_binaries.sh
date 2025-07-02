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

BIN_DIR="$(dirname "$0")/../bin/${OS}/${ARCH}"

# Function to check if a command exists
command_exists() {
    command -v "$1" >/dev/null 2>&1
}

# Function to check if a binary exists in the local bin directory
local_binary_exists() {
    local bin_name=$1
    local local_bin_path="${BIN_DIR}/${bin_name}"
    
    if [ -f "$local_bin_path" ] && [ -x "$local_bin_path" ]; then
        return 0
    fi
    return 1
}

# Function to get binary version
get_binary_version() {
    local bin_path=$1
    local bin_name=$(basename "$bin_path")
    local version_cmd=""
    
    case $bin_name in
        ffmpeg|ffprobe)
            version_cmd="$bin_path -version | head -n 1"
            ;;
        convert|magick)
            version_cmd="$bin_path --version | head -n 1"
            ;;
        soffice)
            version_cmd="$bin_path --version | head -n 1"
            ;;
        unrar|unar)
            version_cmd="$bin_path | head -n 2 | tail -n 1"
            ;;
        7z)
            version_cmd="$bin_path | head -n 3 | tail -n 1"
            ;;
        *)
            version_cmd="$bin_path --version 2>&1 | head -n 1 || echo 'Version not available'"
            ;;
    esac
    
    eval "$version_cmd" 2>/dev/null || echo "Version check failed"
}

# Check required binaries (adjust based on OS)
if [ "$OS" = "darwin" ]; then
    BINARIES=("ffmpeg" "convert" "soffice" "unar" "7z")
else
    BINARIES=("ffmpeg" "convert" "soffice" "unrar" "7z")
fi

echo -e "🔍 Checking for required binaries on ${OS} ${ARCH}...\n"

ALL_AVAILABLE=true
for bin in "${BINARIES[@]}"; do
    echo -n "${bin}: "
    
    if local_binary_exists "$bin"; then
        echo -e "${GREEN}✓ Found in local bin directory${NC}"
        bin_path="${BIN_DIR}/${bin}"
    elif command_exists "$bin"; then
        echo -e "${GREEN}✓ Found in system PATH${NC}"
        bin_path="$(command -v "$bin")"
    else
        echo -e "${RED}✗ Not found${NC}"
        ALL_AVAILABLE=false
        echo
        continue
    fi
    
    # Get and display version
    version=$(get_binary_version "$bin_path")
    echo "  Version: $version"
    
    # Check if the binary is executable
    if [ -x "$bin_path" ]; then
        echo -e "  Executable: ${GREEN}Yes${NC}"
    else
        echo -e "  Executable: ${RED}No${NC} (check permissions)"
        ALL_AVAILABLE=false
    fi
    
    echo
done

# Check if all binaries are available
if $ALL_AVAILABLE; then
    echo -e "${GREEN}✅ All required binaries are available!${NC}"
    exit 0
else
    echo -e "${YELLOW}⚠️  Some required binaries are missing. You can install them using:${NC}"
    echo -e "  make setup-binaries  # Install all required binaries locally"
    echo -e "  # Or install them manually:"
    
    if [ "$OS" = "darwin" ]; then
        echo -e "  # On macOS (using Homebrew):"
        echo -e "  brew install ffmpeg imagemagick p7zip unar"
        echo -e "  # For LibreOffice (soffice), download from: https://www.libreoffice.org/download/download/"
    else
        echo -e "  # On Linux (Debian/Ubuntu):"
        echo -e "  sudo apt-get update && sudo apt-get install -y ffmpeg imagemagick p7zip-full unrar"
        echo -e "  # For LibreOffice (soffice):"
        echo -e "  sudo apt-get install -y libreoffice"
        echo -e "  "
        echo -e "  # On Linux (RHEL/CentOS):"
        echo -e "  # sudo yum install -y ffmpeg ImageMagick p7zip p7zip-plugins unrar"
        echo -e "  # sudo yum install -y libreoffice"
    fi
    
    echo -e "\nAfter installation, you can run this script again to verify the installation."
    exit 1
fi
