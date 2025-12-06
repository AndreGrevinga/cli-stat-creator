# HTMX Migration Guide

**Date:** 2025-12-05
**Status:** Implementation Ready
**Migration Type:** Full Migration with Dual Format Support

## Overview

This guide walks through migrating the Game Score Statistics web interface from vanilla JavaScript with manual DOM manipulation to htmx for server-driven HTML responses. This will reduce JavaScript code by ~85% (from 335 lines to ~50 lines) and simplify the architecture.

**Key Feature:** The migration maintains **full backward compatibility** by supporting both JSON (for API clients) and HTML (for htmx/browsers) responses from the same endpoint, determined by the `Accept` header. This ensures existing API consumers continue to work while enabling modern htmx-driven web interactions.

---

## Phase 1: Setup htmx Infrastructure and Templating

### 1.1: Add htmx to index.html ✓

**Status:** Already completed (htmx script tag is present on line 8)

### 1.2: Create Templates Directory

```bash
mkdir -p web/templates
```

### 1.3: Create internal/templates Package

Create `internal/templates/templates.go`:

```go
// Package templates provides HTML template management and rendering functionality.
// It supports template caching, parsing, and execution for the web interface.
package templates

import (
	"fmt"
	"html/template"
	"io"
	"path/filepath"
	"sync"
)

// Manager handles loading, caching, and rendering of HTML templates.
// It provides thread-safe template caching for improved performance.
type Manager struct {
	cache map[string]*template.Template
	dir   string
	mu    sync.RWMutex
}

// New creates a new template manager that loads templates from the specified directory.
func New(dir string) *Manager {
	return &Manager{
		cache: make(map[string]*template.Template),
		dir:   dir,
	}
}

// Get retrieves a parsed template from cache or loads it from disk.
// Templates are cached after first load for improved performance.
func (m *Manager) Get(name string) (*template.Template, error) {
	m.mu.RLock()
	tmpl, ok := m.cache[name]
	m.mu.RUnlock()

	if ok {
		return tmpl, nil
	}

	// Template not in cache, load it
	m.mu.Lock()
	defer m.mu.Unlock()

	// Double-check in case another goroutine loaded it
	if tmpl, ok := m.cache[name]; ok {
		return tmpl, nil
	}

	// Parse template file
	path := filepath.Join(m.dir, name)
	tmpl, err := template.ParseFiles(path)
	if err != nil {
		return nil, fmt.Errorf("failed to parse template %s: %w", name, err)
	}

	m.cache[name] = tmpl
	return tmpl, nil
}

// Render executes a template with the provided data and writes to the writer.
func (m *Manager) Render(w io.Writer, name string, data interface{}) error {
	tmpl, err := m.Get(name)
	if err != nil {
		return err
	}
	return tmpl.Execute(w, data)
}

// ClearCache removes all cached templates, forcing reload on next access.
// Useful during development when templates are being modified.
func (m *Manager) ClearCache() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.cache = make(map[string]*template.Template)
}
```

---

## Phase 2: Create HTML Templates

### 2.1: Create Overall Statistics Template

Create `web/templates/overall-stats.html`:

```html
<div class="stats-card">
    <h3 class="card-title">Overall Statistics</h3>
    <table class="stats-table">
        <tbody>
            <tr>
                <td style="font-weight: 500;">Total Games</td>
                <td class="stat-value">{{.TotalGamesPlayed}}</td>
            </tr>
            <tr>
                <td style="font-weight: 500;">Total Score</td>
                <td class="stat-value">{{.TotalScore}}</td>
            </tr>
            <tr>
                <td style="font-weight: 500;">Average Score</td>
                <td class="stat-value">{{printf "%.2f" .AverageScore}}</td>
            </tr>
            <tr>
                <td style="font-weight: 500;">Median Score</td>
                <td class="stat-value">{{printf "%.2f" .MedianScore}}</td>
            </tr>
            {{if .ShowDetailed}}
            <tr>
                <td style="font-weight: 500;">Minimum Score</td>
                <td class="stat-value">{{.MinimumScore}}</td>
            </tr>
            <tr>
                <td style="font-weight: 500;">Maximum Score</td>
                <td class="stat-value">{{.MaximumScore}}</td>
            </tr>
            {{end}}
        </tbody>
    </table>
</div>
```

### 2.2: Create Level Statistics Template

Create `web/templates/level-stats.html`:

```html
{{if .LevelStats}}
<div class="stats-card">
    <details open>
        <summary class="card-title">Statistics by Level</summary>
        <table class="stats-table">
            <thead>
                <tr>
                    <th>Level</th>
                    <th>Total Games</th>
                    <th>Total Score</th>
                    <th>Average</th>
                    <th>Median</th>
                    {{if .ShowDetailed}}
                    <th>Min</th>
                    <th>Max</th>
                    {{end}}
                </tr>
            </thead>
            <tbody>
                {{range $level, $stats := .LevelStats}}
                <tr>
                    <td>{{$level}}</td>
                    <td class="stat-value">{{$stats.TotalGamesPlayed}}</td>
                    <td class="stat-value">{{$stats.TotalScore}}</td>
                    <td class="stat-value">{{printf "%.2f" $stats.AverageScore}}</td>
                    <td class="stat-value">{{printf "%.2f" $stats.MedianScore}}</td>
                    {{if $.ShowDetailed}}
                    <td class="stat-value">{{$stats.MinimumScore}}</td>
                    <td class="stat-value">{{$stats.MaximumScore}}</td>
                    {{end}}
                </tr>
                {{end}}
            </tbody>
        </table>
    </details>
</div>
{{end}}
```

### 2.3: Create Player Statistics Template

Create `web/templates/player-stats.html`:

```html
{{if .PlayerStats}}
<div class="stats-card">
    <details open>
        <summary class="card-title">Statistics by Player</summary>
        <table class="stats-table">
            <thead>
                <tr>
                    <th>Player</th>
                    <th>Total Games</th>
                    <th>Total Score</th>
                    <th>Average</th>
                    <th>Median</th>
                    {{if .ShowDetailed}}
                    <th>Min</th>
                    <th>Max</th>
                    {{end}}
                </tr>
            </thead>
            <tbody>
                {{range $player, $stats := .PlayerStats}}
                <tr>
                    <td>{{$player.Name}} (ID: {{$player.ID}})</td>
                    <td class="stat-value">{{$stats.TotalGamesPlayed}}</td>
                    <td class="stat-value">{{$stats.TotalScore}}</td>
                    <td class="stat-value">{{printf "%.2f" $stats.AverageScore}}</td>
                    <td class="stat-value">{{printf "%.2f" $stats.MedianScore}}</td>
                    {{if $.ShowDetailed}}
                    <td class="stat-value">{{$stats.MinimumScore}}</td>
                    <td class="stat-value">{{$stats.MaximumScore}}</td>
                    {{end}}
                </tr>
                {{end}}
            </tbody>
        </table>
    </details>
</div>
{{end}}
```

### 2.4: Create Error Template

Create `web/templates/error.html`:

```html
<section class="error-section">
    <div class="error-message">
        <svg
            xmlns="http://www.w3.org/2000/svg"
            width="24"
            height="24"
            viewBox="0 0 24 24"
            fill="none"
            stroke="currentColor"
            stroke-width="2"
            stroke-linecap="round"
            stroke-linejoin="round"
        >
            <circle cx="12" cy="12" r="10"></circle>
            <line x1="12" y1="8" x2="12" y2="12"></line>
            <line x1="12" y1="16" x2="12.01" y2="16"></line>
        </svg>
        <div>
            <strong>Error</strong>
            <p>{{.Message}}</p>
        </div>
    </div>
</section>
```

### 2.5: Create Results Container Template

Create `web/templates/results.html`:

```html
<section class="results-section">
    <div class="results-header">
        <h2>Analysis Results</h2>
        <button
            id="clear-results"
            class="clear-btn"
            hx-get="/api/clear"
            hx-target="#results-section"
            hx-swap="innerHTML"
        >
            Clear Results
        </button>
    </div>

    {{template "overall-stats.html" .}}
    {{if .LevelStats}}
    {{template "level-stats.html" .}}
    {{end}}
    {{if .PlayerStats}}
    {{template "player-stats.html" .}}
    {{end}}
</section>
```

---

## Phase 3: Update Backend for HTML Rendering

### 3.1: Add ShowDetailed to pipeline.Config

Edit `internal/pipeline/pipeline.go`:

```go
// Config holds configuration options for the statistics processing pipeline.
type Config struct {
	MinScore          int
	MaxScore          int
	FilterByLevel     []int
	CalculateByLevel  bool
	CalculateByPlayer bool
	ShowDetailed      bool  // Add this field
}
```

### 3.2: Create Template Data Structures

Create `internal/handlers/template_data.go`:

```go
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
```

### 3.3: Update handlers/stats.go (Dual JSON/HTML Support)

Edit `internal/handlers/stats.go`:

**Key Change:** Support both HTML (for htmx/browsers) and JSON (for API clients) based on the Accept header. This maintains backward compatibility with existing API consumers while enabling htmx for the web interface.

**Add imports:**
```go
import (
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/pipeline"
	"cli-stat-creator/internal/render"
	"cli-stat-creator/internal/templates"
	"net/http"
	"strconv"
	"strings"
)
```

**Add template manager:**
```go
var tmplManager *templates.Manager

func init() {
	tmplManager = templates.New("web/templates")
}
```

**Update HandleStats function:**
```go
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
		respondWithError(w, r, 400, "Invalid query parameters: "+err.Error(), "INVALID_PARAMS")
		return
	}

	ctx := r.Context()
	logger := logging.FromContext(ctx)
	p := pipeline.New(cfg, pipeline.Filter(cfg))
	results, err := p.Run(ctx, pipeline.ReaderProvider{Reader: file})
	if err != nil {
		logger.Error("Pipeline execution failed", "error", err)
		respondWithError(w, r, 500, "Failed to process file: "+err.Error(), "PROCESSING_ERROR")
		return
	}

	logger.Info("Pipeline completed successfully",
		"overall_count", results.Overall.TotalGamesPlayed,
		"player_count", len(results.ByPlayer),
		"level_count", len(results.ByLevel),
	)

	// Return HTML for browsers/htmx, JSON for API clients
	if wantsHTML(r) {
		renderHTMLResults(w, results, cfg)
	} else {
		render.JSON(w, 200, results)
	}
}
```

**Add helper functions:**
```go
// wantsHTML checks if the client prefers HTML over JSON.
// htmx automatically sends "text/html" in the Accept header.
// This allows the same endpoint to serve both web browsers and API clients.
func wantsHTML(r *http.Request) bool {
	accept := r.Header.Get("Accept")
	// Check if client explicitly wants HTML
	if strings.Contains(accept, "text/html") {
		return true
	}
	// Check for htmx-specific header (extra safety)
	if r.Header.Get("HX-Request") == "true" {
		return true
	}
	return false
}

// respondWithError sends an error response in either HTML or JSON format
// based on the client's Accept header.
func respondWithError(w http.ResponseWriter, r *http.Request, status int, msg, code string) {
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
```

**Update parseQueryParams:**
```go
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

	// Parse detailed flag
	detailed := r.URL.Query().Get("detailed") == "true"

	return pipeline.Config{
		MinScore:          minScore,
		MaxScore:          maxScore,
		FilterByLevel:     levels,
		CalculateByLevel:  includeLevels,
		CalculateByPlayer: includePlayers,
		ShowDetailed:      detailed,
	}, nil
}
```

### 3.4: Add HandleClear Endpoint

Add to `internal/handlers/stats.go`:

```go
// HandleClear returns empty HTML to clear the results section.
// Used by htmx to clear displayed results when the clear button is clicked.
func HandleClear(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.Write([]byte(""))
}
```

### 3.5: Add Route to http-server

Edit `cmd/http-server/main.go`:

```go
func setupRouter(logger *slog.Logger, staticDir string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggerMiddleware(logger))
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// API routes
	r.Post("/api/stats", handlers.HandleStats)
	r.Get("/api/clear", handlers.HandleClear)  // Add this line

	// Serve static files
	fileServer := http.FileServer(http.Dir(staticDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	return r
}
```

---

## Phase 4: Convert Frontend to htmx

### 4.1: Update Form with htmx Attributes

Edit `web/static/index.html`, update the form section:

```html
<form
    id="upload-form"
    hx-post="/api/stats"
    hx-target="#results-section"
    hx-encoding="multipart/form-data"
    hx-indicator="#loading-spinner"
>
    <div class="file-upload">
        <label for="file-input" class="file-label">
            <svg
                xmlns="http://www.w3.org/2000/svg"
                width="24"
                height="24"
                viewBox="0 0 24 24"
                fill="none"
                stroke="currentColor"
                stroke-width="2"
                stroke-linecap="round"
                stroke-linejoin="round"
            >
                <path d="M21 15v4a2 2 0 0 1-2 2H5a2 2 0 0 1-2-2v-4"></path>
                <polyline points="17 8 12 3 7 8"></polyline>
                <line x1="12" y1="3" x2="12" y2="15"></line>
            </svg>
            <span id="file-name">Choose JSON file or drag here</span>
        </label>
        <input
            type="file"
            id="file-input"
            name="file"
            accept=".json"
            required
        />
    </div>

    <div class="filters">
        <div class="filter-group">
            <label for="level-filter">Level Filter</label>
            <input
                type="text"
                id="level-filter"
                name="level"
                placeholder="e.g., 5 or 1-5"
            />
        </div>

        <div class="filter-row">
            <div class="filter-group">
                <label for="min-score">Min Score</label>
                <input
                    type="number"
                    id="min-score"
                    name="min-score"
                    min="0"
                    placeholder="0"
                />
            </div>

            <div class="filter-group">
                <label for="max-score">Max Score</label>
                <input
                    type="number"
                    id="max-score"
                    name="max-score"
                    min="0"
                    placeholder="No limit"
                />
            </div>
        </div>

        <div class="checkboxes">
            <label class="checkbox-label">
                <input
                    type="checkbox"
                    id="detailed"
                    name="detailed"
                    value="true"
                    checked
                />
                <span>Show detailed statistics (min/max)</span>
            </label>
            <label class="checkbox-label">
                <input
                    type="checkbox"
                    id="include-levels"
                    name="include"
                    value="levels"
                    checked
                />
                <span>Include per-level statistics</span>
            </label>
            <label class="checkbox-label">
                <input
                    type="checkbox"
                    id="include-players"
                    name="include"
                    value="players"
                    checked
                />
                <span>Include per-player statistics</span>
            </label>
        </div>
    </div>

    <button type="submit" class="submit-btn">
        <span class="btn-text">Analyze Statistics</span>
        <span id="loading-spinner" class="btn-loader htmx-indicator">
            <svg class="spinner" viewBox="0 0 50 50">
                <circle
                    class="path"
                    cx="25"
                    cy="25"
                    r="20"
                    fill="none"
                    stroke-width="5"
                ></circle>
            </svg>
        </span>
    </button>
</form>
```

### 4.2: Remove Error Section (Handled by Templates)

Remove the standalone error-section from `web/static/index.html`:

```html
<!-- DELETE THIS ENTIRE SECTION -->
<section
    id="error-section"
    class="error-section"
    style="display: none"
>
    ...
</section>
```

### 4.3: Update Results Section

Replace the results section with a simple target div:

```html
<!-- Replace the entire results-section with this: -->
<div id="results-section"></div>
```

### 4.4: Simplify app.js

Replace the entire contents of `web/static/app.js` with:

```javascript
// Keep ONLY drag-and-drop functionality
const fileLabel = document.querySelector('.file-label');
const fileInput = document.getElementById('file-input');
const fileName = document.getElementById('file-name');

// File drag and drop handlers
fileLabel.addEventListener('dragover', (e) => {
    e.preventDefault();
    fileLabel.style.borderColor = 'var(--primary-color)';
    fileLabel.style.background = '#eff6ff';
});

fileLabel.addEventListener('dragleave', () => {
    fileLabel.style.borderColor = 'var(--border-color)';
    fileLabel.style.background = 'var(--bg-secondary)';
});

fileLabel.addEventListener('drop', (e) => {
    e.preventDefault();
    fileLabel.style.borderColor = 'var(--border-color)';
    fileLabel.style.background = 'var(--bg-secondary)';

    const files = e.dataTransfer.files;
    if (files.length > 0) {
        fileInput.files = files;
        fileName.textContent = files[0].name;
    }
});

// File input change handler
fileInput.addEventListener('change', (e) => {
    if (e.target.files.length > 0) {
        fileName.textContent = e.target.files[0].name;
    }
});

// Optional: htmx event listeners for debugging/logging
document.body.addEventListener('htmx:afterSwap', (e) => {
    console.log('Results loaded successfully');
});

document.body.addEventListener('htmx:responseError', (e) => {
    console.error('Request failed:', e.detail);
});
```

**Code removed:** ~280 lines! (fetch, DOM manipulation, table rendering, error handling)

---

## Phase 5: Polish and Enhance

### 5.1: Add htmx CSS Transitions

Add to `web/static/styles.css`:

```css
/* htmx Loading Indicator */
.htmx-indicator {
    display: none;
}

.htmx-request .htmx-indicator {
    display: inline-block;
}

.htmx-request.htmx-indicator {
    display: inline-block;
}

/* htmx Swap Animations */
#results-section {
    transition: opacity 200ms ease-out;
}

.htmx-swapping #results-section {
    opacity: 0;
}

.htmx-settling #results-section {
    opacity: 1;
}

/* Error Section Animations */
.error-section {
    animation: fadeIn 200ms ease-out;
}

@keyframes fadeIn {
    from {
        opacity: 0;
        transform: translateY(-10px);
    }
    to {
        opacity: 1;
        transform: translateY(0);
    }
}
```

### 5.2: Update Button Loader Styles

Update the `.btn-loader` styles in `styles.css`:

```css
.btn-loader {
    display: none; /* Remove this if it conflicts */
}

/* Let htmx control visibility instead */
```

---

## Phase 6: Testing

### Testing Checklist

Test all these scenarios manually:

- [ ] **Dual Format Support (JSON vs HTML)**
  - [ ] Test JSON response with curl (default/no Accept header)
  - [ ] Test JSON response with Accept: application/json header
  - [ ] Test HTML response with Accept: text/html header
  - [ ] Test HTML response through web form (htmx automatic)
  - [ ] Verify both formats return same data (just different serialization)
  - [ ] Test error responses in both JSON and HTML formats

- [ ] **File Upload**
  - [ ] Upload valid JSON file from data/input.json
  - [ ] Upload invalid file (should show error template)
  - [ ] Upload non-JSON file (should show error)
  - [ ] Upload file > 10MB (should show error)

- [ ] **Filters**
  - [ ] Single level filter (e.g., level=3)
  - [ ] Level range filter (e.g., level=1-5)
  - [ ] Min score filter only
  - [ ] Max score filter only
  - [ ] Combined min/max score filters
  - [ ] All filters together

- [ ] **Display Options**
  - [ ] Toggle detailed checkbox (min/max columns should appear/disappear)
  - [ ] Uncheck include levels (level stats should not appear)
  - [ ] Uncheck include players (player stats should not appear)
  - [ ] Uncheck both (only overall stats should appear)

- [ ] **Interactions**
  - [ ] Drag and drop file
  - [ ] Click to select file
  - [ ] Clear results button
  - [ ] Loading spinner appears during processing
  - [ ] Multiple sequential uploads

- [ ] **Error Handling**
  - [ ] Invalid JSON content
  - [ ] Missing required fields
  - [ ] Negative scores (validation error)
  - [ ] Empty player names (validation error)
  - [ ] Network errors

- [ ] **Browser Console**
  - [ ] No JavaScript errors
  - [ ] htmx debug events logging correctly

---

## Phase 7: Documentation

### 7.1: Update CLAUDE.md

Add to the "## Project Overview" section:

```markdown
## Tech Stack
- **Backend**: Go 1.21+ with chi router
- **Frontend**: htmx for dynamic interactions (no build tooling required)
- **Templating**: Go html/template for server-side rendering
- **Styling**: Vanilla CSS
```

Add to the "## Code Style Guidelines" section:

```markdown
### HTML Templates
- Templates are located in `web/templates/`
- Use Go's `html/template` syntax
- Use `{{template "name.html" .}}` for composition
- Pass data structures from handlers via template data types
- All templates auto-escape HTML for security
```

### 7.2: Update README.md

Add to a "## Architecture" section:

```markdown
### Frontend Architecture

The web interface uses **htmx** for dynamic interactions, providing a modern user experience without complex JavaScript frameworks. The architecture follows a server-driven model:

- **Forms** use htmx attributes (`hx-post`, `hx-target`) to submit asynchronously
- **Server** returns HTML fragments (not JSON) rendered via Go templates
- **htmx** swaps the HTML into the page automatically
- **JavaScript** is minimal (~50 lines) and only handles drag-and-drop file selection

This approach:
- Eliminates ~85% of JavaScript code
- Keeps UI logic server-side in Go
- Requires no build step or npm dependencies
- Provides progressive enhancement (works without JS for basic functionality)
```

### 7.3: Document API Compatibility

Add to README.md in the "## HTTP Server" section:

```markdown
### API Compatibility

The `/api/stats` endpoint supports both JSON and HTML responses, making it suitable for both programmatic access and web browser interactions:

**JSON Response (API clients):**
- Default response format when no Accept header is specified
- Set `Accept: application/json` header for explicit JSON response
- Returns structured data with `overall`, `byLevel`, `byPlayer` fields
- Use for scripts, automation, mobile apps, or any programmatic access

**HTML Response (htmx/browsers):**
- Set `Accept: text/html` header (htmx does this automatically)
- Returns HTML fragment ready for insertion into the page
- Used by the web interface for dynamic content updates

**Examples:**

```bash
# JSON for API clients (default)
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json"

# JSON with explicit Accept header
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: application/json"

# HTML for htmx (automatic when using web form)
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: text/html"

# With filters (works for both JSON and HTML)
curl -X POST "http://localhost:8080/api/stats?level=1-5&detailed=true" \
  -F "file=@data/input.json"
```

This dual-format support ensures the API remains flexible and backward-compatible while enabling modern htmx-driven web interactions.
```

---

## Migration Complete!

### Summary of Changes

**Files Created:**
- `web/templates/results.html`
- `web/templates/overall-stats.html`
- `web/templates/level-stats.html`
- `web/templates/player-stats.html`
- `web/templates/error.html`
- `internal/templates/templates.go`
- `internal/handlers/template_data.go`

**Files Modified:**
- `web/static/index.html` (htmx attributes, simplified structure)
- `web/static/app.js` (reduced from 335 to ~50 lines)
- `web/static/styles.css` (added htmx transitions)
- `internal/handlers/stats.go` (dual JSON/HTML rendering with Accept header detection)
- `internal/pipeline/pipeline.go` (added ShowDetailed field)
- `cmd/http-server/main.go` (added /api/clear route)
- `docs/CLAUDE.md` (documentation updates)
- `README.md` (architecture documentation)

**Code Metrics:**
- JavaScript: 335 lines → 50 lines (85% reduction)
- Backend: +~150 lines (template package + rendering logic)
- Templates: +~150 lines (HTML templates)
- Net complexity: Significantly reduced (server-side logic is simpler than client-side DOM manipulation)

**Key Benefits:**
- ✅ **Backward compatible:** Existing JSON API clients continue to work without changes
- ✅ **Dual-purpose endpoint:** Same `/api/stats` serves both web UI and API consumers
- ✅ **Future-proof:** Easy to add mobile apps, CLI tools, or other API clients
- ✅ **Progressive enhancement:** Web interface uses modern htmx while maintaining API access
- ✅ **Simple implementation:** Only one `if/else` check based on Accept header

### Next Steps

1. Follow the phases in order
2. Test thoroughly after each phase
3. Commit after each major phase completes
4. Enjoy your simplified, server-driven architecture!

---

## Troubleshooting

### Common Issues

**Templates not found:**
- Ensure `web/templates/` directory exists
- Check that template manager is initialized with correct path
- Verify template file names match exactly (case-sensitive)

**htmx not swapping content:**
- Check browser console for htmx errors
- Verify `hx-target` selector matches element ID
- Ensure server returns `Content-Type: text/html`
- Check that form has `hx-encoding="multipart/form-data"` for file uploads

**Checkboxes not working:**
- Ensure checkboxes have `value="true"` or `value="levels"` etc.
- Verify backend parses checkbox values correctly
- Remember: unchecked checkboxes don't send values

**Loading indicator not showing:**
- Check `hx-indicator` ID matches spinner element ID
- Verify CSS for `.htmx-indicator` and `.htmx-request .htmx-indicator`
- Test with slow network (DevTools throttling)

**Getting JSON when expecting HTML (or vice versa):**
- Check the `Accept` header in browser DevTools Network tab
- Verify `wantsHTML()` function logic in handlers/stats.go
- htmx should send `HX-Request: true` header
- Test with curl to verify both formats work independently
- Check server logs to see which branch (HTML vs JSON) is executing

### Debugging Tips

Enable htmx event logging:
```javascript
htmx.logAll(); // Add to app.js during development
```

Check htmx requests in DevTools:
- Network tab → check request/response
- Verify `HX-Request: true` header is sent
- Verify HTML is returned (not JSON)

Template rendering errors:
- Check server logs for template parsing errors
- Ensure data structure matches template expectations
- Use `{{printf "%#v" .}}` in templates to debug data

Testing both JSON and HTML formats:
```bash
# Test JSON response (should see JSON structure)
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: application/json" | jq

# Test HTML response (should see HTML tags)
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: text/html"

# Compare both responses have same data
curl -X POST http://localhost:8080/api/stats \
  -F "file=@data/input.json" \
  -H "Accept: application/json" > json_response.txt

# Check server logs for format detection
# Should log which branch (HTML/JSON) was taken
```
