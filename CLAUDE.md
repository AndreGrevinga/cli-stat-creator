# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build/Run/Test Commands
- Build: `go build -o myproject`
- Run: `go run main.go`
- Test: `go test ./...`
- Test single file: `go test -v path/to/file_test.go`
- Format code: `gofmt -w .`

## Code Style Guidelines
- **Imports**: Group standard library imports first, followed by third-party imports
- **Formatting**: Follow Go standard formatting with `gofmt`
- **Types**: Use clear type definitions with descriptive field names and JSON tags
- **Naming**: 
  - Use CamelCase for exported identifiers
  - Use descriptive names for functions and variables
  - Prefix interface names with verb or adjective (e.g., `Reader`)
- **Error Handling**:
  - Always check errors with proper context (e.g., `fmt.Errorf("context: %w", err)`)
  - Return errors instead of logging in functions
- **Comments**: Follow Go standard with `//` for line comments and `/* */` for package documentation
  - Use complete sentences with proper punctuation

## Project Structure
```
.
├── main.go          # Application entry point
├── go.mod           # Go module definition
├── README.md        # Project documentation
└── CLAUDE.md        # AI assistant guidelines
```
