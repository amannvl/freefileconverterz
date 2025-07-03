#!/bin/bash
set -e
set -x

# Debug: Show system information
echo "=== System Information ==="
uname -a
cat /etc/os-release
echo "CPU Info:"
cat /proc/cpuinfo | grep "model name\|cpu cores" || true
echo "========================="

# Debug: Show current directory and contents
echo "Current directory: $(pwd)"
echo "Directory contents:"
ls -la /app/

# Verify binary exists and is executable
echo "Checking for binary..."
if [ ! -f "/app/freefileconverterz" ]; then
    echo "Error: Binary not found at /app/freefileconverterz"
    echo "Current directory: $(pwd)"
    echo "Directory contents:"
    ls -la /app/
    exit 1
fi

# Check binary type and dependencies
echo "=== Binary Information ==="
file /app/freefileconverterz
ldd /app/freefileconverterz || true
echo "========================="

# Verify binary permissions
echo "Binary permissions:"
ls -la /app/freefileconverterz

# Start nginx in the foreground
echo "Starting nginx..."
/usr/sbin/nginx -g 'daemon off;' &

# Start the backend server
echo "Starting backend server..."
cd /app

# Run the binary directly to see any errors
./freefileconverterz || {
    echo "Failed to start the application"
    echo "Trying with strace for more details..."
    strace ./freefileconverterz
    exit 1
}

# Keep the container running if the above fails
tail -f /dev/null
