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
VERSION="v1.0.7"
DOWNLOAD_URL="https://github.com/relayapi/RelayAPI/releases/download/${VERSION}/relayapi-${OS}-${ARCH}.tar.gz"

# Default paths
DEFAULT_DOWNLOAD_DIR="$PWD/relayapi"
INSTALL_DIR="/usr/local/relayapi"

# Ask if user wants to change download directory
echo "📂 Default download directory: $DEFAULT_DOWNLOAD_DIR"
read -p "Do you want to change the download directory? [y/N] " -n 1 -r
echo
DOWNLOAD_DIR="$DEFAULT_DOWNLOAD_DIR"
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "Enter new download directory:"
    read -r NEW_DIR
    if [ ! -z "$NEW_DIR" ]; then
        DOWNLOAD_DIR="$NEW_DIR"
    fi
fi

# Create download directory
mkdir -p "$DOWNLOAD_DIR"
cd "$DOWNLOAD_DIR" || exit 1

echo "📦 Downloading RelayAPI to $DOWNLOAD_DIR ..."
if ! curl -fsSL $DOWNLOAD_URL -o relayapi.tar.gz; then
    echo "❌ Download failed"
    exit 1
fi

echo "📂 Extracting files..."
tar -xzf relayapi.tar.gz
EXTRACT_DIR="relayapi-${OS}-${ARCH}"
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
    echo "🚀 Start service with: relayapi-server"
else
    echo "📁 Files are available in: $PWD/$EXTRACT_DIR"
    echo "💡 You can manually install later by copying files to your preferred location"
fi

# Clean up
cd - > /dev/null
if [[ $REPLY =~ ^[Yy]$ ]]; then
    echo "🧹 Do you want to remove downloaded files? [y/N] "
    read -n 1 -r
    echo
    if [[ $REPLY =~ ^[Yy]$ ]]; then
        rm -rf "$DOWNLOAD_DIR"
        echo "✅ Temporary files cleaned"
    fi
fi 