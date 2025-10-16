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

if ! gpg --batch --list-secret-keys "$GPG_FINGERPRINT" >/dev/null 2>&1; then
  echo "No secret key found for fingerprint ${GPG_FINGERPRINT}. Did you import the key?" >&2
  exit 1
fi

# Clean up any previous build artifacts.
rm -f terraform-provider-slack_v*.zip SHA256SUMS SHA256SUMS.sig terraform-registry-manifest.json

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
if [ -n "${GPG_PASSPHRASE:-}" ]; then
  printf '%s' "${GPG_PASSPHRASE}" | gpg --batch --yes --armor --pinentry-mode loopback \
    --passphrase-fd 0 --local-user "${GPG_FINGERPRINT}" \
    --output SHA256SUMS.sig --detach-sign SHA256SUMS
else
  gpg --batch --yes --armor --pinentry-mode loopback --local-user "${GPG_FINGERPRINT}" \
    --output SHA256SUMS.sig --detach-sign SHA256SUMS
fi

echo "Generating terraform-registry-manifest.json..."
export VERSION GPG_FINGERPRINT
python3 <<'PY'
import json
import os
import subprocess

version = os.environ["VERSION"]
fingerprint = os.environ["GPG_FINGERPRINT"]

public_key = subprocess.run(
    ["gpg", "--armor", "--export", fingerprint],
    check=True,
    capture_output=True,
    text=True,
).stdout

if not public_key.strip():
    raise SystemExit("Unable to export GPG public key for manifest generation.")

fingerprint_output = subprocess.run(
    ["gpg", "--with-colons", "--fingerprint", fingerprint],
    check=True,
    capture_output=True,
    text=True,
).stdout

key_id = None
for line in fingerprint_output.splitlines():
    if line.startswith("fpr:"):
        key_id = line.split(":")[9][-16:]
        break

if not key_id:
    raise SystemExit("Unable to determine GPG key id.")

packages = {}
prefix = f"terraform-provider-slack_v{version}_"

with open("SHA256SUMS", "r", encoding="utf-8") as sums_file:
    for raw_line in sums_file:
        line = raw_line.strip()
        if not line:
            continue
        parts = line.split(None, 1)
        if len(parts) != 2:
            continue
        shasum, filename = parts
        if not filename.endswith(".zip") or not filename.startswith(prefix):
            continue
        target = filename[len(prefix):-4]  # Strip prefix and .zip
        packages[target] = {
            "filename": filename,
            "shasum": shasum,
            "signing_keys": {
                "gpg_public_keys": [
                    {
                        "key_id": key_id,
                        "ascii_armor": public_key,
                    }
                ]
            },
        }

if not packages:
    raise SystemExit("No package entries were discovered while building the manifest.")

manifest = {
    "version": 1,
    "metadata": {
        "protocol_versions": ["6.0"],
    },
    "packages": dict(sorted(packages.items())),
}

with open("terraform-registry-manifest.json", "w", encoding="utf-8") as manifest_file:
    json.dump(manifest, manifest_file, indent=2)
    manifest_file.write("\n")
PY

echo "Artifacts ready:"
ls -1 terraform-provider-slack_v*.zip SHA256SUMS SHA256SUMS.sig terraform-registry-manifest.json
