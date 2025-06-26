# dotdotdot

## Overview

This repository includes a custom Git commit message hook powered by OpenAI. The hook helps generate concise, meaningful commit messages automatically, based on your staged changes.

## Features

- **AI-generated commit messages:** Uses OpenAI to suggest commit messages based on your code changes or file list.
- **Static file awareness:** Recognizes static file changes (CSS, HTML, images) and notes them accordingly.
- **Two modes:**
  - `...` (three dots): Generates a commit message based on the full diff of staged changes.
  - `..` (two dots): Generates a commit message based only on the list of changed files.

## Quick Install

```sh
curl -fsSL https://raw.githubusercontent.com/inem/dotdotdot/main/install.sh | sh
```

```sh
curl -fsSL https://instll.sh/inem/dotdotdot | sh
```


## Manual Installation (Go version)

1. Clone this repository and enter the directory.
2. Build the Go binary:
   ```bash
   make build
   ```
3. Install the hook and binary (this will overwrite `.git/hooks/commit-msg` and `.git/hooks/commit-msg-go`):
   ```bash
   cp _git/hooks/commit-msg .git/hooks/commit-msg
   cp release/commit-msg-go .git/hooks/commit-msg-go
   chmod +x .git/hooks/commit-msg .git/hooks/commit-msg-go
   ```
4. Set your OpenAI API key in the environment:
   ```bash
   export OPENAI_API_KEY=your-api-key
   ```

## Usage

- When committing, use `git commit -m "..."` or `git commit -m ".."`:
  - `...` — AI will analyze the full diff and generate a message.
  - `..` — AI will analyze only the file list and generate a message.
- For other commit messages, the hook does nothing.

## Requirements

- Go 1.21+
- The `github.com/sashabaranov/go-openai` Go module (installed automatically by `make build`)
- An OpenAI API key

## Notes

- The hook logs debug information to `.git/commit-msg-debug.log`.
- If only static files are changed, the commit message will reflect this.
