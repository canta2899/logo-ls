#!/bin/sh

ARCH=$(uname -m)
OS=$(uname)

# Normalize architecture
case "$ARCH" in
  arm64|aarch64)
    ARCH="arm64"
    ;;
  x86_64|amd64)
    ARCH="amd64"
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

# Normalize OS name
case "$OS" in
  Darwin)
    OS="darwin"
    ;;
  Linux)
    OS="linux"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

LATEST_VERSION=$(curl -s https://api.github.com/repos/canta2899/logo-ls/releases/latest | grep '"tag_name"' | cut -d '"' -f 4)

DOWNLOAD_URL="https://github.com/canta2899/logo-ls/releases/download/${LATEST_VERSION}/logo-ls-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"

INSTALL_DIR="$HOME/.local/bin"
mkdir -p "$INSTALL_DIR"

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
  echo "Warning: $INSTALL_DIR is not in your PATH."
  echo "You should add it to your PATH or move logo-ls to a directory that is."
fi

curl -sL "$DOWNLOAD_URL" | tar -xzf - logo-ls && mv logo-ls "$INSTALL_DIR"

echo "✅ logo-ls installed successfully to $INSTALL_DIR"
