# Pipeline Architecture Design

**Date:** 2025-11-18
**Goal:** Implement a streaming pipeline architecture for learning Go concurrency patterns

## Overview

A streaming pipeline where each stage runs in its own goroutine, connected by channels. Each `GameScore` flows through independently:

```
┌────────┐   ┌────────┐   ┌───────────┐   ┌───────────┐   ┌────────┐
│ Source │ → │ Filter │ → │ Transform │ → │ Aggregate │ → │ Output │
└────────┘   └────────┘   └───────────┘   └───────────┘   └────────┘
     ↓            ↓             ↓              ↓              ↓
   chan        chan          chan           chan          Statistics
 GameScore   GameScore     GameScore      GameScore
```

## Design Decisions

- **Data flow:** Item-by-item streaming (not batch) - teaches channel patterns
- **Error handling:** Fail fast with context cancellation - clean and simple
- **Stages:** Source → Filter → Transform → Aggregate → Output

## Package Structure

New package `internal/pipeline` containing:
- `pipeline.go` - Core types and orchestration
- `stages.go` - Individual stage implementations
- `pipeline_test.go` - Tests

## Key Types

```go
// Stage function signature - each stage follows this pattern
type Stage func(ctx context.Context, in <-chan GameScore) <-chan GameScore

// Pipeline orchestrates the full flow
type Pipeline struct {
    stages []Stage
}

// Config holds filter/transform settings
type Config struct {
    FilterByPlayer    []string
    FilterByLevel     []int
    MinScore, MaxScore int
}
```

## Stage Implementations

### Source Stage

Reads from file and emits scores one at a time:

```go
func Source(ctx context.Context, filename string) (<-chan GameScore, error) {
    // Read file once, then stream items
    scores, err := reader.ReadScoresFromFile(filename)
    if err != nil {
        return nil, err
    }

    out := make(chan GameScore)
    go func() {
        defer close(out)
        for _, score := range scores {
            select {
            case out <- score:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out, nil
}
```

### Filter Stage

Passes through only matching scores:

```go
func Filter(cfg Config) Stage {
    return func(ctx context.Context, in <-chan GameScore) <-chan GameScore {
        out := make(chan GameScore)
        go func() {
            defer close(out)
            for score := range in {
                if matchesFilter(score, cfg) {
                    select {
                    case out <- score:
                    case <-ctx.Done():
                        return
                    }
                }
            }
        }()
        return out
    }
}
```

### Transform Stage

Applies transformations (placeholder for future enhancements):

```go
func Transform(cfg Config) Stage {
    return func(ctx context.Context, in <-chan GameScore) <-chan GameScore {
        // Same pattern, apply transformation logic
    }
}
```

### Aggregate Stage

Accumulates all scores to calculate statistics:

```go
func Aggregate(ctx context.Context, in <-chan GameScore) (Statistics, error) {
    var scores []GameScore

    for score := range in {
        select {
        case <-ctx.Done():
            return Statistics{}, ctx.Err()
        default:
            scores = append(scores, score)
        }
    }

    if len(scores) == 0 {
        return Statistics{}, errors.New("no scores after filtering")
    }

    return stats.CalculateStatistics(scores)
}
```

## Pipeline Orchestration

```go
func (p *Pipeline) Run(ctx context.Context, filename string) (Statistics, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    // Start source
    in, err := Source(ctx, filename)
    if err != nil {
        return Statistics{}, err
    }

    // Chain stages
    for _, stage := range p.stages {
        in = stage(ctx, in)
    }

    // Aggregate and return
    return Aggregate(ctx, in)
}
```

## Integration with Existing Code

The pipeline package will use existing packages:
- `stats.GameScore` and `stats.Statistics` types
- `stats.CalculateStatistics()` for the actual math
- `reader.ReadScoresFromFile()` in the Source stage
- `display.RenderStatistics()` for output

Main.go changes from:
```go
scores := reader.ReadScoresFromFile(filepath)
statistics := stats.CalculateStatistics(scores)
display.RenderStatistics(statistics)
```

To:
```go
p := pipeline.New(
    pipeline.Filter(cfg),
    pipeline.Transform(cfg),
)
statistics, err := p.Run(ctx, filepath)
display.RenderStatistics(statistics)
```

## Testing Strategy

Table-driven tests for each stage:

1. **Source tests** - Valid file, missing file, empty file
2. **Filter tests** - By player, by level, by score range, combinations
3. **Transform tests** - Placeholder for now
4. **Aggregate tests** - Normal flow, empty input, context cancellation
5. **Integration test** - Full pipeline end-to-end

Key test pattern for channels:
```go
func TestFilter(t *testing.T) {
    ctx := context.Background()
    in := make(chan GameScore, 3)
    // send test data, close channel
    out := Filter(cfg)(ctx, in)
    // collect results, assert
}
```

## Future Extensions

### File Watching (Phase 2)

Once the pipeline is working, adding streaming from a file watcher:

```go
func WatchSource(ctx context.Context, dir string) <-chan GameScore {
    out := make(chan GameScore)
    go func() {
        defer close(out)
        watcher, _ := fsnotify.NewWatcher()
        watcher.Add(dir)

        for {
            select {
            case event := <-watcher.Events:
                if event.Op&fsnotify.Write != 0 {
                    // Read new/changed file, emit scores
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}
```

### Other Possibilities

- **Fan-out/fan-in** - Multiple filter stages in parallel
- **Metrics** - Track items processed per stage
- **Backpressure** - Buffered channels for flow control

## Learning Goals

This feature teaches:
- Goroutines and channel communication
- Context cancellation patterns
- Pipeline composition
- Functional stage design
- Concurrent testing patterns
