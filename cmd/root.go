package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "rivio",
	Short: "Open-source AI code reviewer that runs on your terms.",
	Long: `Rivio is a local-first AI-powered code review tool designed to help
developers review changes before they reach production.

It analyzes your local Git changes using the AI provider of your choice,
identifies potential bugs, security issues, performance problems, and other
actionable findings, while keeping you in control of what gets published.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}
