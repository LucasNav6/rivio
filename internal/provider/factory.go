// Package provider constructs configured API and CLI model clients.
package provider

import (
	"fmt"
	"strings"

	"github.com/LucasNav6/rivio/internal/config"
	minimaxapi "github.com/LucasNav6/rivio/internal/provider/api/minimax"
	"github.com/LucasNav6/rivio/internal/provider/api/openrouter"
	"github.com/LucasNav6/rivio/internal/provider/cli/antigravity"
	"github.com/LucasNav6/rivio/internal/provider/cli/claude"
	"github.com/LucasNav6/rivio/internal/provider/cli/codex"
	minimaxcli "github.com/LucasNav6/rivio/internal/provider/cli/minimax"
)

// New creates the API or CLI client selected by cfg.
//
// It accepts provider configuration and returns a common Client or an error if
// the integration is unsupported or cannot be initialized.
func New(cfg config.ProviderConfig) (Client, error) {
	// 1. Normalize the configured transport and provider names.
	providerType := strings.ToLower(strings.TrimSpace(cfg.Type))
	providerName := strings.ToLower(strings.TrimSpace(cfg.Name))

	// 2. Construct the selected provider implementation.
	switch providerType {
	case "api":
		switch providerName {
		case "minimax":
			return minimaxapi.New(cfg.APIKey, cfg.Model)
		case "openrouter":
			return openrouter.New(cfg.APIKey, cfg.Model)
		default:
			return nil, fmt.Errorf("unsupported API provider: %q", cfg.Name)
		}
	case "cli":
		switch providerName {
		case "antigravity":
			return antigravity.New(cfg.Model)
		case "claude":
			return claude.New(cfg.Model)
		case "codex":
			return codex.New(cfg.Model)
		case "minimax":
			return minimaxcli.New()
		default:
			return nil, fmt.Errorf("unsupported CLI provider: %q", cfg.Name)
		}
	default:
		return nil, fmt.Errorf("unsupported provider type: %q", cfg.Type)
	}
}
