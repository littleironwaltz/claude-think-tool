package domain_test

import (
	"testing"
	"time"

	"claude-think-tool/internal/domain"
)

func TestToolCreation(t *testing.T) {
	tests := []struct{
		name        string
		toolType    string
		toolName    string
		description string
		wantType    string
		wantName    string
	}{
		{
			name:        "valid tool creation",
			toolType:    "custom",
			toolName:    "think",
			description: "A test tool",
			wantType:    "custom",
			wantName:    "think",
		},
		{
			name:        "empty type",
			toolType:    "",
			toolName:    "think",
			description: "A test tool",
			wantType:    "",
			wantName:    "think",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tool := domain.Tool{
				Type:        tt.toolType,
				Name:        tt.toolName,
				Description: tt.description,
				InputSchema: make(map[string]interface{}),
			}

			if tool.Type != tt.wantType {
				t.Errorf("Tool.Type = %v, want %v", tool.Type, tt.wantType)
			}

			if tool.Name != tt.wantName {
				t.Errorf("Tool.Name = %v, want %v", tool.Name, tt.wantName)
			}
		})
	}
}

func TestConfigValues(t *testing.T) {
	tests := []struct{
		name           string
		apiKey         string
		model          string
		timeout        time.Duration
		maxTokens      int
		outputFormat   string
		verbose        bool
		interactive    bool
		thoughtPrompt  string
		expectedApiKey string
	}{
		{
			name:           "default config",
			apiKey:         "test-key",
			model:          "claude-3-opus-20240229",
			timeout:        30 * time.Second,
			maxTokens:      1024,
			outputFormat:   "text",
			verbose:        false,
			interactive:    false,
			thoughtPrompt:  "",
			expectedApiKey: "test-key",
		},
		{
			name:           "custom config",
			apiKey:         "custom-key",
			model:          "claude-3-sonnet",
			timeout:        60 * time.Second,
			maxTokens:      2048,
			outputFormat:   "json",
			verbose:        true,
			interactive:    true,
			thoughtPrompt:  "Analyze this:",
			expectedApiKey: "custom-key",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := domain.Config{
				APIKey:        tt.apiKey,
				Model:         tt.model,
				Timeout:       tt.timeout,
				MaxTokens:     tt.maxTokens,
				OutputFormat:  tt.outputFormat,
				Verbose:       tt.verbose,
				Interactive:   tt.interactive,
				ThoughtPrompt: tt.thoughtPrompt,
			}

			if config.APIKey != tt.expectedApiKey {
				t.Errorf("Config.APIKey = %v, want %v", config.APIKey, tt.expectedApiKey)
			}

			if config.Model != tt.model {
				t.Errorf("Config.Model = %v, want %v", config.Model, tt.model)
			}

			if config.Timeout != tt.timeout {
				t.Errorf("Config.Timeout = %v, want %v", config.Timeout, tt.timeout)
			}

			if config.MaxTokens != tt.maxTokens {
				t.Errorf("Config.MaxTokens = %v, want %v", config.MaxTokens, tt.maxTokens)
			}
		})
	}
}

func TestThinkResponse(t *testing.T) {
	tests := []struct{
		name           string
		raw            map[string]interface{}
		content        string
		wantContentLen int
	}{
		{
			name:           "empty response",
			raw:            map[string]interface{}{},
			content:        "",
			wantContentLen: 0,
		},
		{
			name: "typical response",
			raw: map[string]interface{}{
				"id":         "msg_123",
				"stop_reason": "end_turn",
			},
			content:        "Test content",
			wantContentLen: 12,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			response := domain.ThinkResponse{
				Raw:     tt.raw,
				Content: tt.content,
			}

			if len(response.Content) != tt.wantContentLen {
				t.Errorf("ThinkResponse.Content length = %v, want %v", len(response.Content), tt.wantContentLen)
			}

			if response.Raw["id"] != tt.raw["id"] {
				t.Errorf("ThinkResponse.Raw[\"id\"] = %v, want %v", response.Raw["id"], tt.raw["id"])
			}
		})
	}
}