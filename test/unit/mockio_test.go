package unit_test

import (
	"os"
	"testing"
)

func TestMockIO(t *testing.T) {
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
	os.Stdout.Write([]byte("Test output\n"))
	stdoutW.Close()

	// Read from stdout
	var output [1024]byte
	n, err = stdoutR.Read(output[:])
	if err != nil {
		t.Fatalf("Failed to read from stdout: %v", err)
	}

	if string(output[:n]) != "Test output\n" {
		t.Errorf("Expected 'Test output\\n', got %q", string(output[:n]))
	}
}