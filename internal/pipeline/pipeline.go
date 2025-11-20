package pipeline

import (
	"cli-stat-creator/internal/stats"
	"context"
)

// Config holds configuration options for filtering game scores in the pipeline.
type Config struct {
	FilterByPlayer     []string
	FilterByLevel      []int
	MinScore, MaxScore int
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
}

// New creates a new Pipeline with the specified stages.
// Stages will be executed in the order they are provided.
func New(stages ...Stage) *Pipeline {
	pipeline := Pipeline{stages: stages}
	return &pipeline
}

// Run executes the pipeline by reading game scores from a file and processing them
// through all configured stages. It returns comprehensive statistics including overall,
// per-level, and per-player results. The context can be used to cancel the operation.
func (p *Pipeline) Run(ctx context.Context, filename string) (Results,
	error) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Start source
	in, err := Source(ctx, filename)
	if err != nil {
		return Results{}, err
	}

	// Chain stages
	for _, stage := range p.stages {
		in = stage(ctx, in)
	}

	// Aggregate and return
	return Aggregate(ctx, in)
}
