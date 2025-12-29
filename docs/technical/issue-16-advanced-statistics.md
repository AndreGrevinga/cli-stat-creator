# Technical Documentation: Advanced Statistical Calculations (Issue #16)

## Overview

This document outlines the technical design for implementing advanced statistical calculations in the cli-stat-creator application. The enhancement adds statistical measures including standard deviation, variance, percentiles, mode, and interquartile range (IQR).

## Current Architecture

The existing statistics system is centered around:
- `GameScore`: Core data structure containing Player, Score, and Level
- `Statistics`: Structure holding calculated statistical measures
- `CalculateStatistics()`: Main function for computing statistics from scores

## Proposed Changes

### 1. Enhanced Statistics Structure

```go
// internal/stats/stats.go

type Statistics struct {
    // Existing fields
    Count      int
    Total      int
    Average    float64
    Median     float64
    Min        int
    Max        int

    // New fields for advanced statistics
    StdDev     float64        // Standard deviation
    Variance   float64        // Variance
    Mode       []int          // Most frequent score(s) - slice to handle multiple modes
    Percentiles Percentiles   // Percentile calculations
    IQR        float64        // Interquartile range
}

type Percentiles struct {
    P25  float64  // 25th percentile (Q1)
    P75  float64  // 75th percentile (Q3)
    P90  float64  // 90th percentile
}
```

### 2. Function Architecture

#### Main Calculation Flow

```go
// internal/stats/stats.go

// CalculateStatistics calculates comprehensive statistics including advanced measures
func CalculateStatistics(scores []GameScore) (Statistics, error) {
    // 1. Validation
    if err := validateScores(scores); err != nil {
        return Statistics{}, err
    }

    // 2. Extract score values for calculations
    values := extractScoreValues(scores)

    // 3. Sort values once for multiple calculations
    sorted := sortScoreValues(values)

    // 4. Calculate basic statistics (existing)
    stats := calculateBasicStats(values, sorted)

    // 5. Calculate advanced statistics (new)
    stats.Variance = calculateVariance(values, stats.Average)
    stats.StdDev = calculateStdDev(stats.Variance)
    stats.Percentiles = calculatePercentiles(sorted)
    stats.IQR = calculateIQR(stats.Percentiles)
    stats.Mode = calculateMode(values)

    return stats, nil
}
```

#### Helper Functions for Advanced Statistics

```go
// internal/stats/advanced.go (new file)

// calculateVariance computes the variance of score values
// Variance = sum((x - mean)²) / n
func calculateVariance(values []int, mean float64) float64 {
    // Implementation details omitted
    // Returns: variance value
}

// calculateStdDev computes standard deviation from variance
// StdDev = √variance
func calculateStdDev(variance float64) float64 {
    // Implementation details omitted
    // Returns: standard deviation value
}

// calculatePercentiles computes requested percentile values from sorted scores
// Uses linear interpolation method for percentile calculation
func calculatePercentiles(sortedValues []int) Percentiles {
    // Implementation details omitted
    // Returns: Percentiles{P25, P75, P90}
}

// calculateIQR computes the interquartile range
// IQR = Q3 - Q1 (75th percentile - 25th percentile)
func calculateIQR(percentiles Percentiles) float64 {
    // Implementation details omitted
    // Returns: IQR value
}

// calculateMode finds the most frequently occurring score(s)
// Returns slice to handle multiple modes (bimodal, multimodal distributions)
func calculateMode(values []int) []int {
    // Implementation details omitted
    // Returns: slice of mode values (empty if no mode exists)
}
```

#### Utility Functions

```go
// internal/stats/utils.go (new file or add to existing)

// extractScoreValues extracts score integers from GameScore slice
func extractScoreValues(scores []GameScore) []int {
    // Implementation details omitted
    // Returns: []int containing only score values
}

// sortScoreValues creates a sorted copy of score values
func sortScoreValues(values []int) []int {
    // Implementation details omitted
    // Returns: sorted copy of values (ascending order)
}

// validateScores performs basic validation on score data
func validateScores(scores []GameScore) error {
    // Implementation details omitted
    // Returns: error if validation fails, nil otherwise
}
```

### 3. Integration with Existing Functions

#### Grouped Statistics Functions

The grouped statistics functions will automatically benefit from the enhanced Statistics structure:

```go
// internal/stats/stats.go

// CalculateStatisticsByLevel calculates advanced statistics for each level
func CalculateStatisticsByLevel(scores []GameScore) (map[int]Statistics, error) {
    grouped := GroupByLevel(scores)

    result := make(map[int]Statistics)
    for level, levelScores := range grouped {
        // Calls enhanced CalculateStatistics internally
        stats, err := CalculateStatistics(levelScores)
        if err != nil {
            return nil, fmt.Errorf("level %d: %w", level, err)
        }
        result[level] = stats
    }

    return result, nil
}

// CalculateStatisticsByPlayer calculates advanced statistics for each player
func CalculateStatisticsByPlayer(scores []GameScore) (map[Player]Statistics, error) {
    grouped := GroupByPlayer(scores)

    result := make(map[Player]Statistics)
    for player, playerScores := range grouped {
        // Calls enhanced CalculateStatistics internally
        stats, err := CalculateStatistics(playerScores)
        if err != nil {
            return nil, fmt.Errorf("player %s: %w", player.Name, err)
        }
        result[player] = stats
    }

    return result, nil
}
```

### 4. Display Layer Updates

#### CLI Display

```go
// internal/display/table.go

// RenderStatistics displays overall statistics with advanced measures
func RenderStatistics(s Statistics) {
    // Create table with existing columns
    // Add new columns for advanced statistics:
    // - Std Dev
    // - Variance
    // - Mode
    // - P25/P75/P90 (in detailed view)
    // - IQR

    // Formatting considerations:
    // - Float values: 2 decimal places
    // - Mode: comma-separated list if multiple modes
    // - Percentiles: show in detailed view only
}

// RenderLevelStatistics displays per-level statistics with advanced measures
func RenderLevelStatistics(statistics map[int]Statistics, detailed bool) {
    // Use RenderGroupedStatistics with enhanced column set
    // Show advanced stats based on detailed flag
}

// RenderPlayerStatistics displays per-player statistics with advanced measures
func RenderPlayerStatistics(statistics map[Player]Statistics, detailed bool) {
    // Use RenderGroupedStatistics with enhanced column set
    // Show advanced stats based on detailed flag
}

// Helper: formatMode handles mode display (multiple values)
func formatMode(modes []int) string {
    // Implementation details omitted
    // Returns: comma-separated string or "N/A" if no mode
}
```

#### HTTP Response Updates

```go
// internal/handlers/template_data.go

type StatisticsData struct {
    // Existing fields
    Count      int     `json:"count"`
    Total      int     `json:"total"`
    Average    float64 `json:"average"`
    Median     float64 `json:"median"`
    Min        int     `json:"min"`
    Max        int     `json:"max"`

    // New fields
    StdDev     float64   `json:"stdDev"`
    Variance   float64   `json:"variance"`
    Mode       []int     `json:"mode"`
    P25        float64   `json:"p25"`
    P75        float64   `json:"p75"`
    P90        float64   `json:"p90"`
    IQR        float64   `json:"iqr"`
}

// convertToStatisticsData maps Stats structure to template data
func convertToStatisticsData(s stats.Statistics) StatisticsData {
    // Map all fields including new advanced statistics
    // Returns: StatisticsData for template rendering
}
```

### 5. Data Flow Diagram

```
Input: []GameScore
    ↓
CalculateStatistics()
    ├── validateScores()
    ├── extractScoreValues() → []int
    ├── sortScoreValues() → []int (sorted)
    │
    ├── calculateBasicStats()
    │   ├── Count, Total, Min, Max
    │   ├── Average
    │   └── Median (from sorted values)
    │
    └── Advanced Statistics
        ├── calculateVariance(values, average) → float64
        ├── calculateStdDev(variance) → float64
        ├── calculatePercentiles(sorted) → Percentiles
        ├── calculateIQR(percentiles) → float64
        └── calculateMode(values) → []int
    ↓
Output: Statistics (enhanced)
    ↓
Display/Render Layer
    ├── CLI: RenderStatistics()
    │   └── Table with additional columns
    │
    └── HTTP: convertToStatisticsData()
        └── JSON/HTML with advanced stats
```

### 6. Testing Strategy

#### Unit Tests

```go
// internal/stats/advanced_test.go (new file)

// Test individual calculation functions
func TestCalculateVariance(t *testing.T) {
    // Test cases:
    // - Empty slice
    // - Single value
    // - Known variance values
    // - Zero variance (all same values)
}

func TestCalculateStdDev(t *testing.T) {
    // Test cases:
    // - Zero variance
    // - Known standard deviation values
    // - Relationship to variance (stddev² = variance)
}

func TestCalculatePercentiles(t *testing.T) {
    // Test cases:
    // - Small datasets (< 4 values)
    // - Even vs odd number of values
    // - Known percentile values
    // - Edge cases (all same values)
}

func TestCalculateIQR(t *testing.T) {
    // Test cases:
    // - Normal distribution
    // - Narrow distribution (small IQR)
    // - Wide distribution (large IQR)
}

func TestCalculateMode(t *testing.T) {
    // Test cases:
    // - No mode (all unique values)
    // - Single mode (unimodal)
    // - Multiple modes (bimodal, multimodal)
    // - All values the same
}
```

#### Integration Tests

```go
// internal/stats/stats_test.go (existing file)

func TestCalculateStatistics_WithAdvancedStats(t *testing.T) {
    // Test that CalculateStatistics correctly populates all fields
    // Verify relationships between stats (e.g., IQR = P75 - P25)
}

func TestCalculateStatisticsByLevel_WithAdvancedStats(t *testing.T) {
    // Test grouped calculations include advanced statistics
}

func TestCalculateStatisticsByPlayer_WithAdvancedStats(t *testing.T) {
    // Test grouped calculations include advanced statistics
}
```

## Implementation Phases

### Phase 1: Core Calculations
1. Create `internal/stats/advanced.go` with calculation functions
2. Add fields to Statistics structure
3. Update `CalculateStatistics()` to call new functions
4. Write comprehensive unit tests

### Phase 2: Display Integration
1. Update CLI display functions to show new statistics
2. Update HTTP response structures and templates
3. Add formatting helpers for mode display
4. Test rendering in both CLI and web interface

### Phase 3: Documentation
1. Update README.md with new statistics descriptions
2. Update CLAUDE.md with new function documentation
3. Add Go doc comments for all new functions
4. Update data format documentation if needed

## Performance Considerations

- **Single Sort**: Perform sorting once and reuse for median, percentiles, and mode calculations
- **In-place Operations**: Use in-place sorting where possible to minimize memory allocation
- **Mode Calculation**: Use hash map for frequency counting (O(n) time complexity)
- **Variance Calculation**: Single pass through data after computing mean

## Backward Compatibility

- All existing API endpoints remain unchanged
- CLI output adds new columns but preserves existing columns
- JSON responses add new fields but retain all existing fields
- No breaking changes to function signatures (Statistics struct is expanded, not modified)

## Error Handling

- Return errors for invalid inputs (empty slices, negative scores)
- Handle edge cases gracefully (insufficient data for percentiles, no mode exists)
- Provide meaningful error messages with context
- Validate preconditions before calculations

## Future Enhancements

- Configurable percentiles (allow users to request custom percentile values)
- Outlier detection using IQR method (values outside Q1 - 1.5×IQR or Q3 + 1.5×IQR)
- Coefficient of variation (CV = stddev / mean)
- Skewness and kurtosis for distribution shape analysis
