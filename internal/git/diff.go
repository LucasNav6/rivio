// Package git reads repository state and changes through the Git executable.
//
// These operations run in the current working directory and accept contexts so
// callers can cancel long-running Git commands. Diff compares the working tree
// with a selected base branch and returns the patch for provider review.
//
// # Comparing Changes
//
//	changes, err := git.Diff(ctx, "main")
//	if err != nil {
//	    return err
//	}
//
// The package does not modify the repository.
package git

import (
	"context"
	"fmt"
	"os/exec"
	"strings"
)

// Diff returns the working-tree changes relative to base.
//
//	changes, err := git.Diff(ctx, "main")
//
// It accepts a command context and base branch name and returns the Git patch,
// or an error if the base is empty or Git cannot compare it.
func Diff(ctx context.Context, base string) (string, error) {
	// 1. Require a base branch before invoking Git.
	if strings.TrimSpace(base) == "" {
		return "", fmt.Errorf("compare Git changes: base branch is required")
	}

	// 2. Ask Git to compare the working tree with the selected base branch.
	cmd := exec.CommandContext(ctx, "git", "diff", base, "--")
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", fmt.Errorf("compare Git changes with base %q: %w: %s", base, err, strings.TrimSpace(string(output)))
	}

	// 3. Return the diff output to the caller.
	return string(output), nil
}
