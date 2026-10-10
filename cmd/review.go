// Package cmd runs Rivio from the command line.
//
// Rivio compares local Git changes with a base branch and sends the resulting
// diff to the configured AI provider. The command package connects the CLI
// flags to configuration loading, Git inspection, provider selection, and
// review output.
//
// # Reviewing Changes
//
// Use the review command from a Git working tree. The configuration file at
// ~/.config/rivio/config.yml supplies the provider and default base branch;
// flags can override the file path and base branch for an individual run.
//
//	rivio review --base main
//
// Each analysis overwrites diff.md and umd.md in ~/.config/rivio. Operational
// events and failures are reported through Rivio's structured logger.
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/LucasNav6/rivio/internal/config"
	"github.com/LucasNav6/rivio/internal/git"
	"github.com/LucasNav6/rivio/internal/prompt"
	"github.com/LucasNav6/rivio/internal/provider"
	"github.com/spf13/cobra"

	"charm.land/log/v2"
)

// baseBranch stores the optional base branch provided to the review command.
var baseBranch string

// configPath stores the path to the Rivio YAML configuration file.
var configPath string

// reviewCmd defines the command that reviews local changes with an AI provider.
var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes with AI",
	Long: `Review code changes against a base branch using your configured AI provider.

Rivio analyzes the Git diff, looking for actionable issues such as bugs,
security vulnerabilities, incorrect behavior, performance problems, and
potential regressions.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return runReview(cmd.Context())
	},
}

// runReview reviews the changes between a Git base branch and the working tree.
//
// The base branch comes from --base when set, otherwise from review.base in the
// configuration file. The selected provider receives the shared review prompt
// containing the Git diff. The sequence diagram and analyzed diff overwrite
// separate files in ~/.config/rivio.
//
//	err := runReview(ctx)
//
// It accepts a command context and returns an error if
// repository detection, configuration loading, provider setup, diff creation,
// review generation, or report writing fails.
func runReview(ctx context.Context) error {
	// 1. Load the Rivio configuration from the configured file path.
	cfg, err := config.Load(configPath)
	if err != nil {
		log.Fatal("It could not load the configuration", "error", err)
		return err
	}

	// 2. Initialize the client selected by transport and provider name.
	client, err := provider.New(cfg.Provider)
	if err != nil {
		log.Fatal("It could not initialize the AI provider", "error", err)
		return err
	}

	// 3. Resolve the base branch from the flag or configuration.
	selectedBase, err := git.ResolveBaseBranch(baseBranch, cfg.Review.Base)
	if err != nil {
		log.Fatal("It could not select the base branch", "error", err)
		return err
	}

	// 4. Detect whether Rivio is running inside a Git repository.
	if err := git.Detect(ctx); err != nil {
		log.Fatal("It could not detect the Git repository", "error", err)
		return err
	}

	// 5. Get the changes between the selected base branch and the working tree.
	changes, err := git.Diff(ctx, selectedBase)
	if err != nil {
		log.Fatal("It could not get the Git changes", "error", err)
		return err
	}
	originBranch, err := git.CurrentBranch(ctx)
	if err != nil {
		log.Fatal("It could not determine the current branch", "error", err)
		return err
	}

	// 6. Gather changed source files so the provider can trace the affected flow.
	codeContext, err := git.Context(ctx, selectedBase)
	if err != nil {
		log.Fatal("It could not read the changed source files", "error", err)
		return err
	}

	// 7. Build the flow-diagram prompt and run it with the selected provider.
	result, err := client.Run(ctx, prompt.Review(changes, codeContext))
	if err != nil {
		log.Fatal("It could not complete the AI code review", "error", err)
		return err
	}

	// 8. Keep only a fenced sequenceDiagram from the provider response.
	diagram, err := extractSequenceDiagram(result)
	if err != nil {
		log.Fatal("The provider response did not contain a sequence diagram", "error", err)
		return err
	}

	// 9. Replace the diff and diagram reports in the user config directory.
	home, err := os.UserHomeDir()
	if err != nil {
		log.Fatal("Rivio could not resolve the user config directory", "error", err)
		return err
	}
	reportDir := filepath.Join(home, ".config", "rivio")
	if err := os.MkdirAll(reportDir, 0700); err != nil {
		log.Fatal("Rivio could not create the user config directory", "error", err)
		return err
	}
	diffReport := fmt.Sprintf("branch origin: %s\nbranch target: %s\n\n### DIFF ANALIZED\n```git diff\n%s\n```\n", originBranch, selectedBase, changes)
	if err := os.WriteFile(filepath.Join(reportDir, "diff.md"), []byte(diffReport), 0600); err != nil {
		log.Fatal("Rivio could not write the diff report", "action", "write_diff_report", "result", "failed", "error", err)
		return err
	}
	sequenceReport := fmt.Sprintf("branch origin: %s\nbranch target: %s\n\n### SEQUENCE DIAGRAM\n%s\n", originBranch, selectedBase, diagram)
	if err := os.WriteFile(filepath.Join(reportDir, "umd.md"), []byte(sequenceReport), 0600); err != nil {
		log.Fatal("Rivio could not write the sequence diagram", "action", "write_sequence_diagram", "result", "failed", "error", err)
		return err
	}
	log.Info("Rivio wrote analysis files", "action", "write_analysis_files", "result", "success", "directory", reportDir)

	return nil
}

func extractSequenceDiagram(response string) (string, error) {
	lines := strings.Split(strings.ReplaceAll(response, "\r\n", "\n"), "\n")
	for start, line := range lines {
		if !strings.HasPrefix(strings.TrimSpace(line), "```") {
			continue
		}
		for end := start + 1; end < len(lines); end++ {
			if strings.TrimSpace(lines[end]) != "```" {
				continue
			}
			block := strings.TrimSpace(strings.Join(lines[start+1:end], "\n"))
			if strings.HasPrefix(block, "sequenceDiagram") {
				return "```mermaid\n" + block + "\n```", nil
			}
			break
		}
	}
	return "", fmt.Errorf("missing fenced Mermaid sequenceDiagram")
}

// init registers review command flags and attaches reviewCmd to rootCmd. It
// accepts no arguments and returns no value.
func init() {
	reviewCmd.Flags().StringVar(&configPath, "config", "~/.config/rivio/config.yml", "Path to Rivio config file")
	reviewCmd.Flags().StringVar(&baseBranch, "base", "", "Base branch to compare against config")
	rootCmd.AddCommand(reviewCmd)
}
