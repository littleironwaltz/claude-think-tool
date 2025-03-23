package usecase

import (
	"context"
	"encoding/json"
	"fmt"
	"os"

	"claude-think-tool/internal/domain"
)

// ThinkService implements the domain.ThinkService interface
type ThinkService struct {
	apiClient domain.APIClient
}

// NewThinkService creates a new instance of ThinkService
func NewThinkService(apiClient domain.APIClient) *ThinkService {
	return &ThinkService{
		apiClient: apiClient,
	}
}

// AnalyzeThought runs a complete tool use cycle with Claude to analyze a thought
func (s *ThinkService) AnalyzeThought(ctx context.Context, thought string, config domain.Config) (*domain.ThinkResponse, error) {
	// Get API key from config or environment variable if not set
	apiKey := config.APIKey
	if apiKey == "" {
		apiKey = os.Getenv("ANTHROPIC_API_KEY")
		if apiKey == "" {
			return nil, fmt.Errorf("API key not found. Set it using the -apikey flag or ANTHROPIC_API_KEY environment variable")
		}
	}

	// Create the think tool
	thinkTool := createThinkTool()
	
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

	// Build initial request
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

	// Send initial request
	initialResp, err := s.apiClient.SendRequest(ctx, initialRequestMap)
	if err != nil {
		return nil, fmt.Errorf("initial request failed: %w", err)
	}

	// Parse the response
	var initialResponseMap map[string]interface{}
	if err := json.Unmarshal(initialResp, &initialResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse initial response: %v", err)
	}

	// Check if Claude wants to use our tool
	stopReason, ok := initialResponseMap["stop_reason"].(string)
	if !ok || stopReason != "tool_use" {
		// Format the response and return it
		return formatThinkResponse(initialResponseMap)
	}

	// Extract tool use information
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

	// Process the tool request - in this case, providing an analysis of the thought
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

	// Prepare follow-up request with tool result
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
	finalResp, err := s.apiClient.SendRequest(ctx, followUpRequestMap)
	if err != nil {
		return nil, fmt.Errorf("follow-up request failed: %w", err)
	}

	// Parse final response
	var finalResponseMap map[string]interface{}
	if err := json.Unmarshal(finalResp, &finalResponseMap); err != nil {
		return nil, fmt.Errorf("failed to parse final response: %v", err)
	}

	// Format the response and return it
	return formatThinkResponse(finalResponseMap)
}

// createThinkTool creates a new instance of the think tool
func createThinkTool() domain.Tool {
	return domain.Tool{
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

// formatThinkResponse converts API response to a ThinkResponse
func formatThinkResponse(responseMap map[string]interface{}) (*domain.ThinkResponse, error) {
	// Extract just the text content from Claude's response
	content, ok := responseMap["content"].([]interface{})
	if !ok {
		return nil, fmt.Errorf("couldn't extract content from response")
	}
	
	var textContent string
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
			textContent += text + "\n"
		}
	}

	return &domain.ThinkResponse{
		Raw:     responseMap,
		Content: textContent,
	}, nil
}