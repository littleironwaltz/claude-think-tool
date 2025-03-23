package interfacelayer

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"time"

	"claude-think-tool/internal/domain"
)

// Version information
const (
	Version = "0.1.0"
)

// CLI handles command line interface functionality
type CLI struct {
	thinkService domain.ThinkService
	fileStorage  domain.FileStorage
	formatter    *Formatter
}

// NewCLI creates a new CLI instance
func NewCLI(service domain.ThinkService, storage domain.FileStorage, formatter *Formatter) *CLI {
	return &CLI{
		thinkService: service,
		fileStorage:  storage,
		formatter:    formatter,
	}
}

// Run executes the CLI application
func (c *CLI) Run() {
	c.runWithExit(true)
}

// TestRun executes the CLI application without exiting the program (for testing)
func (c *CLI) TestRun() {
	c.runWithExit(false)
}

// runWithExit executes the CLI application with option to exit program
func (c *CLI) runWithExit(shouldExit bool) {
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
		c.printVersion()
		if shouldExit {
			os.Exit(0)
		}
		return
	}
	
	// Print help and exit if requested
	if *help {
		c.printHelp()
		if shouldExit {
			os.Exit(0)
		}
		return
	}
	
	// Create config from flags
	config := domain.Config{
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
		thought, err = c.fileStorage.ReadFromFile(*inputFile)
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
		c.runInteractiveMode(ctx, config)
		return
	}
	
	// Process the thought
	response, err := c.thinkService.AnalyzeThought(ctx, thought, config)
	if err != nil {
		log.Fatalf("Think tool call error: %v", err)
	}
	
	// Format the output
	output := c.formatter.FormatOutput(response, config.OutputFormat)
	
	// Write to file or print to console
	if *outputFile != "" {
		if err := c.fileStorage.WriteToFile(*outputFile, output); err != nil {
			log.Fatalf("Error writing output file: %v", err)
		}
		fmt.Printf("Analysis written to %s\n", *outputFile)
	} else {
		fmt.Println(output)
	}
}

// runInteractiveMode handles interactive CLI mode
func (c *CLI) runInteractiveMode(ctx context.Context, config domain.Config) {
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
		response, err := c.thinkService.AnalyzeThought(ctx, input, config)
		if err != nil {
			fmt.Printf("Error: %v\n", err)
			continue
		}
		
		// Format and print the output
		output := c.formatter.FormatOutput(response, config.OutputFormat)
		fmt.Println(output)
	}
	
	fmt.Println("Goodbye!")
}

// printVersion prints the version information
func (c *CLI) printVersion() {
	fmt.Printf("Claude Think Tool v%s\n", Version)
	fmt.Println("A tool for analyzing and verifying thinking processes with Claude")
	fmt.Println("https://github.com/yourusername/claude-think-tool")
}

// printHelp prints usage information
func (c *CLI) printHelp() {
	c.printVersion()
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