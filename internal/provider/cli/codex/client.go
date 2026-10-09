// Package codex runs prompts through the authenticated Codex CLI.
package codex

import (
	"context"
	"fmt"
	"strings"

	"github.com/LucasNav6/rivio/internal/provider/cli/runner"
)

// Client runs prompts with Codex using read-only sandboxing.
type Client struct {
	executable string
	model      string
}

// New locates Codex and creates a configured client.
//
// It accepts an optional model name and returns a client or an error if Codex
// is not installed or cannot be found on PATH.
func New(model string) (*Client, error) {
	// 1. Locate the Codex executable.
	executable, err := runner.Locate("codex")
	if err != nil {
		return nil, err
	}
	return &Client{executable: executable, model: strings.TrimSpace(model)}, nil
}

// Run sends prompt to Codex and returns the generated response.
//
// It accepts a context and prompt and returns generated text or an error if
// the CLI command fails.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Configure Codex for an ephemeral, read-only execution.
	args := []string{"exec", "--sandbox", "read-only", "--ephemeral", "--color", "never"}
	if client.model != "" {
		args = append(args, "--model", client.model)
	}
	args = append(args, "-")

	// 2. Send the prompt on standard input and return Codex output.
	content, err := runner.Run(ctx, client.executable, args, prompt)
	if err != nil {
		return "", fmt.Errorf("Codex: %w", err)
	}
	return content, nil
}
