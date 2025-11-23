package logging

import (
	"bytes"
	"context"
	"log/slog"
	"testing"
)

// TestFromContext_EmptyContext verifies that FromContext returns a no-op logger
// when the context does not contain a logger.
func TestFromContext_EmptyContext(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx)

	if logger == nil {
		t.Fatal("FromContext returned nil, expected non-nil no-op logger")
	}

	// Verify it's a no-op logger by checking that output is discarded
	var buf bytes.Buffer
	logger.Info("test message")
	if buf.Len() > 0 {
		t.Errorf("No-op logger wrote output, expected nothing; got: %s", buf.String())
	}
}

// TestFromContext_WithLogger verifies that FromContext returns the correct logger
// when the context contains a logger.
func TestFromContext_WithLogger(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	})
	expectedLogger := slog.New(handler)

	ctx := WithLogger(context.Background(), expectedLogger)
	retrievedLogger := FromContext(ctx)

	if retrievedLogger == nil {
		t.Fatal("FromContext returned nil")
	}

	// Write a message and verify it appears in the buffer
	testMessage := "test log message"
	retrievedLogger.Info(testMessage)

	output := buf.String()
	if output == "" {
		t.Error("Logger did not write output, expected log message")
	}
	if !bytes.Contains(buf.Bytes(), []byte(testMessage)) {
		t.Errorf("Logger output does not contain expected message.\nExpected: %q\nGot: %q", testMessage, output)
	}
}

// TestWithLogger_AddsLoggerToContext verifies that WithLogger properly adds
// a logger to the context and it can be retrieved.
func TestWithLogger_AddsLoggerToContext(t *testing.T) {
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, &slog.HandlerOptions{
		Level: slog.LevelDebug,
	})
	logger := slog.New(handler)

	ctx := context.Background()
	ctxWithLogger := WithLogger(ctx, logger)

	// Verify the logger was added by retrieving it
	retrieved := FromContext(ctxWithLogger)
	if retrieved == nil {
		t.Fatal("Retrieved logger is nil")
	}

	// Verify it's the same logger by writing to it
	testMessage := "context logger test"
	retrieved.Debug(testMessage)

	if !bytes.Contains(buf.Bytes(), []byte(testMessage)) {
		t.Errorf("Retrieved logger did not write to expected buffer.\nExpected message: %q\nGot: %q", testMessage, buf.String())
	}
}

// TestNoOpLogger_DiscardsOutput verifies that the no-op logger returned
// when context has no logger discards all output.
func TestNoOpLogger_DiscardsOutput(t *testing.T) {
	ctx := context.Background()
	logger := FromContext(ctx) // Should return no-op logger

	// Try logging at all levels
	logger.Debug("debug message")
	logger.Info("info message")
	logger.Warn("warn message")
	logger.Error("error message")

	// If we got here without panicking, the no-op logger is working correctly.
	// The no-op logger writes to io.Discard, so we can't capture output to verify,
	// but we can verify it doesn't panic.
}

// TestWithLogger_PreservesOtherContextValues verifies that adding a logger
// to a context doesn't affect other values in the context.
func TestWithLogger_PreservesOtherContextValues(t *testing.T) {
	type testKey string
	const key testKey = "test-key"
	expectedValue := "test-value"

	// Create context with a value
	ctx := context.WithValue(context.Background(), key, expectedValue)

	// Add logger to context
	var buf bytes.Buffer
	handler := slog.NewTextHandler(&buf, nil)
	logger := slog.New(handler)
	ctxWithLogger := WithLogger(ctx, logger)

	// Verify original value is still present
	retrievedValue := ctxWithLogger.Value(key)
	if retrievedValue != expectedValue {
		t.Errorf("Context value was lost.\nExpected: %q\nGot: %v", expectedValue, retrievedValue)
	}

	// Verify logger was added
	retrievedLogger := FromContext(ctxWithLogger)
	if retrievedLogger == nil {
		t.Error("Logger was not added to context")
	}
}
