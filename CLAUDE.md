# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview
This is a CLI application for analyzing game score statistics. It reads JSON files containing game scores and calculates various statistics including averages, medians, and per-level breakdowns.

## Build/Run/Test Commands

### CLI Application
- Build: `go build -o cli-stat-creator ./cmd/cli-stat-creator`
- Run: `go run ./cmd/cli-stat-creator -i`

### HTTP Server
- Build: `go build -o http-server ./cmd/http-server`
- Run: `./http-server` or `go run ./cmd/http-server`

### Testing and Formatting
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
- `--max-score <value>`: Only include scores <= this value (default: no limit)

### Logging Options
- `-l`, `--log-level <level>`: Set logging level (debug, info, warn, error; default: warn)

## HTTP Server Flags
The HTTP server supports the following command-line flags:

### Server Configuration
- `-p`, `--port <value>`: Server port (default: 8080)
- `-h`, `--host <value>`: Bind address (default: localhost)
- `--static-dir <value>`: Path to static files directory (default: web/static)

### Logging Options
- `-l`, `--log-level <level>`: Set logging level (debug, info, warn, error; default: info)

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
│   ├── cli-stat-creator/
│   │   └── main.go           # CLI application entry point
│   └── http-server/
│       └── main.go           # HTTP server entry point, router setup, middleware
├── internal/
│   ├── stats/
│   │   └── stats.go          # Statistics calculation and game score types
│   ├── reader/
│   │   └── json.go           # JSON file reading functionality
│   ├── display/
│   │   └── table.go          # Table rendering and display functions (CLI)
│   ├── logging/
│   │   └── logging.go        # Context-based structured logging
│   ├── pipeline/
│   │   ├── pipeline.go       # Data processing pipeline and Results type
│   │   ├── stages.go         # Pipeline stage implementations
│   │   ├── params.go         # Level string parsing utilities
│   │   └── source.go         # ScoreProvider interface and implementations
│   ├── handlers/
│   │   ├── stats.go          # HTTP request handlers for stats endpoint
│   │   └── template_data.go  # Data structures for template rendering
│   ├── render/
│   │   └── json.go           # JSON and HTML response rendering
│   └── templates/
│       └── templates.go      # Template loading and caching
├── web/
│   ├── static/               # HTML, CSS, JavaScript for web interface
│   └── templates/            # Go HTML templates for rendering
├── data/                     # Sample input data directory
│   ├── input.json            # Sample game scores in JSON format
│   └── README.md             # Documentation for data structure
├── go.mod                    # Go module definition (cli-stat-creator)
├── README.md                 # Project documentation
└── CLAUDE.md                 # AI assistant guidelines
```

## Key Packages and Functions

### cmd/cli-stat-creator
- `parseFlags(ctx context.Context) (pipeline.Config, error)`: Validates all command-line flags and returns a pipeline configuration
- `parseLogLevel(flag string, envValue string) slog.Leveler`: Parses log level from flag or environment variable
- `setupLogging(logLevelFlag string) (context.Context, error)`: Initializes logging system with specified log level
- `getInputFilepath() string`: Returns input file path based on flags or user input
- `runPipeline(ctx context.Context, cfg pipeline.Config, filepath string) (pipeline.Results, error)`: Executes the data processing pipeline
- `displayResults(results pipeline.Results, detailed, showPlayers, showLevels bool)`: Renders pipeline results to stdout

### internal/stats
- `Player`: Type representing a game player with name and unique ID
- `GameScore`: Type representing individual game score entries (Player, Score, Level)
- `Statistics`: Type holding calculated statistics from score data (totals, min/max, average, median, standard deviation, variance, mode, percentiles, IQR)
- `Percentiles`: Type holding percentile values (P25, P75, P90) computed from scores
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
- `Config`: Configuration options for filtering game scores in the pipeline (FilterByPlayer, FilterByLevel, MinScore, MaxScore, CalculateByLevel, CalculateByPlayer, ShowDetailed)
- `Results`: Contains all calculated statistics from a pipeline run (Overall, ByLevel, ByPlayer)
- `Pipeline`: Represents a series of processing stages that game scores flow through
- `New(config Config, stages ...Stage) *Pipeline`: Creates a new Pipeline with specified configuration and stages
- `(p *Pipeline) Run(ctx context.Context, provider ScoreProvider) (Results, error)`: Executes the pipeline and returns statistics
- `ScoreProvider`: Interface for sources that can provide game scores
- `FileProvider`: Provides game scores by reading from a file path (implements ScoreProvider)
- `ReaderProvider`: Provides game scores by reading from an io.Reader (implements ScoreProvider)
- `ParseLevelString(level string) ([]int, error)`: Parses level flag string into slice of integers (supports single level "5" or range "1-5")
- `Filter(cfg Config) Stage`: Creates a filtering stage based on configuration
- `Source(ctx context.Context, provider ScoreProvider) (<-chan stats.GameScore, error)`: Reads scores from provider into a channel
- `Aggregate(ctx context.Context, in <-chan stats.GameScore, cfg Config) (Results, error)`: Aggregates scores into statistics

### internal/handlers
- `HandleStats(w http.ResponseWriter, r *http.Request)`: Main endpoint handler for file uploads and statistics processing (POST /api/stats)
- `HandleClear(w http.ResponseWriter, r *http.Request)`: Clears results display, returns empty HTML (GET /api/clear, used by htmx)
- `ResultsData`: Type holding all data needed to render the results template
- `StatsWithDetailed`: Type wrapping Statistics with a ShowDetailed flag for template rendering
- `ErrorData`: Type holding error information for rendering error templates

### internal/render
- `JSON(w http.ResponseWriter, status int, data any)`: Writes JSON response with specified status code and data
- `Error(w http.ResponseWriter, status int, message string, code string)`: Writes JSON error response with status, message, and error code
- `FilterResults(results pipeline.Results, overall, levels, players bool) any`: Creates filtered map of pipeline results based on specified flags

### internal/templates
- `Manager`: Template manager with caching for HTML template rendering (thread-safe with sync.RWMutex)
- `New(dir string) *Manager`: Creates new template manager that loads templates from specified directory
- `(m *Manager) Get(name string) (*template.Template, error)`: Retrieves parsed template from cache or loads from disk
- `(m *Manager) Render(w io.Writer, name string, data any) error`: Executes template with provided data and writes to writer
- `(m *Manager) ClearCache()`: Removes all cached templates, forcing reload on next access

## Data Format
Input JSON should contain an array of game score objects with fields:
- `Player` (object): Player information
  - `name` (string): Player name (required, cannot be empty)
  - `id` (int): Unique player identifier (must be positive)
- `Score` (int): Game score value (must be non-negative)
- `Level` (int): Game level (must be positive)
