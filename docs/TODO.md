# CLI Stat Creator - Improvement Plan

## HTTP Server Implementation (Priority: MEDIUM)

These issues apply to the `feature/http-server` branch which is not yet merged to main.

### Before Merge

- [ ] Replace `strings.SplitSeq` with `strings.Split` in `parseIncludeParam`
  - Location: `internal/handlers/stats.go:144`
  - Issue: SplitSeq is Go 1.23+ experimental feature
  - Solution: Use standard Split for broader compatibility

- [ ] Move Content-Type header setting before WriteHeader in `internal/render/json.go`
  - Location: Lines 17, 22 and 38, 44
  - Issue: Headers should be set before WriteHeader() for consistency
  - Solution: Reorder lines to set headers first

- [ ] Add tests for handlers package (`internal/handlers/stats_test.go`)
  - Package currently has no test coverage
  - Test HandleStats with various inputs and query parameters

- [ ] Add tests for render package (`internal/render/json_test.go`)
  - Package currently has no test coverage
  - Test JSON, Error, and FilterResults functions

## High Priority - Bugs & Code Quality

### Main Branch Issues

- [ ] **#45** - Fix duplicate error handling in main.go
  - Location: `cmd/cli-stat-creator/main.go:94-111, 116-121`
  - Issue: Errors are logged both to structured logger AND printed to stdout
  - Solution: Choose one approach (preferably stderr for CLI errors)
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/45

- [ ] **#47** - Fix MaxScore=0 sentinel value (0 is valid score)
  - Location: `internal/pipeline/stages.go:118`
  - Issue: `MaxScore=0` means "no limit" but 0 is a valid game score
  - Solution: Use pointer type `*int` (nil = no limit) or separate enabled flag
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/47

- [ ] **#41** - Add validation for negative levels in ParseLevelString
  - Location: `internal/pipeline/params.go:26-40`
  - Issue: Function accepts negative levels without validation
  - Solution: Add explicit check `if minLevel < 0 || maxLevel < 0 { return error }`
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/41

- [ ] **#53** - Remove dead TODO comment code in main.go
  - Location: `cmd/cli-stat-creator/main.go:25-26`
  - Issue: Commented-out player filter flag code
  - Action: Remove or implement the feature
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/53

- [ ] **#11** - Improve error messages with context and suggestions
  - Current: Generic error messages
  - Goal: Clear context, suggest solutions, include file paths
  - Example: "failed to read file at path/to/file: file not found. Please check..."
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/11

## High Priority - Testing

- [ ] **#34** - Add test cases for edge cases in reader package
  - Missing tests: empty JSON array, malformed JSON, extra fields, large numbers, unicode
  - Location: `internal/reader/json_test.go`
  - Effort: ~30 minutes
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/34

- [ ] **#9 / #54** - Add benchmark tests (duplicate issues)
  - Packages to benchmark: stats, pipeline, reader
  - Key functions: CalculateStatistics, GroupByLevel, Filter stages
  - Track performance regressions in CI
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/9
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/54

## Medium Priority - Performance & Tech Debt

- [ ] **#46** - Use map for O(1) filter lookups instead of O(n) slice iteration
  - Location: `internal/pipeline/stages.go:65-104`
  - Current: Linear search through filter slices
  - Solution: Convert filter slices to maps at start of Filter function
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/46

- [ ] **#52** - Use int64 for totalScore to prevent overflow
  - Location: `internal/stats/stats.go:61`
  - Risk: Low probability but high impact (silent incorrect statistics)
  - Consider: Document limitation or use int64 for large datasets
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/52

- [ ] **#50** - Fix non-hermetic pipeline tests (use fixtures instead of real files)
  - Location: `internal/pipeline/pipeline_test.go:281, 310, 384, 413, 474`
  - Issue: Tests depend on `../../data/input.json`
  - Solution: Use `t.TempDir()` or Go embed for test fixtures
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/50

- [ ] **#49** - Create test fixtures for edge cases
  - Needed fixtures:
    - `data/test-empty.json` - Empty array
    - `data/test-single-entry.json` - Minimal valid data
    - `data/test-edge-cases.json` - Extreme values (zero scores, high scores, level 1, unicode)
    - `data/test-large.json` - 10,000+ records for performance
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/49

- [ ] **#48** - Add end-to-end CLI integration test
  - Missing: Full CLI execution tests, flag combinations, error output, exit codes
  - Location: `cmd/cli-stat-creator/main_test.go`
  - Current: Only tests helper functions (ParseLevelString, parseLogLevel)
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/48

## Medium Priority - Features (Statistics)

- [ ] **#21** - Implement export functionality for multiple formats
  - Formats: JSON, CSV, plain text, markdown
  - Implementation: New package `internal/export` with format-specific exporters
  - CLI flags: `-o, --output` (format), `--output-file` (destination)
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/21

- [ ] **#19** - Add player count and game count per level
  - Current: Only shows average scores by level
  - Add to Statistics struct: `PlayerCountByLevel`, `GameCountByLevel`
  - Display in level breakdown table
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/19

- [ ] **#16** - Add advanced statistical calculations
  - New statistics: standard deviation, variance, percentiles (25th, 75th, 90th), mode, IQR
  - Benefits: Deeper insights, identify outliers and patterns
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/16

## Low Priority - CLI & UX

- [ ] **#32** - Add --version and --help CLI flags
  - Implementation: Use Go's `flag` package
  - Flags: `--version/-v` (print version), `--help/-h` (usage info)
  - Effort: ~30 minutes
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/32

- [ ] **#15** - Add configuration file support
  - Format: YAML or JSON
  - Locations: `~/.cli-stat-creator.yaml`, `.cli-stat-creator.yaml`
  - Options: Default input file, output format, logging level, display settings
  - Priority: CLI flags override config file
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/15

## Low Priority - Advanced Analysis Features

- [ ] **#23** - Add top N players/scores ranking feature
  - Features:
    - Top N by average score, total score, games played
    - Top N individual high scores (leaderboard)
    - Configurable N (default: 10)
  - CLI: `--top=N`, `--sort-by=average|total|games`
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/23

- [ ] **#18** - Add score distribution and histogram feature
  - Create score buckets (0-100, 101-200, etc.)
  - Display as ASCII histogram
  - Show distribution statistics (skewness, kurtosis)
  - GitHub: https://github.com/AndreGrevinga/cli-stat-creator/issues/18

## Completed Items

- [x] **#51** - Add doc comments to Validate() and RenderGroupedStatistics
  - Status: COMPLETED - Both functions now have proper Go doc comments
  - Validate: `internal/stats/stats.go:33-35`
  - RenderGroupedStatistics: `internal/display/table.go:39-47`

- [x] **#44** - Refactor main() into focused functions
  - Status: COMPLETED - See PR #55, commit 6e94ef7
  - Main now broken into: setupLogging, parseFlags, getInputFilepath, runPipeline, displayResults

- [x] Add comprehensive unit tests for core functions
  - ReadScoresFromFile: `internal/reader/json_test.go`
  - CalculateStatistics: `internal/stats/stats_test.go`
  - GroupByLevel: `internal/stats/stats_test.go`

- [x] Add input validation (file extension check, score/level validation)
  - File extension: `internal/reader/json.go`
  - GameScore validation: `internal/stats/stats.go:36-49`

- [x] Add logging support (structured logging with levels)
  - Package: `internal/logging/`
  - Used throughout the application

- [x] Implement CLI argument parsing using flag package
  - Flags for: detailed, default-input, log-level, no-players, no-levels, level, min-score, max-score

- [x] Extract display logic from main into separate presentation functions
  - Package: `internal/display/`

- [x] Refactor code into separate packages
  - Packages: stats, reader, display, logging, pipeline

- [x] Add filtering capabilities (by level, score range)
  - Level filtering: Complete
  - Score range filtering: Complete
  - Player filtering: Backend ready, CLI flag pending (see #53)

- [x] Add per-player statistics
  - Implemented in `internal/stats/stats.go:128-138`

## Notes

- Issues #9 and #54 are duplicates (both about benchmark tests)
- The HTTP server work is on branch `feature/http-server` (not yet merged to main)
- All HTTP server critical issues have been resolved
- Player filter functionality exists in pipeline but CLI flag is not implemented (see #53)
