package patterns

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"net/http"
	"regexp"
	"strings"

	v2 "github.com/SpectoLabs/hoverfly/core/handlers/v2"
	"github.com/vmyroslav/api-test-demo/tools/postprocessor/config"
)

// StaticEndpointProcessor implements the EndpointProcessor interface for static endpoints
type StaticEndpointProcessor struct {
	rules         []config.EndpointRule
	caseSensitive bool
	logger        *slog.Logger
}

// NewStaticEndpointProcessor creates a new static endpoint processor
func NewStaticEndpointProcessor(rules []config.EndpointRule, caseSensitive bool, logger *slog.Logger) *StaticEndpointProcessor {
	return &StaticEndpointProcessor{
		rules:         rules,
		caseSensitive: caseSensitive,
		logger:        logger,
	}
}

// FindMatchingRule checks if a pair matches any endpoint rule
func (p *StaticEndpointProcessor) FindMatchingRule(pair *v2.RequestMatcherResponsePairViewV5) *config.EndpointRule {
	if len(p.rules) == 0 || len(pair.RequestMatcher.Method) == 0 || len(pair.RequestMatcher.Path) == 0 {
		return nil
	}

	method := fmt.Sprintf("%v", pair.RequestMatcher.Method[0].Value)
	path := fmt.Sprintf("%v", pair.RequestMatcher.Path[0].Value)

	p.logger.Debug("Finding matching endpoint",
		"method", method,
		"path", path)

	for _, rule := range p.rules {
		if p.methodMatches(method, rule.Method) && p.pathMatches(path, rule.Path) {
			p.logger.Debug("Found matching endpoint rule",
				"rule_method", rule.Method,
				"rule_path", rule.Path)
			return &rule
		}
	}

	return nil
}

// methodMatches checks if a method matches a rule's method
func (p *StaticEndpointProcessor) methodMatches(method, ruleMethod string) bool {
	if p.caseSensitive {
		return method == ruleMethod
	}
	return strings.EqualFold(method, ruleMethod)
}

// pathMatches checks if a path matches a pattern (including glob patterns)
func (p *StaticEndpointProcessor) pathMatches(path, pattern string) bool {
	// Exact match
	if path == pattern {
		return true
	}

	// Case insensitive match if configured
	if !p.caseSensitive && strings.EqualFold(path, pattern) {
		return true
	}

	// Handle glob patterns
	if strings.Contains(pattern, "*") {
		// Convert glob to regex pattern
		regexPattern := "^" + strings.ReplaceAll(pattern, "*", "[^/]*") + "$"
		re, err := regexp.Compile(regexPattern)
		if err != nil {
			p.logger.Error("Invalid path pattern", "pattern", pattern, "error", err)
			return false
		}
		return re.MatchString(path)
	}

	return false
}

// ApplyRule applies a static response to a pair
func (p *StaticEndpointProcessor) ApplyRule(pair *v2.RequestMatcherResponsePairViewV5, rule *config.EndpointRule) error {
	p.logger.Debug("Applying static response",
		"method", rule.Method,
		"path", rule.Path)

	// Validate JSON format
	var respData map[string]interface{}
	if err := json.Unmarshal([]byte(rule.StaticResponse), &respData); err != nil {
		return fmt.Errorf("invalid JSON in static response: %w", err)
	}

	// Marshal to ensure consistent formatting
	respBody, err := json.Marshal(respData)
	if err != nil {
		return fmt.Errorf("marshaling response data: %w", err)
	}

	// Update the response
	pair.Response.Body = string(respBody)

	if rule.Status > 0 {
		pair.Response.Status = rule.Status
	} else {
		pair.Response.Status = http.StatusOK
	}

	pair.Response.Templated = false

	// Set headers
	if pair.Response.Headers == nil {
		pair.Response.Headers = make(map[string][]string)
	}

	pair.Response.Headers["Content-Type"] = []string{"application/json"}
	pair.Response.Headers["Content-Length"] = []string{fmt.Sprint(len(respBody))}

	// Set matchers to exact
	if len(pair.RequestMatcher.Method) > 0 {
		pair.RequestMatcher.Method[0].Matcher = "exact"
	}

	if len(pair.RequestMatcher.Path) > 0 {
		pair.RequestMatcher.Path[0].Matcher = "exact"
	}

	p.logger.Debug("Static response applied",
		"status", pair.Response.Status,
		"body_length", len(pair.Response.Body))

	return nil
}
