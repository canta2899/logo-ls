#!/bin/sh

if [ "$(uname -m)" = "aarch64" ]; then ARCH="arm64"; else ARCH="amd64"; fi;
OS="$(uname)"

LATEST_VERSION=$(curl -s https://api.github.com/repos/canta2899/logo-ls/releases/latest | grep '"tag_name"' | cut -d '"' -f 4)

DOWNLOAD_URL="https://github.com/canta2899/logo-ls/releases/download/${LATEST_VERSION}/logo-ls-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"

INSTALL_DIR="$HOME/.local/bin"

mkdir -p $INSTALL_DIR

if ! echo "$PATH" | grep -q "$INSTALL_DIR"; then
    echo "Warning: $INSTALL_DIR is not in your PATH."
    echo "You should either add it to your PATH or move logo-ls to a directory that is in your PATH."
fi

curl -sL ${DOWNLOAD_URL} | tar -xzf - logo-ls && mv logo-ls $INSTALL_DIR
