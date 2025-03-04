package patterns

import (
	"fmt"
	"regexp"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// UUIDProcessor handles UUID patterns
type UUIDProcessor struct {
	replaceWith string
}

func NewUUIDProcessor(replaceWith string) *UUIDProcessor {
	return &UUIDProcessor{
		replaceWith: replaceWith,
	}
}

// Match checks if the value is a UUID
func (p *UUIDProcessor) Match(value string) bool {
	uuidRegex := regexp.MustCompile(`^[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}$`)
	return uuidRegex.MatchString(value)
}

// ProcessRequest transforms a request value according to UUID pattern
func (p *UUIDProcessor) ProcessRequest(_ string) (string, bool) {
	if p.replaceWith != "" {
		return p.replaceWith, true
	}
	// Use proper regex pattern for UUIDs that will work with Hoverfly's RegexMatch
	return `[0-9a-fA-F]{8}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{4}-[0-9a-fA-F]{12}`, true
}

// MatcherType returns the matcher type for this pattern
func (p *UUIDProcessor) MatcherType() config.MatcherType {
	if p.replaceWith != "" {
		return config.ExactMatcher
	}

	return config.RegexMatcher
}

// ProcessResponse transforms a response value according to UUID pattern
func (p *UUIDProcessor) ProcessResponse(field string, value string, modifiedFields map[string]bool) (string, bool) {
	if !p.Match(value) {
		return value, false
	}

	// If we have a fixed replacement, use it regardless of modified fields
	if p.replaceWith != "" {
		return p.replaceWith, true
	}

	// If this field was modified in the request and modifiedFields is provided
	if modifiedFields != nil && modifiedFields[field] {
		return fmt.Sprintf("{{ Request.Body 'jsonpath' '$.%s' }}", field), true
	}

	return value, false
}

// HasReplacement returns true if this processor has a fixed replacement value
func (p *UUIDProcessor) HasReplacement() bool {
	return p.replaceWith != ""
}
