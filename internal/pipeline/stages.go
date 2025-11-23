package pipeline

import (
	"cli-stat-creator/internal/logging"
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
	logger := logging.FromContext(ctx)

	logger.Info("Reading file", "filename", filename)
	scores, err := reader.ReadScoresFromFile(ctx, filename)
	if err != nil {
		logger.Error("File read failed",
			"filename", filename,
			"error", err,
		)
		return nil, err
	}
	logger.Info("File read completed",
		"filename", filename,
		"score_count", len(scores),
	)

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

// Filter returns a Stage that filters game scores based on the provided configuration.
// Scores can be filtered by player names, levels, and score ranges (min/max).
// When a filter list is empty, all values for that field pass through.
func Filter(cfg Config) Stage {
	return func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore {
		logger := logging.FromContext(ctx)

		logger.Info("Filter stage started",
			"min_score", cfg.MinScore,
			"max_score", cfg.MaxScore,
			"level_filter_count", len(cfg.FilterByLevel),
			"player_filter_count", len(cfg.FilterByPlayer),
		)

		out := make(chan stats.GameScore)
		go func() {
			defer close(out)
			var inputCount, outputCount, filteredCount int

			for score := range in {
				inputCount++

				// Check level filter
				if len(cfg.FilterByLevel) > 0 {
					levelValid := false
					for _, level := range cfg.FilterByLevel {
						if score.Level == level {
							levelValid = true
							break
						}
					}
					if !levelValid {
						filteredCount++
						logger.Debug("Score filtered out",
							"player", score.Player.Name,
							"level", score.Level,
							"score", score.Score,
							"reason", "level_not_in_filter",
						)
						continue
					}
				}

				// Check player filter
				if len(cfg.FilterByPlayer) > 0 {
					playerValid := false
					for _, player := range cfg.FilterByPlayer {
						if score.Player.Name == player {
							playerValid = true
							break
						}
					}
					if !playerValid {
						filteredCount++
						logger.Debug("Score filtered out",
							"player", score.Player.Name,
							"level", score.Level,
							"score", score.Score,
							"reason", "player_not_in_filter",
						)
						continue
					}
				}

				// Check score range filters
				if score.Score < cfg.MinScore {
					filteredCount++
					logger.Debug("Score filtered out",
						"player", score.Player.Name,
						"level", score.Level,
						"score", score.Score,
						"reason", "score_below_min",
					)
					continue
				}

				if cfg.MaxScore > 0 && score.Score > cfg.MaxScore {
					filteredCount++
					logger.Debug("Score filtered out",
						"player", score.Player.Name,
						"level", score.Level,
						"score", score.Score,
						"reason", "score_above_max",
					)
					continue
				}

				// Score passed all filters
				logger.Debug("Score passed filter",
					"player", score.Player.Name,
					"level", score.Level,
					"score", score.Score,
				)
				outputCount++
				select {
				case out <- score:
				case <-ctx.Done():
					return
				}
			}

			logger.Info("Filter stage completed",
				"input_count", inputCount,
				"output_count", outputCount,
				"filtered_count", filteredCount,
			)
		}()
		return out
	}
}

// Placeholder for the future
/*func Transform(cfg Config) Stage {
return func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore {
}
}*/

// Aggregate collects all game scores from the input channel and calculates comprehensive statistics.
// It returns overall statistics, as well as statistics grouped by level and player.
// The config parameter controls which statistics are calculated.
// Returns an error if no scores are available to process.
func Aggregate(ctx context.Context, in <-chan stats.GameScore, config Config) (Results, error) {
	logger := logging.FromContext(ctx)

	logger.Info("Aggregation started")

	var scores []stats.GameScore

	for score := range in {
		select {
		case <-ctx.Done():
			return Results{}, ctx.Err()
		default:
			scores = append(scores, score)
		}
	}

	if len(scores) == 0 {
		logger.Warn("No scores to process",
			"message", "All scores were filtered out or the input file was empty",
		)
		return Results{}, errors.New("no scores to process")
	}

	logger.Info("Calculating overall statistics", "score_count", len(scores))
	overall, err := stats.CalculateStatistics(scores)
	if err != nil {
		logger.Error("Statistics calculation failed",
			"stage", "overall",
			"error", err,
		)
		return Results{}, err
	}

	var byLevel map[int]stats.Statistics
	if config.CalculateByLevel {
		logger.Info("Calculating level statistics", "enabled", true)
		byLevel, err = stats.CalculateStatisticsByLevel(scores)
		if err != nil {
			logger.Error("Statistics calculation failed",
				"stage", "by_level",
				"error", err,
			)
			return Results{}, err
		}
	} else {
		logger.Info("Calculating level statistics", "enabled", false)
	}

	var byPlayer map[stats.Player]stats.Statistics
	if config.CalculateByPlayer {
		logger.Info("Calculating player statistics", "enabled", true)
		byPlayer, err = stats.CalculateStatisticsByPlayer(scores)
		if err != nil {
			logger.Error("Statistics calculation failed",
				"stage", "by_player",
				"error", err,
			)
			return Results{}, err
		}
	} else {
		logger.Info("Calculating player statistics", "enabled", false)
	}

	logger.Info("Aggregation completed",
		"overall_count", overall.TotalGamesPlayed,
		"by_level_count", len(byLevel),
		"by_player_count", len(byPlayer),
	)

	return Results{
		Overall:  overall,
		ByLevel:  byLevel,
		ByPlayer: byPlayer,
	}, nil
}
