// Package handlers provides HTTP request handlers for the game statistics API.
// It includes handlers for processing game score uploads and returning statistical analysis.
package handlers

import (
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/pipeline"
	"cli-stat-creator/internal/render"
	"net/http"
	"strconv"
	"strings"
)

// HandleStats processes uploaded game score files and returns statistical analysis.
// It accepts a multipart form file upload with a JSON file, processes the scores through
// the pipeline with optional filtering, and returns overall, level, and player statistics.
// Query parameters supported: min-score, max-score, level, include (overall/levels/players/all).
func HandleStats(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		render.Error(w, 413, "File too large (max 10MB)", "FILE_TOO_LARGE")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		render.Error(w, 400, "No file uploaded", "MISSING_FILE")
		return
	}
	defer file.Close()
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".json") {
		render.Error(w, 415, "File must be .json", "INVALID_FILE")
		return
	}
	cfg, err := parseQueryParams(r)
	if err != nil {
		render.Error(w, 400, "Invalid query parameters", "INVALID_PARAMS")
		return
	}
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	p := pipeline.New(cfg, pipeline.Filter(cfg))
	results, err := p.Run(ctx, pipeline.ReaderProvider{Reader: file})
	if err != nil {
		logger.Error("Pipeline execution failed", "error", err)
		render.Error(w, 500, "Failed to process file", "PROCESSING_ERROR")
		return
	}
	logger.Info("Pipeline completed successfully",
		"overall_count", results.Overall.TotalGamesPlayed,
		"player_count", len(results.ByPlayer),
		"level_count", len(results.ByLevel),
	)
	render.JSON(w, 200, results)
}
func parseIncludeParam(include string) (overall, levels, players bool) {
	// Default to all if not specified
	if include == "" {
		return true, true, true
	}

	overall, levels, players = false, false, false
	includeStrings := strings.SplitSeq(include, ",")
	for string := range includeStrings {
		switch string {
		case "overall":
			overall = true
		case "levels":
			levels = true
		case "players":
			players = true
		case "all":
			overall, levels, players = true, true, true
		}
	}
	return overall, levels, players
}

func parseOptionalInt(value string, defaultValue int) (int, error) {
	if value == "" {
		return defaultValue, nil
	}
	result, err := strconv.Atoi(value)
	if err != nil {
		return 0, err
	}
	return result, nil
}

func parseQueryParams(r *http.Request) (pipeline.Config, error) {
	minScoreParam := r.URL.Query().Get("min-score")
	minScore, err := parseOptionalInt(minScoreParam, 0)
	if err != nil {
		return pipeline.Config{}, err
	}
	maxScoreParam := r.URL.Query().Get("max-score")
	maxScore, err := parseOptionalInt(maxScoreParam, 0)
	if err != nil {
		return pipeline.Config{}, err
	}
	levelString := r.URL.Query().Get("level")
	levels, err := pipeline.ParseLevelString(levelString)
	if err != nil {
		return pipeline.Config{}, err
	}

	// Parse include parameter to determine which stats to calculate
	includeParam := r.URL.Query().Get("include")
	_, includeLevels, includePlayers := parseIncludeParam(includeParam)

	return pipeline.Config{
		MinScore:          minScore,
		MaxScore:          maxScore,
		FilterByLevel:     levels,
		CalculateByLevel:  includeLevels,
		CalculateByPlayer: includePlayers,
	}, nil
}
