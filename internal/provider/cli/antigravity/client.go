// Package antigravity runs prompts through the Antigravity CLI.
package antigravity

import (
	"context"
	"fmt"
	"strings"

	"github.com/LucasNav6/rivio/internal/provider/cli/runner"
)

// Client runs prompts with Antigravity in non-interactive sandbox mode.
type Client struct {
	executable string
	model      string
}

// New locates Antigravity and creates a configured client.
//
// It accepts an optional model name and returns a client or an error if the
// Antigravity CLI is unavailable.
func New(model string) (*Client, error) {
	// 1. Locate the Antigravity executable.
	executable, err := runner.Locate("agy")
	if err != nil {
		return nil, err
	}
	return &Client{executable: executable, model: strings.TrimSpace(model)}, nil
}

// Run sends prompt to Antigravity and returns the generated response.
//
// It accepts a context and prompt and returns generated text or an error if
// the CLI command fails.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Configure a single sandboxed headless prompt.
	args := []string{"-p", prompt, "--sandbox"}
	if client.model != "" {
		args = append(args, "--model", client.model)
	}

	// 2. Run Antigravity and return its response.
	content, err := runner.Run(ctx, client.executable, args, "")
	if err != nil {
		return "", fmt.Errorf("Antigravity: %w", err)
	}
	return content, nil
}
