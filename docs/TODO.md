# CLI Stat Creator - Improvement Plan

## HTTP Server Implementation - Code Review Issues (Priority: CRITICAL)

### Critical Issues (Must Fix Before Commit)

#### Documentation
- [ ] Add Go doc comments to ScoreProvider interface and implementations in `internal/pipeline/source.go`
  - ScoreProvider interface (line 10)
  - FileProvider struct (line 14)
  - FileProvider.Scores method (line 18)
  - ReaderProvider struct (line 22)
  - ReaderProvider.Scores method (line 26)
- [ ] Add Go doc comments to all exported functions in `internal/render/json.go`
  - JSON function (line 9)
  - Error function (line 20)
  - FilterResults function (line 39)
- [ ] Add Go doc comments to HandleStats function in `internal/handlers/stats.go` (line 12)

#### Functionality Bugs
- [ ] Fix `parseQueryParams` in `internal/handlers/stats.go` to handle missing/empty query parameters with defaults
  - Lines 66-83: `strconv.Atoi("")` fails when min-score/max-score query params are absent
  - Should default to 0 for missing values
  - Should only error on genuinely malformed input
- [ ] Set `CalculateByLevel` and `CalculateByPlayer` fields in `parseQueryParams` (lines 79-83)
  - Currently missing from pipeline.Config
  - Should default to true or be configurable via query params
- [ ] Fix `HandleStats` error handling to send HTTP response on pipeline failure
  - Line 39: Currently returns without sending response (client hangs)
  - Should call `render.Error(w, 500, "Failed to process file", "PROCESSING_ERROR")`
- [ ] Register `POST /api/stats` route in `setupRouter` in `cmd/http-server/main.go`
  - Route handler exists but is never registered with chi router
- [ ] Start HTTP server with `http.ListenAndServe()` in `main()`
  - Currently setupRouter is called but server never starts listening

### Important Issues (Should Fix Before Commit)

- [ ] Add logging middleware to attach slog logger to request context
  - `handlers/stats.go` line 31 calls `logging.FromContext(ctx)` but logger is never attached
  - Chi middleware.Logger is used but doesn't set up slog in context
- [ ] Fix `parseConfig` in `cmd/http-server/main.go` to populate all fields and handle errors properly
  - Lines 34-42: Error silently returns empty config (port 0)
  - Missing: Host, LogLevel, StaticDir population
  - Flags are declared but never parsed in main()
- [ ] Add file size limit check in `HandleStats` with 413 response
  - Line 13: `r.ParseMultipartForm(10 << 20)` doesn't check error
  - Should return 413 Payload Too Large for files exceeding limit
- [ ] Connect `parseIncludeParam` and `FilterResults` to actual response filtering
  - Both functions exist but are never called
  - Should filter results based on include query parameter

### Minor Issues (Nice to Have)

- [ ] Fix log message in `internal/reader/json.go` line 49 to match design doc
  - Currently: "Invalid score skipped during validation"
  - Should be: "Invalid score during validation"
- [ ] Replace `strings.SplitSeq` with `strings.Split` in `parseIncludeParam`
  - Line 50: SplitSeq is Go 1.23+ experimental feature
  - Use standard Split for broader compatibility
- [ ] Move Content-Type header setting before WriteHeader in `internal/render/json.go`
  - Lines 10-17: Headers should be set before WriteHeader() for consistency
- [ ] Add tests for handlers package (`internal/handlers/stats_test.go`)
  - Package currently has no test coverage
- [ ] Add tests for render package (`internal/render/json_test.go`)
  - Package currently has no test coverage

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
