// Package stats provides types and functions for calculating game score statistics.
// It supports computing aggregates like averages, medians, min/max values,
// and per-level breakdowns from collections of game scores.
package stats

import (
	"errors"
	"slices"
)

type Player struct {
	name string
	id   int
}

// GameScore represents a single game score entry from a player.
// It contains the player name, their score, and the level they were playing at.
type GameScore struct {
	Player       Player
	Score, Level int
}

// Statistics holds the calculated statistics from a collection of game scores.
// It includes totals, aggregates (min, max, average, median), and per-level breakdowns.
type Statistics struct {
	TotalGamesPlayed, TotalScore, MinimumScore, MaximumScore int
	MedianScore, AverageScore                                float64
}

type LevelStatistics struct {
	Level int
}

type PlayerStatistics struct {
	Player string
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
	totalScore := 0
	sortedScores := make([]int, totalGamesPlayed)
	for i, gameScore := range scores {
		score := gameScore.Score
		totalScore += score
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

func GroupByPlayer(scores []GameScore) {
	resultMap := make(map[Player][]GameScore)
	for _, gameScore := range scores {
		slice := append(resultMap[gameScore.Player], gameScore)
		resultMap[gameScore.Player] = slice
	}
}
