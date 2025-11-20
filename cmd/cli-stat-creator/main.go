// Package main implements a CLI application for analyzing game score statistics.
// It reads JSON files containing game scores and calculates various statistics
// including averages, medians, min/max values, and per-level breakdowns.
package main

import (
	"cli-stat-creator/internal/display"
	"cli-stat-creator/internal/pipeline"
	"context"
	"flag"
	"fmt"
)

var (
	detailed     bool
	defaultInput bool
	noPlayers    = flag.Bool("no-players", false, "Hide player statistics")
	noLevels     = flag.Bool("no-levels", false, "Hide level statistics")
)

const defaultInputFile = "data/input.json"

// main is the entry point of the application.
// It reads game scores from a JSON file (data/input.json), calculates statistics,
// and displays the results in formatted tables including overall statistics
// and per-level average scores.
func main() {
	flag.BoolVar(&detailed, "detailed", false, "Show all statistics columns")
	flag.BoolVar(&detailed, "d", false, "Show all statistics columns (shorthand)")
	flag.BoolVar(&defaultInput, "i", false, "Uses the default input.json (shorthand)")
	flag.BoolVar(&defaultInput, "default-input", false, "Uses the default input.json")
	flag.Parse()
	ctx := context.Background()
	config := pipeline.Config{}
	var filepath string
	if defaultInput {
		filepath = defaultInputFile
	} else {
		fmt.Println("Please input the file path to analyze")
		fmt.Scan(&filepath)
	}
	p := pipeline.New(
		pipeline.Filter(config),
	)
	results, err := p.Run(ctx, filepath)
	if err != nil {
		fmt.Println("Error calculating statistics:", err)
		return
	}
	display.RenderStatistics(results.Overall)
	if !*noPlayers {
		fmt.Println("\nPlayer Statistics:")
		display.RenderPlayerStatistics(results.ByPlayer, detailed)
	}
	if !*noLevels {
		fmt.Println("\nLevel Statistics:")
		display.RenderLevelStatistics(results.ByLevel, detailed)
	}
}
