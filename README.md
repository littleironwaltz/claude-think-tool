# Claude Think Tool 🧠

A Go implementation that demonstrates Claude's "think" tool - a method for enhancing structured reasoning during complex problem-solving.

## Overview

The "think" tool allows Claude to add an additional thinking step within its response generation process. This implementation connects to Anthropic's Claude API and demonstrates a complete tool use cycle with a custom "think" tool that analyzes decision-making processes and provides structured feedback.

Based on [Anthropic's think tool concept](https://www.anthropic.com/engineering/claude-think-tool), this project shows how to:
- Define a custom tool that creates space for Claude's structured reasoning
- Process Claude's thinking through a complete API interaction cycle
- Enhance complex decision analysis with additional reasoning steps

## When to Use This Tool

The think tool is particularly effective for:
- Analyzing complex decision-making processes
- Following multi-step policy guidelines
- Making sequential decisions
- Evaluating tool outputs and next actions
- Providing structured analysis with clear recommendations

## Features

- Complete implementation of Claude's custom tool API
- Demonstrates the full tool use cycle with Claude
- Provides structured analysis of thought processes
- Implements proper error handling and context support
- Follows Go best practices and idiomatic patterns
- CLI interface with flexible configuration options
- Interactive mode for conversational thought analysis
- File input/output support for batch processing
- Multiple output formats (text, JSON)
- Custom prompt templates

## How It Works

The think tool works by:
1. Defining a custom tool that accepts a "thought" parameter
2. Sending this tool definition to Claude with a prompt
3. Claude using the tool to analyze the thought process
4. Processing Claude's tool use request
5. Sending an analysis result back to Claude
6. Claude incorporating the analysis into its final response

## Technical Implementation

The tool is defined with this schema:

```json
{
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
}
```

## Requirements

- Go 1.18 or higher
- Valid Anthropic API key

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/claude-think-tool.git
   cd claude-think-tool
   ```

2. Build the executable:
   ```bash
   go build
   ```

## Usage

Before running the tool, make sure to set your Anthropic API key as an environment variable:

```bash
export ANTHROPIC_API_KEY="your-api-key-here"
```

### Basic Usage

Run the program with a custom thought:

```bash
go run main.go "I think we should launch the new product immediately because our competitor just released a similar feature"
```

Or use the default example thought:

```bash
go run main.go
```

### Command Line Options

The tool supports various command-line options:

```
Usage:
  claude-think-tool [options] [thought]

Options:
  -apikey string
        Anthropic API key (default: ANTHROPIC_API_KEY env var)
  -format string
        Output format (text, json) (default "text")
  -help
        Print help information
  -input string
        Input file containing thought to analyze
  -interactive
        Interactive mode
  -max-tokens int
        Maximum tokens in Claude's response (default 1024)
  -model string
        Claude model to use (default "claude-3-7-sonnet-20250219")
  -output string
        Output file for analysis results
  -prompt string
        Custom prompt template (default: "Please analyze the following thought: %s")
  -timeout duration
        API request timeout (default 30s)
  -verbose
        Verbose output mode
  -version
        Print version information
```

### Examples

Analyze a thought and save the result to a file in JSON format:
```bash
go run main.go -output analysis.json -format json "I believe we should launch this feature"
```

Read a thought from a file:
```bash
go run main.go -input thought.txt
```

Use interactive mode for continuous analysis:
```bash
go run main.go -interactive
```

Use a custom prompt template:
```bash
go run main.go -prompt "Critically evaluate this hypothesis:" "Our new marketing strategy will increase conversion rates by 25%"
```

## Example Output

The tool analyzes a thought process and provides feedback on its strengths and weaknesses:

```
## Analysis of Your Launch Reasoning

Your thinking demonstrates a data-driven approach with clear business justification, but contains an important blindspot regarding security.

### Strengths
- **Data-backed decision making**: You've cited specific metrics (23% engagement improvement, 15% faster load times) that support the launch
- **Strategic alignment**: You've connected the launch to Q2 goals, showing business purpose
- **Solution-oriented**: You've proposed a pragmatic approach (limited rollout with parallel testing)

### Concerns
- **Security risk acceptance**: Launching without completed security testing creates significant risk
- **Timing assumption**: The parallel testing approach assumes security issues would be discovered early in the limited rollout
- **Missing contingency**: There's no mentioned plan for handling security issues discovered post-launch

### Recommendation
I'd recommend modifying your approach to prioritize basic security testing before any public release. Consider a more structured phased rollout with clear evaluation criteria at each stage and prepare contingency plans for potential security vulnerabilities that might be discovered.
```

## Project Structure

- `main.go`: Complete implementation of the think tool cycle
- `README.md`: This file with project overview and usage instructions
- `CLAUDE.md`: Build instructions and code style guidelines
- `IMPROVEMENTS.md`: Completed and planned improvements

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## License

MIT License

## References

- [Claude Think Tool - Anthropic Engineering Blog](https://www.anthropic.com/engineering/claude-think-tool)
- [Claude API Documentation](https://docs.anthropic.com/claude/docs/tool-use)
