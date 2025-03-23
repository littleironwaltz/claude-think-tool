package interfacelayer_test

import (
	"encoding/json"
	"strings"
	"testing"

	"claude-think-tool/internal/domain"
	interfacelayer "claude-think-tool/internal/interface"
)

func TestOutputFormats(t *testing.T) {
	// Create test cases
	tests := []struct {
		name           string
		format         string
		response       *domain.ThinkResponse
		expectJSON     bool
		shouldContain  string
		shouldNotContain string
	}{
		{
			name:   "text format",
			format: "text",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
					"content": []map[string]interface{}{
						{"type": "text", "text": "This is a test response"},
					},
				},
				Content: "This is a test response",
			},
			expectJSON:      false,
			shouldContain:   "This is a test response",
			shouldNotContain: "\"id\":",
		},
		{
			name:   "json format",
			format: "json",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
					"content": []map[string]interface{}{
						{"type": "text", "text": "This is a test response"},
					},
				},
				Content: "This is a test response",
			},
			expectJSON:      true,
			shouldContain:   "\"id\": \"msg_123\"",
			shouldNotContain: "",
		},
		{
			name:   "default to json for unknown format",
			format: "unknown",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
				},
				Content: "This is a test response",
			},
			expectJSON:      true,
			shouldContain:   "\"id\": \"msg_123\"",
			shouldNotContain: "",
		},
		{
			name:   "complex json response",
			format: "json",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
					"content": []map[string]interface{}{
						{"type": "text", "text": "This is a test response"},
						{"type": "tool_use", "id": "tu_123", "name": "think"},
					},
					"model": "claude-3-7-sonnet-20250219",
				},
				Content: "This is a test response",
			},
			expectJSON:      true,
			shouldContain:   "\"model\":",
			shouldNotContain: "",
		},
	}

	// Create formatter
	formatter := interfacelayer.NewFormatter()

	// Run tests
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Format the output
			output := formatter.FormatOutput(tt.response, tt.format)

			// Verify JSON parsing if expected
			if tt.expectJSON {
				var jsonObj map[string]interface{}
				err := json.Unmarshal([]byte(output), &jsonObj)
				if err != nil {
					t.Errorf("Expected valid JSON, got error: %v", err)
				}
			}

			// Verify content contains expected string
			if tt.shouldContain != "" && !strings.Contains(output, tt.shouldContain) {
				t.Errorf("Expected output to contain %q, got %q", tt.shouldContain, output)
			}

			// Verify content does not contain unexpected string
			if tt.shouldNotContain != "" && strings.Contains(output, tt.shouldNotContain) {
				t.Errorf("Expected output not to contain %q, but it does", tt.shouldNotContain)
			}
		})
	}
}