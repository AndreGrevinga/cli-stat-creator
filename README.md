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

## Usage

**Display Options:**
- `-i`, `--default-input` - Use default `data/input.json` without prompting
- `-d`, `--detailed` - Show all statistics columns (min/max values)
- `--no-players` / `--no-levels` - Hide specific statistic sections

**Filtering:**
- `--level <value>` - Single level (`5`) or range (`1-5`)
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
├── cmd/cli-stat-creator/    # Application entry point
├── internal/
│   ├── stats/               # Statistics calculation
│   ├── reader/              # JSON file reading
│   ├── display/             # Table rendering
│   ├── logging/             # Structured logging
│   └── pipeline/            # Data processing pipeline
├── data/                    # Sample input data
└── CLAUDE.md                # AI development guidelines
```
