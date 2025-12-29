// Package stats provides types and functions for calculating game score statistics.
// It supports computing aggregates like averages, medians, min/max values,
// and per-level breakdowns from collections of game scores.
package stats

import (
	"errors"
	"slices"
)

// Player represents a game player with a unique identifier.
// It contains the player's name and a numeric ID to distinguish players with the same name.
type Player struct {
	Name string `json:"name"`
	ID   int    `json:"id"`
}

// GameScore represents a single game score entry from a player.
// It contains the player name, their score, and the level they were playing at.
type GameScore struct {
	Player Player `json:"player"`
	Score  int    `json:"score"`
	Level  int    `json:"level"`
}

// Statistics holds the calculated statistics from a collection of game scores.
// It includes totals, aggregates (min, max, average, median), and per-level breakdowns.
type Statistics struct {
	TotalScore                                   int64
	TotalGamesPlayed, MinimumScore, MaximumScore int
	MedianScore, AverageScore                    float64
	StdDev                                       float64
	Variance                                     float64
	Mode                                         []int
	Percentiles                                  Percentiles
	IQR                                          float64
}

type Percentiles struct {
	P25 float64
	P75 float64
	P90 float64
}

// Validate checks that all GameScore fields contain valid values.
// Returns an error if the player name is empty, player ID is not positive,
// score is negative, or level is not positive.
func (g GameScore) Validate() error {
	if g.Player.Name == "" {
		return errors.New("player name cannot be empty")
	}
	if g.Player.ID <= 0 {
		return errors.New("player id must be positive")
	}
	if g.Score < 0 {
		return errors.New("score cannot be negative")
	}
	if g.Level <= 0 {
		return errors.New("level must be positive")
	}
	return nil
}

// CalculateStatistics computes comprehensive statistics from a slice of game scores.
// It calculates total games, total score, minimum, maximum, average, median,
// and average scores grouped by level.
// Returns an error if the scores slice is empty.
func CalculateStatistics(scores []GameScore) (Statistics, error) {
	totalGamesPlayed := len(scores)
	if totalGamesPlayed == 0 {
		return Statistics{}, errors.New("no scores provided")
	}
	totalScore := int64(0)
	sortedScores := make([]int, totalGamesPlayed)
	for i, gameScore := range scores {
		score := gameScore.Score
		totalScore += int64(score)
		sortedScores[i] = score
	}
	slices.Sort(sortedScores)
	var medianScore float64
	if totalGamesPlayed%2 == 0 {
		medianScore = (float64(sortedScores[(totalGamesPlayed/2)-1]) + float64(sortedScores[(totalGamesPlayed/2)])) / 2.
	} else {
		medianScore = float64(sortedScores[(totalGamesPlayed-1)/2])
	}
	minimumScore := sortedScores[0]
	maximumScore := sortedScores[len(sortedScores)-1]
	averageScore := float64(totalScore) / float64(totalGamesPlayed)
	statistic := Statistics{
		TotalGamesPlayed: totalGamesPlayed,
		TotalScore:       totalScore,
		MinimumScore:     minimumScore,
		MaximumScore:     maximumScore,
		MedianScore:      medianScore,
		AverageScore:     averageScore}
	return statistic, nil
}

// GroupByLevel organizes game scores into a map grouped by level.
// The returned map uses level numbers as keys and slices of GameScore as values.
// This is useful for calculating per-level statistics.
func GroupByLevel(scores []GameScore) map[int][]GameScore {
	resultMap := make(map[int][]GameScore)
	for _, gameScore := range scores {
		slice := append(resultMap[gameScore.Level], gameScore)
		resultMap[gameScore.Level] = slice
	}
	return resultMap
}

// GroupByPlayer organizes game scores into a map grouped by player.
// The returned map uses Player structs as keys and slices of GameScore as values.
// This is useful for calculating per-player statistics.
func GroupByPlayer(scores []GameScore) map[Player][]GameScore {
	resultMap := make(map[Player][]GameScore)
	for _, gameScore := range scores {
		slice := append(resultMap[gameScore.Player], gameScore)
		resultMap[gameScore.Player] = slice
	}
	return resultMap
}

// CalculateStatisticsByLevel computes statistics for each level separately.
// Returns a map of level numbers to their statistics, or an error if any level has no scores.
func CalculateStatisticsByLevel(scores []GameScore) (map[int]Statistics, error) {
	result := make(map[int]Statistics)
	for level, levelScores := range GroupByLevel(scores) {
		levelStats, err := CalculateStatistics(levelScores)
		if err != nil {
			return nil, err
		}
		result[level] = levelStats
	}
	return result, nil
}

// CalculateStatisticsByPlayer computes statistics for each player separately.
// Returns a map of players to their statistics, or an error if any player has no scores.
func CalculateStatisticsByPlayer(scores []GameScore) (map[Player]Statistics, error) {
	result := make(map[Player]Statistics)
	for player, playerScores := range GroupByPlayer(scores) {
		playerStats, err := CalculateStatistics(playerScores)
		if err != nil {
			return nil, err
		}
		result[player] = playerStats
	}
	return result, nil
}
