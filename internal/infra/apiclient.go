package infra

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// Constants for Claude API
const (
	AnthropicAPIURL     = "https://api.anthropic.com/v1/messages"
	AnthropicAPIVersion = "2023-06-01"
)

// ClaudeAPIClient implements the domain.APIClient interface
type ClaudeAPIClient struct {
	Client  *http.Client
	APIKey  string
	BaseURL string // Can be overridden for testing
}

// NewClaudeAPIClient creates a new API client for Claude
func NewClaudeAPIClient(client *http.Client, apiKey string) *ClaudeAPIClient {
	return &ClaudeAPIClient{
		Client:  client,
		APIKey:  apiKey,
		BaseURL: AnthropicAPIURL,
	}
}

// SendRequest sends a JSON request to the Claude API
func (c *ClaudeAPIClient) SendRequest(ctx context.Context, requestMap map[string]interface{}) ([]byte, error) {
	requestJSON, err := json.Marshal(requestMap)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", c.BaseURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", c.APIKey)
	req.Header.Set("anthropic-version", AnthropicAPIVersion)

	resp, err := c.Client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("received non-200 response: %d, failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("received non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseData, nil
}