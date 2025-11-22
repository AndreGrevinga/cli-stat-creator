package main

import (
	"testing"
)

func TestParseLevelFlag_EmptyString(t *testing.T) {
	result, err := parseLevelFlag("")
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
			result, err := parseLevelFlag(tt.input)
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
			result, err := parseLevelFlag(tt.input)
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
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := parseLevelFlag(tt.input)
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
			result, err := parseLevelFlag(tt.input)
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
	for i := 0; i < 100; i++ {
		levels[i] = i + 1
	}
	return levels
}
