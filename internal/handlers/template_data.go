package handlers

import "cli-stat-creator/internal/stats"

// ResultsData holds all data needed to render the results template.
type ResultsData struct {
	Overall      StatsWithDetailed
	LevelStats   map[int]stats.Statistics
	PlayerStats  map[stats.Player]stats.Statistics
	ShowDetailed bool
}

// StatsWithDetailed wraps Statistics with a ShowDetailed flag for template rendering.
type StatsWithDetailed struct {
	stats.Statistics
	ShowDetailed bool
}

// ErrorData holds error information for rendering error templates.
type ErrorData struct {
	Message string
	Code    string
}
