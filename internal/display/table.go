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

func RenderLevelStatistics(statistics map[int]stats.Statistics, detailed bool) {
	table := tablewriter.NewTable(os.Stdout)
	var header, data []string
	if detailed {
		header = []string{"Level", "GamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"}
	} else {
		header = []string{"Level", "GamesPlayed", "AverageScore", "MaximumScore"}
	}
	table.Header(header)
	sortedLevels := make([]int, 0, len(statistics))
	for level, _ := range statistics {
		sortedLevels = append(sortedLevels, level)
	}
	slices.Sort(sortedLevels)
	for _, level := range sortedLevels {
		stats := statistics[level]
		if detailed {
			data = []string{
				fmt.Sprintf("%d", level),
				fmt.Sprintf("%d", stats.TotalGamesPlayed), fmt.Sprintf("%d", stats.TotalScore),
				fmt.Sprintf("%.2f", stats.AverageScore), fmt.Sprintf("%.2f", stats.MedianScore),
				fmt.Sprintf("%d", stats.MinimumScore), fmt.Sprintf("%d", stats.MaximumScore),
			}
		} else {
			data = []string{
				fmt.Sprintf("%d", level), fmt.Sprintf("%d", stats.TotalGamesPlayed),
				fmt.Sprintf("%.2f", stats.AverageScore), fmt.Sprintf("%d", stats.MaximumScore)}
		}
		table.Append(data)
	}
	table.Render()
}

func RenderPlayerStatistics(statistics map[stats.Player]stats.Statistics, detailed bool) {
	table := tablewriter.NewTable(os.Stdout)
	var header, data []string
	if detailed {
		header = []string{"Player", "GamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"}
	} else {
		header = []string{"Player", "GamesPlayed", "AverageScore", "MaximumScore"}
	}
	table.Header(header)
	sortedPlayers := make([]stats.Player, 0, len(statistics))
	for player, _ := range statistics {
		sortedPlayers = append(sortedPlayers, player)
	}
	slices.SortFunc(sortedPlayers, func(a, b stats.Player) int {
		if a.Name < b.Name {
			return -1
		} else if a.Name > b.Name {
			return 1
		}
		return 0
	})
	for _, player := range sortedPlayers {
		stats := statistics[player]
		if detailed {
			data = []string{
				player.Name,
				fmt.Sprintf("%d", stats.TotalGamesPlayed), fmt.Sprintf("%d", stats.TotalScore),
				fmt.Sprintf("%.2f", stats.AverageScore), fmt.Sprintf("%.2f", stats.MedianScore),
				fmt.Sprintf("%d", stats.MinimumScore), fmt.Sprintf("%d", stats.MaximumScore),
			}
		} else {
			data = []string{
				player.Name, fmt.Sprintf("%d", stats.TotalGamesPlayed),
				fmt.Sprintf("%.2f", stats.AverageScore), fmt.Sprintf("%d", stats.MaximumScore)}
		}
		table.Append(data)
	}
	table.Render()
}
