// Package provider constructs configured API and CLI model clients.
//
// All implementations satisfy Client, so application code can run the same
// prompt without depending on a provider's transport or command-line syntax.
package provider

import "github.com/LucasNav6/rivio/internal/provider/contract"

// Client is the common contract implemented by every provider integration.
type Client = contract.Client
