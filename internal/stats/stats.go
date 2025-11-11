package stats

import (
	"errors"
	"slices"
)

// GameScore represents a single game score entry from a player.
// It contains the player name, their score, and the level they were playing at.
type GameScore struct {
	Player       string
	Score, Level int
}

// Statistics holds the calculated statistics from a collection of game scores.
// It includes totals, aggregates (min, max, average, median), and per-level breakdowns.
type Statistics struct {
	TotalGamesPlayed, TotalScore, MinimumScore, MaximumScore int
	MedianScore, AverageScore                                float64
	AverageScoreByLevel                                      map[int]float64
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
	var minimumScore int
	var maximumScore int
	sortedScores := make([]int, totalGamesPlayed)
	for i, score := range scores {
		sortedScores[i] = score.Score
	}
	slices.Sort(sortedScores)
	var medianScore float64
	if totalGamesPlayed%2 == 0 {
		medianScore = (float64(sortedScores[(totalGamesPlayed/2)-1]) + float64(sortedScores[(totalGamesPlayed/2)])) / 2.
	} else {
		medianScore = float64(sortedScores[(totalGamesPlayed-1)/2])
	}
	for i, gameScore := range scores {
		score := gameScore.Score
		totalScore += score
		switch {
		case i == 0:
			minimumScore = score
			maximumScore = score
		case score < minimumScore:
			minimumScore = score
		case score > maximumScore:
			maximumScore = score
		}
	}
	averageScore := float64(totalScore) / float64(totalGamesPlayed)
	groupedByLevel := GroupByLevel(scores)
	averageScoreByLevel := make(map[int]float64, len(groupedByLevel))
	for level, innerScores := range groupedByLevel {
		sum := 0
		for _, score := range innerScores {
			sum += score.Score
		}
		levelAverage := float64(sum) / float64(len(innerScores))
		averageScoreByLevel[level] = levelAverage
	}
	statistic := Statistics{
		TotalGamesPlayed:    totalGamesPlayed,
		TotalScore:          totalScore,
		MinimumScore:        minimumScore,
		MaximumScore:        maximumScore,
		MedianScore:         medianScore,
		AverageScore:        averageScore,
		AverageScoreByLevel: averageScoreByLevel}
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
