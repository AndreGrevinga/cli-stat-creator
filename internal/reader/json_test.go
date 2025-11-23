package reader

import (
	"cli-stat-creator/internal/stats"
	"context"
	"os"
	"path/filepath"
	"reflect"
	"testing"
)

func TestReadScoresFromFile(t *testing.T) {
	tests := []struct {
		name      string
		content   string
		want      []stats.GameScore
		wantError bool
	}{
		{
			name: "valid JSON file",
			content: `[
				{
					"player": {"name": "alice", "id": 1},
					"score": 100,
					"level": 1
				},
				{
					"player": {"name": "bob", "id": 2},
					"score": 200,
					"level": 2
				}
			]`,
			want: []stats.GameScore{
				{Player: stats.Player{Name: "alice", ID: 1}, Score: 100, Level: 1},
				{Player: stats.Player{Name: "bob", ID: 2}, Score: 200, Level: 2},
			},
			wantError: false,
		},
		{
			name:      "empty JSON array",
			content:   `[]`,
			want:      []stats.GameScore{},
			wantError: false,
		},
		{
			name: "single game score",
			content: `[
				{
					"player": {"name": "alice", "id": 1},
					"score": 100,
					"level": 1
				}
			]`,
			want: []stats.GameScore{
				{Player: stats.Player{Name: "alice", ID: 1}, Score: 100, Level: 1},
			},
			wantError: false,
		},
		{
			name:      "invalid JSON format",
			content:   `{"invalid": "json"`,
			want:      nil,
			wantError: true,
		},
		{
			name:      "malformed JSON array",
			content:   `[{"player": {"name": "alice", "id": 1}, "score": 100,]`,
			want:      nil,
			wantError: true,
		},
		{
			name: "missing required fields",
			content: `[
				{
					"player": {"name": "alice"},
					"score": 100
				}
			]`,
			want: []stats.GameScore{
				{Player: stats.Player{Name: "alice", ID: 0}, Score: 100, Level: 0},
			},
			wantError: true, // JSON unmarshaling uses zero values for missing fields
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			// Create temporary file with test content
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, "test_scores.json")

			err := os.WriteFile(tmpFile, []byte(test.content), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			// Test the function
			result, err := ReadScoresFromFile(context.Background(), tmpFile)

			// Check error expectation
			if (err != nil) != test.wantError {
				t.Errorf("ReadScoresFromFile() error = %v, wantError %v", err, test.wantError)
				return
			}

			// If we didn't expect an error, verify the result
			if !test.wantError && !reflect.DeepEqual(result, test.want) {
				t.Errorf("ReadScoresFromFile() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestReadScoresFromFile_NonExistentFile(t *testing.T) {
	tmpDir := t.TempDir()
	nonExistentFile := filepath.Join(tmpDir, "does_not_exist.json")

	_, err := ReadScoresFromFile(context.Background(), nonExistentFile)

	if err == nil {
		t.Error("ReadScoresFromFile() expected error for non-existent file, got nil")
	}
}

func TestReadScoresFromFile_FileExtension(t *testing.T) {
	tests := []struct {
		name      string
		filename  string
		wantError bool
	}{
		{
			name:      "lowercase .json",
			filename:  "test.json",
			wantError: false,
		},
		{
			name:      "uppercase .JSON",
			filename:  "test.JSON",
			wantError: false,
		},
		{
			name:      "mixed case .Json",
			filename:  "test.Json",
			wantError: false,
		},
		{
			name:      "wrong extension .txt",
			filename:  "test.txt",
			wantError: true,
		},
		{
			name:      "no extension",
			filename:  "test",
			wantError: true,
		},
		{
			name:      "json in name but wrong extension",
			filename:  "json.txt",
			wantError: true,
		},
	}

	validContent := `[{"player": {"name": "alice", "id": 1}, "score": 100, "level": 1}]`

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			tmpDir := t.TempDir()
			tmpFile := filepath.Join(tmpDir, test.filename)

			err := os.WriteFile(tmpFile, []byte(validContent), 0644)
			if err != nil {
				t.Fatalf("Failed to create test file: %v", err)
			}

			_, err = ReadScoresFromFile(context.Background(), tmpFile)

			if (err != nil) != test.wantError {
				t.Errorf("ReadScoresFromFile(%s) error = %v, wantError %v", test.filename, err, test.wantError)
			}
		})
	}
}

func TestReadScoresFromFile_ComplexData(t *testing.T) {
	tmpDir := t.TempDir()
	tmpFile := filepath.Join(tmpDir, "complex_scores.json")

	content := `[
		{
			"player": {"name": "alice", "id": 1},
			"score": 100,
			"level": 1
		},
		{
			"player": {"name": "bob", "id": 2},
			"score": 200,
			"level": 2
		},
		{
			"player": {"name": "charlie", "id": 3},
			"score": 150,
			"level": 1
		},
		{
			"player": {"name": "alice", "id": 1},
			"score": 175,
			"level": 3
		}
	]`

	err := os.WriteFile(tmpFile, []byte(content), 0644)
	if err != nil {
		t.Fatalf("Failed to create test file: %v", err)
	}

	result, err := ReadScoresFromFile(context.Background(), tmpFile)

	if err != nil {
		t.Fatalf("ReadScoresFromFile() unexpected error: %v", err)
	}

	if len(result) != 4 {
		t.Errorf("ReadScoresFromFile() returned %d scores, want 4", len(result))
	}

	// Verify first and last entries
	expectedFirst := stats.GameScore{Player: stats.Player{Name: "alice", ID: 1}, Score: 100, Level: 1}
	if !reflect.DeepEqual(result[0], expectedFirst) {
		t.Errorf("ReadScoresFromFile() first entry = %v, want %v", result[0], expectedFirst)
	}

	expectedLast := stats.GameScore{Player: stats.Player{Name: "alice", ID: 1}, Score: 175, Level: 3}
	if !reflect.DeepEqual(result[3], expectedLast) {
		t.Errorf("ReadScoresFromFile() last entry = %v, want %v", result[3], expectedLast)
	}
}
