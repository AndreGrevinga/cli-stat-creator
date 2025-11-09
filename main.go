package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameScore struct {
	Player       string
	Score, Level int
}

type Statistics struct {
	TotalGamesPlayed, TotalScore                          int
	AverageScore, MedianScore, MinimumScore, MaximumScore float64
	AverageScoreByLevel                                   map[int]float64
}

func ReadScoresFromFile(filename string) ([]GameScore, error) {
	gameScores := []GameScore{}
	data, err := os.ReadFile(filename)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return nil, err
	}
	json.Unmarshal(data, &gameScores)
	for i := range gameScores {
		fmt.Println("Player: " + gameScores[i].Player)
		fmt.Println("Score: " + fmt.Sprintf("%d", gameScores[i].Score))
		fmt.Println("Level: " + fmt.Sprintf("%d", gameScores[i].Level))
	}
	return gameScores, nil
}

func CalculateStatistics(scores []GameScore) (Statistics, error) {
	return Statistics{}, nil
}

func GroupByLevel(scores []GameScore) map[int][]GameScore {
	return map[int][]GameScore{}
}

func main() {
	ReadScoresFromFile("data/input.json")
}
