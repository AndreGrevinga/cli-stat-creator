// Package main implements a CLI application for analyzing game score statistics.
// It reads JSON files containing game scores and calculates various statistics
// including averages, medians, min/max values, and per-level breakdowns.
package main

import (
	"cli-stat-creator/internal/display"
	"cli-stat-creator/internal/reader"
	"cli-stat-creator/internal/stats"
	"flag"
	"fmt"
)

var (
	detailed     bool
	defaultInput bool
	noPlayers    bool
	noLevels     bool
)

// main is the entry point of the application.
// It reads game scores from a JSON file (configurable via -i flag or user input),
// calculates statistics, and displays the results in formatted tables including
// overall statistics, per-player, and per-level average scores.
func main() {
	flag.BoolVar(&detailed, "detailed", false, "Show all statistics columns")
	flag.BoolVar(&detailed, "d", false, "Show all statistics columns (shorthand)")
	flag.BoolVar(&defaultInput, "i", false, "Use the default input.json (shorthand)")
	flag.BoolVar(&defaultInput, "default-input", false, "Use the default input.json")
	flag.BoolVar(&noPlayers, "no-players", false, "Hide player statistics")
	flag.BoolVar(&noLevels, "no-levels", false, "Hide level statistics")
	flag.Parse()
	var filepath string
	if defaultInput {
		filepath = "data/input.json"
	} else {
		fmt.Println("Please input the file path to analyze")
		fmt.Scan(&filepath)
	}
	gameScores, err := reader.ReadScoresFromFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	var levelGroupedStatistics map[int]stats.Statistics
	if !noLevels {
		levelGroupedStatistics, err = stats.CalculateStatisticsByLevel(gameScores)
		if err != nil {
			fmt.Println("Error calculating level statistics:", err)
			return
		}
	}
	var playerGroupedStatistics map[stats.Player]stats.Statistics
	if !noPlayers {
		playerGroupedStatistics, err = stats.CalculateStatisticsByPlayer(gameScores)
		if err != nil {
			fmt.Println("Error calculating player statistics:", err)
			return
		}
	}
	statistics, err := stats.CalculateStatistics(gameScores)
	if err != nil {
		fmt.Println("Error calculating statistics:", err)
		return
	}
	display.RenderStatistics(statistics)
	if !noPlayers {
		fmt.Println("\nPlayer Statistics:")
		display.RenderPlayerStatistics(playerGroupedStatistics, detailed)
	}
	if !noLevels {
		fmt.Println("\nLevel Statistics:")
		display.RenderLevelStatistics(levelGroupedStatistics, detailed)
	}
}
