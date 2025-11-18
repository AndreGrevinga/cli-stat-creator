// Package display provides functions for rendering statistics and data
// in formatted tables for console output.
package display

import (
	"cli-stat-creator/internal/stats"
	"fmt"
	"os"
	"slices"
	"strconv"

	"github.com/olekukonko/tablewriter"
)

const defaultPrecision = 2

func formatFloat(f float64) string {
	return strconv.FormatFloat(f, 'f', defaultPrecision, 64)
}

// RenderStatistics displays overall game statistics in a formatted table.
// It renders total games played, total score, average score, median score,
// minimum score, and maximum score to stdout.
func RenderStatistics(s stats.Statistics) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{"TotalGamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"})
	table.Append([]string{
		fmt.Sprintf("%d", s.TotalGamesPlayed),
		fmt.Sprintf("%d", s.TotalScore),
		formatFloat(s.AverageScore),
		formatFloat(s.MedianScore),
		fmt.Sprintf("%d", s.MinimumScore),
		fmt.Sprintf("%d", s.MaximumScore),
	})
	table.Render()
}

// RenderGroupedStatistics is a generic function that renders statistics grouped by any comparable key type.
// It accepts a map of statistics, a detailed flag, the column title for the key,
// a sorting function for the keys, and a formatting function to convert keys to strings.
func RenderGroupedStatistics[K comparable](
	statistics map[K]stats.Statistics,
	detailed bool,
	keyColumnTitle string,
	sortKeys func([]K),
	formatKey func(K) string,
) {
	table := tablewriter.NewTable(os.Stdout)

	var header []string
	if detailed {
		header = []string{keyColumnTitle, "GamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"}
	} else {
		header = []string{keyColumnTitle, "GamesPlayed", "AverageScore", "MaximumScore"}
	}
	table.Header(header)

	// Extract keys into a slice
	keys := make([]K, 0, len(statistics))
	for key := range statistics {
		keys = append(keys, key)
	}

	// Sort using the provided sorting function
	sortKeys(keys)

	// Build table rows
	for _, key := range keys {
		s := statistics[key]
		var data []string
		if detailed {
			data = []string{
				formatKey(key),
				fmt.Sprintf("%d", s.TotalGamesPlayed),
				fmt.Sprintf("%d", s.TotalScore),
				formatFloat(s.AverageScore),
				formatFloat(s.MedianScore),
				fmt.Sprintf("%d", s.MinimumScore),
				fmt.Sprintf("%d", s.MaximumScore),
			}
		} else {
			data = []string{
				formatKey(key),
				fmt.Sprintf("%d", s.TotalGamesPlayed),
				formatFloat(s.AverageScore),
				fmt.Sprintf("%d", s.MaximumScore),
			}
		}
		table.Append(data)
	}
	table.Render()
}

// RenderLevelStatistics displays statistics grouped by level in a formatted table.
// It renders level numbers sorted numerically along with their statistics to stdout.
func RenderLevelStatistics(statistics map[int]stats.Statistics, detailed bool) {
	RenderGroupedStatistics(
		statistics,
		detailed,
		"Level",
		func(levels []int) { slices.Sort(levels) },
		func(level int) string { return fmt.Sprintf("%d", level) },
	)
}

// RenderPlayerStatistics displays statistics grouped by player in a formatted table.
// It renders player names sorted alphabetically along with their statistics to stdout.
func RenderPlayerStatistics(statistics map[stats.Player]stats.Statistics, detailed bool) {
	RenderGroupedStatistics(
		statistics,
		detailed,
		"Player",
		func(players []stats.Player) {
			slices.SortFunc(players, func(a, b stats.Player) int {
				if a.Name < b.Name {
					return -1
				} else if a.Name > b.Name {
					return 1
				}
				return 0
			})
		},
		func(player stats.Player) string { return player.Name },
	)
}
