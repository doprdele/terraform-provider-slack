#!/bin/sh
set -eu

# Build and package the Terraform provider for release.

VERSION=${1:-}
if [ -z "$VERSION" ]; then
  echo "Usage: ./build.sh <version>" >&2
  exit 1
fi

if [ -z "${GPG_FINGERPRINT:-}" ]; then
  echo "GPG_FINGERPRINT must be set to sign release artifacts." >&2
  exit 1
fi

if ! command -v gpg >/dev/null 2>&1; then
  echo "gpg is required to sign release artifacts." >&2
  exit 1
fi

# Clean up any previous build artifacts.
rm -f terraform-provider-slack_v*.zip SHA256SUMS SHA256SUMS.sig

build_and_zip() {
  OS=$1
  ARCH=$2

  echo "Building for ${OS}/${ARCH}..."

  BINARY_NAME="terraform-provider-slack_v${VERSION}_${OS}_${ARCH}"
  if [ "$OS" = "windows" ]; then
    BINARY_NAME="${BINARY_NAME}.exe"
  fi

  CGO_ENABLED=0 GOOS="$OS" GOARCH="$ARCH" \
    go build -trimpath -ldflags="-s -w" -o "$BINARY_NAME" .

  zip -9 "terraform-provider-slack_v${VERSION}_${OS}_${ARCH}.zip" "$BINARY_NAME"

  rm -f "$BINARY_NAME"
}

build_and_zip linux amd64
build_and_zip linux arm64
build_and_zip windows amd64
build_and_zip windows arm64
build_and_zip darwin amd64
build_and_zip darwin arm64

echo "Generating SHA256SUMS..."
sha256sum terraform-provider-slack_v*.zip > SHA256SUMS

echo "Signing SHA256SUMS..."
gpg --batch --yes --armor --local-user "${GPG_FINGERPRINT}" \
  --output SHA256SUMS.sig --detach-sign SHA256SUMS

echo "Artifacts ready:"
ls -1 terraform-provider-slack_v*.zip SHA256SUMS SHA256SUMS.sig
