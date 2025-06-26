#!/bin/sh
set -e

HOOK_PATH=".git/hooks/commit-msg"
REPO_URL="https://raw.githubusercontent.com/inem/dotdotdot/main/_git/hooks/commit-msg"

echo "Downloading commit_msg hook..."
curl -fsSL "$REPO_URL" -o "$HOOK_PATH"
chmod +x "$HOOK_PATH"
echo "Hook installed to $HOOK_PATH"

if ! gem list -i openai > /dev/null 2>&1; then
  echo "Ruby gem 'openai' not found. Installing..."
  gem install openai
fi

if [ -z "$OPENAI_API_KEY" ]; then
  echo "WARNING: OPENAI_API_KEY is not set. Please export it in your shell."
fi

echo "Done."