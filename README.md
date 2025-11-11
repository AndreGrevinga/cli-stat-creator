# CLI Stat Creator

A Go CLI application for analyzing game score statistics. Reads game score data from JSON files and calculates comprehensive statistics including averages, medians, min/max values, and per-level breakdowns.

## Features

- Read game scores from JSON files
- Calculate comprehensive statistics:
  - Total games played and total score
  - Average, median, minimum, and maximum scores
  - Average score by level
- Group scores by player level
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

## Input Data Format

The application expects JSON input in the following format:

```json
[
  {
    "Player": "PlayerName",
    "Score": 85,
    "Level": 3
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
- `GameScore`: Type representing individual game score entries (Player, Score, Level)
- `Statistics`: Type holding calculated statistics from score data
- `CalculateStatistics(scores []GameScore) (Statistics, error)`: Calculates comprehensive statistics
- `GroupByLevel(scores []GameScore) map[int][]GameScore`: Groups scores by level for analysis

### internal/reader
- `ReadScoresFromFile(filename string) ([]GameScore, error)`: Reads and parses game scores from JSON file

### internal/display
- `RenderStatistics(s Statistics)`: Renders overall statistics table to stdout
- `RenderLevelBreakdown(s Statistics)`: Renders per-level average scores table to stdout

## License

This project is licensed under the MIT License.
