package pipeline

import (
	"cli-stat-creator/internal/stats"
	"context"
)

type Config struct {
	FilterByPlayer     []string
	FilterByLevel      []int
	MinScore, MaxScore int
}

type Results struct {
	Overall  stats.Statistics
	ByLevel  map[int]stats.Statistics
	ByPlayer map[stats.Player]stats.Statistics
}

type Pipeline struct {
	stages []Stage
}

func New(stages ...Stage) *Pipeline {
	pipeline := Pipeline{stages: stages}
	return &pipeline
}

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
