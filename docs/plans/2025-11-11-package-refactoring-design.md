# Package Refactoring Design

## Overview

Refactor the single-file CLI application into a multi-package structure following Go conventions. This refactoring aims to:
- Teach Go package organization and project structure
- Improve testability by separating concerns
- Prepare the codebase for future growth
- Follow standard Go project layout patterns

## Target Directory Structure

```
.
├── cmd/
│   └── cli-stat-creator/
│       └── main.go           # Application entry point
├── internal/
│   ├── stats/
│   │   └── stats.go          # Core types & calculations
│   ├── display/
│   │   └── table.go          # Rendering logic
│   └── reader/
│       └── json.go           # File I/O
├── data/
│   └── input.json            # Sample data (unchanged)
├── go.mod
├── CLAUDE.md
└── README.md
```

## Package Responsibilities

### internal/stats
**Purpose:** Core business logic and domain types

**Contains:**
- `GameScore` type definition
- `Statistics` type definition
- `CalculateStatistics(scores []GameScore) (Statistics, error)`
- `GroupByLevel(scores []GameScore) map[int][]GameScore`

**Dependencies:** None (pure business logic)

**Why here:** Domain types and calculations are the heart of the application. They should be independent of I/O and presentation concerns.

### internal/reader
**Purpose:** File I/O and data parsing

**Contains:**
- `ReadScoresFromFile(filename string) ([]stats.GameScore, error)`

**Dependencies:** `internal/stats` (for GameScore type)

**Why here:** Separates the concern of "how we get data" from "what we do with data". Makes it easy to add new readers (CSV, XML, API) later.

### internal/display
**Purpose:** Output formatting and rendering

**Contains:**
- `RenderStatistics(s stats.Statistics)` - renders overall statistics table
- `RenderLevelBreakdown(s stats.Statistics)` - renders per-level statistics table

**Dependencies:**
- `internal/stats` (for Statistics type)
- `github.com/olekukonko/tablewriter`

**Why here:** Separates presentation logic from calculation. Table formatting code extracted from main. Split into two functions for better reusability.

### cmd/cli-stat-creator/main.go
**Purpose:** Application entry point and coordination

**Contains:**
- `main()` function
- User input handling
- Error message display
- Coordination of reader → stats → display flow

**Dependencies:** All three internal packages

**Why here:** The `cmd/` directory is a Go convention for executable entry points. Main should be thin - just wiring components together.

## Code Organization Details

### Type Ownership
- `GameScore` and `Statistics` live in `internal/stats`
- They are domain types, not I/O types
- Other packages import and use these types

### Visibility Rules
- All package names are lowercase (Go convention)
- Exported symbols use UpperCamelCase
- Package-private symbols use lowerCamelCase
- Everything in these packages is exported since they're internal

### Import Paths
All imports use the module name from go.mod:
```go
import "cli-stat-creator/internal/stats"
import "cli-stat-creator/internal/reader"
import "cli-stat-creator/internal/display"
```

## Build Process Changes

### Before (single file)
```bash
go build -o cli-stat-creator
go run main.go
```

### After (cmd directory)
```bash
go build -o cli-stat-creator ./cmd/cli-stat-creator
go run ./cmd/cli-stat-creator
```

The `./cmd/cli-stat-creator` path tells Go where to find the main package.

## Testing Strategy

### Package-Level Testing
Each package can be tested independently:
```bash
go test ./internal/stats       # Test calculations in isolation
go test ./internal/reader      # Test file reading with fixtures
go test ./internal/display     # Test table rendering
go test ./...                  # Test everything
```

### Test Files Structure
```
internal/stats/stats_test.go
internal/reader/json_test.go
internal/display/table_test.go
```

### Benefits
- Pure functions in `stats` are easy to test with table-driven tests
- Reader tests can use fixture files in `testdata/`
- Display tests can capture stdout or use io.Writer interface
- No need to test main coordination logic separately

## Why `internal/` Directory?

The `internal/` directory is a special Go convention:
- Packages inside `internal/` cannot be imported by external projects
- Signals "these are implementation details, not public API"
- Useful even for CLI-only tools as it enforces architectural boundaries
- Standard practice in Go projects

## Migration Path

1. Create directory structure
2. Move stats logic to `internal/stats/stats.go`
3. Move reader logic to `internal/reader/json.go`
4. Extract display logic to `internal/display/table.go`
5. Create new `cmd/cli-stat-creator/main.go`
6. Fix all import statements
7. Test build and execution
8. Update documentation

## Future Growth Considerations

This structure supports future additions:

**New readers:**
- `internal/reader/csv.go`
- `internal/reader/api.go`

**New statistics:**
- Add functions to `internal/stats/stats.go`
- Or create `internal/stats/advanced.go`

**New displays:**
- `internal/display/json.go`
- `internal/display/chart.go`

**Multiple commands:**
- `cmd/cli-stat-creator/` (existing)
- `cmd/stat-server/` (future web server)
- Both share `internal/` packages

## Learning Outcomes

By implementing this structure, you'll learn:
- How to organize multi-package Go projects
- Go's visibility rules (capitalization-based exports)
- The `internal/` package convention
- The `cmd/` directory pattern for executables
- How to structure imports and manage dependencies
- Package-level testing strategies
- How Go project layout scales from small to large
