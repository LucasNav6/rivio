// Package contract defines the interface shared by provider integrations.
package contract

import "context"

// Client sends a prompt to a configured model and returns its response.
//
// Implementations hide transport details such as HTTP APIs or local CLI
// processes from the application.
type Client interface {
	Run(ctx context.Context, prompt string) (string, error)
}
