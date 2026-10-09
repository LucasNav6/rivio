// Package cmd runs Rivio from the command line.
//
// Rivio compares local Git changes with a base branch and sends the resulting
// diff to the configured AI provider. The command package connects the CLI
// flags to configuration loading, Git inspection, provider selection, and
// review output.
//
// # Reviewing Changes
//
// Use the review command from a Git working tree. The configuration file
// supplies the provider and default base branch; flags can override the file
// path and base branch for an individual run.
//
//	rivio review --config rivio.yml --base main
//
// Review findings are written to standard output. Operational events and
// failures are reported through Rivio's structured logger.
package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

// rootCmd is the root command that configures and runs the Rivio CLI.
var rootCmd = &cobra.Command{
	Use:   "rivio",
	Short: "Open-source AI code reviewer that runs on your terms.",
	Long: `Rivio is a local-first AI-powered code review tool designed to help
developers review changes before they reach production.

It analyzes your local Git changes using the AI provider of your choice,
identifies potential bugs, security issues, performance problems, and other
actionable findings, while keeping you in control of what gets published.`,
}

// Execute runs the root command. It accepts no arguments and exits the process
// with a non-zero status when command execution fails.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
