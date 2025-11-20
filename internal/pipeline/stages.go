package pipeline

import (
	"cli-stat-creator/internal/reader"
	"cli-stat-creator/internal/stats"
	"context"
	"errors"
)

type Stage func(ctx context.Context, in <-chan stats.GameScore) <-chan stats.GameScore

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
	isScoreValid := (config.MaxScore == 0 || score.Score < config.MaxScore) && score.Score > config.MinScore
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

func Aggregate(ctx context.Context, in <-chan stats.GameScore) (Results, error) {
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
		return Results{}, errors.New("No scores after filtering")
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
