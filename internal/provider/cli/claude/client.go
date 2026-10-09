// Package claude runs prompts through the authenticated Claude Code CLI.
package claude

import (
	"context"
	"fmt"
	"strings"

	"github.com/LucasNav6/rivio/internal/provider/cli/runner"
)

// Client runs prompts with Claude Code without enabling tools.
type Client struct {
	executable string
	model      string
}

// New locates Claude Code and creates a configured client.
//
// It accepts an optional model name and returns a client or an error if the
// Claude Code CLI is unavailable.
func New(model string) (*Client, error) {
	// 1. Locate the Claude Code executable.
	executable, err := runner.Locate("claude")
	if err != nil {
		return nil, err
	}
	return &Client{executable: executable, model: strings.TrimSpace(model)}, nil
}

// Run sends prompt to Claude Code and returns the generated response.
//
// It accepts a context and prompt and returns generated text or an error if
// the CLI command fails.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Configure non-interactive output and disable Claude tools.
	args := []string{"-p", "--tools", ""}
	if client.model != "" {
		args = append(args, "--model", client.model)
	}
	args = append(args, prompt)

	// 2. Run Claude Code and return its response.
	content, err := runner.Run(ctx, client.executable, args, "")
	if err != nil {
		return "", fmt.Errorf("Claude Code: %w", err)
	}
	return content, nil
}
