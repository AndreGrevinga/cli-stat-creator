# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a CLI application for analyzing game score statistics. It reads JSON files containing game scores and calculates various statistics including averages, medians, and per-level breakdowns.

## Build/Run/Test Commands
- Build: `go build -o cli-stat-creator ./cmd/cli-stat-creator`
- Run: `go run ./cmd/cli-stat-creator -i`
- Test: `go test ./...`
- Test single file: `go test -v path/to/file_test.go`
- Format code: `gofmt -w .`

## CLI Flags
The application supports the following command-line flags:

### Display Options
- `-d`, `--detailed`: Show all statistics columns (includes min/max values)
- `-i`, `--default-input`: Use the default `data/input.json` file without prompting
- `--no-players`: Hide player statistics from output
- `--no-levels`: Hide level statistics from output

### Filtering Options
- `--level <value>`: Filter statistics by level (single: "5" or range: "1-5")
- `--players <value>`: Filter statistics by player names (comma-separated: "Alice,Bob")
- `--min-score <value>`: Only include scores >= this value (default: 0)
- `--max-score <value>`: Only include scores <= this value (default: 0)

### Logging Options
- `-l`, `--log-level <level>`: Set logging level (debug, info, warn, error; default: warn)

## Code Style Guidelines
- **Imports**: Use `gofmt` which sorts imports alphabetically within a single block
- **Formatting**: Follow Go standard formatting with `gofmt`
- **Types**: Use clear type definitions with descriptive field names and JSON tags
  - `Player`: Represents a player with name and ID fields
  - `GameScore`: Represents individual game score entries (contains Player struct, Score, Level)
  - `Statistics`: Holds calculated statistics from score data
- **Naming**:
  - Use CamelCase for exported identifiers
  - Use descriptive names for functions and variables
  - Prefix interface names with verb or adjective (e.g., `Reader`)
  - Use short names for short scopes (e.g., `i` for loop indices, `n` for counts)
  - Accept common Go abbreviations: `err`, `ctx`, `buf`, `req`, `resp`, `src`, `dst`, `msg`, `cfg`
  - Avoid abbreviations outside of these standard idioms; prefer full words for clarity
- **Error Handling**:
  - Always check errors with proper context (e.g., `fmt.Errorf("context: %w", err)`)
  - Return errors instead of logging in functions
  - All file operations should return errors for proper handling
- **Comments**: Follow Go standard with `//` for line comments and `/* */` for package documentation
  - Use complete sentences with proper punctuation

## Code Review Guidelines
When reviewing or implementing larger changes:
- **Documentation Updates**: Check if README.md, CLAUDE.md, or other documentation needs to be updated to reflect the changes
- **Go Doc Comments**: Ensure every exported function, type, and struct has appropriate Go doc documentation
  - Doc comments should start with the identifier name
  - Use complete sentences that describe what the function/type does

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
│   ├── display/
│   │   └── table.go          # Table rendering and display functions
│   ├── logging/
│   │   └── logging.go        # Context-based structured logging
│   └── pipeline/
│       ├── pipeline.go       # Data processing pipeline
│       └── stages.go         # Pipeline stage implementations
├── data/                     # Sample input data directory
│   ├── input.json            # Sample game scores in JSON format
│   └── README.md             # Documentation for data structure
├── go.mod                    # Go module definition (cli-stat-creator)
├── README.md                 # Project documentation
└── CLAUDE.md                 # AI assistant guidelines
```

## Key Packages and Functions

### cmd/cli-stat-creator
- `parseLevelFlag(level string) ([]int, error)`: Parses level flag string into slice of integers (supports single level "5" or range "1-5")

### internal/stats
- `Player`: Type representing a game player with name and unique ID
- `GameScore`: Type representing individual game score entries (Player, Score, Level)
- `Statistics`: Type holding calculated statistics from score data
- `(g GameScore) Validate() error`: Validates game score fields (player name not empty, score non-negative, level and player ID positive)
- `CalculateStatistics(scores []GameScore) (Statistics, error)`: Calculates comprehensive statistics including totals, averages, medians
- `GroupByLevel(scores []GameScore) map[int][]GameScore`: Groups scores by level for analysis
- `GroupByPlayer(scores []GameScore) map[Player][]GameScore`: Groups scores by player for analysis
- `CalculateStatisticsByLevel(scores []GameScore) (map[int]Statistics, error)`: Calculates statistics for each level separately
- `CalculateStatisticsByPlayer(scores []GameScore) (map[Player]Statistics, error)`: Calculates statistics for each player separately

### internal/reader
- `ReadScoresFromFile(ctx context.Context, filename string) ([]GameScore, error)`: Reads and parses game scores from JSON file with context-based logging (validates .json extension and all game scores)

### internal/display
- `RenderStatistics(s Statistics)`: Renders overall statistics table to stdout
- `RenderGroupedStatistics[K comparable](...)`: Generic function for rendering statistics grouped by any comparable key type (used by level and player renderers)
- `RenderLevelStatistics(statistics map[int]Statistics, detailed bool)`: Renders per-level statistics table to stdout with optional detailed view
- `RenderPlayerStatistics(statistics map[Player]Statistics, detailed bool)`: Renders per-player statistics table to stdout with optional detailed view

### internal/logging
- `WithLogger(ctx context.Context, logger *slog.Logger) context.Context`: Adds a structured logger to a context
- `FromContext(ctx context.Context) *slog.Logger`: Extracts logger from context, returns no-op logger if not found (graceful degradation)

### internal/pipeline
- Pipeline stages for data processing with structured logging support
- Implements stages for reading, filtering, calculating statistics, and rendering output

## Data Format
Input JSON should contain an array of game score objects with fields:
- `Player` (object): Player information
  - `name` (string): Player name (required, cannot be empty)
  - `id` (int): Unique player identifier (must be positive)
- `Score` (int): Game score value (must be non-negative)
- `Level` (int): Game level (must be positive)
