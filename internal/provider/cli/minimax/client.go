// Package minimax runs prompts through the MiniMax Code CLI.
package minimax

import (
	"context"
	"fmt"

	"github.com/LucasNav6/rivio/internal/provider/cli/runner"
)

// Client runs prompts with MiniMax Code in lightweight mode.
type Client struct {
	executable string
}

// New locates MiniMax Code and creates a client.
//
// It accepts no arguments and returns a client or an error if the mcode
// executable is unavailable.
func New() (*Client, error) {
	// 1. Locate the MiniMax Code executable.
	executable, err := runner.Locate("mcode")
	if err != nil {
		return nil, err
	}
	return &Client{executable: executable}, nil
}

// Run sends prompt to MiniMax Code and returns the generated response.
//
// It accepts a context and prompt and returns generated text or an error if
// the CLI command fails.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Run one lightweight, non-interactive MiniMax Code request.
	content, err := runner.Run(ctx, client.executable, []string{"exec", "--mode", "lightweight", prompt}, "")
	if err != nil {
		return "", fmt.Errorf("MiniMax Code: %w", err)
	}
	return content, nil
}
