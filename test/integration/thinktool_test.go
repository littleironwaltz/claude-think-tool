package integration

import (
	"context"
	"os"
	"testing"
	"time"

	"claude-think-tool/internal/domain"
	"claude-think-tool/internal/usecase"
	"claude-think-tool/test/unit"
	interfacelayer "claude-think-tool/internal/interface"
)

func TestIntegrationWithMocks(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	tests := []struct {
		name           string
		thought        string
		expectError    bool
		mockResponses  [][]byte
		mockErrors     []error
	}{
		{
			name:    "complete flow with tool use",
			thought: "I believe we should launch this feature",
			mockResponses: [][]byte{
				createMockResponse("tool_use", true),  // Initial response with tool use
				createMockResponse("end_turn", false), // Final response after tool result
			},
			mockErrors:  []error{nil, nil},
			expectError: false,
		},
		{
			name:    "direct response without tool use",
			thought: "Simple thought",
			mockResponses: [][]byte{
				createMockResponse("end_turn", false), // Direct response without tool use
			},
			mockErrors:  []error{nil},
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set up mock API client
			mockAPIClient := &unit.MockAPIClient{}
			
			// Configure mock responses
			callCount := 0
			mockAPIClient.SendRequestFunc = func(ctx context.Context, requestMap map[string]interface{}) ([]byte, error) {
				defer func() { callCount++ }()
				if callCount < len(tt.mockResponses) {
					return tt.mockResponses[callCount], tt.mockErrors[callCount]
				}
				t.Errorf("Unexpected call to API client (%d calls already made)", callCount)
				return nil, nil
			}
			
			// Create the service with our mock
			thinkService := usecase.NewThinkService(mockAPIClient)
			
			// Create formatter
			formatter := interfacelayer.NewFormatter()
			
			// Create config
			config := domain.Config{
				APIKey:       "test-api-key",
				Model:        "claude-3-opus-20240229",
				Timeout:      30 * time.Second,
				MaxTokens:    1024,
				OutputFormat: "text",
				Verbose:      false,
			}
			
			// Run the core service
			ctx := context.Background()
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
				}
			}
			
			// Verify all expected API calls were made
			if callCount != len(tt.mockResponses) {
				t.Errorf("Expected %d API calls, got %d", len(tt.mockResponses), callCount)
			}
		})
	}
}

// Helper function to create test responses
func createMockResponse(stopReason string, includeToolUse bool) []byte {
	response, _ := unit.CreateMockAPIResponse(stopReason, includeToolUse)
	return response
}