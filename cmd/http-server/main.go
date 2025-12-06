// Package main implements an HTTP server for the game score statistics API.
// It provides REST endpoints for uploading game score data and retrieving
// statistical analysis including overall, per-level, and per-player statistics.
package main

import (
	"cli-stat-creator/internal/handlers"
	"cli-stat-creator/internal/logging"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

// Config holds the server configuration including network settings,
// logging level, and static file directory path.
type Config struct {
	Port      int
	Host      string
	LogLevel  string
	StaticDir string
}

// loggerMiddleware creates a middleware that attaches the logger to the request context.
// This allows handlers to retrieve the logger using logging.FromContext(ctx).
func loggerMiddleware(logger *slog.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := logging.WithLogger(r.Context(), logger)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// setupRouter creates and configures a chi router with middleware and routes.
// It sets up CORS, logging, recovery, and mounts the API endpoints and static file server.
func setupRouter(logger *slog.Logger, staticDir string) *chi.Mux {
	r := chi.NewRouter()
	r.Use(middleware.RequestID)
	r.Use(loggerMiddleware(logger)) // Attach logger to context
	r.Use(middleware.Logger)
	r.Use(cors.Handler(cors.Options{
		// AllowedOrigins:   []string{"https://foo.com"}, // Use this to allow specific origin hosts
		AllowedOrigins: []string{"https://*", "http://*"},
		// AllowOriginFunc:  func(r *http.Request, origin string) bool { return true },
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	r.Use(middleware.Recoverer)
	r.Use(middleware.URLFormat)

	// API routes
	r.Post("/api/stats", handlers.HandleStats)
	r.Get("/api/clear", handlers.HandleClear)

	// Serve static files
	fileServer := http.FileServer(http.Dir(staticDir))
	r.Get("/*", func(w http.ResponseWriter, r *http.Request) {
		fileServer.ServeHTTP(w, r)
	})

	return r
}

// parseConfig parses command-line flags into a Config struct.
// It validates the port flag and uses sensible defaults for invalid values.
func parseConfig() Config {
	port, err := strconv.ParseInt(portFlag, 10, 0)
	if err != nil {
		slog.Warn("Invalid port flag, using default 8080", "value", portFlag, "error", err)
		port = 8080
	}
	return Config{
		Port:      int(port),
		Host:      *hostFlag,
		LogLevel:  logLevelFlag,
		StaticDir: *staticDirFlag,
	}
}

var (
	portFlag      string
	hostFlag      = flag.String("host", "localhost", "Bind address")
	logLevelFlag  string
	staticDirFlag = flag.String("static-dir", "web/static", "Path to static files directory")
)

func main() {
	flag.StringVar(&portFlag, "port", "8080", "Server port")
	flag.StringVar(&portFlag, "p", "8080", "Server port (shorthand)")
	flag.StringVar(&logLevelFlag, "log-level", "info", "Logging level")
	flag.StringVar(&logLevelFlag, "l", "info", "Logging level (shorthand)")
	flag.Parse()
	config := parseConfig()
	router := setupRouter(slog.Default(), config.StaticDir)
	URL := fmt.Sprintf("%s:%d", config.Host, config.Port)
	slog.Info("Starting server", "url", URL, "static_dir", config.StaticDir)
	err := http.ListenAndServe(URL, router)
	if err != nil {
		slog.Error("Failed to start server", "error", err)
	}
}
