// Package display provides functions for rendering statistics and data
// in formatted tables for console output.
package display

import (
	"cli-stat-creator/internal/stats"
	"fmt"
	"os"
	"slices"

	"github.com/olekukonko/tablewriter"
)

// RenderStatistics displays overall game statistics in a formatted table.
// It renders total games played, total score, average score, median score,
// minimum score, and maximum score to stdout.
func RenderStatistics(s stats.Statistics) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{"TotalGamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"})
	table.Append([]string{
		fmt.Sprintf("%d", s.TotalGamesPlayed), fmt.Sprintf("%d", s.TotalScore),
		fmt.Sprintf("%.2f", s.AverageScore), fmt.Sprintf("%.2f", s.MedianScore),
		fmt.Sprintf("%d", s.MinimumScore), fmt.Sprintf("%d", s.MaximumScore),
	})
	table.Render()
}

// RenderLevelBreakdown displays average scores grouped by level in a formatted table.
// It renders each level number along with its corresponding average score to stdout.
// The order of levels displayed is non-deterministic due to map iteration.
func RenderLevelBreakdown(s stats.Statistics) {
	table := tablewriter.NewTable(os.Stdout)

	levels := make([]int, 0, len(s.AverageScoreByLevel))
	for level := range s.AverageScoreByLevel {
		levels = append(levels, level)
	}
	slices.Sort(levels)
	table.Header([]string{"Level", "Average Score"})
	for _, level := range levels {
		averageScore := s.AverageScoreByLevel[level]
		table.Append([]string{fmt.Sprintf("%d", level), fmt.Sprintf("%.2f", averageScore)})
	}
	table.Render()
}
