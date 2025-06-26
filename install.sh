#!/bin/sh
set -e

HOOK_PATH=".git/hooks/commit-msg"
REPO_URL="https://raw.githubusercontent.com/inem/dotdotdot/refs/heads/main/_git/hooks/commit-msg-go"

if [ ! -d .git/hooks ]; then
  echo "Error: .git/hooks directory not found. Please run this script from the root of your git repository."
  exit 1
fi

echo "Downloading commit_msg hook..."
curl -fsSL "$REPO_URL" -o "$HOOK_PATH"
chmod +x "$HOOK_PATH"
echo "Hook installed to $HOOK_PATH"

if [ -z "$OPENAI_API_KEY" ]; then
  echo "WARNING: OPENAI_API_KEY is not set. Please export it in your shell."
fi

echo "Done."