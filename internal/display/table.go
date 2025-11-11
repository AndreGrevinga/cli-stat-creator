package display

import (
	"cli-stat-creator/internal/stats"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
)

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

func RenderLevelBreakdown(s stats.Statistics) {
	table := tablewriter.NewTable(os.Stdout)
	table.Header([]string{"Level", "Average Score"})
	for level, averageScore := range s.AverageScoreByLevel {
		table.Append([]string{fmt.Sprintf("%d", level), fmt.Sprintf("%.2f", averageScore)})
	}
	table.Render()
}
