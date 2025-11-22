# CLI Stat Creator - Improvement Plan

## Testing (Priority: High)

- [ ] Add comprehensive unit tests for `ReadScoresFromFile` (valid/invalid JSON, missing files, empty files)
- [ ] Add unit tests for `CalculateStatistics` (edge cases: empty scores, single score, even/odd counts)
- [ ] Add unit tests for `GroupByLevel` function
- [ ] Add benchmark tests for performance-critical functions

## Input Validation & Error Handling

- [ ] Add input validation (file extension check, negative scores, invalid levels)
- [ ] Add better error messages with context and suggestions
- [ ] Add logging support (structured logging with levels)

## CLI Improvements

- [ ] Implement CLI argument parsing using flag package or cobra (file path, output format)
- [ ] Remove commented code at `main.go:102` or implement actual user input
- [ ] Add configuration file support (YAML/JSON) for default settings

## Code Refactoring

- [ ] Extract display logic from main into separate presentation functions
- [ ] Refactor code into separate packages (statistics, io, display, cli)

## New Statistics Features

- [ ] Add more statistics: standard deviation, percentiles (25th, 75th, 90th), mode, variance
- [ ] Add per-player statistics (individual averages, best/worst scores, games played)
- [ ] Add score distribution/histogram feature
- [ ] Add player count and game count per level to statistics
- [ ] Implement median score calculation per level (currently only average)

## Data Management Features

- [ ] Implement export functionality (JSON, CSV, plain text output formats)
- [ ] Add filtering capabilities (by player, level, score range)
- [ ] Add top N players/scores feature with ranking
