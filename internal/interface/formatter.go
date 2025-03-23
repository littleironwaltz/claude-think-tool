package interfacelayer

import (
	"encoding/json"
	"fmt"

	"claude-think-tool/internal/domain"
)

// Formatter handles formatting of responses
type Formatter struct{}

// NewFormatter creates a new formatter
func NewFormatter() *Formatter {
	return &Formatter{}
}

// FormatOutput formats the response according to the specified format
func (f *Formatter) FormatOutput(response *domain.ThinkResponse, format string) string {
	switch format {
	case "json":
		jsonBytes, err := json.MarshalIndent(response.Raw, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error formatting JSON: %v", err)
		}
		return string(jsonBytes)
	case "text":
		// Just return the extracted text content
		return response.Content
	default:
		// Default to JSON format
		jsonBytes, err := json.MarshalIndent(response.Raw, "", "  ")
		if err != nil {
			return fmt.Sprintf("Error formatting output: %v", err)
		}
		return string(jsonBytes)
	}
}