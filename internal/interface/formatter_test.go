package interfacelayer_test

import (
	"encoding/json"
	"strings"
	"testing"

	"claude-think-tool/internal/domain"
	interfacelayer "claude-think-tool/internal/interface"
)

func TestFormatter_FormatOutput(t *testing.T) {
	tests := []struct{
		name            string
		response        *domain.ThinkResponse
		format          string
		expectJSON      bool
		expectedContent string
	}{
		{
			name: "text format",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
					"content": []map[string]interface{}{
						{"type": "text", "text": "This is a test response"},
					},
				},
				Content: "This is a test response",
			},
			format:          "text",
			expectJSON:      false,
			expectedContent: "This is a test response",
		},
		{
			name: "json format",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
					"content": []map[string]interface{}{
						{"type": "text", "text": "This is a test response"},
					},
				},
				Content: "This is a test response",
			},
			format:          "json",
			expectJSON:      true,
			expectedContent: "id",
		},
		{
			name: "default format falls back to json",
			response: &domain.ThinkResponse{
				Raw: map[string]interface{}{
					"id": "msg_123",
				},
				Content: "This is a test response",
			},
			format:          "unknown",
			expectJSON:      true,
			expectedContent: "id",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			formatter := interfacelayer.NewFormatter()
			output := formatter.FormatOutput(tt.response, tt.format)
			
			if tt.expectJSON {
				// Verify it's valid JSON
				var jsonObj map[string]interface{}
				err := json.Unmarshal([]byte(output), &jsonObj)
				if err != nil {
					t.Errorf("Expected valid JSON, got error: %v", err)
					return
				}
				
				// Verify it contains expected field
				if _, ok := jsonObj[tt.expectedContent]; !ok {
					t.Errorf("Expected JSON to contain %q", tt.expectedContent)
				}
			} else {
				// For text format, just check content
				if !strings.Contains(output, tt.expectedContent) {
					t.Errorf("Expected output to contain %q, got %q", tt.expectedContent, output)
				}
			}
		})
	}
}