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

Then run the program:

```bash
go run main.go
```

The program will:
1. Send an initial request to Claude with a sample thought about launching a feature
2. Process Claude's tool use request
3. Send a structured analysis back to Claude
4. Display Claude's final analysis incorporating the thinking step

## Example Output

The tool analyzes a thought process and provides feedback on its strengths and weaknesses:

```
## Analysis of the Feature Launch Decision

### Strengths in the Reasoning
- **Data-Driven Approach**: The decision is supported by concrete metrics
- **Strategic Alignment**: The proposal directly addresses quarterly goals
- **Solution-Oriented**: Offers a compromise rather than just delay

### Concerns
- **Security Testing Gap**: Security testing appears to be treated as secondary
- **Risk Assessment**: The thought minimizes potential security risks
- **Process Sequence**: Security testing should precede user exposure

### Recommendation
This thinking demonstrates good business awareness but potentially undervalues security 
considerations. A stronger approach would be to:
1. Complete essential security testing before any user exposure
2. Design a phased rollout with clear success metrics
3. Develop a contingency plan for security issues discovered during the limited release
```

## Project Structure

- `main.go`: Complete implementation of the think tool cycle
- `README.md`: This file with project overview and usage instructions

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
