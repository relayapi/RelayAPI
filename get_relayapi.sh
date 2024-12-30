#!/bin/bash

# Detect system type and architecture
OS=$(uname -s | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

# Convert architecture names
case $ARCH in
    x86_64)
        ARCH="amd64"
        ;;
    aarch64)
        ARCH="arm64"
        ;;
    armv7l)
        ARCH="arm"
        ;;
esac

# Version and download URL
VERSION="v1.0.0"
DOWNLOAD_URL="https://github.com/relayapi/RelayAPI/releases/download/${VERSION}/relayapi-${OS}-${ARCH}.tar.gz"
INSTALL_DIR="/usr/local/relayapi"
EXTRACT_DIR="relayapi-${OS}-${ARCH}"

# Create temporary directory
TMP_DIR=$(mktemp -d)
cd $TMP_DIR

echo "📦 Downloading RelayAPI..."
if ! curl -fsSL $DOWNLOAD_URL -o relayapi.tar.gz; then
    echo "❌ Download failed"
    exit 1
fi

echo "📂 Extracting files..."
tar -xzf relayapi.tar.gz
echo "✅ Files extracted to: $PWD/$EXTRACT_DIR"

# Ask for installation
read -p "Do you want to install RelayAPI to $INSTALL_DIR? [y/N] " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🔧 Installing RelayAPI..."
    # Create installation directory
    sudo mkdir -p $INSTALL_DIR
    sudo cp -r $EXTRACT_DIR/* $INSTALL_DIR/

    # Create symlink
    sudo ln -sf $INSTALL_DIR/relayapi-server /usr/local/bin/relayapi-server

    echo "✅ RelayAPI installed successfully!"
    echo "📝 Config file location: $INSTALL_DIR/config.json"
    echo "🚀 Start service: relayapi-server"
else
    echo "📁 Files are available in: $PWD/$EXTRACT_DIR"
    echo "💡 You can manually install later by copying files to your preferred location"
fi

# Clean up temporary files
cd - > /dev/null
rm -rf $TMP_DIR 