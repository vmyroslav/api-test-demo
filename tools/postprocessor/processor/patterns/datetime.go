package patterns

import (
	"fmt"
	"time"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// DatetimeProcessor handles date/time patterns
type DatetimeProcessor struct {
	formats     []string
	replaceWith string
}

// NewDatetimeProcessor creates a new datetime processor
func NewDatetimeProcessor(formats []string, replaceWith string) *DatetimeProcessor {
	return &DatetimeProcessor{
		formats:     formats,
		replaceWith: replaceWith,
	}
}

// ProcessResponse transforms a response value according to datetime pattern
func (p *DatetimeProcessor) ProcessResponse(field string, value string, modifiedFields map[string]bool) (string, bool) {
	if !p.Match(value) {
		return value, false
	}

	// If this field was modified in the request and no fixed replacement
	if modifiedFields[field] && p.replaceWith == "" {
		return fmt.Sprintf("{{ Request.Body 'jsonpath' '$.%s' }}", field), true
	}

	if p.replaceWith != "" {
		return p.replaceWith, true
	}

	return value, false
}

// Match checks if the value matches any of the datetime formats
func (p *DatetimeProcessor) Match(value string) bool {
	for _, format := range p.formats {
		if _, err := time.Parse(format, value); err == nil {
			return true
		}
	}
	return false
}

// ProcessRequest transforms a request value according to datetime pattern
func (p *DatetimeProcessor) ProcessRequest(_ string) (string, bool) {
	if p.replaceWith != "" {
		return p.replaceWith, true
	}

	// Use a proper regex pattern for datetime values
	// [^"]* matches any characters except double quotes
	// This is safer than .* because it won't accidentally match beyond the field
	return `[^"]*`, true
}

// MatcherType returns the matcher type for this pattern
func (p *DatetimeProcessor) MatcherType() config.MatcherType {
	if p.replaceWith != "" {
		return config.ExactMatcher
	}

	return config.RegexMatcher
}

// HasReplacement returns true if this processor has a fixed replacement value
func (p *DatetimeProcessor) HasReplacement() bool {
	return p.replaceWith != ""
}
