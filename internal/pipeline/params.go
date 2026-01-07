package pipeline

import (
	"fmt"
	"strconv"
	"strings"
)

// parseLevelString parses the level flag string into a slice of level integers.
// It accepts either a single level (e.g., "5") or a range (e.g., "1-5").
// Returns an empty slice if the input string is empty.
// Returns an error if the input contains non-numeric values or is malformed.
func ParseLevelString(level string) ([]int, error) {
	var levels []int
	var err error
	var minLevel, maxLevel int64

	if level == "" {
		return []int{}, nil
	}

	subStrings := strings.Split(level, "-")
	if len(subStrings) > 2 {
		return nil, fmt.Errorf("invalid level format: expected single level or range (e.g., '5' or '1-5'), got '%s'", level)
	}
	minLevel, err = strconv.ParseInt(subStrings[0], 10, 0)
	if err != nil {
		return nil, err
	}
	if minLevel <= 0 {
		return nil, fmt.Errorf("invalid level: levels must be positive, got %d", minLevel)
	}
	if len(subStrings) == 1 {
		maxLevel = minLevel
	} else {
		maxLevel, err = strconv.ParseInt(subStrings[1], 10, 0)
		if err != nil {
			return nil, err
		}
		if maxLevel <= 0 {
			return nil, fmt.Errorf("invalid level: levels must be positive, got %d", maxLevel)
		}
	}
	if minLevel > maxLevel {
		return nil, fmt.Errorf("invalid level range: min (%d) cannot be greater than max (%d)", minLevel, maxLevel)
	}
	if maxLevel-minLevel > 1000 { // reasonable upper bound
		return nil, fmt.Errorf("level range too large: maximum 1000 levels allowed")
	}
	levels = make([]int, 0, maxLevel-minLevel+1)
	for i := minLevel; i <= maxLevel; i++ {
		levels = append(levels, int(i))
	}
	return levels, nil
}
