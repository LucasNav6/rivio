// Package git reads repository state and changes through the Git executable.
//
// These operations run in the current working directory and accept contexts so
// callers can cancel long-running Git commands. Detect verifies that the
// current directory is inside a working tree before review begins.
//
// # Checking a Working Tree
//
//	if err := git.Detect(ctx); err != nil {
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

// Detect verifies that the current directory is inside a Git working tree.
//
//	if err := git.Detect(ctx); err != nil {
//	    return err
//	}
//
// It accepts a context for the Git command and returns an error if Git cannot
// confirm that the directory belongs to a working tree.
func Detect(ctx context.Context) error {
	// 1. Ask Git whether the current directory is inside a working tree.
	cmd := exec.CommandContext(ctx, "git", "rev-parse", "--is-inside-work-tree")

	output, err := cmd.Output()
	if err != nil {
		return fmt.Errorf("detect git repository: %w", err)
	}

	// 2. Accept the directory only when Git confirms it is a working tree.
	if strings.TrimSpace(string(output)) != "true" {
		return fmt.Errorf("current directory is not inside a Git repository")
	}

	return nil
}
