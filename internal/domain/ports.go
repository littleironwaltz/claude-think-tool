package domain

import "context"

// ThinkService defines the interface for the core thinking analysis service
type ThinkService interface {
	AnalyzeThought(ctx context.Context, thought string, config Config) (*ThinkResponse, error)
}

// APIClient defines the interface for Claude API interaction
type APIClient interface {
	SendRequest(ctx context.Context, requestMap map[string]interface{}) ([]byte, error)
}

// FileStorage defines the interface for file operations
type FileStorage interface {
	ReadFromFile(filePath string) (string, error)
	WriteToFile(filePath string, content string) error
}