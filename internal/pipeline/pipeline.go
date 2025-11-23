package pipeline

import (
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/stats"
	"context"
)

// Config holds configuration options for filtering game scores in the pipeline.
type Config struct {
	FilterByPlayer     []string
	FilterByLevel      []int
	MinScore, MaxScore int
	CalculateByLevel   bool
	CalculateByPlayer  bool
}

// Results contains all calculated statistics from a pipeline run.
// It includes overall statistics, as well as statistics grouped by level and player.
type Results struct {
	Overall  stats.Statistics
	ByLevel  map[int]stats.Statistics
	ByPlayer map[stats.Player]stats.Statistics
}

// Pipeline represents a series of processing stages that game scores flow through.
// Stages are executed sequentially in a concurrent, channel-based architecture.
type Pipeline struct {
	stages []Stage
	config Config
}

// New creates a new Pipeline with the specified configuration and stages.
// Stages will be executed in the order they are provided.
func New(config Config, stages ...Stage) *Pipeline {
	pipeline := Pipeline{
		stages: stages,
		config: config,
	}
	return &pipeline
}

// Run executes the pipeline by reading game scores from a file and processing them
// through all configured stages. It returns comprehensive statistics including overall,
// per-level, and per-player results. The context can be used to cancel the operation.
func (p *Pipeline) Run(ctx context.Context, filename string) (Results,
	error) {
	logger := logging.FromContext(ctx)
	logger.Info("Pipeline started", "stage_count", len(p.stages))

	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	in, err := Source(ctx, filename)
	if err != nil {
		logger.Error("Pipeline failed at source stage", "error", err)
		return Results{}, err
	}

	for _, stage := range p.stages {
		in = stage(ctx, in)
	}

	results, err := Aggregate(ctx, in, p.config)
	if err != nil {
		logger.Error("Pipeline failed at aggregation stage", "error", err)
		return Results{}, err
	}

	logger.Info("Pipeline completed successfully",
		"overall_count", results.Overall.TotalGamesPlayed,
		"by_level_count", len(results.ByLevel),
		"by_player_count", len(results.ByPlayer),
	)

	return results, nil
}
