// Package reader provides functionality for reading and parsing game score data
// from various input sources. Currently supports JSON file format.
package reader

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"cli-stat-creator/internal/stats"
)

// ReadScoresFromFile reads game scores from a JSON file and returns them as a slice of GameScore.
// The function expects the JSON file to contain an array of objects with Player, Score, and Level fields.
// Returns an error if the file cannot be read or if the JSON format is invalid.
func ReadScoresFromFile(filename string) ([]stats.GameScore, error) {
	var gameScores []stats.GameScore
	if !strings.HasSuffix(strings.ToLower(filename), ".json") {
		return nil, fmt.Errorf("file '%s' must have .json extension", filename)
	}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, fmt.Errorf("failed to read file %s: %w", filename, err)
	}
	if err := json.Unmarshal(data, &gameScores); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	for i, score := range gameScores {
		if err := score.Validate(); err != nil {
			return nil, fmt.Errorf("invalid score at index %d: %w", i, err)
		}
	}
	return gameScores, nil
}
