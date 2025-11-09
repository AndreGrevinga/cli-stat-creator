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
git clone https://github.com/yourusername/cli-stat-creator.git
cd cli-stat-creator
```

### Running the Application

```bash
go run main.go
```

The application will read game scores from `data/input.json` by default.

### Building

```bash
go build -o cli-stat-creator
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
├── main.go          # Application entry point with score analysis functions
├── go.mod           # Go module definition
├── data/            # Sample input data
│   ├── input.json   # Sample game scores
│   └── README.md    # Data structure documentation
├── README.md        # Project documentation
└── CLAUDE.md        # AI assistant guidelines
```

## Functions

- `ReadScoresFromFile(filename string)`: Reads game scores from a JSON file
- `CalculateStatistics(scores []GameScore)`: Calculates comprehensive statistics from game scores
- `GroupByLevel(scores []GameScore)`: Groups scores by player level for level-based analysis

## License

This project is licensed under the MIT License.
