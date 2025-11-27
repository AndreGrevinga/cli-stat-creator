// Package main implements a CLI application for analyzing game score statistics.
// It reads JSON files containing game scores and calculates various statistics
// including averages, medians, min/max values, and per-level breakdowns.
package main

import (
	"cli-stat-creator/internal/display"
	"cli-stat-creator/internal/logging"
	"cli-stat-creator/internal/pipeline"
	"context"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"strconv"
	"strings"
)

var (
	detailedFlag     bool
	defaultInputFlag bool
	logLevelFlag     string
	noPlayersFlag    = flag.Bool("no-players", false, "Hide player statistics")
	noLevelsFlag     = flag.Bool("no-levels", false, "Hide level statistics")
	// Todo: Implement player filter flag
	//playerString   = flag.String("player", "", "Only show statistics for given player")
	levelFlag    = flag.String("level", "", "Only show statistics for given levels")
	minScoreFlag = flag.String("min-score", "", "Minimum score filter")
	maxScoreFlag = flag.String("max-score", "", "Maximum score filter")
	logLevelMap  = map[string]slog.Leveler{
		"debug": slog.LevelDebug,
		"info":  slog.LevelInfo,
		"warn":  slog.LevelWarn,
		"error": slog.LevelError,
	}
)

const (
	defaultInputFile = "data/input.json"
	defaultLogLevel  = slog.LevelWarn
)

// parseLevelFlag parses the level flag string into a slice of level integers.
// It accepts either a single level (e.g., "5") or a range (e.g., "1-5").
// Returns an empty slice if the input string is empty.
// Returns an error if the input contains non-numeric values or is malformed.
func parseLevelFlag(level string) ([]int, error) {
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
	if len(subStrings) == 1 {
		maxLevel = minLevel
	} else {
		maxLevel, err = strconv.ParseInt(subStrings[1], 10, 0)
		if err != nil {
			return nil, err
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

// parseLogLevel parses a log level from a flag or environment variable.
// The flag parameter takes precedence over envValue. Returns slog.LevelWarn
// as the default for invalid or empty inputs.
// Valid levels: "debug", "info", "warn", "error" (case-insensitive).
func parseLogLevel(flag string, envValue string) slog.Leveler {
	var levelString string
	if flag == "" {
		levelString = envValue
	} else {
		levelString = flag
	}

	// Normalize input: trim whitespace and convert to lowercase
	levelString = strings.ToLower(strings.TrimSpace(levelString))

	level, ok := logLevelMap[levelString]
	if !ok {
		if levelString != "" {
			fmt.Fprintf(os.Stderr, "Warning: invalid log level '%s', defaulting to 'warn'\n", levelString)
		}
		level = defaultLogLevel
	}
	return level
}

func setupLogging(logLevelFlag string) (context.Context, error) {
	logLevel := parseLogLevel(logLevelFlag, os.Getenv("LOG_LEVEL"))
	handler := slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level:     logLevel,
		AddSource: logLevel == slog.LevelDebug, // Add file:line for DEBUG
	})
	logger := slog.New(handler)
	ctx := context.Background()
	ctx = logging.WithLogger(ctx, logger)
	logger.Info("Application started")
	return ctx, nil
}

// parseFlags validates all command-line flags and returns a pipeline configuration.
// Returns an error if any flag values are invalid.
func parseFlags(ctx context.Context) (pipeline.Config, error) {
	logger := logging.FromContext(ctx)
	var minScore, maxScore int64
	var err error
	if *minScoreFlag == "" {
		minScore = 0
	} else {
		minScore, err = strconv.ParseInt(*minScoreFlag, 10, 0)
		if err != nil {
			logger.Error("Invalid min-score flag value",
				"value", *minScoreFlag,
				"error", err,
			)
			return pipeline.Config{}, err
		}
	}
	if *maxScoreFlag == "" {
		maxScore = 0
	} else {
		maxScore, err = strconv.ParseInt(*maxScoreFlag, 10, 0)
		if err != nil {
			logger.Error("Invalid max-score flag value",
				"value", *maxScoreFlag,
				"error", err,
			)
			return pipeline.Config{}, err
		}
	}
	levels, err := parseLevelFlag(*levelFlag)
	if err != nil {
		logger.Error("Invalid level flag value",
			"value", *levelFlag,
			"error", err,
		)
		return pipeline.Config{}, err
	}
	logger.Debug("Parsed command-line flags",
		"detailed", detailedFlag,
		"default_input", defaultInputFlag,
		"no_players", *noPlayersFlag,
		"no_levels", *noLevelsFlag,
		"level", *levelFlag,
		"min_score", *minScoreFlag,
		"max_score", *maxScoreFlag,
		"log_level", logLevelFlag,
	)
	return pipeline.Config{
		CalculateByLevel:  !*noLevelsFlag,
		CalculateByPlayer: !*noPlayersFlag,
		MinScore:          int(minScore),
		MaxScore:          int(maxScore),
		FilterByLevel:     levels,
	}, nil
}

// getInputFilepath returns the input file path based on command-line flags.
// If the default-input flag is set, returns the default file path.
// Otherwise, prompts the user to enter a file path.
func getInputFilepath() string {
	if defaultInputFlag {
		return defaultInputFile
	}
	fmt.Println("Please input the file path to analyze")
	var filepath string
	fmt.Scan(&filepath)
	return filepath
}

// runPipeline executes the data processing pipeline with the given configuration.
// It creates a pipeline with filtering stages and returns the aggregated results.
func runPipeline(ctx context.Context, cfg pipeline.Config, filepath string) (pipeline.Results, error) {
	logger := logging.FromContext(ctx)
	p := pipeline.New(cfg, pipeline.Filter(cfg))
	logger.Info("Pipeline execution started",
		"filepath", filepath,
		"calculate_by_level", cfg.CalculateByLevel,
		"calculate_by_player", cfg.CalculateByPlayer,
	)
	results, err := p.Run(ctx, filepath)
	if err != nil {
		logger.Error("Pipeline execution failed",
			"filepath", filepath,
			"error", err,
		)
		return pipeline.Results{}, err
	}
	logger.Info("Pipeline completed successfully",
		"overall_count", results.Overall.TotalGamesPlayed,
		"player_count", len(results.ByPlayer),
		"level_count", len(results.ByLevel),
	)
	return results, nil
}

// displayResults renders the pipeline results to stdout.
// It always displays overall statistics, and conditionally shows
// player and level statistics based on the provided flags.
func displayResults(results pipeline.Results, detailed, showPlayers, showLevels bool) {
	display.RenderStatistics(results.Overall)
	if showPlayers {
		fmt.Println("\nPlayer Statistics:")
		display.RenderPlayerStatistics(results.ByPlayer, detailed)
	}
	if showLevels {
		fmt.Println("\nLevel Statistics:")
		display.RenderLevelStatistics(results.ByLevel, detailed)
	}
}

// main is the entry point of the application.
// It reads game scores from a JSON file (data/input.json), calculates statistics,
// and displays the results in formatted tables including overall statistics
// and per-level average scores.
func main() {
	flag.BoolVar(&detailedFlag, "detailed", false, "Show all statistics columns")
	flag.BoolVar(&detailedFlag, "d", false, "Show all statistics columns (shorthand)")
	flag.BoolVar(&defaultInputFlag, "i", false, "Uses the default input.json (shorthand)")
	flag.BoolVar(&defaultInputFlag, "default-input", false, "Uses the default input.json")
	flag.StringVar(&logLevelFlag, "log-level", "", "Level of logging")
	flag.StringVar(&logLevelFlag, "l", "", "Level of logging (shorthand)")
	flag.Parse()

	ctx, err := setupLogging(logLevelFlag)
	if err != nil {
		fmt.Println("Error setting up logging:", err)
		return
	}

	config, err := parseFlags(ctx)
	if err != nil {
		fmt.Println("Error parsing flags:", err)
		return
	}

	filepath := getInputFilepath()

	results, err := runPipeline(ctx, config, filepath)
	if err != nil {
		return
	}

	displayResults(results, detailedFlag, !*noPlayersFlag, !*noLevelsFlag)
}
