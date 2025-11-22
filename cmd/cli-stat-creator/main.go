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
	"strconv"
	"strings"
)

var (
	detailed     bool
	defaultInput bool
	noPlayers    = flag.Bool("no-players", false, "Hide player statistics")
	noLevels     = flag.Bool("no-levels", false, "Hide level statistics")
	// Todo: Implement player filter flag
	//playerString   = flag.String("player", "", "Only show statistics for given player")
	levelString    = flag.String("level", "", "Only show statistics for given levels")
	minScoreString = flag.String("min-score", "", "Minimum score filter")
	maxScoreString = flag.String("max-score", "", "Maximum score filter")
)

const defaultInputFile = "data/input.json"

// parseLevelFlag parses the level flag string into a slice of level integers.
// It accepts either a single level (e.g., "5") or a range (e.g., "1-5").
// Returns an empty slice if the input string is empty.
// Returns an error if the input contains non-numeric values or is malformed.
func parseLevelFlag(level string) ([]int, error) {
	var levels []int
	var err error
	var minLevel, maxLevel int64

	if level == "" {
		levels = make([]int, 0)
		return levels, nil
	}

	subStrings := strings.Split(level, "-")
	if len(subStrings) > 2 {
		return nil, fmt.Errorf("invalid level format: expected single level or range (e.g., '5' or '1-5'), got '%s'", level)
	}
	minLevel, err = strconv.ParseInt(subStrings[0], 0, 0)
	if err != nil {
		return nil, err
	}
	if len(subStrings) == 1 {
		maxLevel = minLevel
	} else {
		maxLevel, err = strconv.ParseInt(subStrings[1], 0, 0)
		if err != nil {
			return nil, err
		}
	}
	if minLevel > maxLevel {
		return nil, fmt.Errorf("invalid level range: min (%d) cannot be greater than max (%d)", minLevel, maxLevel)
	}
	if maxLevel-minLevel > 1000 { // reasonable upper bound
		return nil, fmt.Errorf("level range too large: maximum 1000 levels allowed")
	}
	levels = make([]int, 0, maxLevel-minLevel+1)
	for i := minLevel; i <= maxLevel; i++ {
		levels = append(levels, int(i))
	}
	return levels, nil
}

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
	var minScore, maxScore int64
	var err error
	if *minScoreString == "" {
		minScore = 0
	} else {
		minScore, err = strconv.ParseInt(*minScoreString, 0, 0)
		if err != nil {
			fmt.Printf("Error: invalid min-score value '%s': %v\n", *minScoreString, err)
			return
		}
	}
	if *maxScoreString == "" {
		maxScore = 0
	} else {
		maxScore, err = strconv.ParseInt(*maxScoreString, 0, 0)
		if err != nil {
			fmt.Printf("Error: invalid max-score value '%s': %v\n", *maxScoreString, err)
			return
		}
	}
	levels, err := parseLevelFlag(*levelString)
	if err != nil {
		fmt.Printf("Error: invalid level value '%s': %v\n", *levelString, err)
		return
	}
	config := pipeline.Config{
		CalculateByLevel:  !*noLevels,
		CalculateByPlayer: !*noPlayers,
		MinScore:          int(minScore),
		MaxScore:          int(maxScore),
		FilterByLevel:     levels,
	}
	var filepath string
	if defaultInput {
		filepath = defaultInputFile
	} else {
		fmt.Println("Please input the file path to analyze")
		fmt.Scan(&filepath)
	}
	p := pipeline.New(
		config,
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
