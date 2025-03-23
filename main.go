package main

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"time"
)

// Claude API constants
const (
	AnthropicAPIURL     = "https://api.anthropic.com/v1/messages"
	AnthropicAPIVersion = "2023-06-01"
)

// RunCompleteTool runs a complete tool use cycle with Claude
func RunCompleteTool(ctx context.Context, thought string) (map[string]interface{}, error) {
	// Get API key from environment variable
	apiKey := os.Getenv("ANTHROPIC_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("environment variable 'ANTHROPIC_API_KEY' is not set")
	}

	// Define our tool in raw JSON
	rawToolJSON := `{
		"type": "custom",
		"name": "think",
		"description": "A tool to analyze and verify thinking processes",
		"input_schema": {
			"type": "object",
			"properties": {
				"thought": {
					"type": "string",
					"description": "The thought content to be analyzed and verified"
				}
			},
			"required": ["thought"]
		}
	}`

	var toolMap map[string]interface{}
	if err := json.Unmarshal([]byte(rawToolJSON), &toolMap); err != nil {
		return nil, fmt.Errorf("failed to parse tool JSON: %w", err)
	}

	// STEP 1: Send initial request
	fmt.Println("STEP 1: Sending initial request to Claude...")

	initialRequestMap := map[string]interface{}{
		"model":      "claude-3-7-sonnet-20250219",
		"max_tokens": 1024,
		"messages": []map[string]interface{}{
			{
				"role":    "user",
				"content": fmt.Sprintf("Please analyze the following thought: %s", thought),
			},
		},
		"tools": []interface{}{toolMap},
	}

	// Create HTTP client
	client := &http.Client{}

	// Send initial request
	initialResp, err := sendJSONRequest(ctx, client, initialRequestMap, apiKey)
	if err != nil {
		return nil, fmt.Errorf("initial request failed: %w", err)
	}

	fmt.Println("Initial response:", string(initialResp))

	// Parse the response
	var initialResponseMap map[string]interface{}
	if err := json.Unmarshal(initialResp, &initialResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse initial response: %v", err)
	}

	// Check if Claude wants to use our tool
	stopReason, ok := initialResponseMap["stop_reason"].(string)
	if !ok || stopReason != "tool_use" {
		fmt.Println("Claude didn't request to use the tool. Returning initial response.")
		return initialResponseMap, nil
	}

	// STEP 2: Extract tool use information
	fmt.Println("STEP 2: Extracting tool use information...")

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

	fmt.Printf("Found tool use: ID=%s, Name=%s\n", toolUseID, toolName)

	// STEP 3: Process the tool request (this would be our actual tool logic)
	fmt.Println("STEP 3: Processing tool request...")

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
	fmt.Println("STEP 4: Sending follow-up request with tool result...")

	followUpRequestMap := map[string]interface{}{
		"model":      "claude-3-7-sonnet-20250219",
		"max_tokens": 1024,
		"messages": []map[string]interface{}{
			// Original user message
			{
				"role":    "user",
				"content": fmt.Sprintf("Please analyze the following thought: %s", thought),
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

	fmt.Println("Final response from Claude:", string(finalResp))

	// Parse final response
	var finalResponseMap map[string]interface{}
	if err := json.Unmarshal(finalResp, &finalResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse final response: %v", err)
	}

	fmt.Println("Complete tool use cycle successful!")
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

func main() {
	// Create a context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	thought := "I believe we should launch the new feature next week because our testing shows it improves user engagement by 23% and reduces load times by 15%, which addresses our Q2 goals. The only concern is that we haven't completed security testing, but I think we can do that in parallel during a limited rollout."

	// Use the tool cycle implementation
	response, err := RunCompleteTool(ctx, thought)
	if err != nil {
		log.Fatalf("think tool call error: %v", err)
	}

	// Pretty print the response
	prettyJSON, err := json.MarshalIndent(response, "", "  ")
	if err != nil {
		log.Fatalf("failed to format response: %v", err)
	}

	fmt.Println("Tool call successful:")
	fmt.Println(string(prettyJSON))
}
