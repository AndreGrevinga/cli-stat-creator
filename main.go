// Package main implements a CLI application for analyzing game score statistics.
// It reads JSON files containing game scores and calculates various statistics
// including averages, medians, min/max values, and per-level breakdowns.
package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
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

// ReadScoresFromFile reads game scores from a JSON file and returns them as a slice of GameScore.
// The function expects the JSON file to contain an array of objects with Player, Score, and Level fields.
// Returns an error if the file cannot be read or if the JSON format is invalid.
func ReadScoresFromFile(filename string) ([]GameScore, error) {
	gameScores := []GameScore{}
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(data, &gameScores); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}
	return gameScores, nil
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
	averageScoreByLevel := make(map[int]float64, totalGamesPlayed)
	sortedScores := make([]GameScore, totalGamesPlayed)
	copy(sortedScores, scores)
	sort.Slice(sortedScores, func(i, j int) bool {
		return sortedScores[i].Score < sortedScores[j].Score
	})
	var medianScore float64
	if totalGamesPlayed%2 == 0 {
		medianScore = (float64(sortedScores[(totalGamesPlayed/2)-1].Score) + float64(sortedScores[(totalGamesPlayed/2)].Score)) / 2.
	} else {
		medianScore = float64(sortedScores[(totalGamesPlayed-1)/2].Score)
	}
	for i, gameScore := range scores {
		score := gameScore.Score
		totalScore += score
		if i == 0 {
			minimumScore = score
			maximumScore = score
		}
		if score < minimumScore {
			minimumScore = score
		}
		if score > maximumScore {
			maximumScore = score
		}

	}
	averageScore := float64(totalScore) / float64(totalGamesPlayed)
	groupedByLevel := GroupByLevel(scores)
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
	resultMap := make(map[int][]GameScore, len(scores))
	for _, gameScore := range scores {
		slice := append(resultMap[gameScore.Level], gameScore)
		resultMap[gameScore.Level] = slice
	}
	return resultMap
}

// main is the entry point of the application.
// It reads game scores from a JSON file (data/input.json), calculates statistics,
// and displays the results in formatted tables including overall statistics
// and per-level average scores.
func main() {
	fmt.Println("Please input the file path to analyze")
	var filepath string
	//fmt.Scan(&filepath)
	filepath = "data/input.json"
	var gameScores []GameScore
	gameScores, err := ReadScoresFromFile(filepath)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}
	Statistics, err := CalculateStatistics(gameScores)
	if err != nil {
		fmt.Println("Error calculating statistics:", err)
		return
	}
	dataTable := tablewriter.NewTable(os.Stdout)
	dataTable.Header([]string{"TotalGamesPlayed", "TotalScore", "AverageScore", "MedianScore", "MinimumScore", "MaximumScore"})
	dataTable.Append([]string{
		fmt.Sprintf("%d", Statistics.TotalGamesPlayed), fmt.Sprintf("%d", Statistics.TotalScore),
		fmt.Sprintf("%.2f", Statistics.AverageScore), fmt.Sprintf("%.2f", Statistics.MedianScore),
		fmt.Sprintf("%d", Statistics.MinimumScore), fmt.Sprintf("%d", Statistics.MaximumScore),
	})
	dataTable.Render()
	groupedDataTable := tablewriter.NewTable(os.Stdout)
	groupedDataTable.Header([]string{"Level", "Average Score"})
	for level, averageScore := range Statistics.AverageScoreByLevel {
		groupedDataTable.Append([]string{fmt.Sprintf("%d", level), fmt.Sprintf("%.2f", averageScore)})
	}
	groupedDataTable.Render()
}
