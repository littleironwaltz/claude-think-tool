package main

import (
	"bytes"
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Version information
const (
	Version = "0.1.0"
)

// Claude API constants
const (
	AnthropicAPIURL     = "https://api.anthropic.com/v1/messages"
	AnthropicAPIVersion = "2023-06-01"
)

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

// Tool represents a Claude custom tool definition
type Tool struct {
	Type        string                 `json:"type"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	InputSchema map[string]interface{} `json:"input_schema"`
}

// NewThinkTool creates a new instance of the think tool
func NewThinkTool() Tool {
	return Tool{
		Type:        "custom",
		Name:        "think",
		Description: "A tool to analyze and verify thinking processes",
		InputSchema: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"thought": map[string]interface{}{
					"type":        "string",
					"description": "The thought content to be analyzed and verified",
				},
			},
			"required": []string{"thought"},
		},
	}
}

// RunCompleteTool runs a complete tool use cycle with Claude
func RunCompleteTool(ctx context.Context, thought string, config Config) (map[string]interface{}, error) {
	// Get API key from config or environment variable if not set
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("API key not found. Set it using the -apikey flag or ANTHROPIC_API_KEY environment variable")
		}
	}

	// Create the think tool
	thinkTool := NewThinkTool()
	
	// Convert to map for API request
	var toolMap map[string]interface{}
	toolBytes, err := json.Marshal(thinkTool)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal tool: %w", err)
	}
	if err := json.Unmarshal(toolBytes, &toolMap); err != nil {
		return nil, fmt.Errorf("failed to convert tool to map: %w", err)
	}

	// Prepare the user prompt
	userPrompt := thought
	if config.ThoughtPrompt != "" {
		userPrompt = fmt.Sprintf("%s %s", config.ThoughtPrompt, thought)
	} else {
		userPrompt = fmt.Sprintf("Please analyze the following thought: %s", thought)
	}

	// STEP 1: Send initial request
	if config.Verbose {
		fmt.Println("STEP 1: Sending initial request to Claude...")
	}

	initialRequestMap := map[string]interface{}{
		"model":      config.Model,
		"max_tokens": config.MaxTokens,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": userPrompt,
			},
		},
		"tools": []interface{}{toolMap},
	}

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: config.Timeout,
	}

	// Send initial request
	initialResp, err := sendJSONRequest(ctx, client, initialRequestMap, apiKey)
	if err != nil {
		return nil, fmt.Errorf("initial request failed: %w", err)
	}

	if config.Verbose {
		fmt.Println("Initial response:", string(initialResp))
	}

	// Parse the response
	var initialResponseMap map[string]interface{}
	if err := json.Unmarshal(initialResp, &initialResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse initial response: %v", err)
	}

	// Check if Claude wants to use our tool
	stopReason, ok := initialResponseMap["stop_reason"].(string)
	if !ok || stopReason != "tool_use" {
		if config.Verbose {
			fmt.Println("Claude didn't request to use the tool. Returning initial response.")
		}
		return initialResponseMap, nil
	}

	// STEP 2: Extract tool use information
	if config.Verbose {
		fmt.Println("STEP 2: Extracting tool use information...")
	}

	content, ok := initialResponseMap["content"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("content field missing or invalid")
	}

	var toolUseID string
	var toolName string

	for _, item := range content {
		block, ok := item.(map[string]interface{})
		if !ok {
			continue
		}

		blockType, ok := block["type"].(string)
		if !ok || blockType != "tool_use" {
			continue
		}

		toolUseID, _ = block["id"].(string)
		toolName, _ = block["name"].(string)
		break
	}

	if toolUseID == "" || toolName == "" {
		return nil, fmt.Errorf("couldn't find valid tool use block")
	}

	if config.Verbose {
		fmt.Printf("Found tool use: ID=%s, Name=%s\n", toolUseID, toolName)
	}

	// STEP 3: Process the tool request (this would be our actual tool logic)
	if config.Verbose {
		fmt.Println("STEP 3: Processing tool request...")
	}

	// TODO: Replace with actual thought analysis logic
	toolResult := `I've analyzed the thought about launching a new feature. Here are my observations:

Strengths:
- Quantitative data supports benefits (23% engagement, 15% load time)
- Aligns with Q2 goals
- Shows consideration of both benefits and risks

Concerns:
- Incomplete security testing is a significant risk
- Parallel security testing during rollout might identify issues too late
- No mention of rollback plan if security issues are found

Recommendation:
- Complete at least basic security testing before any release
- Consider a phased rollout approach with clear metrics for each phase
- Prepare a contingency plan for security issues`

	// STEP 4: Send follow-up request with tool result
	if config.Verbose {
		fmt.Println("STEP 4: Sending follow-up request with tool result...")
	}

	followUpRequestMap := map[string]interface{}{
		"model":      config.Model,
		"max_tokens": config.MaxTokens,
		"messages": []map[string]interface{}{
			// Original user message
			{
				"role":    "user",
				"content": userPrompt,
			},
			// Assistant's response with tool use
			{
				"role":    "assistant",
				"content": content,
			},
			// Our tool result
			{
				"role": "user",
				"content": []map[string]interface{}{
					{
						"type":        "tool_result",
						"tool_use_id": toolUseID,
						"content":     toolResult,
					},
				},
			},
		},
	}

	// Send follow-up request
	finalResp, err := sendJSONRequest(ctx, client, followUpRequestMap, apiKey)
	if err != nil {
		return nil, fmt.Errorf("follow-up request failed: %w", err)
	}

	if config.Verbose {
		fmt.Println("Final response from Claude:", string(finalResp))
	}

	// Parse final response
	var finalResponseMap map[string]interface{}
	if err := json.Unmarshal(finalResp, &finalResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse final response: %v", err)
	}

	if config.Verbose {
		fmt.Println("Complete tool use cycle successful!")
	}
	return finalResponseMap, nil
}

// sendJSONRequest sends a JSON request to the Claude API
func sendJSONRequest(ctx context.Context, client *http.Client, requestMap map[string]interface{}, apiKey string) ([]byte, error) {
	requestJSON, err := json.Marshal(requestMap)
	if err != nil {
		return nil, fmt.Errorf("failed to serialize request: %w", err)
	}

	req, err := http.NewRequestWithContext(ctx, "POST", AnthropicAPIURL, bytes.NewBuffer(requestJSON))
	if err != nil {
		return nil, fmt.Errorf("failed to create HTTP request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("x-api-key", apiKey)
	req.Header.Set("anthropic-version", AnthropicAPIVersion)

	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("HTTP request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		bodyBytes, readErr := io.ReadAll(resp.Body)
		if readErr != nil {
			return nil, fmt.Errorf("received non-200 response: %d, failed to read body: %w", resp.StatusCode, readErr)
		}
		return nil, fmt.Errorf("received non-200 response: %d, body: %s", resp.StatusCode, string(bodyBytes))
	}

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	return responseData, nil
}

// formatOutput formats the response according to the specified format
func formatOutput(response map[string]interface{}, format string) string {
	switch format {
	case "json":
		jsonBytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error formatting JSON: %v", err)
		}
		return string(jsonBytes)
	case "text":
		// Extract just the text content from Claude's response
		content, ok := response["content"].([]interface{})
		if !ok {
			return "Error: couldn't extract content from response"
		}
		
		var result string
		for _, item := range content {
			block, ok := item.(map[string]interface{})
			if !ok {
				continue
			}
			
			blockType, ok := block["type"].(string)
			if !ok || blockType != "text" {
				continue
			}
			
			text, ok := block["text"].(string)
			if ok {
				result += text + "\n"
			}
		}
		return result
	default:
		// Default to JSON format
		jsonBytes, err := json.MarshalIndent(response, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error formatting output: %v", err)
		}
		return string(jsonBytes)
	}
}

// readFromFile reads content from a file
func readFromFile(filePath string) (string, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return "", fmt.Errorf("failed to read file: %w", err)
	}
	return string(data), nil
}

// writeToFile writes content to a file
func writeToFile(filePath string, content string) error {
	err := os.WriteFile(filePath, []byte(content), 0644)
	if err != nil {
		return fmt.Errorf("failed to write file: %w", err)
	}
	return nil
}

// printVersion prints the version information
func printVersion() {
	fmt.Printf("Claude Think Tool v%s\n", Version)
	fmt.Println("A tool for analyzing and verifying thinking processes with Claude")
	fmt.Println("https://github.com/yourusername/claude-think-tool")
}

// printHelp prints usage information
func printHelp() {
	printVersion()
	fmt.Println("\nUsage:")
	fmt.Println("  claude-think-tool [options] [thought]")
	fmt.Println("\nOptions:")
	flag.PrintDefaults()
	fmt.Println("\nExamples:")
	fmt.Println("  claude-think-tool \"I believe we should launch the feature next week\"")
	fmt.Println("  claude-think-tool -input thoughts.txt -output analysis.json -format json")
	fmt.Println("  claude-think-tool -interactive")
	fmt.Println("\nDocumentation:")
	fmt.Println("  For full documentation, visit: https://github.com/yourusername/claude-think-tool")
}

func main() {
	// Define command line flags
	apiKey := flag.String("apikey", "", "Anthropic API key (default: ANTHROPIC_API_KEY env var)")
	model := flag.String("model", "claude-3-7-sonnet-20250219", "Claude model to use")
	timeout := flag.Duration("timeout", 30*time.Second, "API request timeout")
	maxTokens := flag.Int("max-tokens", 1024, "Maximum tokens in Claude's response")
	inputFile := flag.String("input", "", "Input file containing thought to analyze")
	outputFile := flag.String("output", "", "Output file for analysis results")
	outputFormat := flag.String("format", "text", "Output format (text, json)")
	verbose := flag.Bool("verbose", false, "Verbose output mode")
	interactive := flag.Bool("interactive", false, "Interactive mode")
	version := flag.Bool("version", false, "Print version information")
	help := flag.Bool("help", false, "Print help information")
	thoughtPrompt := flag.String("prompt", "", "Custom prompt template (default: \"Please analyze the following thought: %s\")")
	
	flag.Parse()

	// Print version and exit if requested
	if *version {
		printVersion()
		return
	}
	
	// Print help and exit if requested
	if *help {
		printHelp()
		return
	}
	
	// Create config from flags
	config := Config{
		APIKey:        *apiKey,
		Model:         *model,
		Timeout:       *timeout,
		MaxTokens:     *maxTokens,
		OutputFormat:  *outputFormat,
		Verbose:       *verbose,
		Interactive:   *interactive,
		ThoughtPrompt: *thoughtPrompt,
	}
	
	// Default thought
	defaultThought := "I believe we should launch the new feature next week because our testing shows it improves user engagement by 23% and reduces load times by 15%, which addresses our Q2 goals. The only concern is that we haven't completed security testing, but I think we can do that in parallel during a limited rollout."
	
	// Determine the thought to analyze
	var thought string
	
	if *inputFile != "" {
		// Read thought from file
		var err error
		thought, err = readFromFile(*inputFile)
		if err != nil {
			log.Fatalf("Error reading input file: %v", err)
		}
	} else if flag.NArg() > 0 {
		// Use first non-flag argument as thought
		thought = flag.Arg(0)
	} else if !*interactive {
		// Use default thought if not in interactive mode
		thought = defaultThought
	}
	
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), config.Timeout)
	defer cancel()
	
	// Handle interactive mode
	if *interactive {
		fmt.Println("Claude Think Tool Interactive Mode")
		fmt.Println("Type 'exit' or 'quit' to exit")
		fmt.Println("Enter a thought to analyze:")
		
		for {
			fmt.Print("> ")
			var input string
			scanner := bytes.NewBuffer(nil)
			if _, err := io.Copy(scanner, os.Stdin); err != nil {
				log.Fatalf("Error reading input: %v", err)
			}
			input = scanner.String()
			
			if input == "exit" || input == "quit" {
				break
			}
			
			// Process the thought
			response, err := RunCompleteTool(ctx, input, config)
			if err != nil {
				fmt.Printf("Error: %v\n", err)
				continue
			}
			
			// Format and print the output
			output := formatOutput(response, config.OutputFormat)
			fmt.Println(output)
		}
		
		fmt.Println("Goodbye!")
		return
	}
	
	// Process the thought
	response, err := RunCompleteTool(ctx, thought, config)
	if err != nil {
		log.Fatalf("Think tool call error: %v", err)
	}
	
	// Format the output
	output := formatOutput(response, config.OutputFormat)
	
	// Write to file or print to console
	if *outputFile != "" {
		if err := writeToFile(*outputFile, output); err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}
		fmt.Printf("Analysis written to %s\n", *outputFile)
	} else {
		fmt.Println(output)
	}
}
