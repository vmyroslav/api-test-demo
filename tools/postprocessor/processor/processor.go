package processor

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"strings"
	"time"
)

// PostProcessor defines interface for processing Hoverfly simulations
type PostProcessor interface {
	Process(simulation *Simulation) error
}

// Pattern defines a pattern to match and how to handle it
type Pattern struct {
	Match   func(string) bool
	Replace func(string) string
}

// DefaultProcessor implements Processor interface with default pattern matching logic
type DefaultProcessor struct {
	patterns []Pattern
}

// NewDefaultProcessor creates new DefaultProcessor with default patterns
func NewDefaultProcessor() *DefaultProcessor {
	fakePattern := "fake-api-test"

	return &DefaultProcessor{
		patterns: []Pattern{
			{
				// fake-api-test pattern
				Match: func(s string) bool {
					return strings.Contains(s, fakePattern)
				},
				Replace: func(s string) string {
					// Replace last 5 characters with "*" as we use this pattern as randomization
					return strings.Replace(s, s[len(s)-5:], "*", 1)
				},
			},
			{
				// timestamp pattern
				Match: func(s string) bool {
					_, err := time.Parse(time.RFC3339, s)
					return err == nil
				},
				Replace: func(s string) string {
					return "*"
				},
			},
		},
	}
}

func (p *DefaultProcessor) Process(simulation *Simulation) error {
	for i := range simulation.Data.Pairs {
		pair := &simulation.Data.Pairs[i]

		if len(pair.Request.Body) > 0 && pair.Request.Body[0].Value != "" {
			newBody, modifiedFields, err := p.processJSONBody(pair.Request.Body[0].Value)
			if err != nil {
				slog.Error("Error processing JSON body", "err", err, "body", pair.Request.Body[0].Value)
				continue
			}

			if len(modifiedFields) > 0 {
				pair.Request.Body[0].Value = newBody
				pair.Request.Body[0].Matcher = "glob"

				// Process response
				if (pair.Response.Status == http.StatusOK || pair.Response.Status == http.StatusCreated) && strings.Contains(pair.Response.Body, "{") {
					var respData map[string]any

					if err = json.Unmarshal([]byte(pair.Response.Body), &respData); err == nil {
						for field := range modifiedFields {
							// Template only the fields that were modified in request
							if _, exists := respData[field]; exists {
								// Replace response field with templated value, propagating the request field to the response
								respData[field] = fmt.Sprintf("{{ Request.Body 'jsonpath' '$.%s' }}", field)
							}
						}
						if newRespBody, err := json.Marshal(respData); err == nil {
							pair.Response.Body = string(newRespBody)
							pair.Response.Templated = true
						}
					}
				}
			}
		}
	}

	return nil
}

// processJSONBody processes JSON body and returns modified body and modified fields
func (p *DefaultProcessor) processJSONBody(body string) (string, map[string]bool, error) {
	var data map[string]any

	if err := json.Unmarshal([]byte(body), &data); err != nil {
		return body, nil, err
	}

	modifiedFields := make(map[string]bool)
	modified := false

	for field, value := range data {
		if strVal, ok := value.(string); ok {
			for _, pattern := range p.patterns {
				if pattern.Match(strVal) {
					data[field] = pattern.Replace(strVal)
					modifiedFields[field] = true
					modified = true
					break
				}
			}
		}
	}

	if !modified {
		return body, nil, nil
	}

	newBody, err := json.Marshal(data)
	if err != nil {
		return body, nil, err
	}

	return string(newBody), modifiedFields, nil
}

type NullProcessor struct{}

func (p *NullProcessor) Process(_ *Simulation) error {
	return nil
}
