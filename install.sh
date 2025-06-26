#!/bin/sh
set -e

HOOK_PATH=".git/hooks/commit-msg"
BIN_PATH=".git/hooks/commit-msg-go"
HOOK_URL="https://raw.githubusercontent.com/inem/dotdotdot/main/_git/hooks/commit-msg"
BIN_URL="https://github.com/inem/dotdotdot/releases/latest/download/commit-msg-go"

if [ ! -d .git/hooks ]; then
  echo "Error: .git/hooks directory not found. Please run this script from the root of your git repository."
  exit 1
fi

echo "Downloading shell wrapper hook..."
curl -fsSL "$HOOK_URL" -o "$HOOK_PATH"
chmod +x "$HOOK_PATH"
echo "Hook installed to $HOOK_PATH"

echo "Downloading Go binary..."
curl -fsSL "$BIN_URL" -o "$BIN_PATH"
chmod +x "$BIN_PATH"
echo "Binary installed to $BIN_PATH"

if [ -z "$OPENAI_API_KEY" ]; then
  echo "WARNING: OPENAI_API_KEY is not set. Please export it in your shell."
fi

echo "Done."