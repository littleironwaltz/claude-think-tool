package infra

import (
	"fmt"
	"os"
)

// FileStorage implements the domain.FileStorage interface
type FileStorage struct{}

// NewFileStorage creates a new file storage implementation
func NewFileStorage() *FileStorage {
	return &FileStorage{}
}

// ReadFromFile reads content from a file
func (fs *FileStorage) ReadFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// WriteToFile writes content to a file
func (fs *FileStorage) WriteToFile(filePath string, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}