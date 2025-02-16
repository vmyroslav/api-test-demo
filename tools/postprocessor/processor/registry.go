package processor

import (
	"fmt"
	"log/slog"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/processor/patterns"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// PatternProcessorFactory is a function that creates a pattern processor from a pattern configuration
type PatternProcessorFactory func(pattern config.Pattern, logger *slog.Logger) PatternProcessor

// EndpointProcessorFactory is a function that creates an endpoint processor from a list of endpoint rules
type EndpointProcessorFactory func(rules []config.EndpointRule, caseSensitive bool, logger *slog.Logger) EndpointProcessor

// Registry manages the registration and creation of both pattern and endpoint processors
type Registry struct {
	patternFactories  map[config.PatternType]PatternProcessorFactory
	endpointFactories map[string]EndpointProcessorFactory
	logger            *slog.Logger
}

// NewRegistry creates and returns a registry with all the standard processors registered
func NewRegistry(logger *slog.Logger) *Registry {
	registry := &Registry{
		patternFactories:  make(map[config.PatternType]PatternProcessorFactory),
		endpointFactories: make(map[string]EndpointProcessorFactory),
		logger:            logger,
	}

	// Register the standard pattern processors
	registry.RegisterPatternProcessor(config.UUIDPattern, func(pattern config.Pattern, logger *slog.Logger) PatternProcessor {
		return patterns.NewUUIDProcessor(pattern.ReplaceWith)
	})

	registry.RegisterPatternProcessor(config.DatetimePattern, func(pattern config.Pattern, logger *slog.Logger) PatternProcessor {
		return patterns.NewDatetimeProcessor(pattern.Formats, pattern.ReplaceWith)
	})

	registry.RegisterPatternProcessor(config.PrefixPattern, func(pattern config.Pattern, logger *slog.Logger) PatternProcessor {
		return patterns.NewPrefixProcessor(pattern.Pattern, pattern.Length, pattern.ReplaceWith)
	})

	// Register the standard endpoint processors
	registry.RegisterEndpointProcessor("static", func(rules []config.EndpointRule, caseSensitive bool, logger *slog.Logger) EndpointProcessor {
		return patterns.NewStaticEndpointProcessor(rules, caseSensitive, logger)
	})

	return registry
}

// RegisterPatternProcessor registers a pattern processor factory for a pattern type
func (r *Registry) RegisterPatternProcessor(patternType config.PatternType, factory PatternProcessorFactory) {
	r.patternFactories[patternType] = factory
	r.logger.Debug("Registered pattern processor factory", "pattern_type", patternType)
}

// RegisterEndpointProcessor registers an endpoint processor factory with the given name
func (r *Registry) RegisterEndpointProcessor(name string, factory EndpointProcessorFactory) {
	r.endpointFactories[name] = factory
	r.logger.Debug("Registered endpoint processor factory", "name", name)
}

// CreatePatternProcessor creates a pattern processor for the given pattern
func (r *Registry) CreatePatternProcessor(pattern config.Pattern) (PatternProcessor, error) {
	factory, exists := r.patternFactories[pattern.Type]
	if !exists {
		return nil, fmt.Errorf("no processor factory registered for pattern type: %s", pattern.Type)
	}

	processor := factory(pattern, r.logger)
	r.logger.Debug("Created pattern processor", "pattern_type", pattern.Type)

	return processor, nil
}

// CreatePatternProcessors creates pattern processors for all the patterns in the config
func (r *Registry) CreatePatternProcessors(ps []config.Pattern) ([]PatternProcessor, error) {
	var processors []PatternProcessor

	for _, pattern := range ps {
		processor, err := r.CreatePatternProcessor(pattern)
		if err != nil {
			return nil, fmt.Errorf("creating processor for pattern type %s: %w", pattern.Type, err)
		}
		processors = append(processors, processor)
	}

	r.logger.Debug("Created pattern processors", "count", len(processors))

	return processors, nil
}

// CreateEndpointProcessor creates an endpoint processor using the factory with the given name
func (r *Registry) CreateEndpointProcessor(name string, rules []config.EndpointRule, caseSensitive bool) EndpointProcessor {
	factory, exists := r.endpointFactories[name]
	if !exists {
		r.logger.Warn("No endpoint processor factory registered with name", "name", name)
		return nil
	}

	processor := factory(rules, caseSensitive, r.logger)
	r.logger.Debug("Created endpoint processor", "name", name)

	return processor
}
