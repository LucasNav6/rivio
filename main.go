// Package main starts Rivio's command-line application.
//
// The executable delegates command parsing and execution to package cmd. Run
// rivio review from a Git working tree to analyze local changes with the
// configured AI provider.
package main

import "github.com/LucasNav6/rivio/cmd"

// main starts Rivio's command-line interface. It accepts no arguments and
// returns no value; command execution errors are handled by the command package.
func main() {
	cmd.Execute()
}
