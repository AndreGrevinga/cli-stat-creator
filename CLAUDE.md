# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a CLI application for analyzing game score statistics. It reads JSON files containing game scores and calculates various statistics including averages, medians, and per-level breakdowns.

## Build/Run/Test Commands
- Build: `go build -o cli-stat-creator ./cmd/cli-stat-creator`
- Run: `go run ./cmd/cli-stat-creator`
- Test: `go test ./...`
- Test single file: `go test -v path/to/file_test.go`
- Format code: `gofmt -w .`

## Code Style Guidelines
- **Imports**: Group standard library imports first, followed by third-party imports
- **Formatting**: Follow Go standard formatting with `gofmt`
- **Types**: Use clear type definitions with descriptive field names and JSON tags
  - `GameScore`: Represents individual game score entries (Player, Score, Level)
  - `Statistics`: Holds calculated statistics from score data
- **Naming**:
  - Use CamelCase for exported identifiers
  - Use descriptive names for functions and variables
  - Prefix interface names with verb or adjective (e.g., `Reader`)
- **Error Handling**:
  - Always check errors with proper context (e.g., `fmt.Errorf("context: %w", err)`)
  - Return errors instead of logging in functions
  - All file operations should return errors for proper handling
- **Comments**: Follow Go standard with `//` for line comments and `/* */` for package documentation
  - Use complete sentences with proper punctuation

## Project Structure
```
.
├── cmd/
│   └── cli-stat-creator/
│       └── main.go           # Application entry point
├── internal/
│   ├── stats/
│   │   └── stats.go          # Statistics calculation and game score types
│   ├── reader/
│   │   └── json.go           # JSON file reading functionality
│   └── display/
│       └── table.go          # Table rendering and display functions
├── data/                     # Sample input data directory
│   ├── input.json            # Sample game scores in JSON format
│   └── README.md             # Documentation for data structure
├── go.mod                    # Go module definition (cli-stat-creator)
├── README.md                 # Project documentation
└── CLAUDE.md                 # AI assistant guidelines
```

## Key Packages and Functions

### internal/stats
- `GameScore`: Type representing individual game score entries (Player, Score, Level)
- `Statistics`: Type holding calculated statistics from score data
- `CalculateStatistics(scores []GameScore) (Statistics, error)`: Calculates comprehensive statistics including totals, averages, medians, and per-level breakdowns
- `GroupByLevel(scores []GameScore) map[int][]GameScore`: Groups scores by level for analysis

### internal/reader
- `ReadScoresFromFile(filename string) ([]GameScore, error)`: Reads and parses game scores from JSON file

### internal/display
- `RenderStatistics(s Statistics)`: Renders overall statistics table to stdout
- `RenderLevelBreakdown(s Statistics)`: Renders per-level average scores table to stdout

## Data Format
Input JSON should contain an array of game score objects with fields:
- `Player` (string): Player name
- `Score` (int): Game score value
- `Level` (int): Player level
