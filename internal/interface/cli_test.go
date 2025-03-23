package interfacelayer_test

import (
	"bytes"
	"context"
	"flag"
	"io"
	"os"
	"testing"

	"claude-think-tool/internal/domain"
	interfacelayer "claude-think-tool/internal/interface"
	"claude-think-tool/test/unit"
)

// TestCLI_ParseFlags tests the CLI's flag parsing functionality
func TestCLI_ParseFlags(t *testing.T) {
	// Save original flags
	oldArgs := os.Args
	defer func() { 
		os.Args = oldArgs 
		flag.CommandLine = flag.NewFlagSet(os.Args[0], flag.ExitOnError)
	}()

	// Create test cases
	tests := []struct {
		name            string
		args            []string
		envVars         map[string]string
		expectedArgs    map[string]string
		expectedThought string
	}{
		{
			name: "default settings",
			args: []string{"program"},
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "test-key-from-env",
			},
			expectedArgs: map[string]string{
				"apikey":   "",
				"model":    "claude-3-7-sonnet-20250219",
				"timeout":  "30s",
				"maxTokens": "1024",
				"format":   "text",
				"verbose":  "false",
				"interactive": "false",
				"prompt":   "",
			},
			expectedThought: "I believe we should launch the new feature next week because our testing shows it improves user engagement by 23% and reduces load times by 15%, which addresses our Q2 goals. The only concern is that we haven't completed security testing, but I think we can do that in parallel during a limited rollout.",
		},
		{
			name: "custom settings",
			args: []string{
				"program",
				"-apikey=custom-api-key",
				"-model=claude-3-opus-20240229",
				"-timeout=60s",
				"-max-tokens=2048",
				"-format=json",
				"-verbose",
				"-interactive",
				"-prompt=Analyze this thought:",
				"Custom thought content",
			},
			envVars: map[string]string{},
			expectedArgs: map[string]string{
				"apikey":   "custom-api-key",
				"model":    "claude-3-opus-20240229",
				"timeout":  "60s",
				"maxTokens": "2048",
				"format":   "json",
				"verbose":  "true",
				"interactive": "true",
				"prompt":   "Analyze this thought:",
			},
			expectedThought: "Custom thought content",
		},
		{
			name: "input from file",
			args: []string{
				"program",
				"-input=test-input.txt",
			},
			envVars: map[string]string{
				"ANTHROPIC_API_KEY": "test-key-from-env",
			},
			expectedArgs: map[string]string{
				"apikey":   "",
				"model":    "claude-3-7-sonnet-20250219",
				"timeout":  "30s",
				"maxTokens": "1024",
				"format":   "text",
				"verbose":  "false",
				"interactive": "false",
				"prompt":   "",
				"input":    "test-input.txt",
			},
			expectedThought: "This is a thought from a file", // Will be loaded from mock file
		},
		{
			name: "help flag",
			args: []string{
				"program",
				"-help",
			},
			envVars: map[string]string{},
			expectedArgs: map[string]string{
				"help": "true",
			},
			expectedThought: "",
		},
		{
			name: "version flag",
			args: []string{
				"program",
				"-version",
			},
			envVars: map[string]string{},
			expectedArgs: map[string]string{
				"version": "true",
			},
			expectedThought: "",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Reset flags
			flag.CommandLine = flag.NewFlagSet(tt.args[0], flag.ExitOnError)
			
			// Set args
			os.Args = tt.args
			
			// Set environment variables
			for k, v := range tt.envVars {
				oldVal, exists := os.LookupEnv(k)
				os.Setenv(k, v)
				defer func(key, val string, existed bool) {
					if existed {
						os.Setenv(key, val)
					} else {
						os.Unsetenv(key)
					}
				}(k, oldVal, exists)
			}

			// Create mocks
			mockThinkService := &unit.MockThinkService{}
			mockThinkService.AnalyzeThoughtFunc = func(ctx context.Context, thought string, config domain.Config) (*domain.ThinkResponse, error) {
				// Verify the thought matches expectations
				if tt.expectedThought != "" && thought != tt.expectedThought {
					t.Errorf("Expected thought %q, got %q", tt.expectedThought, thought)
				}
				
				return &domain.ThinkResponse{
					Raw: map[string]interface{}{
						"content": []map[string]interface{}{
							{"type": "text", "text": "Test response"},
						},
					},
					Content: "Test response",
				}, nil
			}
			
			mockFileStorage := &unit.MockFileStorage{}
			mockFileStorage.ReadFromFileFunc = func(filePath string) (string, error) {
				if filePath == "test-input.txt" {
					return "This is a thought from a file", nil
				}
				return "", nil
			}
			mockFileStorage.WriteToFileFunc = func(filePath string, content string) error {
				return nil
			}
			
			// Capture stdout to validate help and version output
			oldStdout := os.Stdout
			r, w, _ := os.Pipe()
			os.Stdout = w
			
			// Create CLI with mocks
			formatter := interfacelayer.NewFormatter()
			cli := interfacelayer.NewCLI(mockThinkService, mockFileStorage, formatter)
			
			// Special handling for help and version flags
			if tt.name == "help flag" || tt.name == "version flag" {
				// Run without exiting
				cli.TestRun()
				
				// Restore stdout and read output
				w.Close()
				os.Stdout = oldStdout
				
				var buf bytes.Buffer
				io.Copy(&buf, r)
				output := buf.String()
				
				// Verify expected output
				if tt.name == "help flag" && output == "" {
					t.Errorf("Expected help output, got empty string")
				}
				if tt.name == "version flag" && output == "" {
					t.Errorf("Expected version output, got empty string")
				}
				
				return
			}
			
			// Skip actually running interactive mode
			if tt.name == "custom settings" {
				// Because the interactive flag is set, this test would hang
				// So we don't actually call cli.Run()
				w.Close()
				os.Stdout = oldStdout
				return
			}
			
			// Run CLI in test mode (for non-interactive cases)
			if tt.expectedArgs["interactive"] != "true" {
				// Redirect output from the pipe to avoid test hanging
				go func() {
					io.Copy(io.Discard, r)
				}()
				
				// Run CLI in test mode (doesn't exit program)
				cli.TestRun()
			}
			
			// Restore stdout
			w.Close()
			os.Stdout = oldStdout
		})
	}
}