package config

// MatcherType defines how Hoverfly should compare incoming requests with stored patterns.
// see https://docs.hoverfly.io/en/latest/pages/reference/hoverfly/request_matchers.html#request-matchers
type MatcherType string

const (
	ExactMatcher MatcherType = "exact" // Direct string comparison
	RegexMatcher MatcherType = "regex" // Regular expression matching
	GlobMatcher  MatcherType = "glob"  // Glob pattern matching (e.g., "*" wildcards)
)

func (mt MatcherType) String() string {
	return string(mt)
}
