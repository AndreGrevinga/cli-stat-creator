// Package logging provides context-based logger management for structured logging.
// It uses Go's log/slog library and propagates loggers through context.Context.
//
// The package provides two main functions:
//   - WithLogger: Adds a logger to a context
//   - FromContext: Retrieves a logger from a context, returning a no-op logger if not present
//
// This design ensures graceful degradation when logger is not available in context.
// The no-op logger discards all log messages, preventing nil pointer panics while
// maintaining the same interface as a configured logger.
package logging

import (
	"context"
	"io"
	"log/slog"
)

type contextKey string

const loggerKey contextKey = "logger"

// FromContext extracts logger from context, returns no-op logger if not found
func FromContext(ctx context.Context) *slog.Logger {
	if logger, ok := ctx.Value(loggerKey).(*slog.Logger); ok {
		return logger
	}
	return slog.New(slog.NewTextHandler(io.Discard, nil))
}

// WithLogger adds logger to context
func WithLogger(ctx context.Context, logger *slog.Logger) context.Context {
	return context.WithValue(ctx, loggerKey, logger)
}
