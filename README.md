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
curl -fsSL https://raw.githubusercontent.com/inem/dotdotdot/refs/heads/main/install.sh | sh

## Regular Installation

1. Copy the `commit_msg` script to your repository's `.git/hooks/` directory as `commit-msg`.
2. Make sure the script is executable:
   ```bash
   chmod +x .git/hooks/commit-msg
   ```
3. Set your OpenAI API key in the environment:
   ```bash
   export OPENAI_API_KEY=your-api-key
   ```

## Usage

- When committing, use `git commit -m "..."` or `git commit -m ".."`:
  - `...` — AI will analyze the full diff and generate a message.
  - `..` — AI will analyze only the file list and generate a message.
- For other commit messages, the hook does nothing.

## Requirements

- Ruby
- The `openai` Ruby gem (`gem install openai`)
- An OpenAI API key

## Notes

- The hook logs debug information to `.git/commit-msg-debug.log`.
- If only static files are changed, the commit message will reflect this.