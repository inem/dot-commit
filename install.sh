#!/bin/sh
set -e

HOOK_PATH=".git/hooks/commit-msg"
BIN_PATH=".git/hooks/commit-msg-go"
SCRIPT_SRC="_git/hooks/commit-msg"
BIN_SRC="release/commit-msg-go"

if [ ! -d .git/hooks ]; then
  echo "Error: .git/hooks directory not found. Please run this script from the root of your git repository."
  exit 1
fi

echo "Installing shell wrapper hook..."
cp "$SCRIPT_SRC" "$HOOK_PATH"
chmod +x "$HOOK_PATH"
echo "Hook installed to $HOOK_PATH"

echo "Copying Go binary to .git/hooks..."
cp "$BIN_SRC" "$BIN_PATH"
chmod +x "$BIN_PATH"
echo "Binary installed to $BIN_PATH"

if [ -z "$OPENAI_API_KEY" ]; then
  echo "WARNING: OPENAI_API_KEY is not set. Please export it in your shell."
fi

echo "Done."