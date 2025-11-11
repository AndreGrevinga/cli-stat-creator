// Package main implements a CLI application for analyzing game score statistics.
// It reads JSON files containing game scores and calculates various statistics
// including averages, medians, min/max values, and per-level breakdowns.
package main

import (
	"cli-stat-creator/internal/display"
	"cli-stat-creator/internal/reader"
	"cli-stat-creator/internal/stats"
	"fmt"
)

// main is the entry point of the application.
// It reads game scores from a JSON file (data/input.json), calculates statistics,
// and displays the results in formatted tables including overall statistics
// and per-level average scores.
func main() {
	fmt.Println("Please input the file path to analyze")
	var filepath string
	//fmt.Scan(&filepath)
	filepath = "data/input.json"
	var gameScores []stats.GameScore
	gameScores, err := reader.ReadScoresFromFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	statistics, err := stats.CalculateStatistics(gameScores)
	if err != nil {
		fmt.Println("Error calculating statistics:", err)
		return
	}
	display.RenderStatistics(statistics)
	display.RenderLevelBreakdown(statistics)
}
