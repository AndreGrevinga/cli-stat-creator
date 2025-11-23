# Structured Logging Design

**Issue:** [#12 Add structured logging support](https://github.com/AndreGrevinga/cli-stat-creator/issues/12)
**Date:** 2025-11-23

## Overview and Architecture

**Goal:** Add structured logging to the CLI application using Go's standard `log/slog` library, making it easy to troubleshoot issues while keeping the default output clean.

**Key Decisions:**
- **Library:** `log/slog` from the standard library (Go 1.21+)
- **Output:** Logs to stderr, statistics to stdout (preserves Unix conventions)
- **Default Level:** WARN (silent during normal operation)
- **Configuration:** CLI flag `--log-level` / `-l` with `LOG_LEVEL` environment variable fallback
- **Logger Distribution:** Pass logger via `context.Context` through the pipeline

**Architecture Changes:**
1. **Main setup:** Create and configure logger in `main()`, add to context
2. **Pipeline integration:** Extract logger from context in each pipeline stage
3. **Error handling:** Add structured logging to existing error paths
4. **Configuration:** Add log level flag parsing and logger initialization

**Log Levels:**
- `DEBUG`: Per-score processing, detailed pipeline operations
- `INFO`: File operations, calculation summaries, pipeline stage transitions
- `WARN`: Invalid scores, empty results, unusual conditions
- `ERROR`: File failures, parsing errors, calculation errors

## Implementation Details

### Logger Initialization (in `main()`)

```go
// Parse log level from flag or environment variable
logLevel := parseLogLevel(*logLevelFlag, os.Getenv("LOG_LEVEL"))

// Create handler (Text for readability)
handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
    Level: logLevel,
    AddSource: logLevel == slog.LevelDebug, // Add file:line for DEBUG
})

// Create logger and add to context
logger := slog.New(handler)
ctx = logging.WithLogger(ctx, logger)
```

### Context Key (new internal/logging package)

```go
package logging

import (
    "context"
    "io"
    "log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

// FromContext extracts logger from context, returns no-op logger if not found
func FromContext(ctx context.Context) *slog.Logger {
    if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
        return logger
    }
    return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
    return context.WithValue(ctx, loggerKey, logger)
}
```

### Pipeline Logging Examples

```go
// In Source()
logger := logging.FromContext(ctx)
logger.Info("reading file", "filename", filename)
scores, err := reader.ReadScoresFromFile(filename)
if err != nil {
    logger.Error("file read failed", "filename", filename, "error", err)
    return nil, err
}
logger.Info("file read completed", "filename", filename, "score_count", len(scores))
```

## Logging Points by Component

### cmd/cli-stat-creator/main.go
- INFO: Application started with configuration details
- DEBUG: Flag values parsed
- ERROR: Invalid flag values (already printing to stdout, add structured log)
- INFO: Pipeline execution started
- ERROR: Pipeline execution failed
- INFO: Pipeline completed successfully with result counts

### internal/pipeline/stages.go

**Source function:**
- INFO: Reading file (filename)
- INFO: File read completed (filename, score_count)
- ERROR: File read failed (filename, error)
- DEBUG: Streaming scores to channel

**Filter stage:**
- INFO: Filter stage started (config summary)
- DEBUG: Score filtered out (player, level, score, reason)
- DEBUG: Score passed filter (player, level, score)
- INFO: Filter stage completed (input_count, output_count, filtered_count)

**Aggregate function:**
- INFO: Aggregation started
- WARN: No scores to process (after filtering or from empty file)
- INFO: Calculating overall statistics (score_count)
- INFO: Calculating level statistics (enabled/disabled)
- INFO: Calculating player statistics (enabled/disabled)
- ERROR: Statistics calculation failed (error)
- INFO: Aggregation completed (overall, by_level_count, by_player_count)

### internal/reader/json.go
- DEBUG: File extension validated
- WARN: Invalid score skipped during validation (player, score, level, validation_error)
- ERROR: JSON parsing failed (filename, error)

## Testing and Edge Cases

### Testing Approach

1. **Logger behavior tests:**
   - Test `FromContext()` returns no-op logger when context has no logger
   - Test `FromContext()` returns correct logger when present
   - Test log level parsing with valid/invalid inputs

2. **Integration tests:**
   - Capture stderr output during pipeline runs
   - Verify expected log messages appear at correct levels
   - Test different log levels produce appropriate output

3. **Unit tests update:**
   - Existing tests should still pass (no-op logger when no logger in context)
   - Add context with test logger where needed for assertions

### Edge Cases

1. **No logger in context:** Return no-op logger (logs to io.Discard) - graceful degradation
2. **Invalid log level flag:** Default to WARN, log warning about invalid value
3. **Empty LOG_LEVEL env var:** Ignore, use flag or default
4. **Both flag and env var set:** Flag takes precedence
5. **Concurrent logging:** `slog` is safe for concurrent use (pipeline stages run concurrently)

### Handler Format Choice
- **TextHandler**: Human-readable for terminal use (recommended for CLI)
- **JSONHandler**: Machine-parseable but harder to read
- Start with TextHandler, can make configurable later if needed

### Performance Considerations
- Logging at DEBUG level could impact performance with large files
- No-op logger has zero overhead when log level filters out messages
- `slog` is designed for high performance with zero allocations for disabled levels when using key-value pairs

**Important:** Use slog's key-value pairs for zero-allocation disabled logs:
```go
// GOOD - zero allocation if DEBUG is disabled
logger.Debug("processing score", "score", score)

// BAD - allocates even if DEBUG is disabled
logger.Debug(fmt.Sprintf("processing score %d", score))
```

## Files and Package Structure

### New Files

**internal/logging/logging.go**
- `FromContext(ctx) *slog.Logger` - Extract logger from context
- `WithLogger(ctx, logger) context.Context` - Add logger to context
- Context key definition
- Helper for creating no-op logger

### Modified Files

**cmd/cli-stat-creator/main.go**
- Add `logLevel` flag (string, default "warn")
- Add `parseLogLevel()` function
- Initialize logger with TextHandler to stderr
- Add logger to context before pipeline.Run()
- Add INFO log for app start and completion
- Add ERROR logs for failures

**internal/pipeline/stages.go**
- Import `internal/logging`
- Add logging to `Source()` - file read operations
- Add logging to `Filter()` - filter operations and results
- Add logging to `Aggregate()` - calculation operations

**internal/pipeline/pipeline.go**
- Add INFO log at pipeline start in `Run()`
- Add INFO/ERROR logs for results

**internal/reader/json.go**
- Add WARN logs for invalid scores during validation
- Add DEBUG log for file extension check

### CLI Flag

- Long form: `--log-level <level>` or `-l <level>`
- Values: `debug`, `info`, `warn`, `error`
- Default: `warn`
- Environment variable: `LOG_LEVEL`

### Example Usage

```bash
# Default (WARN level)
./cli-stat-creator -i

# With DEBUG logging
./cli-stat-creator -i --log-level debug

# Using environment variable
LOG_LEVEL=info ./cli-stat-creator -i

# Redirect logs to file, keep statistics on screen
./cli-stat-creator -i 2>app.log
```

## Implementation Notes

### Log Message Format

Use structured key-value pairs consistently:
- `filename` - for file paths
- `score_count` - for number of scores
- `filter_type` - for filter operations
- `level` - for game level (context, not log level)
- `player` - for player name/info
- `error` - for error details
- `input_count`, `output_count`, `filtered_count` - for filter results

### Error Handling

Maintain existing error returns, but add structured logging:
```go
if err != nil {
    logger.Error("operation failed", "context", value, "error", err)
    return err // Still return the error
}
```

Don't replace error handling with logging - they serve different purposes.
