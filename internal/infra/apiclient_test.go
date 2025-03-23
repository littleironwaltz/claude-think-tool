package infra_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"claude-think-tool/internal/infra"
)

func TestClaudeAPIClient_SendRequest(t *testing.T) {
	tests := []struct {
		name           string
		requestData    map[string]interface{}
		serverResponse map[string]interface{}
		serverStatus   int
		expectError    bool
	}{
		{
			name: "successful request",
			requestData: map[string]interface{}{
				"model": "claude-3-opus-20240229",
				"messages": []map[string]interface{}{
					{"role": "user", "content": "Hello"},
				},
			},
			serverResponse: map[string]interface{}{
				"id":   "msg_123",
				"type": "message",
				"role": "assistant",
				"content": []map[string]interface{}{
					{"type": "text", "text": "Hello, how can I help you?"},
				},
			},
			serverStatus: http.StatusOK,
			expectError:  false,
		},
		{
			name: "server error",
			requestData: map[string]interface{}{
				"model": "invalid-model",
				"messages": []map[string]interface{}{
					{"role": "user", "content": "Hello"},
				},
			},
			serverResponse: map[string]interface{}{
				"error": map[string]interface{}{
					"type":    "invalid_request_error",
					"message": "Model not found",
				},
			},
			serverStatus: http.StatusBadRequest,
			expectError:  true,
		},
	}

	// Create test API client factory that allows overriding the URL
	createTestClient := func(url string) *infra.ClaudeAPIClient {
		client := &http.Client{Timeout: 10 * time.Second}
		apiClient := &infra.ClaudeAPIClient{
			Client:  client,
			APIKey:  "test-api-key",
			BaseURL: url,
		}
		return apiClient
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create a test server
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				// Check headers
				if r.Header.Get("x-api-key") != "test-api-key" {
					t.Errorf("Expected x-api-key header, got %s", r.Header.Get("x-api-key"))
				}
				
				// Set status code and response
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(tt.serverStatus)
				json.NewEncoder(w).Encode(tt.serverResponse)
			}))
			defer server.Close()

			// Create API client, using test server URL
			apiClient := createTestClient(server.URL)

			// Call the API
			ctx := context.Background()
			resp, err := apiClient.SendRequest(ctx, tt.requestData)

			// Check expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}

			// Verify response can be parsed
			var respMap map[string]interface{}
			if err := json.Unmarshal(resp, &respMap); err != nil {
				t.Errorf("Failed to parse response: %v", err)
				return
			}

			// Verify response content
			if respMap["id"] != tt.serverResponse["id"] {
				t.Errorf("Expected id %v, got %v", tt.serverResponse["id"], respMap["id"])
			}
		})
	}
}