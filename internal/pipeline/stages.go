package pipeline

import (
	"cli-stat-creator/internal/reader"
	"cli-stat-creator/internal/stats"
	"context"
	"errors"
)

// Stage represents a processing step in the pipeline that transforms a stream of game scores.
// Each stage receives an input channel and returns an output channel.
type Stage func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore

// Source reads game scores from a JSON file and streams them through a channel.
// It returns immediately with the channel, while scores are sent concurrently in a goroutine.
// The operation can be cancelled via the context.
func Source(ctx context.Context, filename string) (<-chan stats.GameScore, error) {
	// Read file once, then stream items
	scores, err := reader.ReadScoresFromFile(filename)
	if err != nil {
		return nil, err
	}

	out := make(chan stats.GameScore)
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

func matchesFilter(score stats.GameScore, config Config) bool {
	isLevelValid := len(config.FilterByLevel) == 0
	isPlayerValid := len(config.FilterByPlayer) == 0
	isScoreValid := (config.MaxScore == 0 || score.Score <= config.MaxScore) && score.Score >= config.MinScore
	for _, level := range config.FilterByLevel {
		if score.Level == level {
			isLevelValid = true
		}
	}
	for _, player := range config.FilterByPlayer {
		if score.Player.Name == player {
			isPlayerValid = true
		}
	}
	return isLevelValid && isPlayerValid && isScoreValid
}

// Filter returns a Stage that filters game scores based on the provided configuration.
// Scores can be filtered by player names, levels, and score ranges (min/max).
// When a filter list is empty, all values for that field pass through.
func Filter(cfg Config) Stage {
	return func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore {
		out := make(chan stats.GameScore)
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

//Placeholder for the future
/*func Transform(cfg Config) Stage {
return func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore {
}
}*/

// Aggregate collects all game scores from the input channel and calculates comprehensive statistics.
// It returns overall statistics, as well as statistics grouped by level and player.
// Returns an error if no scores are available to process.
func Aggregate(ctx context.Context, in <-chan stats.GameScore) (Results, error) {
	var scores []stats.GameScore

	for score := range in {
		scores = append(scores, score)
	}

	if len(scores) == 0 {
		return Results{}, errors.New("no scores to process")
	}
	overall, err := stats.CalculateStatistics(scores)
	if err != nil {
		return Results{}, err
	}
	byLevel, err := stats.CalculateStatisticsByLevel(scores)
	if err != nil {
		return Results{}, err
	}
	byPlayer, err := stats.CalculateStatisticsByPlayer(scores)
	if err != nil {
		return Results{}, err
	}

	return Results{
		Overall:  overall,
		ByLevel:  byLevel,
		ByPlayer: byPlayer,
	}, nil
}
