package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var reviewCmd = &cobra.Command{
	Use:   "review",
	Short: "Review code changes with AI",
	Long: `Review code changes against a base branch using your configured AI provider.

Rivio analyzes the Git diff, looking for actionable issues such as bugs,
security vulnerabilities, incorrect behavior, performance problems, and
potential regressions.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("review called")
	},
}

func init() {
	rootCmd.AddCommand(reviewCmd)
}
