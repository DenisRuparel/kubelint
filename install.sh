#!/bin/bash

set -e

BINARY_NAME="kubelint"
DOWNLOAD_URL="https://github.com/DenisRuparel/kubelint/releases/latest/download/kubelint-linux-amd64"

echo "Installing latest KubeLint release..."

# Download binary
curl -fsSL "${DOWNLOAD_URL}" -o "${BINARY_NAME}"

# Make executable
chmod +x "${BINARY_NAME}"

# Move to system path
sudo mv "${BINARY_NAME}" /usr/local/bin/${BINARY_NAME}

echo ""
echo "KubeLint installed successfully!"
echo ""

# Verify installation
kubelint version