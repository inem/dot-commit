#!/bin/sh
set -e

HOOK_PATH=".git/hooks/commit-msg"
BIN_PATH=".git/hooks/dot-commit"
HOOK_URL="https://raw.githubusercontent.com/inem/dot-commit/main/_git/hooks/commit-msg"

# Detect OS and ARCH
OS=$(uname | tr '[:upper:]' '[:lower:]')
ARCH=$(uname -m)

case "$ARCH" in
  x86_64|amd64)
    ARCH=amd64
    ;;
  arm64|aarch64)
    ARCH=arm64
    ;;
  *)
    echo "Unsupported architecture: $ARCH"
    exit 1
    ;;
esac

case "$OS" in
  linux|darwin)
    ;;
  msys*|cygwin*|mingw*|windowsnt|windows)
    OS=windows
    BIN_PATH=".git/hooks/dot-commit.exe"
    ;;
  *)
    echo "Unsupported OS: $OS"
    exit 1
    ;;
esac

BIN_URL="https://github.com/inem/dot-commit/releases/latest/download/dot-commit-${OS}-${ARCH}"

if [ ! -d .git/hooks ]; then
  echo "Error: .git/hooks directory not found. Please run this script from the root of your git repository."
  exit 1
fi

echo "Downloading shell wrapper hook..."
curl -fsSL "$HOOK_URL" -o "$HOOK_PATH"
chmod +x "$HOOK_PATH"
echo "Hook installed to $HOOK_PATH"

echo "Downloading Go binary for $OS-$ARCH..."
curl -fsSL "$BIN_URL" -o "$BIN_PATH"
chmod +x "$BIN_PATH"
echo "Binary installed to $BIN_PATH"

if [ -z "$OPENAI_API_KEY" ]; then
  echo "WARNING: OPENAI_API_KEY is not set. Please export it in your shell."
fi

echo "Done."