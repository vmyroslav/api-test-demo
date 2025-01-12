package config

type Config interface {
	Routes() []Route
	Patterns() []Pattern
}

type Route interface {
	Path() string
	Methods() []string
	Request() Request
}

type Request interface {
	Headers() map[string][]string
	Body() string
}

type Response interface {
	Headers() map[string][]string
	Body() string
	TemplateMatchedFields() bool
}

type Pattern interface {
	Type() string
	Formats() []string
	Pattern() *string
	ReplaceWith() string
}
