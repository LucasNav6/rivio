// Package config loads the settings Rivio uses for code reviews.
//
// The YAML configuration selects both a transport and a provider name, then
// supplies provider-specific model settings and a default Git base branch.
// API credentials are used only by API integrations and are never included in
// logs. CLI integrations use the corresponding locally authenticated tool.
//
// # Configuration File
//
// A configuration file has provider and review sections. MiniMax over its API
// and MiniMax through its CLI are selected as different integrations:
//
//	provider:
//	  type: api
//	  name: minimax
//	  model: MiniMax-M2.5
//	  api_key: your-api-key
//	review:
//	  base: main
//
// OpenRouter uses the same API fields, while CLI integrations use type: cli
// and a name such as claude, codex, antigravity, or minimax.
//
// Load parses this structure into Config for the command and provider packages.
package config

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"gopkg.in/yaml.v3"
)

// Config contains all settings loaded from a Rivio YAML configuration file.
//
// It groups AI provider options with review-specific settings.
type Config struct {
	Provider ProviderConfig `yaml:"provider"`
	Review   ReviewConfig   `yaml:"review"`
}

// ProviderConfig selects a provider integration and supplies its settings.
//
// Type distinguishes API from CLI integrations. Name selects the provider
// within that transport; Model and APIKey are used where supported.
type ProviderConfig struct {
	Type   string `yaml:"type"`
	Name   string `yaml:"name"`
	Model  string `yaml:"model"`
	APIKey string `yaml:"api_key"`
}

// ReviewConfig contains settings that control which changes Rivio reviews.
//
// Its base branch identifies the comparison point for the working tree.
type ReviewConfig struct {
	Base string `yaml:"base"`
}

// Load reads and parses a Rivio YAML configuration file.
//
//	cfg, err := config.Load("~/.config/rivio/config.yml")
//
// It accepts a file path and returns the decoded Config or an error if the
// file cannot be read or its contents are not valid YAML.
func Load(path string) (*Config, error) {
	if path == "~" || strings.HasPrefix(path, "~/") {
		home, err := os.UserHomeDir()
		if err != nil {
			return nil, fmt.Errorf("resolve config path %q: %w", path, err)
		}
		path = filepath.Join(home, strings.TrimPrefix(path, "~/"))
	}

	// 1. Read the configuration file from disk.
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read config %q: %w", path, err)
	}

	// 2. Decode the YAML data into the Rivio configuration structure.
	var cfg Config
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, fmt.Errorf("parse config %q: %w", path, err)
	}

	return &cfg, nil
}
