package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/processor"
)

// ProcessorType is a custom type for processor options
type ProcessorType string

const (
	DefaultProcessorType ProcessorType = "default"
	NullProcessorType    ProcessorType = "null"
)

func main() {
	slog.SetDefault(slog.New(slog.NewJSONHandler(os.Stdout, nil)))

	var (
		inputFile     string
		proc          processor.PostProcessor
		processorType ProcessorType = NullProcessorType
	)

	flag.StringVar(&inputFile, "file", "", "Path to the Hoverfly simulation file")
	flag.Var(&processorType, "processor", "Processor type to use (default or null)")
	flag.Parse()

	if inputFile == "" {
		slog.Error("Error: -file flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Read the simulation file
	data, err := os.ReadFile(inputFile)
	if err != nil {
		slog.Error("Error reading file", "err", err)
		os.Exit(1)
	}

	var simulation processor.Simulation
	if err = json.Unmarshal(data, &simulation); err != nil {
		slog.Error("Error unmarshaling simulation", "err", err)
		os.Exit(1)
	}

	// Select the processor based on the flag
	switch processorType {
	case DefaultProcessorType:
		proc = processor.NewDefaultProcessor()
	case NullProcessorType:
		proc = &processor.NullProcessor{}
	default:
		slog.Error("Unknown processor type", "type", processorType)
		os.Exit(1)
	}

	// Process the simulation
	if err = proc.Process(&simulation); err != nil {
		slog.Error("Error processing simulation", "err", err)
		os.Exit(1)
	}

	// Create temporary file for processed output
	dir := filepath.Dir(inputFile)
	tmpFile := filepath.Join(dir, "processed_"+filepath.Base(inputFile))

	// Write processed simulation
	output, err := json.MarshalIndent(simulation, "", "  ")
	if err != nil {
		slog.Error("Error marshaling processed simulation", "err", err)
		os.Exit(1)
	}

	if err = os.WriteFile(tmpFile, output, 0o644); err != nil {
		slog.Error("Error writing processed file", "err", err)
		os.Exit(1)
	}

	if err = os.Rename(tmpFile, inputFile); err != nil {
		slog.Error("Error replacing original file", "err", err)
		os.Exit(1)
	}

	slog.Info(fmt.Sprintf("Processed simulation saved to %s", inputFile))
}

// String returns the string representation of the ProcessorType
func (p *ProcessorType) String() string {
	return string(*p)
}

// Set sets the ProcessorType based on the input string
func (p *ProcessorType) Set(value string) error {
	switch value {
	case string(DefaultProcessorType), string(NullProcessorType):
		*p = ProcessorType(value)

		return nil
	default:
		return fmt.Errorf("invalid processor type: %s", value)
	}
}
