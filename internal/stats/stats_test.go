package stats

import (
	"reflect"
	"testing"
)

func TestGameScore_Validate(t *testing.T) {
	tests := []struct {
		name      string
		score     GameScore
		wantError bool
		errorMsg  string
	}{
		{
			name:      "valid game score",
			score:     GameScore{Player: Player{Name: "alice", ID: 1}, Score: 100, Level: 1},
			wantError: false,
		},
		{
			name:      "empty player name",
			score:     GameScore{Player: Player{Name: "", ID: 1}, Score: 100, Level: 1},
			wantError: true,
			errorMsg:  "player name cannot be empty",
		},
		{
			name:      "negative score",
			score:     GameScore{Player: Player{Name: "alice", ID: 1}, Score: -1, Level: 1},
			wantError: true,
			errorMsg:  "score cannot be negative",
		},
		{
			name:      "zero score is valid",
			score:     GameScore{Player: Player{Name: "alice", ID: 1}, Score: 0, Level: 1},
			wantError: false,
		},
		{
			name:      "zero level",
			score:     GameScore{Player: Player{Name: "alice", ID: 1}, Score: 100, Level: 0},
			wantError: true,
			errorMsg:  "level must be positive",
		},
		{
			name:      "negative level",
			score:     GameScore{Player: Player{Name: "alice", ID: 1}, Score: 100, Level: -1},
			wantError: true,
			errorMsg:  "level must be positive",
		},
		{
			name:      "zero player ID",
			score:     GameScore{Player: Player{Name: "alice", ID: 0}, Score: 100, Level: 1},
			wantError: true,
			errorMsg:  "player id must be positive",
		},
		{
			name:      "negative player ID",
			score:     GameScore{Player: Player{Name: "alice", ID: -1}, Score: 100, Level: 1},
			wantError: true,
			errorMsg:  "player id must be positive",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			err := test.score.Validate()

			if (err != nil) != test.wantError {
				t.Errorf("Validate() error = %v, wantError %v", err, test.wantError)
				return
			}

			if test.wantError && err.Error() != test.errorMsg {
				t.Errorf("Validate() error = %q, want %q", err.Error(), test.errorMsg)
			}
		})
	}
}

func TestGroupByPlayer(t *testing.T) {
	tests := []struct {
		name  string
		input []GameScore
		want  map[Player][]GameScore
	}{
		{
			name:  "empty input",
			input: []GameScore{},
			want:  map[Player][]GameScore{},
		},
		{
			name: "single player",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "alex", ID: 1}, Score: 200, Level: 2},
			},
			want: map[Player][]GameScore{
				{Name: "alex", ID: 1}: {
					{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
					{Player: Player{Name: "alex", ID: 1}, Score: 200, Level: 2},
				},
			},
		},
		{
			name: "multiple players",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "alex", ID: 1}, Score: 150, Level: 2},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 1},
			},
			want: map[Player][]GameScore{
				{Name: "alex", ID: 1}: {
					{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
					{Player: Player{Name: "alex", ID: 1}, Score: 150, Level: 2},
				},
				{Name: "bob", ID: 2}: {
					{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				},
				{Name: "charlie", ID: 3}: {
					{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 1},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GroupByPlayer(test.input)

			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("GroupByPlayer() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestGroupByLevel(t *testing.T) {
	tests := []struct {
		name  string
		input []GameScore
		want  map[int][]GameScore
	}{
		{
			name:  "empty input",
			input: []GameScore{},
			want:  map[int][]GameScore{},
		},
		{
			name: "single level",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
			},
			want: map[int][]GameScore{
				1: {
					{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
					{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				},
			},
		},
		{
			name: "multiple levels",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 2},
				{Player: Player{Name: "charlie", ID: 3}, Score: 150, Level: 1},
				{Player: Player{Name: "diana", ID: 4}, Score: 300, Level: 3},
			},
			want: map[int][]GameScore{
				1: {
					{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
					{Player: Player{Name: "charlie", ID: 3}, Score: 150, Level: 1},
				},
				2: {
					{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 2},
				},
				3: {
					{Player: Player{Name: "diana", ID: 4}, Score: 300, Level: 3},
				},
			},
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result := GroupByLevel(test.input)

			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("GroupByLevel() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestCalculateStatistics(t *testing.T) {
	tests := []struct {
		name           string
		input          []GameScore
		wantStatistics Statistics
		wantError      bool
	}{
		{
			name:           "empty input",
			input:          []GameScore{},
			wantStatistics: Statistics{},
			wantError:      true,
		},
		{
			name:  "single game score",
			input: []GameScore{{Player: Player{Name: "alex", ID: 1}, Score: 10, Level: 2}},
			wantStatistics: Statistics{
				TotalGamesPlayed: 1,
				TotalScore:       10,
				MinimumScore:     10,
				MaximumScore:     10,
				MedianScore:      10,
				AverageScore:     10},
			wantError: false,
		},
		{
			name: "odd number of scores - median calculation",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 2},
			},
			wantStatistics: Statistics{
				TotalGamesPlayed: 3,
				TotalScore:       600,
				MinimumScore:     100,
				MaximumScore:     300,
				MedianScore:      200,
				AverageScore:     200},
			wantError: false,
		},
		{
			name: "even number of scores - median calculation",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 2},
				{Player: Player{Name: "diana", ID: 4}, Score: 400, Level: 2},
			},
			wantStatistics: Statistics{
				TotalGamesPlayed: 4,
				TotalScore:       1000,
				MinimumScore:     100,
				MaximumScore:     400,
				MedianScore:      250, // (200 + 300) / 2
				AverageScore:     250},
			wantError: false,
		},
		{
			name: "duplicate scores",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 150, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 150, Level: 1},
				{Player: Player{Name: "charlie", ID: 3}, Score: 150, Level: 2},
			},
			wantStatistics: Statistics{
				TotalGamesPlayed: 3,
				TotalScore:       450,
				MinimumScore:     150,
				MaximumScore:     150,
				MedianScore:      150,
				AverageScore:     150},
			wantError: false,
		},
		{
			name: "unsorted input scores",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 500, Level: 3},
				{Player: Player{Name: "bob", ID: 2}, Score: 100, Level: 1},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 2},
				{Player: Player{Name: "diana", ID: 4}, Score: 200, Level: 1},
				{Player: Player{Name: "eve", ID: 5}, Score: 400, Level: 2},
			},
			wantStatistics: Statistics{
				TotalGamesPlayed: 5,
				TotalScore:       1500,
				MinimumScore:     100,
				MaximumScore:     500,
				MedianScore:      300, // sorted: [100, 200, 300, 400, 500]
				AverageScore:     300},
			wantError: false,
		},
		{
			name: "large dataset with mixed values",
			input: []GameScore{
				{Player: Player{Name: "p1", ID: 1}, Score: 50, Level: 1},
				{Player: Player{Name: "p2", ID: 2}, Score: 75, Level: 1},
				{Player: Player{Name: "p3", ID: 3}, Score: 100, Level: 2},
				{Player: Player{Name: "p4", ID: 4}, Score: 125, Level: 2},
				{Player: Player{Name: "p5", ID: 5}, Score: 150, Level: 3},
				{Player: Player{Name: "p6", ID: 6}, Score: 175, Level: 3},
				{Player: Player{Name: "p7", ID: 7}, Score: 200, Level: 4},
			},
			wantStatistics: Statistics{
				TotalGamesPlayed: 7,
				TotalScore:       875,
				MinimumScore:     50,
				MaximumScore:     200,
				MedianScore:      125, // sorted middle value
				AverageScore:     125},
			wantError: false,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CalculateStatistics(test.input)

			// Check if we got an error when we expected one
			if (err != nil) != test.wantError {
				t.Errorf("CalculateStatistics() error = %v, wantErr %v", err, test.wantError)
				return
			}

			// If we didn't expect an error, verify the score
			if !test.wantError && !reflect.DeepEqual(result, test.wantStatistics) {
				t.Errorf("CalculateStatistics() score = %v, want %v", result, test.wantStatistics)
			}
		})
	}
}

func TestCalculateStatisticsByLevel(t *testing.T) {
	tests := []struct {
		name      string
		input     []GameScore
		want      map[int]Statistics
		wantError bool
	}{
		{
			name:      "empty input",
			input:     []GameScore{},
			want:      map[int]Statistics{},
			wantError: false,
		},
		{
			name: "single level with multiple scores",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 1},
			},
			want: map[int]Statistics{
				1: {
					TotalGamesPlayed: 3,
					TotalScore:       600,
					MinimumScore:     100,
					MaximumScore:     300,
					MedianScore:      200,
					AverageScore:     200,
				},
			},
			wantError: false,
		},
		{
			name: "multiple levels with varying counts",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 2},
				{Player: Player{Name: "charlie", ID: 3}, Score: 150, Level: 1},
				{Player: Player{Name: "diana", ID: 4}, Score: 300, Level: 3},
				{Player: Player{Name: "eve", ID: 5}, Score: 250, Level: 2},
				{Player: Player{Name: "frank", ID: 6}, Score: 350, Level: 3},
			},
			want: map[int]Statistics{
				1: {
					TotalGamesPlayed: 2,
					TotalScore:       250,
					MinimumScore:     100,
					MaximumScore:     150,
					MedianScore:      125,
					AverageScore:     125,
				},
				2: {
					TotalGamesPlayed: 2,
					TotalScore:       450,
					MinimumScore:     200,
					MaximumScore:     250,
					MedianScore:      225,
					AverageScore:     225,
				},
				3: {
					TotalGamesPlayed: 2,
					TotalScore:       650,
					MinimumScore:     300,
					MaximumScore:     350,
					MedianScore:      325,
					AverageScore:     325,
				},
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CalculateStatisticsByLevel(test.input)

			if (err != nil) != test.wantError {
				t.Errorf("CalculateStatisticsByLevel() error = %v, wantErr %v", err, test.wantError)
				return
			}

			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("CalculateStatisticsByLevel() = %v, want %v", result, test.want)
			}
		})
	}
}

func TestCalculateStatisticsByPlayer(t *testing.T) {
	tests := []struct {
		name      string
		input     []GameScore
		want      map[Player]Statistics
		wantError bool
	}{
		{
			name:      "empty input",
			input:     []GameScore{},
			want:      map[Player]Statistics{},
			wantError: false,
		},
		{
			name: "single player with multiple scores",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "alex", ID: 1}, Score: 200, Level: 2},
				{Player: Player{Name: "alex", ID: 1}, Score: 300, Level: 3},
			},
			want: map[Player]Statistics{
				{Name: "alex", ID: 1}: {
					TotalGamesPlayed: 3,
					TotalScore:       600,
					MinimumScore:     100,
					MaximumScore:     300,
					MedianScore:      200,
					AverageScore:     200,
				},
			},
			wantError: false,
		},
		{
			name: "multiple players with varying counts",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "alex", ID: 1}, Score: 150, Level: 2},
				{Player: Player{Name: "charlie", ID: 3}, Score: 300, Level: 1},
				{Player: Player{Name: "bob", ID: 2}, Score: 250, Level: 2},
			},
			want: map[Player]Statistics{
				{Name: "alex", ID: 1}: {
					TotalGamesPlayed: 2,
					TotalScore:       250,
					MinimumScore:     100,
					MaximumScore:     150,
					MedianScore:      125,
					AverageScore:     125,
				},
				{Name: "bob", ID: 2}: {
					TotalGamesPlayed: 2,
					TotalScore:       450,
					MinimumScore:     200,
					MaximumScore:     250,
					MedianScore:      225,
					AverageScore:     225,
				},
				{Name: "charlie", ID: 3}: {
					TotalGamesPlayed: 1,
					TotalScore:       300,
					MinimumScore:     300,
					MaximumScore:     300,
					MedianScore:      300,
					AverageScore:     300,
				},
			},
			wantError: false,
		},
		{
			name: "players with same name but different IDs",
			input: []GameScore{
				{Player: Player{Name: "alex", ID: 1}, Score: 100, Level: 1},
				{Player: Player{Name: "alex", ID: 2}, Score: 200, Level: 1},
				{Player: Player{Name: "alex", ID: 1}, Score: 150, Level: 2},
			},
			want: map[Player]Statistics{
				{Name: "alex", ID: 1}: {
					TotalGamesPlayed: 2,
					TotalScore:       250,
					MinimumScore:     100,
					MaximumScore:     150,
					MedianScore:      125,
					AverageScore:     125,
				},
				{Name: "alex", ID: 2}: {
					TotalGamesPlayed: 1,
					TotalScore:       200,
					MinimumScore:     200,
					MaximumScore:     200,
					MedianScore:      200,
					AverageScore:     200,
				},
			},
			wantError: false,
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			result, err := CalculateStatisticsByPlayer(test.input)

			if (err != nil) != test.wantError {
				t.Errorf("CalculateStatisticsByPlayer() error = %v, wantErr %v", err, test.wantError)
				return
			}

			if !reflect.DeepEqual(result, test.want) {
				t.Errorf("CalculateStatisticsByPlayer() = %v, want %v", result, test.want)
			}
		})
	}
}
