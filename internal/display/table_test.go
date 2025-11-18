package display

import (
	"bytes"
	"cli-stat-creator/internal/stats"
	"io"
	"os"
	"strings"
	"testing"
)

// captureOutput captures stdout during function execution
func captureOutput(f func()) string {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	f()

	w.Close()
	os.Stdout = old

	var buf bytes.Buffer
	io.Copy(&buf, r)
	return buf.String()
}

func TestRenderStatistics(t *testing.T) {
	tests := []struct {
		name       string
		statistics stats.Statistics
		wantFields []string
	}{
		{
			name: "renders all statistics fields",
			statistics: stats.Statistics{
				TotalGamesPlayed: 10,
				TotalScore:       1000,
				AverageScore:     100.50,
				MedianScore:      95.75,
				MinimumScore:     50,
				MaximumScore:     150,
			},
			wantFields: []string{"10", "1000", "100.50", "95.75", "50", "150"},
		},
		{
			name: "renders zero statistics",
			statistics: stats.Statistics{
				TotalGamesPlayed: 0,
				TotalScore:       0,
				AverageScore:     0.00,
				MedianScore:      0.00,
				MinimumScore:     0,
				MaximumScore:     0,
			},
			wantFields: []string{"0", "0", "0.00", "0.00", "0", "0"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(func() {
				RenderStatistics(test.statistics)
			})

			// Verify all expected fields are present in output
			for _, field := range test.wantFields {
				if !strings.Contains(output, field) {
					t.Errorf("RenderStatistics() output missing field %q\nGot: %s", field, output)
				}
			}

			// Verify table headers are present
			if !strings.Contains(output, "TOTAL GAMES PLAYED") {
				t.Error("RenderStatistics() missing header 'TOTAL GAMES PLAYED'")
			}
			if !strings.Contains(output, "AVERAGE SCORE") {
				t.Error("RenderStatistics() missing header 'AVERAGE SCORE'")
			}
		})
	}
}

func TestRenderLevelStatistics(t *testing.T) {
	tests := []struct {
		name       string
		statistics map[int]stats.Statistics
		detailed   bool
		wantFields []string
	}{
		{
			name: "renders single level summary mode",
			statistics: map[int]stats.Statistics{
				1: {
					TotalGamesPlayed: 5,
					TotalScore:       500,
					AverageScore:     100.00,
					MedianScore:      95.00,
					MinimumScore:     80,
					MaximumScore:     120,
				},
			},
			detailed:   false,
			wantFields: []string{"1", "5", "100.00", "120"},
		},
		{
			name: "renders multiple levels sorted numerically",
			statistics: map[int]stats.Statistics{
				3: {
					TotalGamesPlayed: 3,
					AverageScore:     150.00,
					MaximumScore:     180,
				},
				1: {
					TotalGamesPlayed: 5,
					AverageScore:     100.00,
					MaximumScore:     120,
				},
				2: {
					TotalGamesPlayed: 4,
					AverageScore:     125.00,
					MaximumScore:     150,
				},
			},
			detailed:   false,
			wantFields: []string{"1", "2", "3", "100.00", "125.00", "150.00"},
		},
		{
			name: "renders detailed mode with all columns",
			statistics: map[int]stats.Statistics{
				1: {
					TotalGamesPlayed: 5,
					TotalScore:       500,
					AverageScore:     100.00,
					MedianScore:      95.00,
					MinimumScore:     80,
					MaximumScore:     120,
				},
			},
			detailed:   true,
			wantFields: []string{"1", "5", "500", "100.00", "95.00", "80", "120"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(func() {
				RenderLevelStatistics(test.statistics, test.detailed)
			})

			// Verify all expected fields are present
			for _, field := range test.wantFields {
				if !strings.Contains(output, field) {
					t.Errorf("RenderLevelStatistics() output missing field %q\nGot: %s", field, output)
				}
			}

			// Verify header is present
			if !strings.Contains(output, "LEVEL") {
				t.Error("RenderLevelStatistics() missing header 'LEVEL'")
			}
		})
	}
}

func TestRenderPlayerStatistics(t *testing.T) {
	tests := []struct {
		name       string
		statistics map[stats.Player]stats.Statistics
		detailed   bool
		wantFields []string
	}{
		{
			name: "renders single player summary mode",
			statistics: map[stats.Player]stats.Statistics{
				{Name: "alice", ID: 1}: {
					TotalGamesPlayed: 5,
					TotalScore:       500,
					AverageScore:     100.00,
					MedianScore:      95.00,
					MinimumScore:     80,
					MaximumScore:     120,
				},
			},
			detailed:   false,
			wantFields: []string{"alice", "5", "100.00", "120"},
		},
		{
			name: "renders multiple players sorted alphabetically",
			statistics: map[stats.Player]stats.Statistics{
				{Name: "charlie", ID: 3}: {
					TotalGamesPlayed: 3,
					AverageScore:     150.00,
					MaximumScore:     180,
				},
				{Name: "alice", ID: 1}: {
					TotalGamesPlayed: 5,
					AverageScore:     100.00,
					MaximumScore:     120,
				},
				{Name: "bob", ID: 2}: {
					TotalGamesPlayed: 4,
					AverageScore:     125.00,
					MaximumScore:     150,
				},
			},
			detailed:   false,
			wantFields: []string{"alice", "bob", "charlie"},
		},
		{
			name: "renders detailed mode with all columns",
			statistics: map[stats.Player]stats.Statistics{
				{Name: "alice", ID: 1}: {
					TotalGamesPlayed: 5,
					TotalScore:       500,
					AverageScore:     100.00,
					MedianScore:      95.00,
					MinimumScore:     80,
					MaximumScore:     120,
				},
			},
			detailed:   true,
			wantFields: []string{"alice", "5", "500", "100.00", "95.00", "80", "120"},
		},
		{
			name: "handles players with same name but different IDs",
			statistics: map[stats.Player]stats.Statistics{
				{Name: "alex", ID: 1}: {
					TotalGamesPlayed: 3,
					AverageScore:     100.00,
					MaximumScore:     110,
				},
				{Name: "alex", ID: 2}: {
					TotalGamesPlayed: 2,
					AverageScore:     90.00,
					MaximumScore:     95,
				},
			},
			detailed:   false,
			wantFields: []string{"alex", "3", "2", "100.00", "90.00"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(func() {
				RenderPlayerStatistics(test.statistics, test.detailed)
			})

			// Verify all expected fields are present
			for _, field := range test.wantFields {
				if !strings.Contains(output, field) {
					t.Errorf("RenderPlayerStatistics() output missing field %q\nGot: %s", field, output)
				}
			}

			// Verify header is present
			if !strings.Contains(output, "PLAYER") {
				t.Error("RenderPlayerStatistics() missing header 'PLAYER'")
			}
		})
	}
}

func TestRenderGroupedStatistics(t *testing.T) {
	tests := []struct {
		name       string
		testFunc   func()
		wantFields []string
	}{
		{
			name: "generic function works with int keys",
			testFunc: func() {
				statistics := map[int]stats.Statistics{
					1: {
						TotalGamesPlayed: 5,
						AverageScore:     100.00,
						MaximumScore:     120,
					},
				}
				RenderGroupedStatistics(
					statistics,
					false,
					"TestKey",
					func(keys []int) {}, // no-op sort for single element
					func(key int) string { return "1" },
				)
			},
			wantFields: []string{"TEST KEY", "5", "100.00", "120"},
		},
		{
			name: "generic function works with Player keys",
			testFunc: func() {
				statistics := map[stats.Player]stats.Statistics{
					{Name: "test", ID: 1}: {
						TotalGamesPlayed: 3,
						AverageScore:     90.00,
						MaximumScore:     100,
					},
				}
				RenderGroupedStatistics(
					statistics,
					false,
					"TestPlayer",
					func(keys []stats.Player) {}, // no-op sort for single element
					func(key stats.Player) string { return key.Name },
				)
			},
			wantFields: []string{"TEST PLAYER", "test", "3", "90.00", "100"},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			output := captureOutput(test.testFunc)

			for _, field := range test.wantFields {
				if !strings.Contains(output, field) {
					t.Errorf("RenderGroupedStatistics() output missing field %q\nGot: %s", field, output)
				}
			}
		})
	}
}
