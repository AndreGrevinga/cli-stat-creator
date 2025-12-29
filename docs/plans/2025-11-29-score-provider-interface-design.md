# ScoreProvider Interface Design Document

**Date:** 2025-11-29
**Status:** Ready for Implementation
**Purpose:** Enable `Pipeline.Run()` to accept data from multiple sources (files, io.Reader) via a unified interface

## Overview

Refactor the pipeline to use a `ScoreProvider` interface instead of accepting a filename directly. This allows both the CLI and HTTP server to use the same `Pipeline.Run()` method with different data sources.

## Current Architecture

```
Pipeline.Run(ctx, filename string)
    → Source(ctx, filename)
        → reader.ReadScoresFromFile(ctx, filename)
            → os.ReadFile(filename)
```

The pipeline is tightly coupled to filesystem-based input.

## Target Architecture

```
Pipeline.Run(ctx, provider ScoreProvider)
    → Source(ctx, provider)
        → provider.Scores(ctx)
            → FileProvider: reader.ReadScoresFromFile()
            → ReaderProvider: reader.ReadScoresFromReader()
```

## Implementation Steps

### Step 1: Extend the reader package

**File:** `internal/reader/json.go`

Add a new function that reads from `io.Reader`:

```go
import "io"

// ReadScoresFromReader reads and parses game scores from an io.Reader.
func ReadScoresFromReader(ctx context.Context, r io.Reader) ([]stats.GameScore, error) {
    logger := logging.FromContext(ctx)

    data, err := io.ReadAll(r)
    if err != nil {
        return nil, fmt.Errorf("failed to read data: %w", err)
    }

    var gameScores []stats.GameScore
    if err := json.Unmarshal(data, &gameScores); err != nil {
        return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
    }

    for i, score := range gameScores {
        if err := score.Validate(); err != nil {
            logger.Warn("Invalid score during validation",
                "index", i,
                "player", score.Player.Name,
                "error", err,
            )
            return nil, fmt.Errorf("invalid score at index %d: %w", i, err)
        }
    }

    return gameScores, nil
}
```

**Optional refactor:** Update `ReadScoresFromFile` to use the new function internally:

```go
func ReadScoresFromFile(ctx context.Context, filename string) ([]stats.GameScore, error) {
    logger := logging.FromContext(ctx)

    if !strings.HasSuffix(strings.ToLower(filename), ".json") {
        return nil, fmt.Errorf("file '%s' must have .json extension", filename)
    }
    logger.Debug("File extension validated", "filename", filename)

    file, err := os.Open(filename)
    if err != nil {
        return nil, fmt.Errorf("failed to open file %s: %w", filename, err)
    }
    defer file.Close()

    return ReadScoresFromReader(ctx, file)
}
```

---

### Step 2: Create ScoreProvider interface

**File:** `internal/pipeline/source.go` (new file)

```go
package pipeline

import (
    "context"
    "io"

    "cli-stat-creator/internal/reader"
    "cli-stat-creator/internal/stats"
)

// ScoreProvider abstracts the source of game scores for the pipeline.
type ScoreProvider interface {
    Scores(ctx context.Context) ([]stats.GameScore, error)
}

// FileProvider reads scores from a file path.
type FileProvider struct {
    Filename string
}

// Scores reads game scores from the configured file.
func (f FileProvider) Scores(ctx context.Context) ([]stats.GameScore, error) {
    return reader.ReadScoresFromFile(ctx, f.Filename)
}

// ReaderProvider reads scores from an io.Reader.
type ReaderProvider struct {
    Reader io.Reader
}

// Scores reads game scores from the configured reader.
func (r ReaderProvider) Scores(ctx context.Context) ([]stats.GameScore, error) {
    return reader.ReadScoresFromReader(ctx, r.Reader)
}
```

---

### Step 3: Update Source function

**File:** `internal/pipeline/stages.go`

Change the `Source` function signature to accept `ScoreProvider`:

```go
// Source reads game scores from a provider and streams them through a channel.
func Source(ctx context.Context, provider ScoreProvider) (<-chan stats.GameScore, error) {
    logger := logging.FromContext(ctx)

    logger.Info("Reading scores from provider")
    scores, err := provider.Scores(ctx)
    if err != nil {
        return nil, err
    }
    logger.Info("Scores loaded", "score_count", len(scores))

    out := make(chan stats.GameScore)
    go func() {
        defer close(out)
        logger.Debug("Streaming scores to channel", "count", len(scores))
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

---

### Step 4: Update Pipeline.Run

**File:** `internal/pipeline/pipeline.go`

Change the method signature:

```go
// Run executes the pipeline using the provided score source.
func (p *Pipeline) Run(ctx context.Context, provider ScoreProvider) (Results, error) {
    logger := logging.FromContext(ctx)
    logger.Info("Pipeline started", "stage_count", len(p.stages))

    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    in, err := Source(ctx, provider)  // Changed from Source(ctx, filename)
    if err != nil {
        return Results{}, err
    }

    for _, stage := range p.stages {
        in = stage(ctx, in)
    }

    results, err := Aggregate(ctx, in, p.config)
    if err != nil {
        return Results{}, err
    }

    logger.Info("Pipeline completed successfully",
        "overall_count", results.Overall.TotalGamesPlayed,
        "by_level_count", len(results.ByLevel),
        "by_player_count", len(results.ByPlayer),
    )

    return results, nil
}
```

---

### Step 5: Update CLI caller

**File:** `cmd/cli-stat-creator/main.go`

Find where `p.Run()` is called and update:

```go
// Before
results, err := p.Run(ctx, filename)

// After
results, err := p.Run(ctx, pipeline.FileProvider{Filename: filename})
```

---

### Step 6: Use in HTTP handler

**File:** `internal/handlers/stats.go` (when you create it)

```go
file, header, err := r.FormFile("file")
if err != nil {
    // handle error
}
defer file.Close()

// Validate .json extension from header.Filename

p := pipeline.New(config, pipeline.Filter(config))
results, err := p.Run(ctx, pipeline.ReaderProvider{Reader: file})
```

---

## Files Summary

| File | Change |
|------|--------|
| `internal/reader/json.go` | Add `ReadScoresFromReader()` |
| `internal/pipeline/source.go` | New file: `ScoreProvider`, `FileProvider`, `ReaderProvider` |
| `internal/pipeline/stages.go` | Update `Source()` signature |
| `internal/pipeline/pipeline.go` | Update `Run()` signature |
| `cmd/cli-stat-creator/main.go` | Use `FileProvider{Filename: ...}` |

## Usage Examples

**CLI:**
```go
p := pipeline.New(config, pipeline.Filter(config))
results, err := p.Run(ctx, pipeline.FileProvider{Filename: "data/input.json"})
```

**HTTP:**
```go
p := pipeline.New(config, pipeline.Filter(config))
results, err := p.Run(ctx, pipeline.ReaderProvider{Reader: uploadedFile})
```

**Testing (with mock data):**
```go
type MockProvider struct {
    Scores []stats.GameScore
}

func (m MockProvider) Scores(ctx context.Context) ([]stats.GameScore, error) {
    return m.Scores, nil
}

// In test
p := pipeline.New(config)
results, err := p.Run(ctx, MockProvider{Scores: testData})
```
