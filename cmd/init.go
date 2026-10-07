package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Configure Rivio for this repository",
	Long: `Initialize Rivio in the current Git repository.

Rivio guides you through the initial setup, including selecting an AI
provider, choosing a model, configuring the default base branch, and setting
review preferences.

The configuration is stored in a .rivio.yml file that can be committed and
shared with your team. API keys and other secrets are never stored in the
repository configuration.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("init called")
	},
}

func init() {
	rootCmd.AddCommand(initCmd)
}
