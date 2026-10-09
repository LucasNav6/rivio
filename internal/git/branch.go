// Package git reads repository state and changes through the Git executable.
//
// These operations run in the current working directory and accept contexts so
// callers can cancel long-running Git commands. Branch helpers resolve the
// comparison branch used when Rivio collects changes for review.
//
// # Selecting a Base Branch
//
// An explicit branch value takes precedence over the configured default:
//
//	base, err := git.ResolveBaseBranch(flagValue, configValue)
//	if err != nil {
//		return err
//	}
//
// The package does not modify the repository.
package git

import (
	"fmt"
	"strings"
)

// ResolveBaseBranch selects the base branch used to compare local changes.
//
//	base, err := git.ResolveBaseBranch(flagValue, configValue)
//
// It accepts the optional command-line branch and configured default, and
// returns the first non-empty value or an error if neither is provided.
func ResolveBaseBranch(flagValue, configValue string) (string, error) {
	// 1. Prefer the branch supplied explicitly on the command line.
	if strings.TrimSpace(flagValue) != "" {
		return strings.TrimSpace(flagValue), nil
	}

	// 2. Fall back to the branch configured for reviews.
	if strings.TrimSpace(configValue) != "" {
		return strings.TrimSpace(configValue), nil
	}

	return "", fmt.Errorf("base branch is required: set --base or review.base")
}
