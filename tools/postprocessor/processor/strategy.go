package processor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"strings"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
)

// ProcessingStrategy defines methods for processing requests and responses
type ProcessingStrategy interface {
	// Process handles both the request and response processing for a pair
	Process(pair *v2.RequestMatcherResponsePairViewV5) error

	// ProcessRequest processes the request matchers in a pair
	// Returns a map of modified fields and an error if one occurred
	ProcessRequest(pair *v2.RequestMatcherResponsePairViewV5) (map[string]bool, error)

	// ProcessResponse processes the response in a pair
	// Takes a map of modified fields from the request processing
	ProcessResponse(pair *v2.RequestMatcherResponsePairViewV5, modifiedFields map[string]bool) error
}

// DefaultProcessingStrategy implements the standard processing logic
type DefaultProcessingStrategy struct {
	patternProcessors []PatternProcessor
	logger            *slog.Logger
}

// NewDefaultProcessingStrategy creates a new default processing strategy
func NewDefaultProcessingStrategy(patternProcessors []PatternProcessor, logger *slog.Logger) *DefaultProcessingStrategy {
	return &DefaultProcessingStrategy{
		patternProcessors: patternProcessors,
		logger:            logger,
	}
}

// Process handles both the request and response processing for a pair
func (s *DefaultProcessingStrategy) Process(pair *v2.RequestMatcherResponsePairViewV5) error {
	// First try to process the request
	modifiedFields, err := s.ProcessRequest(pair)
	if err != nil {
		return err
	}

	// Then process the response
	if len(modifiedFields) > 0 {
		// If we modified fields in the request, process response with them
		if err = s.ProcessResponse(pair, modifiedFields); err != nil {
			return err
		}

		return nil
	}

	// If we didn't modify any request fields (e.g., no request body or no matches),
	// try to process the response independently
	if err = s.processResponseForRequestsWithoutBody(pair); err != nil {
		return err
	}

	return nil
}

// ProcessRequest processes the request matchers in a pair
func (s *DefaultProcessingStrategy) ProcessRequest(pair *v2.RequestMatcherResponsePairViewV5) (map[string]bool, error) {
	if len(pair.RequestMatcher.Body) == 0 {
		return nil, nil
	}

	bodyMatcher := &pair.RequestMatcher.Body[0]
	bodyStr := fmt.Sprintf("%v", bodyMatcher.Value)

	// Skip if body is empty or not JSON
	if bodyStr == "" || !strings.HasPrefix(strings.TrimSpace(bodyStr), "{") {
		return nil, nil
	}

	var bodyData map[string]interface{}
	if err := json.Unmarshal([]byte(bodyStr), &bodyData); err != nil {
		return nil, fmt.Errorf("unmarshaling request body: %w", err)
	}

	modifiedFields := make(map[string]bool)
	modified := false
	newBody := make(map[string]interface{})

	// Process each field with pattern processors
	for field, value := range bodyData {
		if strVal, ok := value.(string); ok {
			processed := false

			// Try each processor until one matches
			for _, processor := range s.patternProcessors {
				if processor.Match(strVal) {
					newValue, changed := processor.ProcessRequest(strVal)
					if changed {
						newBody[field] = newValue
						modifiedFields[field] = true
						bodyMatcher.Matcher = processor.MatcherType().String()

						s.logger.Debug("Request field processed",
							"field", field,
							"pattern", processor.MatcherType(),
							"original", strVal,
							"modified", newValue)

						processed = true
						modified = true
						break
					}
				}
			}

			if !processed {
				newBody[field] = strVal
			}
		} else {
			newBody[field] = value
		}
	}

	// Update body if modified
	if modified {
		newBodyBytes, err := json.Marshal(newBody)
		if err != nil {
			return nil, fmt.Errorf("marshaling modified body: %w", err)
		}

		bodyMatcher.Value = string(newBodyBytes)
	}

	return modifiedFields, nil
}

// ProcessResponse processes the response in a pair
func (s *DefaultProcessingStrategy) ProcessResponse(pair *v2.RequestMatcherResponsePairViewV5, modifiedFields map[string]bool) error {
	// Skip if response body is empty or not JSON
	if pair.Response.Body == "" || !strings.HasPrefix(strings.TrimSpace(pair.Response.Body), "{") {
		return nil
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(pair.Response.Body), &respData); err != nil {
		return fmt.Errorf("unmarshaling response body: %w", err)
	}

	modified := false
	newRespData := make(map[string]interface{})

	// Process each field in the response
	for field, value := range respData {
		if strVal, ok := value.(string); ok {
			processed := false

			// Try each processor until one matches
			for _, processor := range s.patternProcessors {
				if processor.Match(strVal) {
					newValue, changed := processor.ProcessResponse(field, strVal, modifiedFields)
					if changed {
						newRespData[field] = newValue

						s.logger.Debug("Response field processed",
							"field", field,
							"original", strVal,
							"modified", newValue)

						processed = true
						modified = true
						break
					}
				}
			}

			// If no processor handled it but the field was in modified request fields,
			// apply default templating
			if !processed && modifiedFields[field] {
				newRespData[field] = fmt.Sprintf("{{ Request.Body 'jsonpath' '$.%s' }}", field)
				modified = true

				s.logger.Debug("Response field templated (default)",
					"field", field)
			} else if !processed {
				newRespData[field] = strVal
			}
		} else {
			newRespData[field] = value
		}
	}

	// Update response if modified
	if modified {
		newRespBody, err := json.Marshal(newRespData)
		if err != nil {
			return fmt.Errorf("marshaling modified response: %w", err)
		}

		pair.Response.Body = string(newRespBody)
		pair.Response.Templated = true

		s.logger.Debug("Response body updated",
			"templated", true,
			"body_length", len(pair.Response.Body))
	}

	return nil
}

// processResponseForRequestsWithoutBody handles cases like GET requests that have no request body
// but still need response processing for fields like UUIDs, timestamps, etc.
// This ensures pattern processors are applied to responses even when there's no request context.
func (s *DefaultProcessingStrategy) processResponseForRequestsWithoutBody(pair *v2.RequestMatcherResponsePairViewV5) error {
	// Skip if response body is empty or not JSON
	if pair.Response.Body == "" || !strings.HasPrefix(strings.TrimSpace(pair.Response.Body), "{") {
		return nil
	}

	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(pair.Response.Body), &respData); err != nil {
		return fmt.Errorf("unmarshaling response body: %w", err)
	}

	modified := false
	newRespData := make(map[string]interface{})

	// Process each field in the response
	for field, value := range respData {
		if strVal, ok := value.(string); ok {
			processed := false

			// Try each processor
			for _, processor := range s.patternProcessors {
				if processor.Match(strVal) {
					// Process this field regardless of whether it was in the request
					newValue, changed := processor.ProcessResponse(field, strVal, nil)
					if changed {
						newRespData[field] = newValue
						s.logger.Debug("Response field processed (independent)",
							"field", field,
							"original", strVal,
							"modified", newValue)
						processed = true
						modified = true
						break
					}
				}
			}

			if !processed {
				newRespData[field] = strVal
			}
		} else {
			newRespData[field] = value
		}
	}

	// Update response if modified
	if modified {
		newRespBody, err := json.Marshal(newRespData)
		if err != nil {
			return fmt.Errorf("marshaling modified response: %w", err)
		}

		pair.Response.Body = string(newRespBody)
		pair.Response.Templated = true // Mark as templated

		s.logger.Debug("Response body updated independently",
			"templated", true,
			"body_length", len(pair.Response.Body))
	}

	return nil
}
