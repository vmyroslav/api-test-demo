package patterns

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// PrefixProcessor handles prefix patterns
type PrefixProcessor struct {
	prefix       string
	randomLength int
	replaceWith  string
	// New field to track resource name
	includeResourceName bool
	resourceName        string
}

// NewPrefixProcessor creates a new prefix processor
func NewPrefixProcessor(prefix string, randomLength int, replaceWith string) *PrefixProcessor {
	return &PrefixProcessor{
		prefix:              prefix,
		randomLength:        randomLength,
		replaceWith:         replaceWith,
		includeResourceName: false, // Default to false
		resourceName:        "",    // Default to empty
	}
}

func (p *PrefixProcessor) Match(value string) bool {
	// Check if the value starts with the prefix
	return strings.HasPrefix(value, p.prefix)
}

func (p *PrefixProcessor) ProcessRequest(value string) (string, bool) {
	if p.replaceWith != "" {
		return p.replaceWith, true
	}

	// Properly escape any regex special characters in the prefix
	escapedPrefix := regexp.QuoteMeta(p.prefix)

	// Check if we detected a resource name embedded in the value
	if resourceName := p.detectResourceName(value); resourceName != "" {
		p.includeResourceName = true
		p.resourceName = resourceName

		// Pattern with resource name included
		return fmt.Sprintf("%s%s-[a-zA-Z0-9]{%d}", escapedPrefix, resourceName, p.randomLength), true
	}

	// Standard pattern without resource name
	return fmt.Sprintf("%s[a-zA-Z0-9]{%d}", escapedPrefix, p.randomLength), true
}

// detectResourceName tries to extract a resource name from the value (if it exists)
// For example, if value is "fake-api-test-activity-abcde", it would extract "activity"
func (p *PrefixProcessor) detectResourceName(value string) string {
	// Skip if value doesn't start with prefix
	if !strings.HasPrefix(value, p.prefix) {
		return ""
	}

	// Remove the prefix
	remainder := value[len(p.prefix):]

	// Check if there are at least two parts (resource name and random part)
	parts := strings.Split(remainder, "-")
	if len(parts) >= 2 {
		// The resource name is all parts except the last one (which is the random part)
		return strings.Join(parts[:len(parts)-1], "-")
	}

	return ""
}

func (p *PrefixProcessor) MatcherType() config.MatcherType {
	if p.replaceWith != "" {
		return config.ExactMatcher
	}
	return config.RegexMatcher
}

func (p *PrefixProcessor) ProcessResponse(field string, value string, modifiedFields map[string]bool) (string, bool) {
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

// HasReplacement returns true if this processor has a fixed replacement value
func (p *PrefixProcessor) HasReplacement() bool {
	return p.replaceWith != ""
}
