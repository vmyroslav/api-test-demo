package processor

import (
	"fmt"
	"log/slog"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// PatternProcessor is responsible for processing a specific pattern type
type PatternProcessor interface {
	// Match checks if the value matches this pattern type
	Match(value string) bool

	// MatcherType returns the Hoverfly matcher type to use (regex, exact, glob)
	MatcherType() config.MatcherType

	// ProcessRequest transforms a request value according to this pattern
	// Returns the transformed value and whether it was modified
	ProcessRequest(value string) (string, bool)

	// ProcessResponse transforms a response value according to this pattern
	// field: the field name in the JSON
	// value: the field value to process
	// modifiedFields: map of all fields modified in the request
	// Returns the transformed value and whether it was modified
	ProcessResponse(field string, value string, modifiedFields map[string]bool) (string, bool)

	// HasReplacement returns true if this processor has a fixed replacement value
	HasReplacement() bool
}

// EndpointProcessor processes static endpoint rules
type EndpointProcessor interface {
	// FindMatchingRule checks if a pair matches any endpoint rule
	FindMatchingRule(pair *v2.RequestMatcherResponsePairViewV5) *config.EndpointRule

	// ApplyRule applies a static response to a pair
	ApplyRule(pair *v2.RequestMatcherResponsePairViewV5, rule *config.EndpointRule) error
}

// PostProcessor is the main orchestrator for processing Hoverfly simulations
type PostProcessor struct {
	config             *config.Config
	endpointProcessor  EndpointProcessor
	processingStrategy ProcessingStrategy
	logger             *slog.Logger
}

// New creates a new processor with the appropriate processors and strategies
func New(cfg *config.Config, logger *slog.Logger) *PostProcessor {
	// Create a unified processor registry with standard processors
	processorRegistry := NewRegistry(logger)

	// Create pattern processors directly from config
	patternProcessors, err := processorRegistry.CreatePatternProcessors(cfg.Patterns)
	if err != nil {
		logger.Error("Error creating pattern processors", "error", err)
		// Create an empty slice in case of error
		patternProcessors = []PatternProcessor{}
	}

	// Create endpoint processor if there are endpoint rules
	var endpointProc EndpointProcessor
	if len(cfg.Endpoints) > 0 {
		endpointProc = processorRegistry.CreateEndpointProcessor("static", cfg.Endpoints, cfg.Settings.CaseSensitive)
	}

	// Create the default processing strategy with the slice of processors
	processingStrategy := NewDefaultProcessingStrategy(patternProcessors, logger)

	return &PostProcessor{
		config:             cfg,
		endpointProcessor:  endpointProc,
		processingStrategy: processingStrategy,
		logger:             logger,
	}
}

// Process implements the Processor interface
func (p *PostProcessor) Process(simulation *v2.SimulationViewV5) error {
	if simulation == nil {
		return fmt.Errorf("nil simulation provided")
	}

	p.logger.Debug("Processing simulation",
		"pairs_count", len(simulation.RequestResponsePairs))

	for i := range simulation.RequestResponsePairs {
		pair := &simulation.RequestResponsePairs[i]

		// 1. Check if this is a static endpoint rule
		if p.endpointProcessor != nil {
			if rule := p.endpointProcessor.FindMatchingRule(pair); rule != nil {
				if err := p.endpointProcessor.ApplyRule(pair, rule); err != nil {
					return fmt.Errorf("applying static response for pair %d: %w", i, err)
				}
				continue
			}
		}

		// 2. Process the pair
		if err := p.processingStrategy.Process(pair); err != nil {
			return fmt.Errorf("processing pair %d: %w", i, err)
		}
	}

	p.logger.Debug("Simulation processing completed")

	return nil
}
