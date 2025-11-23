# HTTP Server Design Document

**Date:** 2025-11-23
**Status:** Design Complete
**Author:** Architecture Design Session

## Overview

This document describes the design for adding an HTTP server to the cli-stat-creator project. The server will provide a web interface for analyzing game score statistics, complementing the existing CLI tool.

## High-Level Architecture

The HTTP server will be implemented as a separate binary (`cmd/http-server`) that reuses all existing statistics calculation logic. The architecture maintains clean separation between:

- **Backend (Go + chi)**: JSON API that accepts file uploads, processes game scores through the existing pipeline, and returns statistics as JSON
- **Frontend (HTML/CSS/JS)**: Static files served by chi that provide a file upload form and dynamically render statistics tables

### Key Architectural Decisions

1. **Stateless requests**: Each HTTP request uploads a JSON file, processes it, and returns results. No persistent storage or session management.

2. **Code reuse**: The `internal/stats` and `internal/pipeline` packages remain unchanged. They're pure domain logic that works identically for both CLI and HTTP interfaces.

3. **Separation of concerns**: New packages handle HTTP-specific functionality without modifying existing business logic.

## API Design

### Primary Endpoint: POST /api/stats

Single endpoint that handles file uploads and returns statistics based on query parameters.

**Request Format:**
- Method: POST
- Content-Type: multipart/form-data
- File field name: "file"
- File must be JSON format

**Query Parameters (all optional):**
- `include`: Comma-separated list of sections to include in response
  - Values: `overall`, `levels`, `players`, or `all` (default: `all`)
  - Example: `?include=overall,levels`
- `level`: Filter by level (single or range)
  - Example: `?level=5` or `?level=1-5`
- `min-score`: Minimum score filter (integer)
- `max-score`: Maximum score filter (integer)
- `detailed`: Include min/max values in response (boolean, default: false)

**Response Format (JSON):**
```json
{
  "overall": {
    "totalGamesPlayed": 100,
    "totalScore": 8500,
    "averageScore": 85.0,
    "medianScore": 82.0,
    "minimumScore": 45,
    "maximumScore": 100
  },
  "byLevel": {
    "1": { /* Statistics object */ },
    "2": { /* Statistics object */ }
  },
  "byPlayer": {
    "player-key": { /* Statistics object */ }
  }
}
```

Fields are omitted based on `include` parameter.

**Error Response Format:**
```json
{
  "error": "description of what went wrong",
  "code": "ERROR_CODE"
}
```

### Static File Endpoint: GET /

Serves static HTML/CSS/JS files from `web/static/` directory. The root path serves `index.html`.

## Project Structure

```
.
├── cmd/
│   ├── cli-stat-creator/    (existing - unchanged)
│   └── http-server/          (new)
│       └── main.go
├── internal/
│   ├── stats/                (existing - unchanged)
│   ├── reader/               (existing - unchanged)
│   ├── pipeline/             (existing - unchanged)
│   ├── logging/              (existing - unchanged)
│   ├── display/              (existing - unchanged)
│   ├── handlers/             (new)
│   │   └── stats.go
│   └── render/               (new)
│       └── json.go
└── web/
    └── static/               (new)
        ├── index.html
        ├── styles.css
        └── app.js
```

## Backend Components

### cmd/http-server/main.go

Entry point for the HTTP server binary.

**Responsibilities:**
- Parse server configuration from CLI flags
- Set up structured logging with slog
- Initialize chi router with middleware
- Register routes and handlers
- Start HTTP server with graceful shutdown

**Configuration Flags:**
- `--port` or `-p`: Server port (default: 8080)
- `--host`: Bind address (default: localhost)
- `--log-level` or `-l`: Logging level (debug, info, warn, error; default: info)
- `--static-dir`: Path to static files directory (default: ./web/static)

**Chi Router Setup:**

Middleware stack:
- `chi.Logger`: Request logging middleware
- `chi.Recoverer`: Panic recovery middleware
- Custom logging middleware to add request ID to context
- CORS middleware (optional, for local development)

**Routes:**
- `POST /api/stats` → handlers.HandleStats
- `GET /*` → chi.FileServer (serves static files)

**Graceful Shutdown:**

The server listens for interrupt signals (SIGINT, SIGTERM) and gracefully shuts down, waiting for active requests to complete with a timeout (30 seconds).

**Function Signatures:**
```go
type Config struct {
    Port      int
    Host      string
    LogLevel  string
    StaticDir string
}

func setupRouter(logger *slog.Logger, staticDir string) *chi.Mux
func parseConfig() Config
```

### internal/handlers/stats.go

Chi HTTP handlers for the API endpoint.

**Responsibilities:**
- Parse multipart form data to extract uploaded file
- Parse query parameters into `pipeline.Config` structure
- Stream uploaded file to pipeline for processing
- Convert pipeline results to JSON via render package
- Handle errors with appropriate HTTP status codes

**Function Signatures:**
```go
func HandleStats(w http.ResponseWriter, r *http.Request)
func parseQueryParams(r *http.Request) (pipeline.Config, error)
func parseIncludeParam(include string) (overall, levels, players bool)
```

### internal/render/json.go

JSON response formatting and error handling.

**Responsibilities:**
- Marshal statistics results to JSON
- Format error responses with consistent structure
- Set appropriate HTTP headers and status codes
- Filter results based on include parameter

**Function Signatures:**
```go
func JSON(w http.ResponseWriter, status int, data interface{})
func Error(w http.ResponseWriter, status int, message string, code string)
func filterResults(results pipeline.Results, overall, levels, players bool) interface{}
```

## Frontend Components

### web/static/index.html

Single-page application with the following sections:
- **Header**: Title and brief description
- **Upload Form**: File input, filter controls, and submit button
- **Results Display**: Three collapsible/expandable sections for overall, per-level, and per-player statistics

**Form Controls:**
- File input (accept .json files)
- Level filter (text input for single level or range)
- Min/max score filters (number inputs)
- Checkboxes for: "Show detailed stats", "Include levels", "Include players"
- Submit button

### web/static/styles.css

Simple, clean styling with:
- Responsive layout (works on mobile and desktop)
- Table styling for statistics display
- Form styling with clear labels and inputs
- Loading state indicators
- Error message styling

### web/static/app.js

JavaScript handles all client-side logic.

**Core Functionality:**
```javascript
function handleFormSubmit(event)
function fetchStatistics(formData, queryParams)
function renderResults(data)
function renderStatisticsTable(stats, container)
function renderLevelTable(levelStats)
function renderPlayerTable(playerStats)
function showError(message)
function clearResults()
```

**Table Rendering:**

JavaScript creates HTML tables dynamically with columns matching the CLI output (Total Games, Total Score, Average, Median, etc.). When detailed mode is enabled, Min/Max columns are included.

## Error Handling

### Error Categories

**Client Errors (4xx):**
- 400 Bad Request: Invalid file format, malformed JSON, invalid query parameters, missing required fields
- 413 Payload Too Large: Uploaded file exceeds size limit (configurable, default: 10MB)
- 415 Unsupported Media Type: Non-JSON file uploaded
- 422 Unprocessable Entity: Valid JSON but fails GameScore validation (negative scores, empty player names, etc.)

**Server Errors (5xx):**
- 500 Internal Server Error: Unexpected errors during processing, pipeline failures

### Validation Flow

Request validation happens in this order:
1. **File Upload Validation** (handlers): Check file exists, size limit, content type
2. **Query Parameter Validation** (handlers): Parse and validate level ranges, min/max scores
3. **JSON Structure Validation** (reader): Unmarshal JSON, check structure
4. **GameScore Validation** (stats): Existing Validate() method checks business rules
5. **Pipeline Execution** (pipeline): Process data, catch calculation errors

### Error Response Structure

All errors return consistent JSON format with descriptive messages and error codes for client-side handling.

**Error Codes:**
- `INVALID_FILE`: File upload issues
- `INVALID_JSON`: JSON parsing failures
- `VALIDATION_ERROR`: GameScore validation failures
- `INVALID_PARAMETERS`: Query parameter issues
- `PROCESSING_ERROR`: Pipeline execution errors

## Logging Strategy

### Request-Scoped Logging

The HTTP server leverages the existing `internal/logging` package for structured logging with request-scoped context.

**Logging Flow:**

1. **Request Start**: Middleware generates unique request ID, creates logger with request metadata (method, path, request ID), adds to context via `logging.WithLogger()`
2. **Handler Execution**: Handlers extract logger via `logging.FromContext()`, log important events
3. **Pipeline Processing**: Pipeline already uses context-based logging, automatically inherits request-scoped logger
4. **Request Complete**: Middleware logs response status, duration, bytes sent

**Log Fields:**
- `request_id`: Unique identifier for request correlation
- `method`: HTTP method (POST, GET)
- `path`: Request path
- `remote_addr`: Client IP address
- `user_agent`: Client user agent (optional)
- `duration_ms`: Request processing time
- `status`: HTTP response status code
- `error`: Error message (if applicable)

**Log Levels:**
- **Debug**: Query parameter parsing, file size, pipeline configuration
- **Info**: Request started, request completed, statistics calculated
- **Warn**: Validation failures, client errors (4xx)
- **Error**: Server errors (5xx), unexpected failures

## Testing Strategy

### Unit Tests

**internal/handlers/stats_test.go:**
- Test query parameter parsing with valid/invalid inputs
- Test include parameter parsing (all, combinations, invalid values)
- Test error handling for missing files, invalid formats
- Use `httptest.NewRequest()` and `httptest.NewRecorder()`

**internal/render/json_test.go:**
- Test JSON marshaling of Statistics, pipeline.Results
- Test error response formatting
- Test result filtering based on include flags
- Verify correct HTTP status codes and headers

### Integration Tests

**cmd/http-server/server_test.go:**
- Start test server with `httptest.NewServer()`
- Test complete request flow: upload file → process → receive JSON response
- Test with sample data from `data/input.json`
- Test all query parameter combinations
- Test error scenarios (invalid JSON, validation failures)
- Verify response matches expected statistics

### Frontend Testing

Manual testing approach for JavaScript frontend:
- Test file upload with valid/invalid files
- Test all filter combinations
- Test error display
- Test table rendering with different data sets
- Verify responsive behavior on different screen sizes

### Existing Tests

All existing tests in `internal/stats`, `internal/pipeline`, `internal/reader` remain unchanged and continue to pass.

## Build and Deployment

### Build Commands

```bash
# Build CLI (existing)
go build -o cli-stat-creator ./cmd/cli-stat-creator

# Build HTTP server (new)
go build -o http-server ./cmd/http-server

# Build both
go build -o cli-stat-creator ./cmd/cli-stat-creator && \
go build -o http-server ./cmd/http-server
```

### Running the Server

```bash
# Start with defaults (localhost:8080, info logging)
./http-server

# Start with custom configuration
./http-server --port 3000 --host 0.0.0.0 --log-level debug
```

The server will log:
- Server starting message with host:port
- Static files directory location
- Ready to accept requests

### Dependencies

New dependencies to add via `go get`:
- `github.com/go-chi/chi/v5` - HTTP router
- `github.com/go-chi/cors` - CORS middleware (optional, for development)

### Documentation Updates

**CLAUDE.md:**
- Add http-server build command
- Add http-server run command and flags
- Document new internal packages (handlers, render)
- Add API endpoint documentation

**README.md:**
- Add HTTP Server section explaining web interface
- Document server startup and configuration
- Add example curl commands for API testing
- Add screenshot/description of web interface (after implementation)
- Keep existing CLI documentation intact

## Implementation Phases

### Phase 1: Design Documentation
- Write validated design document ✓
- Commit design document to git

### Phase 2: Backend Foundation
- Add chi dependency
- Create `internal/render/json.go` with JSON response helpers
- Create `internal/handlers/stats.go` with HTTP handlers
- Create `cmd/http-server/main.go` with server setup

### Phase 3: Frontend
- Create `web/static/index.html` with upload form
- Create `web/static/styles.css` for styling
- Create `web/static/app.js` for API interaction

### Phase 4: Testing
- Write unit tests for handlers package
- Write unit tests for render package
- Write integration tests for API
- Manual testing of web interface

### Phase 5: Documentation
- Update CLAUDE.md with http-server commands
- Update README.md with HTTP server section
- Update go.mod documentation if needed

## Summary

This design adds an HTTP server to the cli-stat-creator project while maintaining complete backward compatibility with the existing CLI. The architecture leverages the existing, well-tested business logic and adds a thin HTTP layer with a simple web interface. Both the CLI and HTTP server can coexist, providing users flexibility in how they interact with the statistics calculator.
