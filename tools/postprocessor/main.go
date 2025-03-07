package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

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

	// Format with indentation for readability
	processedData, err := json.MarshalIndent(simulation, "", "  ")
	if err != nil {
		logger.Error("error serializing processed simulation", "error", err)
		os.Exit(1)
	}

	// Determine output filename based on override setting
	outputFile := inputFile
	if !cfg.Settings.Override {
		// Generate a new output filename with timestamp
		ext := filepath.Ext(inputFile)
		baseName := strings.TrimSuffix(inputFile, ext)
		timestamp := time.Now().Format("20060102_150405")
		outputFile = fmt.Sprintf("%s_processed_%s%s", baseName, timestamp, ext)
	}

	if err = os.WriteFile(outputFile, processedData, 0o644); err != nil {
		logger.Error("Error writing processed simulation", "error", err)
		os.Exit(1)
	}

	logger.Info("Simulation processed successfully", "output", outputFile)
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
