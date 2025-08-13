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

	// Modes:
	// "."  -> full diff including static files (new default)
	// "-"  -> file list only (replaces old "..")
	// "#"  -> diff excluding static files (old "." behavior)
	if text != "." && text != "-" && text != "#" {
		logDebug(logFile, "message is not one of '.', '-', '#'; exiting")
		return
	}

	logDebug(logFile, fmt.Sprintf("proceeding with AI generation for '%s'", text))

	var body string
	var mode string

	switch text {
	case "-":
		mode = "file_list"
		logDebug(logFile, "using file list only mode (-)")
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
	case "#":
		mode = "diff_excluding_static"
		logDebug(logFile, "using diff excluding static files mode (#)")
		// old behavior: exclude static files from diff, append list at end
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
		staticCmd := exec.Command("git", "diff", "--cached", "--name-only")
		staticOutput, err := staticCmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git name-only command failed: %v", err))
			return
		}
		body = string(diffOutput) + "\nSTATIC FILES CHANGED:\n" + string(staticOutput)
		if len(body) > 8000 {
			body = body[:8000]
		}
	case ".":
		mode = "full_diff"
		logDebug(logFile, "using full diff including static files mode (.)")
		diffCmd := exec.Command("git", "diff", "--cached", "--binary")
		diffOutput, err := diffCmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git diff command failed: %v", err))
			return
		}
		// Also capture name-only list for potential static awareness (not excluding here)
		nameOnlyCmd := exec.Command("git", "diff", "--cached", "--name-only")
		nameOnlyOutput, err := nameOnlyCmd.Output()
		if err != nil {
			logDebug(logFile, fmt.Sprintf("git name-only command failed: %v", err))
			return
		}
		body = string(diffOutput) + "\nFILES CHANGED:\n" + string(nameOnlyOutput)
		if len(body) > 8000 {
			body = body[:8000]
		}
	}

	// Send to OpenAI
	apiKey := os.Getenv("OPENAI_API_KEY")
	if apiKey == "" {
		logDebug(logFile, "OPENAI_API_KEY not set")
		return
	}

	client := openai.NewClient(apiKey)

	var prompt string
	if mode == "file_list" {
		prompt = fmt.Sprintf(`You are a senior software engineer crafting high–quality, intention‑revealing Git commit titles.

TASK:
Infer the underlying intent / purpose of the staged changes from ONLY the file list below and output ONE concise commit title (<=70 chars).

GUIDELINES:
- Focus on WHY / intent, not a mechanical list of filenames.
- Use imperative mood (Add, Fix, Refactor, Improve, Optimize, Remove, Adjust, Document, Configure, Harden).
- Infer probable category: feature, fix, refactor, perf, tests, docs, build, chore, style.
- Avoid bland "update" unless nothing else fits.
- Do NOT include quotes, punctuation at end, or code fences.
- Do NOT just say "update X file" unless that is genuinely all.
- If only static/front-end assets (css/html/images) changed: summarize purpose (e.g. "Refine layout spacing", "Update hero images", "Adjust responsive styles").
- If mostly tests added: "Add tests for ..." (or similar).
- If looks like a rename/move: "Rename X to Y" (if clear).
- Prefer describing effect (e.g. "Fix panic on empty config", "Improve cache invalidation logic").
- No trailing period. Single line only.

FILE CHANGES:
%s

Commit message:`, body)
	} else { // diff modes
		prompt = fmt.Sprintf(`You are a senior software engineer writing a single Git commit title.

TASK:
Analyze the following staged diff (and file list at the end) and produce ONE concise, intention‑revealing commit message line (<=70 chars) capturing the primary purpose (WHY) of the change set.

DIFF & CONTEXT:
%s

INSTRUCTIONS:
1. Derive the intent (feature addition, bug fix, refactor, performance improvement, reliability, error handling, security, configuration, cleanup, tests, docs).
2. Prefer the most meaningful unifying theme over listing multiple trivial edits.
3. Use imperative mood (Add, Fix, Refactor, Improve, Optimize, Remove, Harden, Document).
4. Avoid generic "update", "change", "misc", "adjust" if a more specific verb fits.
5. Do not list filenames unless a rename/move is itself the core change.
6. If mainly cleanup (whitespace / formatting): "Format code" or "Apply go fmt" (whichever fits).
7. If adding error handling / validation: describe the protected aspect ("Handle nil X", "Validate Y before Z").
8. If improving performance: state what was optimized ("Optimize query building", etc.).
9. If only static asset/style changes: summarize design intent (e.g. "Refine button spacing", "Update hero images").
10. No trailing period, no quotes, no code fences. Output ONLY the commit title.

Commit message:`, body)
	}

	logDebug(logFile, "sending request to OpenAI...")

	resp, err := client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: openai.GPT4oMini,
			Messages: []openai.ChatCompletionMessage{{
				Role:    openai.ChatMessageRoleUser,
				Content: prompt,
			}},
			MaxTokens:   40,
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
	if idx := strings.IndexAny(newMsg, "\r\n"); idx >= 0 {
		newMsg = strings.TrimSpace(newMsg[:idx])
	}
	if len(newMsg) > 70 {
		newMsg = strings.TrimSpace(newMsg[:70])
	}
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
