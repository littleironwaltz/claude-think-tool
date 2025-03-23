package domain

import "time"

// Tool represents a Claude custom tool definition
type Tool struct {
	Type         string                 `json:"type"`
	Name         string                 `json:"name"`
	Description  string                 `json:"description"`
	InputSchema  map[string]interface{} `json:"input_schema"`
}

// Config holds application configuration
type Config struct {
	APIKey        string
	Model         string
	Timeout       time.Duration
	MaxTokens     int
	OutputFormat  string
	Verbose       bool
	Interactive   bool
	ThoughtPrompt string
}

// ThinkResponse represents the structured response from a thought analysis
type ThinkResponse struct {
	Raw     map[string]interface{}
	Content string
}