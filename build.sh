#!/bin/sh
set -e

VERSION=$1
if [ -z "$VERSION" ]; then
  echo "Version not provided"
  exit 1
fi

echo "Building version $VERSION"

rm -f terraform-provider-slack*.zip

GOOS=linux GOARCH=amd64 go build -o "terraform-provider-slack_v${VERSION}_linux_amd64"
GOOS=windows GOARCH=amd64 go build -o "terraform-provider-slack_v${VERSION}_windows_amd64.exe"
GOOS=darwin GOARCH=amd64 go build -o "terraform-provider-slack_v${VERSION}_darwin_amd64"

zip "terraform-provider-slack_v${VERSION}_linux_amd64.zip" "terraform-provider-slack_v${VERSION}_linux_amd64"
zip "terraform-provider-slack_v${VERSION}_windows_amd64.zip" "terraform-provider-slack_v${VERSION}_windows_amd64.exe"
zip "terraform-provider-slack_v${VERSION}_darwin_amd64.zip" "terraform-provider-slack_v${VERSION}_darwin_amd64"

rm terraform-provider-slack_v*_*