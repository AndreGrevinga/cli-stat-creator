# Per-Player and Per-Level Statistics Rendering Design

## Overview

Add rendering functions for per-player and per-level statistics with configurable detail levels controlled via CLI flags.

## Architecture

**Two new rendering functions** in `internal/display/table.go`:
- `RenderPlayerStatistics(stats map[Player]Statistics, detailed bool)`
- `RenderLevelStatistics(stats map[int]Statistics, detailed bool)`

**CLI flag parsing** in `cmd/cli-stat-creator/main.go`:
- Use Go's standard `flag` package
- `--detailed` or `-d` boolean flag (default: false)
- `--no-players` boolean flag (default: false)
- `--no-levels` boolean flag (default: false)

The main function checks these flags and conditionally calls the render functions.

## Table Layouts

### Summary View (detailed=false)

Performance-focused metrics: Games Played, Average Score, Max Score.

**Player Statistics:**
```
| Player Name | Games Played | Avg Score | Max Score |
|-------------|--------------|-----------|-----------|
| Alice       | 15           | 850.50    | 1200      |
| Bob         | 12           | 720.00    | 980       |
```

**Level Statistics:**
```
| Level | Games Played | Avg Score | Max Score |
|-------|--------------|-----------|-----------|
| 1     | 20           | 650.25    | 900       |
| 2     | 18           | 780.00    | 1100      |
```

### Detailed View (detailed=true)

All statistics columns.

**Player Statistics:**
```
| Player Name | Games | Total   | Avg    | Median | Min | Max  |
|-------------|-------|---------|--------|--------|-----|------|
| Alice       | 15    | 12758   | 850.53 | 845.00 | 520 | 1200 |
| Bob         | 12    | 8640    | 720.00 | 715.50 | 450 | 980  |
```

**Level Statistics:**
```
| Level | Games | Total   | Avg    | Median | Min | Max  |
|-------|-------|---------|--------|--------|-----|------|
| 1     | 20    | 13005   | 650.25 | 640.00 | 320 | 900  |
| 2     | 18    | 14040   | 780.00 | 775.50 | 450 | 1100 |
```

## CLI Flag Integration

### Flag Definitions

```go
var (
    detailed  = flag.Bool("detailed", false, "Show all statistics columns")
    noPlayers = flag.Bool("no-players", false, "Hide player statistics")
    noLevels  = flag.Bool("no-levels", false, "Hide level statistics")
)
```

Short form for detailed flag:
```go
flag.BoolVar(detailed, "d", false, "Show all statistics columns (shorthand)")
```

### Main Function Flow

1. `flag.Parse()` at start of main
2. Read file and calculate all statistics (as currently done)
3. Render overall statistics (always shown)
4. If `!*noPlayers` then call `RenderPlayerStatistics(playerStats, *detailed)`
5. If `!*noLevels` then call `RenderLevelStatistics(levelStats, *detailed)`

### Usage Examples

- `./cli-stat-creator` - Shows all tables, summary view
- `./cli-stat-creator -d` - Shows all tables, detailed view
- `./cli-stat-creator --detailed` - Shows all tables, detailed view
- `./cli-stat-creator --no-levels` - Shows overall + player stats only
- `./cli-stat-creator --no-players` - Shows overall + level stats only
- `./cli-stat-creator -d --no-levels` - Detailed player stats only

## Implementation Details

### Sorting

**Player sorting:**
- Extract player names into a slice
- Sort alphabetically with `slices.Sort()`
- Iterate in sorted order to build table rows

**Level sorting:**
- Extract level integers into a slice
- Sort numerically with `slices.Sort()`
- Iterate in sorted order to build table rows

### Table Section Headers

Add a blank line and title before each table for clarity:

```go
fmt.Println("\nPlayer Statistics:")
// render player table

fmt.Println("\nLevel Statistics:")
// render level table
```

### Error Handling

The render functions don't need to return errors since they format already-calculated data. All error handling happens during data loading and statistics calculation in main.

## Files to Modify

1. `internal/display/table.go` - Add `RenderPlayerStatistics` and `RenderLevelStatistics` functions
2. `cmd/cli-stat-creator/main.go` - Add flag parsing and conditional rendering calls
