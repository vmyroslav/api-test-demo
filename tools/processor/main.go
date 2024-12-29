package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
)

func main() {
	var inputFile string
	flag.StringVar(&inputFile, "file", "", "Path to the Hoverfly simulation file")
	flag.Parse()

	if inputFile == "" {
		fmt.Println("Error: -file flag is required")
		flag.Usage()
		os.Exit(1)
	}

	// Create temporary file for processed output
	dir := filepath.Dir(inputFile)
	tmpFile := filepath.Join(dir, "processed_"+filepath.Base(inputFile))

	// Process the simulation
	if err := ProcessSimulation(inputFile, tmpFile); err != nil {
		fmt.Printf("Error processing simulation: %v\n", err)
		os.Exit(1)
	}

	// Replace original file with processed one
	if err := os.Rename(tmpFile, inputFile); err != nil {
		fmt.Printf("Error replacing original file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Successfully processed simulation file:", inputFile)
}
