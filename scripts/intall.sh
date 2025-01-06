#!/bin/bash

set -e

# Function to get the latest release tag from GitHub
get_latest_release() {
  curl --silent "https://api.github.com/repos/NethermindEth/sedge/releases/latest" |   # Get latest release from GitHub API
    grep '"tag_name":' |                                                              # Get tag line
    sed -E 's/.*"([^"]+)".*/\1/'                                                      # Extract version number
}

# Detect OS and architecture
OS="$(uname -s)"
ARCH="$(uname -m)"

# Normalize OS and ARCH
if [[ "$OS" == "Linux" ]]; then
  PLATFORM="linux"
elif [[ "$OS" == "Darwin" ]]; then
  PLATFORM="darwin"
else
  echo "Unsupported OS: $OS"
  exit 1
fi

if [[ "$ARCH" == "x86_64" ]]; then
  ARCH="amd64"
elif [[ "$ARCH" == "arm64" || "$ARCH" == "aarch64" ]]; then
  ARCH="arm64"
else
  echo "Unsupported architecture: $ARCH"
  exit 1
fi

# Get the latest version
VERSION=$(get_latest_release)
echo "Latest version is $VERSION"

# Construct the download URL
BINARY_NAME="sedge-${VERSION}-${PLATFORM}-${ARCH}"
DOWNLOAD_URL="https://github.com/NethermindEth/sedge/releases/download/${VERSION}/${BINARY_NAME}"

# Download the binary
echo "Downloading ${BINARY_NAME} from ${DOWNLOAD_URL}..."
curl -L -o sedge "${DOWNLOAD_URL}"

# Make it executable
chmod +x sedge

# Move to /usr/local/bin or another directory in PATH
echo "Installing sedge to /usr/local/bin..."
sudo mv sedge /usr/local/bin/

echo "Installation complete. You can now run 'sedge version' from the command line."
