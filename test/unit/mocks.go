package unit

import (
	"context"
	"encoding/json"
	"errors"

	"claude-think-tool/internal/domain"
)

// MockAPIClient implements domain.APIClient for testing
type MockAPIClient struct {
	SendRequestFunc func(ctx context.Context, requestMap map[string]interface{}) ([]byte, error)
}

// SendRequest calls the mocked function
func (m *MockAPIClient) SendRequest(ctx context.Context, requestMap map[string]interface{}) ([]byte, error) {
	return m.SendRequestFunc(ctx, requestMap)
}

// MockFileStorage implements domain.FileStorage for testing
type MockFileStorage struct {
	ReadFromFileFunc func(filePath string) (string, error)
	WriteToFileFunc  func(filePath string, content string) error
}

// ReadFromFile calls the mocked function
func (m *MockFileStorage) ReadFromFile(filePath string) (string, error) {
	return m.ReadFromFileFunc(filePath)
}

// WriteToFile calls the mocked function
func (m *MockFileStorage) WriteToFile(filePath string, content string) error {
	return m.WriteToFileFunc(filePath, content)
}

// MockThinkService implements domain.ThinkService for testing
type MockThinkService struct {
	AnalyzeThoughtFunc func(ctx context.Context, thought string, config domain.Config) (*domain.ThinkResponse, error)
}

// AnalyzeThought calls the mocked function
func (m *MockThinkService) AnalyzeThought(ctx context.Context, thought string, config domain.Config) (*domain.ThinkResponse, error) {
	return m.AnalyzeThoughtFunc(ctx, thought, config)
}

// Helper function to create mock Claude API responses
func CreateMockAPIResponse(stopReason string, includeToolUse bool) ([]byte, error) {
	content := []map[string]interface{}{}
	
	if includeToolUse {
		content = append(content, map[string]interface{}{
			"type": "tool_use",
			"id":   "tu_123",
			"name": "think",
		})
	} else {
		content = append(content, map[string]interface{}{
			"type": "text",
			"text": "This is a test response",
		})
	}
	
	response := map[string]interface{}{
		"id":          "msg_123",
		"type":        "message",
		"role":        "assistant",
		"content":     content,
		"stop_reason": stopReason,
		"model":       "claude-3-opus-20240229",
	}
	
	return json.Marshal(response)
}

// Common errors
var (
	ErrNotFound = errors.New("not found")
	ErrAPIError = errors.New("API error")
)