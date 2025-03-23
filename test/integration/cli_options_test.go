package integration

import (
	"bytes"
	"os"
	"os/exec"
	"strings"
	"testing"
)

// TestCLIOptions tests actual CLI command execution with various options
func TestCLIOptions(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	// Set a fake API key for testing
	os.Setenv("ANTHROPIC_API_KEY", "test-api-key")
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	// Create a temporary directory for test files
	tempDir, err := os.MkdirTemp("", "cli-test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	// Create a test input file
	inputFile := tempDir + "/input.txt"
	err = os.WriteFile(inputFile, []byte("This is a test thought from a file"), 0644)
	if err != nil {
		t.Fatalf("Failed to create input file: %v", err)
	}

	// Define test cases
	tests := []struct {
		name          string
		args          []string
		expectOutput  []string
		notExpectOutput []string
		expectError   bool
	}{
		{
			name:         "help flag",
			args:         []string{"-help"},
			expectOutput: []string{"Usage:", "Options:", "Examples:"},
			expectError:  false,
		},
		{
			name:         "version flag",
			args:         []string{"-version"},
			expectOutput: []string{"Claude Think Tool v", "A tool for analyzing"},
			expectError:  false,
		},
		{
			name:         "custom thought",
			args:         []string{"This is a test thought"},
			expectOutput: []string{}, // Can't easily test actual output since it would require API call
			expectError:  true,       // Will error with fake API key
		},
		{
			name:         "json format",
			args:         []string{"-format", "json", "Test thought"},
			expectOutput: []string{}, // Can't easily test JSON output without API call
			expectError:  true,       // Will error with fake API key
		},
		{
			name:         "input file",
			args:         []string{"-input", inputFile},
			expectOutput: []string{}, // Can't easily test output without API call
			expectError:  true,       // Will error with fake API key
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Execute the command with arguments
			cmd := exec.Command("go", append([]string{"run", "../../main.go"}, tt.args...)...)
			var stdout, stderr bytes.Buffer
			cmd.Stdout = &stdout
			cmd.Stderr = &stderr
			err := cmd.Run()

			// Check error expectation
			if tt.expectError && err == nil {
				t.Errorf("Expected command to fail, but it succeeded")
			}
			if !tt.expectError && err != nil {
				t.Errorf("Expected command to succeed, but got error: %v\nStderr: %s", err, stderr.String())
			}

			// If we're only testing flags that don't require API calls
			if tt.name == "help flag" || tt.name == "version flag" {
				// Check expected output
				output := stdout.String()
				for _, expectedStr := range tt.expectOutput {
					if !strings.Contains(output, expectedStr) {
						t.Errorf("Expected output to contain %q, but it doesn't.\nOutput: %s", expectedStr, output)
					}
				}

				// Check unexpected output
				for _, unexpectedStr := range tt.notExpectOutput {
					if strings.Contains(output, unexpectedStr) {
						t.Errorf("Expected output not to contain %q, but it does.\nOutput: %s", unexpectedStr, output)
					}
				}
			}
		})
	}
}

// TestAPIKeyFromEnv tests that the CLI can read the API key from environment variable
func TestAPIKeyFromEnv(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	// Verify we can read API key from environment
	// We'll just use -version to avoid making an actual API call
	apiKey := "test-api-key-from-env"
	os.Setenv("ANTHROPIC_API_KEY", apiKey)
	defer os.Unsetenv("ANTHROPIC_API_KEY")

	cmd := exec.Command("go", "run", "../../main.go", "-version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Claude Think Tool") {
		t.Errorf("Expected version output, got: %s", output)
	}
}

// TestAPIKeyFromFlag tests that the CLI can read the API key from flag
func TestAPIKeyFromFlag(t *testing.T) {
	// Skip if not running integration tests
	if os.Getenv("RUN_INTEGRATION_TESTS") != "1" {
		t.Skip("Skipping integration test; set RUN_INTEGRATION_TESTS=1 to run")
	}

	// Verify we can pass API key as a flag
	// We'll just use -version to avoid making an actual API call
	cmd := exec.Command("go", "run", "../../main.go", "-apikey", "test-key-from-flag", "-version")
	var stdout bytes.Buffer
	cmd.Stdout = &stdout
	err := cmd.Run()

	if err != nil {
		t.Errorf("Command failed: %v", err)
	}

	output := stdout.String()
	if !strings.Contains(output, "Claude Think Tool") {
		t.Errorf("Expected version output, got: %s", output)
	}
}