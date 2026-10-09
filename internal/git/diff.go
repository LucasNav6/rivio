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
	"os"
	"os/exec"
	"path/filepath"
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

// Context returns source files from the current repository checkout.
//
// It accepts a command context and reads tracked text files from the checkout,
// returning their paths and contents or an error if Git cannot list the files.
func Context(ctx context.Context) (string, error) {
	cmd := exec.CommandContext(ctx, "git", "ls-files", "-z")
	output, err := cmd.Output()
	if err != nil {
		return "", fmt.Errorf("list repository files: %w", err)
	}
	paths := strings.Split(strings.TrimSuffix(string(output), "\x00"), "\x00")

	changed, err := exec.CommandContext(ctx, "git", "diff", "--name-only", "--diff-filter=ACMR", "HEAD", "--").Output()
	if err != nil {
		return "", fmt.Errorf("list changed repository files: %w", err)
	}
	for _, path := range strings.Split(strings.TrimSpace(string(changed)), "\n") {
		if path != "" {
			paths = append(paths, path)
		}
	}
	untracked, err := exec.CommandContext(ctx, "git", "ls-files", "--others", "--exclude-standard", "-z").Output()
	if err != nil {
		return "", fmt.Errorf("list untracked repository files: %w", err)
	}
	paths = append(paths, strings.Split(strings.TrimSuffix(string(untracked), "\x00"), "\x00")...)

	var files strings.Builder
	seen := make(map[string]struct{}, len(paths))
	for _, path := range paths {
		if path == "" {
			continue
		}
		cleanPath := filepath.Clean(path)
		if filepath.IsAbs(cleanPath) || cleanPath == ".." || strings.HasPrefix(cleanPath, ".."+string(filepath.Separator)) || !sourceFile(cleanPath) {
			continue
		}
		if _, ok := seen[cleanPath]; ok {
			continue
		}
		seen[cleanPath] = struct{}{}

		content, err := os.ReadFile(cleanPath)
		if err != nil {
			continue
		}

		fmt.Fprintf(&files, "\n--- %s ---\n%s", cleanPath, content)
	}

	return files.String(), nil
}

func sourceFile(path string) bool {
	switch strings.ToLower(filepath.Ext(path)) {
	case ".go", ".js", ".jsx", ".ts", ".tsx", ".py", ".java", ".kt", ".rs", ".rb", ".php", ".cs", ".c", ".h", ".cpp", ".hpp", ".swift", ".scala", ".ex", ".exs", ".erl", ".hrl", ".vue", ".svelte", ".sql", ".sh", ".yaml", ".yml", ".toml", ".json", ".xml", ".md":
		return true
	default:
		return false
	}
}
