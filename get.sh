#!/bin/sh

if [ "$(uname -m)" = "aarch64" ]; then ARCH="arm64"; else ARCH="amd64"; fi;
OS="$(uname)"
LATEST_VERSION=$(curl -s https://api.github.com/repos/canta2899/logo-ls/releases/latest | grep '"tag_name"' | cut -d '"' -f 4)
DOWNLOAD_URL="https://github.com/canta2899/logo-ls/releases/download/${LATEST_VERSION}/logo-ls-${LATEST_VERSION}-${OS}-${ARCH}.tar.gz"

curl -sL ${DOWNLOAD_URL} | tar -xzf - logo-ls && sudo mv logo-ls /usr/local/bin
