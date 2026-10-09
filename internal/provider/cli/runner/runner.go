// Package runner executes provider command-line clients with a supplied prompt.
package runner

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Locate returns the executable path for a command available on PATH.
//
// It accepts a command name and returns its resolved path or an error if the
// executable is unavailable.
func Locate(name string) (string, error) {
	// 1. Resolve the executable using the current process PATH.
	executable, err := exec.LookPath(name)
	if err != nil {
		return "", fmt.Errorf("%s CLI is not installed or not on PATH: %w", name, err)
	}
	return executable, nil
}

// Run executes a CLI command with prompt either in its arguments or stdin.
//
// It accepts the context, executable path, argument list, and optional stdin
// prompt and returns standard output or an execution or empty-output error.
func Run(ctx context.Context, executable string, args []string, stdin string) (string, error) {
	// 1. Start the provider command with the supplied execution context.
	command := exec.CommandContext(ctx, executable, args...)
	if stdin != "" {
		command.Stdin = strings.NewReader(stdin)
	}

	// 2. Capture the provider response and reject failed or empty runs.
	output, err := command.Output()
	if err != nil {
		return "", fmt.Errorf("run provider CLI: %w", err)
	}
	content := strings.TrimSpace(string(output))
	if content == "" {
		return "", fmt.Errorf("provider CLI returned no content")
	}

	return content, nil
}
