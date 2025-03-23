package infra_test

import (
	"os"
	"path/filepath"
	"testing"

	"claude-think-tool/internal/infra"
)

func TestFileStorage(t *testing.T) {
	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "filestorage_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create the file storage implementation
	storage := infra.NewFileStorage()

	t.Run("successful read and write", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "test_success.txt")
		content := "test content"

		// Write the file
		err := storage.WriteToFile(filePath, content)
		if err != nil {
			t.Fatalf("Failed to write file: %v", err)
		}

		// Read the file
		readContent, err := storage.ReadFromFile(filePath)
		if err != nil {
			t.Fatalf("Failed to read file: %v", err)
		}

		// Verify content
		if readContent != content {
			t.Errorf("Expected content %q, got %q", content, readContent)
		}
	})

	t.Run("read nonexistent file", func(t *testing.T) {
		filePath := filepath.Join(tempDir, "nonexistent.txt")
		
		// Try to read a nonexistent file
		_, err := storage.ReadFromFile(filePath)
		if err == nil {
			t.Errorf("Expected error reading nonexistent file, got nil")
		}
	})

	t.Run("write to invalid location", func(t *testing.T) {
		// Create a directory instead of a file
		dirPath := filepath.Join(tempDir, "test_dir")
		err := os.Mkdir(dirPath, 0755)
		if err != nil {
			t.Fatalf("Failed to create directory: %v", err)
		}

		// Try to write to a path that's a directory
		err = storage.WriteToFile(dirPath, "test content")
		if err == nil {
			t.Errorf("Expected error writing to directory path, got nil")
		}
	})
}