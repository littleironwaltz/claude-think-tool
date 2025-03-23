package usecase_test

import (
	"context"
	"errors"
	"testing"
	"time"

	"claude-think-tool/internal/domain"
	"claude-think-tool/internal/usecase"
	"claude-think-tool/test/unit"
)

func TestAnalyzeThought(t *testing.T) {
	tests := []struct {
		name           string
		thought        string
		config         domain.Config
		mockResponses  [][]byte
		mockErrors     []error
		expectError    bool
		expectedErrMsg string
	}{
		{
			name:    "successful analysis without tool use",
			thought: "Test thought",
			config: domain.Config{
				APIKey:       "test-key",
				Model:        "test-model",
				Timeout:      30 * time.Second,
				MaxTokens:    1024,
				OutputFormat: "text",
			},
			mockResponses: [][]byte{
				createMockResponse("end_turn", false),
			},
			mockErrors:     []error{nil},
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name:    "successful analysis with tool use",
			thought: "Test thought requiring tool",
			config: domain.Config{
				APIKey:       "test-key",
				Model:        "test-model",
				Timeout:      30 * time.Second,
				MaxTokens:    1024,
				OutputFormat: "text",
			},
			mockResponses: [][]byte{
				createMockResponse("tool_use", true),
				createMockResponse("end_turn", false),
			},
			mockErrors:     []error{nil, nil},
			expectError:    false,
			expectedErrMsg: "",
		},
		{
			name:    "api error on initial request",
			thought: "Test thought",
			config: domain.Config{
				APIKey:       "test-key",
				Model:        "test-model",
				Timeout:      30 * time.Second,
				MaxTokens:    1024,
				OutputFormat: "text",
			},
			mockResponses:  [][]byte{nil},
			mockErrors:     []error{unit.ErrAPIError},
			expectError:    true,
			expectedErrMsg: "initial request failed: API error",
		},
		{
			name:    "api error on follow-up request",
			thought: "Test thought requiring tool",
			config: domain.Config{
				APIKey:       "test-key",
				Model:        "test-model",
				Timeout:      30 * time.Second,
				MaxTokens:    1024,
				OutputFormat: "text",
			},
			mockResponses: [][]byte{
				createMockResponse("tool_use", true),
				nil,
			},
			mockErrors:     []error{nil, unit.ErrAPIError},
			expectError:    true,
			expectedErrMsg: "follow-up request failed: API error",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Setup mock API client
			mockAPIClient := &unit.MockAPIClient{}
			
			// Configure the mock to return different responses for sequential calls
			callCount := 0
			mockAPIClient.SendRequestFunc = func(ctx context.Context, requestMap map[string]interface{}) ([]byte, error) {
				defer func() { callCount++ }()
				if callCount < len(tt.mockResponses) {
					return tt.mockResponses[callCount], tt.mockErrors[callCount]
				}
				return nil, errors.New("unexpected call to SendRequest")
			}
			
			// Create service with mock
			service := usecase.NewThinkService(mockAPIClient)
			
			// Call the service
			ctx := context.Background()
			response, err := service.AnalyzeThought(ctx, tt.thought, tt.config)
			
			// Check error expectations
			if tt.expectError {
				if err == nil {
					t.Errorf("Expected error but got nil")
					return
				}
				if err.Error() != tt.expectedErrMsg {
					t.Errorf("Expected error message %q, got %q", tt.expectedErrMsg, err.Error())
				}
				return
			}
			
			// If we don't expect an error, but got one
			if err != nil {
				t.Errorf("Unexpected error: %v", err)
				return
			}
			
			// Verify response is not nil
			if response == nil {
				t.Errorf("Expected non-nil response")
				return
			}
		})
	}
}

// Helper function to create test responses
func createMockResponse(stopReason string, includeToolUse bool) []byte {
	response, _ := unit.CreateMockAPIResponse(stopReason, includeToolUse)
	return response
}