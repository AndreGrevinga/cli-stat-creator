// Package handlers provides HTTP request handlers for the game statistics API.
// It includes handlers for processing game score uploads and returning statistical analysis.
package handlers

import (
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/pipeline"
	"cli-stat-creator/internal/render"
	"cli-stat-creator/internal/templates"
	"net/http"
	"strconv"
	"strings"
)

var tmplManager *templates.Manager

/*func init() {
tmplManager := templates.New("web/templates")
}*/

// HandleStats processes uploaded game score files and returns statistical analysis.
// It accepts a multipart form file upload with a JSON file, processes the scores through
// the pipeline with optional filtering, and returns results as either HTML (for htmx/browsers)
// or JSON (for API clients) based on the Accept header.
func HandleStats(w http.ResponseWriter, r *http.Request) {
	// Parse multipart form with 10MB limit
	err := r.ParseMultipartForm(10 << 20)
	if err != nil {
		respondWithError(w, r, 413, "File too large (max 10MB)", "FILE_TOO_LARGE")
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		respondWithError(w, r, 400, "No file uploaded", "MISSING_FILE")
		return
	}
	defer file.Close()
	filename := header.Filename
	if !strings.HasSuffix(strings.ToLower(filename), ".json") {
		respondWithError(w, r, 415, "File must be .json", "INVALID_FILE")
		return
	}
	cfg, err := parseQueryParams(r)
	if err != nil {
		respondWithError(w, r, 400, "Invalid query parameters", "INVALID_PARAMS")
		return
	}
	ctx := r.Context()
	logger := logging.FromContext(ctx)
	p := pipeline.New(cfg, pipeline.Filter(cfg))
	results, err := p.Run(ctx, pipeline.ReaderProvider{Reader: file})
	if err != nil {
		logger.Error("Pipeline execution failed", "error", err)
		respondWithError(w, r, 500, "Failed to process file", "PROCESSING_ERROR")
		return
	}
	logger.Info("Pipeline completed successfully",
		"overall_count", results.Overall.TotalGamesPlayed,
		"player_count", len(results.ByPlayer),
		"level_count", len(results.ByLevel),
	)
	if wantsHTML(r) {
		renderHTMLResults(w, results, cfg)
	} else {
		render.JSON(w, 200, results)
	}
}

// HandleClear returns empty HTML to clear the results section.
// Used by htmx to clear displayed results when the clear button is clicked.
func HandleClear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(""))
}

// wantsHTML checks if the client prefers HTML over JSON.
// htmx automatically sends "text/html" in the Accept header.
// This allows the same endpoint to serve both web browsers and API clients.
func wantsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")

	if strings.Contains(accept, "text/html") {
		return true
	}

	if r.Header.Get("HX-Request") == "true" {
		return true
	}
	return false
}

// respondWithError sends an error response in either HTML or JSON format
// based on the client's Accept header.
func respondWithError(w http.ResponseWriter, r *http.Request, status int,
	msg, code string) {
	if wantsHTML(r) {
		renderHTMLError(w, status, msg, code)
	} else {
		render.Error(w, status, msg, code)
	}
}

// renderHTMLResults renders the statistics results as HTML using templates.
func renderHTMLResults(w http.ResponseWriter, results pipeline.Results, cfg pipeline.Config) {
	data := ResultsData{
		Overall: StatsWithDetailed{
			Statistics:   results.Overall,
			ShowDetailed: cfg.ShowDetailed,
		},
		LevelStats:   results.ByLevel,
		PlayerStats:  results.ByPlayer,
		ShowDetailed: cfg.ShowDetailed,
	}
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	err := tmplManager.Render(w, "results.html", data)
	if err != nil {
		http.Error(w, "Template rendering error: "+err.Error(), 500)
	}
}

// renderHTMLError renders an error message as HTML using the error template.
func renderHTMLError(w http.ResponseWriter, status int, msg, code string) {
	data := ErrorData{
		Message: msg,
		Code:    code,
	}

	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(status)
	err := tmplManager.Render(w, "error.html", data)
	if err != nil {
		// Fallback to plain text if template fails
		http.Error(w, msg, status)
	}
}

// parseIncludeParam parses the include parameter from the URL query.
func parseIncludeParam(include string) (overall, levels, players bool) {
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

	includeParam := r.URL.Query().Get("include")
	_, includeLevels, includePlayers := parseIncludeParam(includeParam)

	detailed := r.URL.Query().Get("detaled") == "true"

	return pipeline.Config{
		MinScore:          minScore,
		MaxScore:          maxScore,
		FilterByLevel:     levels,
		CalculateByLevel:  includeLevels,
		CalculateByPlayer: includePlayers,
		ShowDetailed:      detailed,
	}, nil
}
