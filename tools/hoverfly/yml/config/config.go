package config

import "net/http"

// Config represents the main configuration interface.
type Config interface {
	Settings() Settings
	Patterns() []Pattern
	Rules() []Rule
}

// Settings represents configuration settings such as debug mode or case sensitivity.
type Settings interface {
	CaseSensitive() bool
	Debug() bool
}

// Pattern represents a general transformation rule.
type Pattern interface {
	Type() string
	Formats() []string
	Pattern() *string
	KeepChars() *int
	ReplaceWith() string
}

// Rule represents an API route and its associated rules.
type Rule interface {
	Path() string
	Methods() []string
	Request() Request
	Response() Response
}

// Request represents the configuration of a request.
type Request interface {
	Body() RequestBody
	Headers() http.Header
	AdditionalConfig() map[string]any
}

// RequestBody represents specific body rules in a request.
type RequestBody interface {
	Fields() map[string]FieldRule
}

// Response represents the configuration of a response.
type Response interface {
	TemplateMatchedFields() bool
}

// FieldRule defines rules for a specific field in a request body.
type FieldRule interface {
	Type() string
	Pattern() *string
	KeepChars() *int
	Formats() []string
	ReplaceWith() string
}
