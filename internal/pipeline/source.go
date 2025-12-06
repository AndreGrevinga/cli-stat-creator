// Package pipeline extension for score provider implementations.
// This file contains interfaces and implementations for providing game scores
// from various sources (files, readers, etc.) to the pipeline.
package pipeline

import (
	"cli-stat-creator/internal/reader"
	"cli-stat-creator/internal/stats"
	"context"
	"io"
)

// ScoreProvider is an interface for sources that can provide game scores.
// Implementations include file-based and reader-based providers.
type ScoreProvider interface {
	Scores(ctx context.Context) ([]stats.GameScore, error)
}

// FileProvider provides game scores by reading from a file path.
type FileProvider struct {
	Filename string
}

// Scores reads and returns game scores from the file specified in the provider.
func (f FileProvider) Scores(ctx context.Context) ([]stats.GameScore, error) {
	return reader.ReadScoresFromFile(ctx, f.Filename)
}

// ReaderProvider provides game scores by reading from an io.Reader.
// This is useful for processing uploaded data or in-memory data.
type ReaderProvider struct {
	Reader io.Reader
}

// Scores reads and returns game scores from the reader in the provider.
func (r ReaderProvider) Scores(ctx context.Context) ([]stats.GameScore, error) {
	return reader.ReadScoresFromReader(ctx, r.Reader)
}
