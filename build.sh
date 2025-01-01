#!/bin/bash

# Exit immediately if a command exits with a non-zero status
set -e

# Version file
VERSION_FILE="version.txt"

# Load current version from file or create if not exists
if [ -f "$VERSION_FILE" ]; then
    CURRENT_VERSION=$(cat "$VERSION_FILE")
else
    CURRENT_VERSION="v1.0.0"
    echo $CURRENT_VERSION > "$VERSION_FILE"
fi

# Calculate next version (but don't write it yet)
MAJOR=$(echo $CURRENT_VERSION | cut -d. -f1)
MINOR=$(echo $CURRENT_VERSION | cut -d. -f2)
PATCH=$(echo $CURRENT_VERSION | cut -d. -f3)
PATCH=$((PATCH + 1))
NEXT_VERSION="${MAJOR}.${MINOR}.${PATCH}"

echo "📦 Current version: $CURRENT_VERSION"
echo "🔄 Next version will be: $NEXT_VERSION"

CURRENT_VERSION=$NEXT_VERSION

# Load GitHub token from .env
if [ ! -f ".env" ]; then
    echo "❌ Error: .env file not found"
    exit 1
fi
source .env
if [ -z "$GITHUB_TOKEN" ]; then
    echo "❌ Error: GITHUB_TOKEN not found in .env"
    exit 1
fi

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

# 更新版本号
echo "📝 Updating version number in main.go..."
if [[ "$OSTYPE" == "darwin"* ]]; then
    # macOS
    sed -i '' "s/const Version = \".*\"/const Version = \"$CURRENT_VERSION\"/" cmd/server/main.go
else
    # Linux
    sed -i "s/const Version = \".*\"/const Version = \"$CURRENT_VERSION\"/" cmd/server/main.go
fi

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

        # For Linux systems, include the service file
        if [ "$OS" = "linux" ]; then
            cp ../relayapi.service.linux.sh $PACKAGE_DIR/
            chmod +x $PACKAGE_DIR/relayapi.service.linux.sh
        fi

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

# Create GitHub release
echo "📤 Creating GitHub release for version $CURRENT_VERSION..."

# Create git tag
git tag $CURRENT_VERSION
git push 

# Create GitHub release
RELEASE_NOTES="RelayAPI Release $CURRENT_VERSION"
RELEASE_FILES=()
for OS in "${SYSTEMS[@]}"; do
    for ARCH in "${ARCHITECTURES[@]}"; do
        RELEASE_FILES+=("$RELEASE_DIR/relayapi-${OS}-${ARCH}.tar.gz")
    done
done

# Create release using GitHub API
RELEASE_DATA="{\"tag_name\":\"$CURRENT_VERSION\",\"name\":\"Release $CURRENT_VERSION\",\"body\":\"$RELEASE_NOTES\",\"draft\":false,\"prerelease\":false}"
RELEASE_RESPONSE=$(curl -X POST -H "Authorization: token $GITHUB_TOKEN" \
    -H "Accept: application/vnd.github.v3+json" \
    -d "$RELEASE_DATA" \
    "https://api.github.com/repos/relayapi/RelayAPI/releases")

# Get release ID from response
RELEASE_ID=$(echo $RELEASE_RESPONSE | jq -r .id)

# Upload release assets
for FILE in "${RELEASE_FILES[@]}"; do
    FILENAME=$(basename $FILE)
    echo "Uploading $FILENAME..."
    curl -X POST -H "Authorization: token $GITHUB_TOKEN" \
        -H "Content-Type: application/octet-stream" \
        --data-binary @"$FILE" \
        "https://uploads.github.com/repos/relayapi/RelayAPI/releases/$RELEASE_ID/assets?name=$FILENAME"
done

# All steps completed successfully, now update version numbers
echo "📝 Updating version numbers..."

# Update version.txt with next version
echo $NEXT_VERSION > "$VERSION_FILE"

# Update version in get_relayapi.sh
sed -i '' "s/VERSION=\".*\"/VERSION=\"$CURRENT_VERSION\"/" get_relayapi.sh

echo "✨ Release $CURRENT_VERSION published successfully!"
echo "🔄 Version updated to $NEXT_VERSION" 