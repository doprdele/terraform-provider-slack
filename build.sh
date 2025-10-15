#!/bin/sh
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Version argument not provided."
  exit 1
fi

echo "Building binaries for version $VERSION"

# Build for different platforms
GOOS=linux GOARCH=amd64 go build -o terraform-provider-slack_v${VERSION}_linux_amd64
GOOS=windows GOARCH=amd64 go build -o terraform-provider-slack_v${VERSION}_windows_amd64.exe
GOOS=darwin GOARCH=amd64 go build -o terraform-provider-slack_v${VERSION}_darwin_amd64

# Create zip archives
zip terraform-provider-slack_v${VERSION}_linux_amd64.zip terraform-provider-slack_v${VERSION}_linux_amd64
zip terraform-provider-slack_v${VERSION}_windows_amd64.zip terraform-provider-slack_v${VERSION}_windows_amd64.exe
zip terraform-provider-slack_v${VERSION}_darwin_amd64.zip terraform-provider-slack_v${VERSION}_darwin_amd64
