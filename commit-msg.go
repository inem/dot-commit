package main

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
)

var staticExt = []string{"css", "scss", "html", "htm", "png", "jpg", "jpeg", "gif", "svg"}

func isStaticFile(filename string) bool {
	ext := strings.ToLower(filepath.Ext(filename))
	if len(ext) > 0 {
		ext = ext[1:] // remove dot
	}
	for _, staticExtension := range staticExt {
		if ext == staticExtension {
			return true
		}
	}
	return false
}

func logDebug(logFile *os.File, message string) {
	if logFile != nil {
		fmt.Fprintf(logFile, "DEBUG: %s\n", message)
	}
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println("usage: commit-msg <path>")
		os.Exit(1)
	}

	msgFile := os.Args[1]
	content, err := os.ReadFile(msgFile)
	if err != nil {
		fmt.Printf("Error reading file: %v\n", err)
		os.Exit(1)
	}

	text := strings.TrimSpace(string(content))

	// Debug log
	logFile, _ := os.OpenFile(".git/commit-msg-debug.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if logFile != nil {
		defer logFile.Close()
		fmt.Fprintf(logFile, "=== %s ===\n", time.Now().Format(time.RFC3339))
		logDebug(logFile, fmt.Sprintf("commit message = '%s'", text))
	}

	// Only process if message is "..." or ".."
	if text != "..." && text != ".." {
		logDebug(logFile, "message is not '...' or '..', exiting")
		return
	}

	logDebug(logFile, fmt.Sprintf("proceeding with AI generation for '%s'", text))

	var body string

	if text == ".." {
		// 1. For ".." - only file list mode
		logDebug(logFile, "using file list only mode")

		cmd := exec.Command("git", "diff", "--cached", "--name-only")
		output, err := cmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git command failed: %v", err))
			return
		}

		files := strings.Split(strings.TrimSpace(string(output)), "\n")
		var staticFiles, nonStaticFiles []string

		for _, file := range files {
			file = strings.TrimSpace(file)
			if file == "" {
				continue
			}
			if isStaticFile(file) {
				staticFiles = append(staticFiles, file)
			} else {
				nonStaticFiles = append(nonStaticFiles, file)
			}
		}

		body = "FILES CHANGED:\n" + strings.Join(nonStaticFiles, "\n")
		if len(staticFiles) > 0 {
			body += "\nSTATIC FILES CHANGED:\n" + strings.Join(staticFiles, "\n")
		}
	} else {
		// 2. For "..." - full diff mode
		logDebug(logFile, "using full diff mode")

		// Exclude static files from diff
		excludeArgs := []string{"diff", "--cached", "--binary"}
		for _, ext := range staticExt {
			excludeArgs = append(excludeArgs, fmt.Sprintf(":(exclude)*.%s", ext))
		}

		diffCmd := exec.Command("git", excludeArgs...)
		diffOutput, err := diffCmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git diff command failed: %v", err))
			return
		}

		// Get static files list
		staticCmd := exec.Command("git", "diff", "--cached", "--name-only")
		staticOutput, err := staticCmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git name-only command failed: %v", err))
			return
		}

		body = string(diffOutput) + "\nSTATIC FILES CHANGED:\n" + string(staticOutput)

		// Trim to 8KB
		if len(body) > 8000 {
			body = body[:8000]
		}
	}

	// 2. Send to OpenAI
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logDebug(logFile, "OPENAI_API_KEY not set")
		return
	}

	client := openai.NewClient(apiKey)

	var prompt string
	if text == ".." {
		prompt = fmt.Sprintf("Generate a short (≤70 chars) git commit message based on file list.\nIf only static files changed (css/html/images) - note this.\n%s", body)
	} else {
		prompt = fmt.Sprintf("Write a git commit message (max 70 chars) describing these code changes:\n%s\n\nMessage:", body)
	}

	logDebug(logFile, "sending request to OpenAI...")

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{
				{
					Role:    openai.ChatMessageRoleUser,
					Content: prompt,
				},
			},
			MaxTokens:   20,
			Temperature: 0.3,
		},
	)

	if err != nil {
		logDebug(logFile, fmt.Sprintf("OpenAI API error: %v", err))
		return
	}

	if len(resp.Choices) == 0 {
		logDebug(logFile, "no choices in OpenAI response")
		return
	}

	newMsg := strings.TrimSpace(resp.Choices[0].Message.Content)
	logDebug(logFile, fmt.Sprintf("got response: '%s'", newMsg))

	if newMsg != "" {
		err = os.WriteFile(msgFile, []byte(newMsg+"\n"), 0644)
		if err != nil {
			logDebug(logFile, fmt.Sprintf("file write error: %v", err))
			return
		}
		logDebug(logFile, "wrote new message to file")
	}
}