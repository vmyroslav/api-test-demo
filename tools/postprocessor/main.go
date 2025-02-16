package main

import (
	"encoding/json"
	"flag"
	"log/slog"
	"os"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/processor"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

func main() {
	var (
		inputFile  string
		configFile string
		logLevel   string
	)

	flag.StringVar(&inputFile, "file", "", "Path to the Hoverfly simulation file")
	flag.StringVar(&configFile, "config", "", "Path to the processor configuration file")
	flag.StringVar(&logLevel, "log-level", "error", "Log level (debug, info, warn, error)")
	flag.Parse()

	if inputFile == "" {
		slog.Error("missing required argument", "argument", "file")
		flag.Usage()
		os.Exit(1)
	}

	if configFile == "" {
		slog.Error("missing required argument", "argument", "config")
		flag.Usage()
		os.Exit(1)
	}

	// Load configuration
	cfg, err := config.LoadConfig(configFile)
	if err != nil {
		slog.Error("Error loading configuration", "error", err)
		os.Exit(1)
	}

	// Set up logging with appropriate level
	logger := setupLogger(cfg.Settings.Debug)

	// Read and parse the simulation file
	simData, err := os.ReadFile(inputFile)
	if err != nil {
		logger.Error("error reading simulation file", "error", err)
		os.Exit(1)
	}

	var simulation v2.SimulationViewV5
	if err = json.Unmarshal(simData, &simulation); err != nil {
		logger.Error("error parsing simulation", "error", err)
		os.Exit(1)
	}

	// Create and run processor
	proc := processor.New(cfg, logger)
	if err = proc.Process(&simulation); err != nil {
		logger.Error("Error processing simulation", "error", err)
		os.Exit(1)
	}

	// Write processed simulation to temporary file first
	tmpFile := inputFile + ".tmp"

	// Format with indentation for readability
	processedData, err := json.MarshalIndent(simulation, "", "  ")
	if err != nil {
		logger.Error("error serializing processed simulation", "error", err)
		os.Exit(1)
	}

	if err = os.WriteFile(tmpFile, processedData, 0o644); err != nil {
		logger.Error("Error writing processed simulation", "error", err)
		os.Exit(1)
	}

	// Replace original file with processed one
	if err = os.Rename(tmpFile, inputFile); err != nil {
		logger.Error("Error replacing original file", "error", err)
		os.Exit(1)
	}

	logger.Info("Simulation processed successfully", "output", inputFile)
}

func setupLogger(isDebug bool) *slog.Logger {
	logLevel := slog.LevelError

	if isDebug {
		logLevel = slog.LevelDebug
	}

	handler := slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: logLevel,
	})

	l := slog.New(handler)

	slog.SetDefault(l)

	return l
}
