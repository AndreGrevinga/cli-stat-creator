package main

import (
	"cli-stat-creator/internal/pipeline"
	"log/slog"
	"testing"
)

func TestParseLevelFlag_EmptyString(t *testing.T) {
	result, err := pipeline.ParseLevelString("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(result) != 0 {
		t.Errorf("expected empty slice, got %v", result)
	}
}

func TestParseLevelFlag_SingleLevel(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"level 1", "1", []int{1}},
		{"level 5", "5", []int{5}},
		{"level 10", "10", []int{10}},
		{"level 100", "100", []int{100}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pipeline.ParseLevelString(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, tt.expected[i], val)
				}
			}
		})
	}
}

func TestParseLevelFlag_Range(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []int
	}{
		{"range 1-3", "1-3", []int{1, 2, 3}},
		{"range 5-8", "5-8", []int{5, 6, 7, 8}},
		{"range 1-1", "1-1", []int{1}},
		{"range 10-15", "10-15", []int{10, 11, 12, 13, 14, 15}},
		{"range 20-22", "20-22", []int{20, 21, 22}},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pipeline.ParseLevelString(tt.input)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, tt.expected[i], val)
				}
			}
		})
	}
}

func TestParseLevelFlag_InvalidInput(t *testing.T) {
	tests := []struct {
		name  string
		input string
	}{
		{"non-numeric", "abc"},
		{"invalid range start", "abc-5"},
		{"invalid range end", "1-xyz"},
		{"space in number", "1 5"},
		{"negative number", "-5"},
		{"negative range", "-5--3"},
		{"multiple dashes", "1-5-10"},
		{"reverse range", "10-5"},
		{"too large range", "1-1002"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := pipeline.ParseLevelString(tt.input)
			if err == nil {
				t.Error("expected error, got nil")
			}
		})
	}
}

func TestParseLevelFlag_EdgeCases(t *testing.T) {
	tests := []struct {
		name        string
		input       string
		expectError bool
		expected    []int
	}{
		{"zero level", "0", false, []int{0}},
		{"range with zero", "0-2", false, []int{0, 1, 2}},
		{"large range", "1-100", false, make100Levels()},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := pipeline.ParseLevelString(tt.input)
			if tt.expectError {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if len(result) != len(tt.expected) {
				t.Fatalf("expected length %d, got %d", len(tt.expected), len(result))
			}
			for i, val := range result {
				if val != tt.expected[i] {
					t.Errorf("at index %d: expected %d, got %d", i, tt.expected[i], val)
				}
			}
		})
	}
}

// Helper function to generate a slice of 100 levels (1-100)
func make100Levels() []int {
	levels := make([]int, 100)
	for i := range 100 {
		levels[i] = i + 1
	}
	return levels
}

// TestParseLogLevel_ValidLevels tests that all valid log level strings
// are correctly parsed to their corresponding slog.Level values.
func TestParseLogLevel_ValidLevels(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		envValue string
		expected slog.Level
	}{
		{"debug from flag", "debug", "", slog.LevelDebug},
		{"info from flag", "info", "", slog.LevelInfo},
		{"warn from flag", "warn", "", slog.LevelWarn},
		{"error from flag", "error", "", slog.LevelError},
		{"DEBUG uppercase from flag", "DEBUG", "", slog.LevelDebug},
		{"Info mixed case from flag", "Info", "", slog.LevelInfo},
		{"WARN uppercase from flag", "WARN", "", slog.LevelWarn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.flag, tt.envValue)
			if result.Level() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Level())
			}
		})
	}
}

// TestParseLogLevel_InvalidLevels tests that invalid log level strings
// default to WARN level.
func TestParseLogLevel_InvalidLevels(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		envValue string
	}{
		{"invalid flag", "invalid", ""},
		{"random text", "xyz", ""},
		{"numeric value", "123", ""},
		{"empty both", "", ""},
		{"whitespace only", "  ", ""},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.flag, tt.envValue)
			if result.Level() != slog.LevelWarn {
				t.Errorf("expected WARN for invalid input, got %v", result.Level())
			}
		})
	}
}

// TestParseLogLevel_FlagPrecedence tests that the flag parameter
// takes precedence over the environment variable.
func TestParseLogLevel_FlagPrecedence(t *testing.T) {
	tests := []struct {
		name     string
		flag     string
		envValue string
		expected slog.Level
	}{
		{"flag overrides env", "debug", "error", slog.LevelDebug},
		{"flag debug, env info", "debug", "info", slog.LevelDebug},
		{"flag error, env debug", "error", "debug", slog.LevelError},
		{"empty flag uses env", "", "info", slog.LevelInfo},
		{"empty flag with invalid env", "", "invalid", slog.LevelWarn},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.flag, tt.envValue)
			if result.Level() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Level())
			}
		})
	}
}

// TestParseLogLevel_EnvironmentVariable tests that the environment variable
// is used when the flag is empty.
func TestParseLogLevel_EnvironmentVariable(t *testing.T) {
	tests := []struct {
		name     string
		envValue string
		expected slog.Level
	}{
		{"env debug", "debug", slog.LevelDebug},
		{"env info", "info", slog.LevelInfo},
		{"env warn", "warn", slog.LevelWarn},
		{"env error", "error", slog.LevelError},
		{"env with spaces", "  info  ", slog.LevelInfo},
		{"env uppercase", "ERROR", slog.LevelError},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel("", tt.envValue)
			if result.Level() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Level())
			}
		})
	}
}

// TestParseLogLevel_CaseInsensitive tests that log level parsing
// is case-insensitive.
func TestParseLogLevel_CaseInsensitive(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"lowercase debug", "debug", slog.LevelDebug},
		{"uppercase DEBUG", "DEBUG", slog.LevelDebug},
		{"mixed case Debug", "Debug", slog.LevelDebug},
		{"mixed case DeBuG", "DeBuG", slog.LevelDebug},
		{"lowercase info", "info", slog.LevelInfo},
		{"uppercase INFO", "INFO", slog.LevelInfo},
		{"mixed case InFo", "InFo", slog.LevelInfo},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input, "")
			if result.Level() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Level())
			}
		})
	}
}

// TestParseLogLevel_Whitespace tests that leading and trailing whitespace
// is properly trimmed before parsing.
func TestParseLogLevel_Whitespace(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected slog.Level
	}{
		{"leading space", " debug", slog.LevelDebug},
		{"trailing space", "info ", slog.LevelInfo},
		{"both spaces", " warn ", slog.LevelWarn},
		{"multiple spaces", "  error  ", slog.LevelError},
		{"tab characters", "\tdebug\t", slog.LevelDebug},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseLogLevel(tt.input, "")
			if result.Level() != tt.expected {
				t.Errorf("expected %v, got %v", tt.expected, result.Level())
			}
		})
	}
}
