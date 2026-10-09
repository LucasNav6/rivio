// Package chat sends prompts to OpenAI-compatible chat completion APIs.
package chat

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

// Client sends prompts to a configured chat completion endpoint.
type Client struct {
	endpoint   string
	apiKey     string
	model      string
	httpClient *http.Client
}

// request is the JSON payload sent to a chat completion endpoint.
type request struct {
	Model    string    `json:"model"`
	Messages []message `json:"messages"`
}

// message contains one role and prompt or response in the chat protocol.
type message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

// response contains the completion choices returned by a chat endpoint.
type response struct {
	Choices []struct {
		Message message `json:"message"`
	} `json:"choices"`
}

// New creates a client for an OpenAI-compatible chat completion endpoint.
//
// It accepts an endpoint URL, API key, and model name and returns a configured
// client or an error if any required value is empty.
func New(endpoint, apiKey, model string) (*Client, error) {
	// 1. Normalize the endpoint, credential, and model.
	endpoint = strings.TrimSpace(endpoint)
	apiKey = strings.TrimSpace(apiKey)
	model = strings.TrimSpace(model)
	if endpoint == "" || apiKey == "" || model == "" {
		return nil, fmt.Errorf("chat API endpoint, API key, and model are required")
	}

	// 2. Create the configured HTTP client.
	return &Client{
		endpoint: endpoint,
		apiKey:   apiKey,
		model:    model,
		httpClient: &http.Client{
			Timeout: 120 * time.Second,
		},
	}, nil
}

// Run sends prompt to the configured model and returns its response.
//
// It accepts a context and prompt and returns generated text or an error if
// the HTTP request or response cannot be processed.
func (client *Client) Run(ctx context.Context, prompt string) (string, error) {
	// 1. Require prompt text before sending an API request.
	prompt = strings.TrimSpace(prompt)
	if prompt == "" {
		return "", fmt.Errorf("chat API prompt is required")
	}

	// 2. Encode the prompt as a chat completion request.
	body, err := json.Marshal(request{
		Model:    client.model,
		Messages: []message{{Role: "user", Content: prompt}},
	})
	if err != nil {
		return "", fmt.Errorf("encode chat request: %w", err)
	}

	// 3. Create and send the authenticated HTTP request.
	httpRequest, err := http.NewRequestWithContext(ctx, http.MethodPost, client.endpoint, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("create chat request: %w", err)
	}
	httpRequest.Header.Set("Authorization", "Bearer "+client.apiKey)
	httpRequest.Header.Set("Content-Type", "application/json")

	httpResponse, err := client.httpClient.Do(httpRequest)
	if err != nil {
		return "", fmt.Errorf("send chat request: %w", err)
	}
	defer httpResponse.Body.Close()

	// 4. Reject unsuccessful responses and decode the first completion.
	if httpResponse.StatusCode < http.StatusOK || httpResponse.StatusCode >= http.StatusMultipleChoices {
		return "", fmt.Errorf("chat request failed with HTTP %d", httpResponse.StatusCode)
	}
	var completion response
	if err := json.NewDecoder(io.LimitReader(httpResponse.Body, 16<<20)).Decode(&completion); err != nil {
		return "", fmt.Errorf("decode chat response: %w", err)
	}
	if len(completion.Choices) == 0 || strings.TrimSpace(completion.Choices[0].Message.Content) == "" {
		return "", fmt.Errorf("chat response did not contain generated content")
	}

	// 5. Return the validated model response.
	return strings.TrimSpace(completion.Choices[0].Message.Content), nil
}
