#!/bin/sh
set -e

# This script is called by semantic-release to build and package the provider.

# Get the version from the command line argument.
VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Usage: ./build.sh <version>"
  exit 1
fi

# Clean up any previous build artifacts.
rm -f terraform-provider-slack*.zip SHA256SUMS SHA256SUMS.sig

# A function to build and zip for a given OS and architecture.
build_and_zip() {
  OS=$1
  ARCH=$2
  echo "Building for $OS/$ARCH..."

  # Set the output binary name.
  BINARY_NAME="terraform-provider-slack_v${VERSION}_${OS}_${ARCH}"
  if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
  fi

  # Build the binary.
  GOOS=$OS GOARCH=$ARCH go build -o "$BINARY_NAME"

  # Create the zip file.
  zip "terraform-provider-slack_v${VERSION}_${OS}_${ARCH}.zip" "$BINARY_NAME"

  # Clean up the binary.
  rm "$BINARY_NAME"
}

# Build for all target platforms.
build_and_zip linux amd64
build_and_zip linux arm64
build_and_zip windows amd64
build_and_zip windows arm64
build_and_zip darwin amd64
build_and_zip darwin arm64

# Generate the SHA256SUMS file.
echo "Generating SHA256SUMS..."
sha256sum terraform-provider-slack_v*.zip > SHA256SUMS
