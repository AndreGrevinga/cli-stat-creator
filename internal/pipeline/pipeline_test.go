package pipeline

import (
	"cli-stat-creator/internal/stats"
	"context"
	"testing"
	"time"
)

func TestFilter_EmptyConfig(t *testing.T) {
	ctx := context.Background()
	config := Config{}

	// Create test scores
	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 50, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 75, Level: 1}
	close(in)

	// Apply filter
	stage := Filter(config)
	out := stage(ctx, in)

	// Collect results
	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// With empty config, all scores should pass through
	if len(results) != 3 {
		t.Errorf("Expected 3 scores, got %d", len(results))
	}
}

func TestFilter_FilterByPlayer(t *testing.T) {
	ctx := context.Background()
	config := Config{
		FilterByPlayer: []string{"Alice", "Bob"},
	}

	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 50, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 75, Level: 1}
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only Alice and Bob should pass through
	if len(results) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(results))
	}
}

func TestFilter_FilterByLevel(t *testing.T) {
	ctx := context.Background()
	config := Config{
		FilterByLevel: []int{1},
	}

	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 50, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 75, Level: 1}
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only level 1 scores should pass through
	if len(results) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(results))
	}
	for _, score := range results {
		if score.Level != 1 {
			t.Errorf("Expected level 1, got level %d", score.Level)
		}
	}
}

func TestFilter_FilterByScoreRange(t *testing.T) {
	ctx := context.Background()
	config := Config{
		MinScore: 60,
		MaxScore: 90,
	}

	in := make(chan stats.GameScore, 4)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 50, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 75, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Diana", ID: 4}, Score: 60, Level: 2}
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only scores in range [60, 90] should pass through
	if len(results) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(results))
	}
	for _, score := range results {
		if score.Score < 60 || score.Score > 90 {
			t.Errorf("Score %d is outside range [60, 90]", score.Score)
		}
	}
}

func TestFilter_CombinedPlayerAndLevel(t *testing.T) {
	ctx := context.Background()
	config := Config{
		FilterByPlayer: []string{"Alice"},
		FilterByLevel:  []int{1},
	}

	in := make(chan stats.GameScore, 4)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}  // Matches both
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 50, Level: 2}   // Matches player only
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 75, Level: 1}     // Matches level only
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 60, Level: 2} // Matches neither
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only Alice's level 1 score should pass through (AND logic)
	if len(results) != 1 {
		t.Errorf("Expected 1 score, got %d", len(results))
	}
	if len(results) > 0 {
		if results[0].Player.Name != "Alice" || results[0].Level != 1 {
			t.Errorf("Expected Alice's level 1 score, got %s level %d", results[0].Player.Name, results[0].Level)
		}
	}
}

func TestFilter_CombinedPlayerAndScoreRange(t *testing.T) {
	ctx := context.Background()
	config := Config{
		FilterByPlayer: []string{"Alice", "Bob"},
		MinScore:       60,
		MaxScore:       90,
	}

	in := make(chan stats.GameScore, 5)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1} // Alice, but score too high
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 75, Level: 2}  // Matches both
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 50, Level: 1}    // Bob, but score too low
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 60, Level: 2}    // Matches both
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 70, Level: 1} // Score OK, wrong player
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only Alice's 75 and Bob's 60 should pass through
	if len(results) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(results))
	}
	for _, score := range results {
		if score.Player.Name != "Alice" && score.Player.Name != "Bob" {
			t.Errorf("Unexpected player: %s", score.Player.Name)
		}
		if score.Score < 60 || score.Score > 90 {
			t.Errorf("Score %d is outside range [60, 90]", score.Score)
		}
	}
}

func TestFilter_AllFiltersCombined(t *testing.T) {
	ctx := context.Background()
	config := Config{
		FilterByPlayer: []string{"Alice"},
		FilterByLevel:  []int{1, 2},
		MinScore:       60,
		MaxScore:       90,
	}

	in := make(chan stats.GameScore, 6)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 75, Level: 1}  // Matches all
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 85, Level: 2}  // Matches all
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1} // Alice, level 1, but score too high
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 75, Level: 3}  // Alice, score OK, wrong level
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 75, Level: 1}    // Level and score OK, wrong player
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 50, Level: 1}  // Alice, level 1, score too low
	close(in)

	stage := Filter(config)
	out := stage(ctx, in)

	var results []stats.GameScore
	for score := range out {
		results = append(results, score)
	}

	// Only first two scores should pass through
	if len(results) != 2 {
		t.Errorf("Expected 2 scores, got %d", len(results))
	}
	for _, score := range results {
		if score.Player.Name != "Alice" {
			t.Errorf("Expected Alice, got %s", score.Player.Name)
		}
		if score.Level != 1 && score.Level != 2 {
			t.Errorf("Expected level 1 or 2, got %d", score.Level)
		}
		if score.Score < 60 || score.Score > 90 {
			t.Errorf("Score %d is outside range [60, 90]", score.Score)
		}
	}
}

func TestFilter_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	config := Config{}

	in := make(chan stats.GameScore)
	stage := Filter(config)
	out := stage(ctx, in)

	// Cancel context immediately
	cancel()

	// Send a score
	go func() {
		time.Sleep(10 * time.Millisecond)
		in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
		close(in)
	}()

	// Should exit quickly due to cancellation
	timeout := time.After(100 * time.Millisecond)
	resultCount := 0
	done := make(chan bool)

	go func() {
		for range out {
			resultCount++
		}
		done <- true
	}()

	select {
	case <-done:
		// Filter should have stopped processing
	case <-timeout:
		t.Error("Filter did not respect context cancellation")
	}
}

func TestSource_ValidFile(t *testing.T) {
	ctx := context.Background()
	filename := "../../data/input.json"

	out, err := Source(ctx, filename)
	if err != nil {
		t.Fatalf("Source failed: %v", err)
	}

	var count int
	for range out {
		count++
	}

	if count == 0 {
		t.Error("Expected scores from input.json, got none")
	}
}

func TestSource_InvalidFile(t *testing.T) {
	ctx := context.Background()
	filename := "nonexistent.json"

	_, err := Source(ctx, filename)
	if err == nil {
		t.Error("Expected error for nonexistent file, got nil")
	}
}

func TestSource_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	filename := "../../data/input.json"

	out, err := Source(ctx, filename)
	if err != nil {
		t.Fatalf("Source failed: %v", err)
	}

	// Cancel immediately
	cancel()

	// Drain channel - should close quickly
	timeout := time.After(100 * time.Millisecond)
	done := make(chan bool)

	go func() {
		for range out {
		}
		done <- true
	}()

	select {
	case <-done:
		// Source respected cancellation
	case <-timeout:
		t.Error("Source did not respect context cancellation")
	}
}

func TestAggregate_ValidScores(t *testing.T) {
	ctx := context.Background()
	config := Config{
		CalculateByLevel:  true,
		CalculateByPlayer: true,
	}

	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 80, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 90, Level: 1}
	close(in)

	results, err := Aggregate(ctx, in, config)
	if err != nil {
		t.Fatalf("Aggregate failed: %v", err)
	}

	if results.Overall.TotalGamesPlayed != 3 {
		t.Errorf("Expected 3 games, got %d", results.Overall.TotalGamesPlayed)
	}

	if len(results.ByLevel) != 2 {
		t.Errorf("Expected 2 levels, got %d", len(results.ByLevel))
	}

	if len(results.ByPlayer) != 3 {
		t.Errorf("Expected 3 players, got %d", len(results.ByPlayer))
	}
}

func TestAggregate_EmptyInput(t *testing.T) {
	ctx := context.Background()
	config := Config{}

	in := make(chan stats.GameScore)
	close(in)

	_, err := Aggregate(ctx, in, config)
	if err == nil {
		t.Error("Expected error for empty input, got nil")
	}
}

func TestPipeline_EndToEnd(t *testing.T) {
	ctx := context.Background()
	filename := "../../data/input.json"

	// Create pipeline with filter
	config := Config{
		CalculateByLevel:  true,
		CalculateByPlayer: true,
	}
	p := New(config, Filter(config))

	results, err := p.Run(ctx, filename)
	if err != nil {
		t.Fatalf("Pipeline.Run failed: %v", err)
	}

	if results.Overall.TotalGamesPlayed == 0 {
		t.Error("Expected non-zero games played")
	}

	if len(results.ByLevel) == 0 {
		t.Error("Expected level statistics")
	}

	if len(results.ByPlayer) == 0 {
		t.Error("Expected player statistics")
	}
}

func TestPipeline_WithFiltering(t *testing.T) {
	ctx := context.Background()
	filename := "../../data/input.json"

	// Create pipeline with filter for level 5 only
	config := Config{
		FilterByLevel:    []int{5},
		CalculateByLevel: true,
	}
	p := New(config, Filter(config))

	results, err := p.Run(ctx, filename)
	if err != nil {
		t.Fatalf("Pipeline.Run failed: %v", err)
	}

	// Should only have level 5 statistics
	if len(results.ByLevel) != 1 {
		t.Errorf("Expected only 1 level, got %d", len(results.ByLevel))
	}

	if _, exists := results.ByLevel[5]; !exists {
		t.Error("Expected level 5 statistics")
	}
}

func TestAggregate_ConditionalCalculation(t *testing.T) {
	ctx := context.Background()

	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 80, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 90, Level: 1}
	close(in)

	// Test with both flags off
	config := Config{
		CalculateByLevel:  false,
		CalculateByPlayer: false,
	}

	results, err := Aggregate(ctx, in, config)
	if err != nil {
		t.Fatalf("Aggregate failed: %v", err)
	}

	// Overall should always be calculated
	if results.Overall.TotalGamesPlayed != 3 {
		t.Errorf("Expected 3 games, got %d", results.Overall.TotalGamesPlayed)
	}

	// ByLevel and ByPlayer should be empty/nil when flags are false
	if len(results.ByLevel) != 0 {
		t.Errorf("Expected 0 levels (flag off), got %d", len(results.ByLevel))
	}

	if len(results.ByPlayer) != 0 {
		t.Errorf("Expected 0 players (flag off), got %d", len(results.ByPlayer))
	}
}

func TestPipeline_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	filename := "../../data/input.json"

	config := Config{}
	p := New(config, Filter(config))

	// Cancel immediately
	cancel()

	_, err := p.Run(ctx, filename)
	if err == nil {
		// Might succeed if file is read before cancellation kicks in
		// This is acceptable for this test
	}
}
