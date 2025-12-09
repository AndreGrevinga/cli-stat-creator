// Package reader provides functionality for reading and parsing game score data
// from various input sources. Currently supports JSON file format.
package reader

import (
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/stats"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
)

// ReadScoresFromFile reads game scores from a JSON file and returns them as a slice of GameScore.
// The function expects the JSON file to contain an array of objects with Player, Score, and Level fields.
// Returns an error if the file cannot be read or if the JSON format is invalid.
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

// ReadScoresFromReader reads game scores from a reader and returns them as a slice of GameScore.
// The function expects the reader to contain a JSON array of objects with Player, Score, and Level fields.
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
				"score", score.Score,
				"level", score.Level,
				"validation_error", err,
			)
			return nil, fmt.Errorf("invalid score at index %d: %w", i, err)
		}
	}
	return gameScores, nil
}
