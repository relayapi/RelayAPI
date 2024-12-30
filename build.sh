#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Version number
VERSION="v1.0.0"

# Supported systems and architectures
SYSTEMS=("linux" "darwin" "windows")
ARCHITECTURES=("amd64" "arm64")

# Build directories
BUILD_DIR="build"
RELEASE_DIR="release"
SERVER_DIR="server"

# Check if server directory exists
if [ ! -d "$SERVER_DIR" ]; then
    echo "❌ Error: server directory not found"
    exit 1
fi

# Change to server directory
echo "📂 Changing to server directory..."
cd $SERVER_DIR || exit 1

# Check if config files exist
if [ ! -f "config.json" ] || [ ! -f "default.rai" ]; then
    echo "❌ Error: config.json or default.rai not found in server directory"
    exit 1
fi

# Clean old build files
echo "🧹 Cleaning build directories..."
rm -rf ../$BUILD_DIR ../$RELEASE_DIR
mkdir -p ../$BUILD_DIR ../$RELEASE_DIR

# Build for different platforms
for OS in "${SYSTEMS[@]}"; do
    for ARCH in "${ARCHITECTURES[@]}"; do
        echo "🔨 Building for $OS/$ARCH..."
        
        # Set output filename
        if [ "$OS" = "windows" ]; then
            OUTPUT="../$BUILD_DIR/relayapi-server.exe"
        else
            OUTPUT="../$BUILD_DIR/relayapi-server"
        fi

        # Build with error checking
        if ! GOOS=$OS GOARCH=$ARCH go build -o $OUTPUT ./cmd/server; then
            echo "❌ Build failed for $OS/$ARCH"
            exit 1
        fi

        # Create release package
        PACKAGE_NAME="relayapi-${OS}-${ARCH}"
        PACKAGE_DIR="../$BUILD_DIR/$PACKAGE_NAME"
        mkdir -p $PACKAGE_DIR

        # Copy files to release directory
        cp $OUTPUT $PACKAGE_DIR/
        cp config.json $PACKAGE_DIR/config.json
        cp default.rai $PACKAGE_DIR/default.rai

        # Create package
        pushd ../$BUILD_DIR > /dev/null
        if ! tar -czf "../$RELEASE_DIR/$PACKAGE_NAME.tar.gz" $PACKAGE_NAME; then
            echo "❌ Package creation failed for $OS/$ARCH"
            popd > /dev/null
            exit 1
        fi
        popd > /dev/null

        echo "✅ $OS/$ARCH build completed"
    done
done

# Return to original directory and clean up
cd ..
rm -rf $BUILD_DIR

echo "🎉 All platform builds completed!"
echo "📦 Release packages are in $RELEASE_DIR directory" 