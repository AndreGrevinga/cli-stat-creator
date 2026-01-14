# CLI Stat Creator

> **📚 Learning Project Notice**
>
> This project is primarily a learning exercise for exploring Go development with AI assistance.
>
> **Development Approach:**
> - **Implementation code (~80%):** Written by human developer to learn Go hands-on
> - **Tests, documentation, and architecture:** Primarily designed with [Claude Code](https://claude.com/claude-code) assistance
> - AI serves as a mentor for design decisions, best practices, and comprehensive testing strategies
> - All AI-generated code is reviewed, understood, and adapted before integration
>
> **Disclaimer:**
> This is an educational project focused on learning Go and exploring AI-assisted development workflows. The codebase represents a learning journey rather than production-ready software. Use at your own discretion.

---

A Go CLI application for analyzing game score statistics. Reads game score data from JSON files and calculates comprehensive statistics including averages, medians, min/max values, and per-level breakdowns.

## Features

- Read game scores from JSON files with player information (name and ID)
- Calculate comprehensive statistics:
  - Total games played and total score
  - Average, median, minimum, and maximum scores
  - Average score by level
- Group scores by player level or by player
- Structured logging with configurable verbosity levels (debug, info, warn, error)
- JSON-based data format for easy integration

## Quick Start

**Prerequisites:** Go 1.25.3 or higher

```bash
# Clone and navigate
git clone https://github.com/AndreGrevinga/cli-stat-creator.git
cd cli-stat-creator

# Run directly
go run ./cmd/cli-stat-creator -i

# Or build and run
go build -o cli-stat-creator ./cmd/cli-stat-creator
./cli-stat-creator -i

# Run tests
go test ./...
```

## HTTP Server

The project includes a web interface for analyzing game score statistics. The HTTP server provides both a user-friendly web UI and a REST API for programmatic access.

**Build and Run:**
```bash
# Build the server
go build -o http-server ./cmd/http-server

# Run with defaults (localhost:8080)
./http-server

# Or run directly
go run ./cmd/http-server

# Custom port and host
./http-server --port 3000
./http-server --host 0.0.0.0 --port 8080
```

**Available Flags:**
- `-p`, `--port` - Server port (default: 8080)
- `-h`, `--host` - Bind address (default: localhost)
- `-l`, `--log-level` - Logging level: debug, info, warn, error (default: info)
- `--static-dir` - Path to static files directory (default: web/static)

**Web Interface Usage:**

Navigate to `http://localhost:8080` in your browser. The interface allows you to:
- Upload JSON files via drag-and-drop or file picker
- Apply filters (level range, min/max score)
- Toggle detailed statistics view
- View results organized by overall, level, and player statistics

**API Usage:**

The server exposes a REST API at `/api/stats` that accepts JSON files and returns statistics:

```bash
# Basic file upload
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: application/json"

# With filtering parameters
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -F "level=2-4" \
  -F "min-score=200" \
  -F "detailed=true" \
  -H "Accept: application/json"

# Example JSON response
{
  "overall": {
    "totalGamesPlayed": 100,
    "totalScore": 8500,
    "averageScore": 85.0,
    "medianScore": 82.0,
    "minimumScore": 45,
    "maximumScore": 100
  },
  "byLevel": {
    "1": { /* Statistics */ },
    "2": { /* Statistics */ }
  },
  "byPlayer": {
    "Alice (1)": { /* Statistics */ }
  }
}
```

**Query Parameters:**
- `level` - Filter by level (e.g., `"5"` or `"1-5"`)
- `min-score`, `max-score` - Score range filters
- `detailed` - Include min/max values in response
- `include-levels`, `include-players` - Control which statistics sections to include

## Usage

**Display Options:**
- `-i`, `--default-input` - Use default `data/input.json` without prompting
- `-d`, `--detailed` - Show all statistics columns (min/max values)
- `--no-players` / `--no-levels` - Hide specific statistic sections

**Filtering:**
- `--level <value>` - Single level (`5`) or range (`1-5`)
- `--players <value>` - Filter by player names (comma-separated: `Alice,Bob`)
- `--min-score <value>` / `--max-score <value>` - Filter by score range

**Logging:**
- `-l`, `--log-level <level>` - Set verbosity (`debug`, `info`, `warn`, `error`)

**Examples:**
```bash
# Detailed stats with default input
./cli-stat-creator -i -d

# Filter levels 2-4 with minimum score of 200
./cli-stat-creator -i --level 2-4 --min-score 200

# Show only overall stats with debug logging
./cli-stat-creator -i --no-players --no-levels --log-level debug
```

## Input Data Format

The application expects JSON input in the following format:

```json
[
  {
    "player": {
      "name": "Alice",
      "id": 1
    },
    "score": 85,
    "level": 3
  }
]
```

See `data/README.md` for more details about the sample data structure.

## Project Structure

```
cli-stat-creator/
├── cmd/
│   ├── cli-stat-creator/    # CLI application entry point
│   └── http-server/         # HTTP server entry point
├── internal/
│   ├── stats/               # Statistics calculation
│   ├── reader/              # JSON file reading
│   ├── display/             # Table rendering (CLI)
│   ├── logging/             # Structured logging
│   ├── pipeline/            # Data processing pipeline
│   ├── handlers/            # HTTP request handlers
│   ├── render/              # Response rendering (JSON/HTML)
│   └── templates/           # Template management
├── web/
│   ├── static/              # HTML, CSS, JavaScript
│   └── templates/           # Go HTML templates
├── data/                    # Sample input data
└── CLAUDE.md                # AI development guidelines
```
