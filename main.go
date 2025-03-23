package main

import (
	"net/http"
	"time"

	"claude-think-tool/internal/infra"
	interfacelayer "claude-think-tool/internal/interface"
	"claude-think-tool/internal/usecase"
)

func main() {
	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Initialize infrastructure
	apiClient := infra.NewClaudeAPIClient(httpClient, "")
	fileStorage := infra.NewFileStorage()

	// Initialize use cases
	thinkService := usecase.NewThinkService(apiClient)

	// Initialize interface layer
	formatter := interfacelayer.NewFormatter()
	cli := interfacelayer.NewCLI(thinkService, fileStorage, formatter)

	// Run the application
	cli.Run()
}