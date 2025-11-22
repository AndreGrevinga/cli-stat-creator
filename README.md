# CLI Stat Creator

A Go CLI application for analyzing game score statistics. Reads game score data from JSON files and calculates comprehensive statistics including averages, medians, min/max values, and per-level breakdowns.

## Features

- Read game scores from JSON files with player information (name and ID)
- Calculate comprehensive statistics:
  - Total games played and total score
  - Average, median, minimum, and maximum scores
  - Average score by level
- Group scores by player level or by player
- JSON-based data format for easy integration

## Getting Started

### Prerequisites
- Go 1.25.3 or higher

### Installation

Clone the repository:
```bash
git clone https://github.com/AndreGrevinga/cli-stat-creator.git
cd cli-stat-creator
```

### Running the Application

```bash
go run ./cmd/cli-stat-creator
```

The application will read game scores from `data/input.json` by default.

### Building

```bash
go build -o cli-stat-creator ./cmd/cli-stat-creator
./cli-stat-creator
```

### Testing

```bash
go test ./...
```

## CLI Flags

The application supports several command-line flags for customizing output and filtering data:

### Display Options
- `-d`, `--detailed`: Show all statistics columns (includes min/max values)
- `-i`, `--default-input`: Use the default `data/input.json` file without prompting
- `--no-players`: Hide player statistics from output
- `--no-levels`: Hide level statistics from output

### Filtering Options
- `--level <value>`: Filter statistics by level
  - Single level: `--level 5`
  - Range: `--level 1-5` (includes levels 1, 2, 3, 4, 5)
- `--min-score <value>`: Only include scores greater than or equal to this value
- `--max-score <value>`: Only include scores less than or equal to this value

### Examples

```bash
# Use default input file with detailed statistics
./cli-stat-creator -i -d

# Show only level 3 statistics
./cli-stat-creator -i --level 3

# Show statistics for levels 1-5
./cli-stat-creator -i --level 1-5

# Filter scores between 100 and 500
./cli-stat-creator -i --min-score 100 --max-score 500

# Show only overall stats (hide player and level breakdowns)
./cli-stat-creator -i --no-players --no-levels

# Combine filters: levels 2-4 with scores 200+
./cli-stat-creator -i --level 2-4 --min-score 200 -d
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
- `Player`: Type representing a player with name and ID fields
- `GameScore`: Type representing individual game score entries (contains Player struct, Score, Level)
- `Statistics`: Type holding calculated statistics from score data
- `CalculateStatistics(scores []GameScore) (Statistics, error)`: Calculates comprehensive statistics
- `GroupByLevel(scores []GameScore) map[int][]GameScore`: Groups scores by level for analysis
- `GroupByPlayer(scores []GameScore) map[Player][]GameScore`: Groups scores by player for analysis

### internal/reader
- `ReadScoresFromFile(filename string) ([]GameScore, error)`: Reads and parses game scores from JSON file

### internal/display
- `RenderStatistics(s Statistics)`: Renders overall statistics table to stdout
- `RenderLevelBreakdown(s Statistics)`: Renders per-level average scores table to stdout

## License

This project is licensed under the MIT License.
