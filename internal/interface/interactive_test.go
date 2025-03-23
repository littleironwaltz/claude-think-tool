package interfacelayer_test

import (
	"bufio"
	"context"
	"os"
	"testing"
	"time"

	"claude-think-tool/internal/domain"
	interfacelayer "claude-think-tool/internal/interface"
	"claude-think-tool/test/unit"
)

// TestMockStdinStdout is a helper test to ensure we can mock stdin/stdout
func TestMockStdinStdout(t *testing.T) {
	// Save original stdin and stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	// Create a pipe for stdin
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdin pipe: %v", err)
	}
	os.Stdin = r

	// Create a pipe for stdout
	stdoutR, stdoutW, err := os.Pipe()
	if err != nil {
		t.Fatalf("Failed to create stdout pipe: %v", err)
	}
	os.Stdout = stdoutW

	// Write to stdin
	go func() {
		defer w.Close()
		w.Write([]byte("test input\n"))
	}()

	// Read from stdin
	var input [1024]byte
	n, err := os.Stdin.Read(input[:])
	if err != nil {
		t.Fatalf("Failed to read from stdin: %v", err)
	}

	if string(input[:n]) != "test input\n" {
		t.Errorf("Expected 'test input\\n', got %q", string(input[:n]))
	}

	// Write to stdout
	os.Stdout.Write([]byte("test output\n"))
	stdoutW.Close()

	// Read from stdout
	var output [1024]byte
	n, err = stdoutR.Read(output[:])
	if err != nil {
		t.Fatalf("Failed to read from stdout: %v", err)
	}

	if string(output[:n]) != "test output\n" {
		t.Errorf("Expected 'test output\\n', got %q", string(output[:n]))
	}
}

func TestInteractiveModeWithScannerInput(t *testing.T) {
	// Save original stdin and stdout
	oldStdin := os.Stdin
	oldStdout := os.Stdout
	defer func() {
		os.Stdin = oldStdin
		os.Stdout = oldStdout
	}()

	// Create pipes
	stdinReader, stdinWriter, _ := os.Pipe()
	stdoutReader, stdoutWriter, _ := os.Pipe()
	
	// Redirect stdin and stdout
	os.Stdin = stdinReader
	os.Stdout = stdoutWriter

	// Create mock dependencies
	mockService := &unit.MockThinkService{}
	mockFileStorage := &unit.MockFileStorage{}
	formatter := interfacelayer.NewFormatter()
	
	// Set up input and expected thoughts
	inputPrompts := []string{
		"thought 1",
		"thought 2",
		"exit",
	}
	
	// Set up mock service to handle each thought
	callCount := 0
	mockService.AnalyzeThoughtFunc = func(ctx context.Context, thought string, config domain.Config) (*domain.ThinkResponse, error) {
		callCount++
		expectedThoughts := []string{"thought 1", "thought 2"}
		
		if callCount <= len(expectedThoughts) && thought != expectedThoughts[callCount-1] {
			t.Errorf("Expected thought %q for call %d, got %q", expectedThoughts[callCount-1], callCount, thought)
		}
		
		return &domain.ThinkResponse{
			Raw: map[string]interface{}{
				"content": []map[string]interface{}{
					{"type": "text", "text": "Response for: " + thought},
				},
			},
			Content: "Response for: " + thought,
		}, nil
	}

	// Create CLI
	cli := interfacelayer.NewCLI(mockService, mockFileStorage, formatter)
	
	// Run interactive mode in a goroutine
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	
	config := domain.Config{
		APIKey:        "test-key",
		Model:         "test-model",
		MaxTokens:     100,
		OutputFormat:  "text",
	}
	
	// Run the interactive mode in a separate goroutine
	done := make(chan bool)
	go func() {
		cli.RunInteractiveMode(ctx, config)
		done <- true
	}()
	
	// Write inputs to stdin with small delays
	go func() {
		// Let CLI print its welcome message
		time.Sleep(100 * time.Millisecond)
		
		// Feed each input prompt with a small delay
		for _, prompt := range inputPrompts {
			stdinWriter.Write([]byte(prompt + "\n"))
			time.Sleep(100 * time.Millisecond)
		}
		stdinWriter.Close()
	}()
	
	// Read output
	go func() {
		scanner := bufio.NewScanner(stdoutReader)
		for scanner.Scan() {
			// Just consume the output
		}
	}()
	
	// Wait for interactive mode to finish
	select {
	case <-done:
		// Test passes if we get here
	case <-time.After(5 * time.Second):
		t.Fatal("Test timed out")
	}
	
	// Verify the correct number of calls were made
	if callCount != 2 {
		t.Errorf("Expected 2 calls to AnalyzeThought, got %d", callCount)
	}
	
	// Close stdout to allow output reader to complete
	stdoutWriter.Close()
}