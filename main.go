package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"sort"

	"github.com/olekukonko/tablewriter"
)

type GameScore struct {
	Player       string
	Score, Level int
}

type Statistics struct {
	TotalGamesPlayed, TotalScore, MinimumScore, MaximumScore int
	MedianScore, AverageScore                                float64
	AverageScoreByLevel                                      map[int]float64
}

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

func GroupByLevel(scores []GameScore) map[int][]GameScore {
	resultMap := make(map[int][]GameScore, len(scores))
	for _, gameScore := range scores {
		slice := append(resultMap[gameScore.Level], gameScore)
		resultMap[gameScore.Level] = slice
	}
	return resultMap
}

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
