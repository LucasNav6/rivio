// Package openrouter connects prompts to OpenRouter's chat completion API.
package openrouter

import (
	"context"
	"fmt"

	"github.com/LucasNav6/rivio/internal/provider/api/chat"
)

const endpoint = "https://openrouter.ai/api/v1/chat/completions"

// Client runs prompts through the OpenRouter API.
type Client struct {
	chatClient *chat.Client
}

// New creates an OpenRouter API client for the configured model.
//
// It accepts an API key and model name and returns a client or an error if
// either setting is missing.
func New(apiKey, model string) (*Client, error) {
	// 1. Configure the shared chat completion transport for OpenRouter.
	chatClient, err := chat.New(endpoint, apiKey, model)
	if err != nil {
		return nil, fmt.Errorf("initialize OpenRouter API client: %w", err)
	}

	return &Client{chatClient: chatClient}, nil
}

// Run sends prompt to OpenRouter and returns the generated response.
//
// It accepts a context and prompt and returns generated text or an error if
// the API request fails.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	return client.chatClient.Run(ctx, prompt)
}
