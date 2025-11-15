package stats

import (
	"reflect"
	"testing"
)

func TestGroupByLevel(t *testing.T) {

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
			input: []GameScore{{Player: Player{name: "alex", id: 1}, Score: 10, Level: 2}},
			wantStatistics: Statistics{
				TotalGamesPlayed: 1,
				TotalScore:       10,
				MinimumScore:     10,
				MaximumScore:     10,
				MedianScore:      10,
				AverageScore:     10},
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
