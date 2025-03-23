package integration

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"claude-think-tool/internal/domain"
	"claude-think-tool/internal/infra"
	interfacelayer "claude-think-tool/internal/interface"
	"claude-think-tool/internal/usecase"
)

func TestE2E(t *testing.T) {
	// Skip if not running E2E tests or missing API key
	if os.Getenv("RUN_E2E_TESTS") != "1" {
		t.Skip("Skipping E2E test; set RUN_E2E_TESTS=1 to run")
	}

	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		t.Skip("Skipping E2E test; ANTHROPIC_API_KEY environment variable not set")
	}

	// Define test cases
	tests := []struct {
		name        string
		thought     string
		timeout     time.Duration
		expectError bool
	}{
		{
			name:        "simple thought analysis",
			thought:     "I believe this is a good idea because it aligns with our goals.",
			timeout:     60 * time.Second,
			expectError: false,
		},
		{
			name:        "complex thought with risks",
			thought:     "I believe we should launch this feature next week because it improves user engagement, but we haven't done security testing yet.",
			timeout:     60 * time.Second,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Create real http client with timeout
			httpClient := &http.Client{
				Timeout: tt.timeout,
			}

			// Create real API client
			apiClient := infra.NewClaudeAPIClient(httpClient, apiKey)
			
			// Create the service
			thinkService := usecase.NewThinkService(apiClient)
			
			// Create formatter
			formatter := interfacelayer.NewFormatter()
			
			// Create config
			config := domain.Config{
				APIKey:       apiKey,
				Model:        "claude-3-opus-20240229",
				Timeout:      tt.timeout,
				MaxTokens:    1024,
				OutputFormat: "text",
				Verbose:      true,
			}
			
			// Run the core service
			ctx := context.Background()
			ctx, cancel := context.WithTimeout(ctx, tt.timeout)
			defer cancel()
			
			response, err := thinkService.AnalyzeThought(ctx, tt.thought, config)
			
			// Check error expectations
			if tt.expectError && err == nil {
				t.Errorf("Expected error, got nil")
				return
			}
			if !tt.expectError && err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			// If we expect success, verify we have a valid response
			if !tt.expectError {
				if response == nil {
					t.Errorf("Expected non-nil response")
					return
				}
				
				// Format the output and verify it's not empty
				output := formatter.FormatOutput(response, config.OutputFormat)
				if output == "" {
					t.Errorf("Expected non-empty formatted output")
					return
				}
				
				t.Logf("Received response: %s", output)
			}
		})
	}
}