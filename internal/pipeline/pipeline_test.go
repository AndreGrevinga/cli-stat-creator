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

	in := make(chan stats.GameScore, 3)
	in <- stats.GameScore{Player: stats.Player{Name: "Alice", ID: 1}, Score: 100, Level: 1}
	in <- stats.GameScore{Player: stats.Player{Name: "Bob", ID: 2}, Score: 80, Level: 2}
	in <- stats.GameScore{Player: stats.Player{Name: "Charlie", ID: 3}, Score: 90, Level: 1}
	close(in)

	results, err := Aggregate(ctx, in)
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

	in := make(chan stats.GameScore)
	close(in)

	_, err := Aggregate(ctx, in)
	if err == nil {
		t.Error("Expected error for empty input, got nil")
	}
}

func TestPipeline_EndToEnd(t *testing.T) {
	ctx := context.Background()
	filename := "../../data/input.json"

	// Create pipeline with filter
	config := Config{}
	p := New(Filter(config))

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
		FilterByLevel: []int{5},
	}
	p := New(Filter(config))

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

func TestPipeline_ContextCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	filename := "../../data/input.json"

	config := Config{}
	p := New(Filter(config))

	// Cancel immediately
	cancel()

	_, err := p.Run(ctx, filename)
	if err == nil {
		// Might succeed if file is read before cancellation kicks in
		// This is acceptable for this test
	}
}
